package adhoctestimages

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-logr/logr"
	operatorsv1alpha1 "github.com/operator-framework/api/pkg/operators/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	runtimeclient "sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	operatorVersionPollInterval = 5 * time.Second
	operatorVersionTimeout      = 10 * time.Minute
)

// ExtractOperatorVersionFromImage parses a container image reference and returns
// the operator name (with common test suffixes stripped) and the version tag.
// Example: "quay.io/org/rbac-permissions-operator-e2e:v0.1.459-g3fa5c0d"
//
//	returns ("rbac-permissions-operator", "v0.1.459-g3fa5c0d")
func ExtractOperatorVersionFromImage(image string) (string, string) {
	if image == "" {
		return "", ""
	}

	var tag string
	if idx := strings.LastIndex(image, ":"); idx != -1 {
		tag = image[idx+1:]
		image = image[:idx]
	}

	if idx := strings.LastIndex(image, "/"); idx != -1 {
		image = image[idx+1:]
	}

	suffixes := []string{"-e2e", "-test", "-tests", "-harness"}
	for _, suffix := range suffixes {
		if strings.HasSuffix(image, suffix) {
			image = strings.TrimSuffix(image, suffix)
			break
		}
	}

	return image, tag
}

// ExtractCommitFromTag extracts the short commit hash from a version tag.
// Tags like "v0.1.459-g3fa5c0d" contain a git describe suffix where "g" prefix
// indicates a git commit hash. Returns the commit hash without the "g" prefix,
// or the full tag if no commit hash is found.
func ExtractCommitFromTag(tag string) string {
	if tag == "" {
		return ""
	}

	parts := strings.Split(tag, "-")
	for _, part := range parts {
		if strings.HasPrefix(part, "g") && len(part) >= 7 {
			return part[1:]
		}
	}

	return tag
}

// CSVMatchesVersion checks if a ClusterServiceVersion name contains the expected
// version tag or commit hash. CSV names follow the pattern:
// "operator-name.v0.1.459-g3fa5c0d"
func CSVMatchesVersion(csvName, expectedTag string) bool {
	if csvName == "" || expectedTag == "" {
		return false
	}

	if strings.Contains(csvName, expectedTag) {
		return true
	}

	commit := ExtractCommitFromTag(expectedTag)
	if commit != expectedTag && commit != "" {
		return strings.Contains(csvName, commit)
	}

	return false
}

// WaitForOperatorVersion polls the cluster until the operator's CSV matches the
// expected version from the e2e image tag, or until the timeout is reached.
func WaitForOperatorVersion(ctx context.Context, logger logr.Logger, restConfig *rest.Config, operatorName, expectedTag string) error {
	if expectedTag == "" || expectedTag == "latest" {
		logger.Info("skipping operator version check — no specific version in image tag", "operator", operatorName, "tag", expectedTag)
		return nil
	}

	logger.Info("waiting for operator version on cluster",
		"operator", operatorName,
		"expectedTag", expectedTag,
		"timeout", operatorVersionTimeout,
	)

	client, err := runtimeclient.New(restConfig, runtimeclient.Options{})
	if err != nil {
		return fmt.Errorf("failed to create runtime client: %w", err)
	}

	kubeClient, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return fmt.Errorf("failed to create kubernetes client: %w", err)
	}

	namespace := findOperatorNamespace(ctx, kubeClient, operatorName)
	if namespace == "" {
		logger.Info("could not determine operator namespace, searching all namespaces", "operator", operatorName)
	}

	var lastCSVNames []string

	err = wait.PollUntilContextTimeout(ctx, operatorVersionPollInterval, operatorVersionTimeout, false, func(ctx context.Context) (bool, error) {
		csvList := &operatorsv1alpha1.ClusterServiceVersionList{}
		listOpts := []runtimeclient.ListOption{}
		if namespace != "" {
			listOpts = append(listOpts, runtimeclient.InNamespace(namespace))
		}

		if err := client.List(ctx, csvList, listOpts...); err != nil {
			logger.Error(err, "failed to list CSVs, retrying")
			return false, nil
		}

		lastCSVNames = nil
		for _, csv := range csvList.Items {
			if !strings.Contains(strings.ToLower(csv.Name), strings.ToLower(operatorName)) {
				continue
			}
			lastCSVNames = append(lastCSVNames, csv.Name)

			if CSVMatchesVersion(csv.Name, expectedTag) && csv.Status.Phase == operatorsv1alpha1.CSVPhaseSucceeded {
				logger.Info("operator version matched",
					"csv", csv.Name,
					"phase", csv.Status.Phase,
					"expectedTag", expectedTag,
				)
				return true, nil
			}
			logger.Info("CSV found but version or phase doesn't match yet",
				"csv", csv.Name,
				"phase", csv.Status.Phase,
				"expectedTag", expectedTag,
			)
		}

		return false, nil
	})
	if err != nil {
		return fmt.Errorf("operator %q did not reach expected version %q within %v (last seen CSVs: %v): %w",
			operatorName, expectedTag, operatorVersionTimeout, lastCSVNames, err)
	}

	return nil
}

// findOperatorNamespace searches for a namespace matching the operator name.
// Operators typically live in namespaces like "openshift-<operator-name>".
func findOperatorNamespace(ctx context.Context, kubeClient kubernetes.Interface, operatorName string) string {
	candidates := []string{
		"openshift-" + operatorName,
		operatorName,
	}

	for _, ns := range candidates {
		if _, err := kubeClient.CoreV1().Namespaces().Get(ctx, ns, metav1.GetOptions{}); err == nil {
			return ns
		}
	}

	return ""
}
