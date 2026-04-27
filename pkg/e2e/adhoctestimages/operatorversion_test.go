package adhoctestimages

import (
	"github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = ginkgo.Describe("ExtractOperatorVersionFromImage", func() {
	ginkgo.DescribeTable("parsing",
		func(image, wantOperator, wantTag string) {
			gotOperator, gotTag := ExtractOperatorVersionFromImage(image)
			Expect(gotOperator).To(Equal(wantOperator))
			Expect(gotTag).To(Equal(wantTag))
		},
		ginkgo.Entry("full image with tag",
			"quay.io/openshift-qe/rbac-permissions-operator-e2e:v0.1.459-g3fa5c0d",
			"rbac-permissions-operator", "v0.1.459-g3fa5c0d"),
		ginkgo.Entry("image with test suffix",
			"quay.io/org/my-service-test:v1.2.3",
			"my-service", "v1.2.3"),
		ginkgo.Entry("image with latest tag",
			"quay.io/org/my-operator-e2e:latest",
			"my-operator", "latest"),
		ginkgo.Entry("image without tag",
			"quay.io/org/my-operator-e2e",
			"my-operator", ""),
		ginkgo.Entry("simple image",
			"simple:v1",
			"simple", "v1"),
		ginkgo.Entry("empty image",
			"", "", ""),
		ginkgo.Entry("no suffix to strip",
			"quay.io/org/osd-example-operator:v2.0.0",
			"osd-example-operator", "v2.0.0"),
		ginkgo.Entry("harness suffix",
			"quay.io/org/my-operator-harness:abc123",
			"my-operator", "abc123"),
	)
})

var _ = ginkgo.Describe("ExtractCommitFromTag", func() {
	ginkgo.DescribeTable("extraction",
		func(tag, want string) {
			Expect(ExtractCommitFromTag(tag)).To(Equal(want))
		},
		ginkgo.Entry("git describe format", "v0.1.459-g3fa5c0d", "3fa5c0d"),
		ginkgo.Entry("longer commit hash", "v0.1.455-gc8fe2a1", "c8fe2a1"),
		ginkgo.Entry("no commit hash", "v1.2.3", "v1.2.3"),
		ginkgo.Entry("plain commit", "abc1234", "abc1234"),
		ginkgo.Entry("empty", "", ""),
		ginkgo.Entry("latest", "latest", "latest"),
	)
})

var _ = ginkgo.Describe("CSVMatchesVersion", func() {
	ginkgo.DescribeTable("matching",
		func(csvName, expectedTag string, want bool) {
			Expect(CSVMatchesVersion(csvName, expectedTag)).To(Equal(want))
		},
		ginkgo.Entry("exact tag match",
			"rbac-permissions-operator.v0.1.459-g3fa5c0d", "v0.1.459-g3fa5c0d", true),
		ginkgo.Entry("commit hash match",
			"rbac-permissions-operator.v0.1.459-g3fa5c0d", "v0.1.459-g3fa5c0d", true),
		ginkgo.Entry("stale version",
			"rbac-permissions-operator.v0.1.455-gc8fe2a1", "v0.1.459-g3fa5c0d", false),
		ginkgo.Entry("empty csv name",
			"", "v0.1.459-g3fa5c0d", false),
		ginkgo.Entry("empty expected tag",
			"rbac-permissions-operator.v0.1.459-g3fa5c0d", "", false),
		ginkgo.Entry("commit only match",
			"my-operator.v1.0.0-gabcdef1", "v1.0.0-gabcdef1", true),
	)
})
