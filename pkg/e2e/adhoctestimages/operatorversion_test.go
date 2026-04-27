package adhoctestimages

import (
	"testing"
)

func TestExtractOperatorVersionFromImage(t *testing.T) {
	tests := []struct {
		name         string
		image        string
		wantOperator string
		wantTag      string
	}{
		{
			name:         "full image with tag",
			image:        "quay.io/openshift-qe/rbac-permissions-operator-e2e:v0.1.459-g3fa5c0d",
			wantOperator: "rbac-permissions-operator",
			wantTag:      "v0.1.459-g3fa5c0d",
		},
		{
			name:         "image with test suffix",
			image:        "quay.io/org/my-service-test:v1.2.3",
			wantOperator: "my-service",
			wantTag:      "v1.2.3",
		},
		{
			name:         "image with latest tag",
			image:        "quay.io/org/my-operator-e2e:latest",
			wantOperator: "my-operator",
			wantTag:      "latest",
		},
		{
			name:         "image without tag",
			image:        "quay.io/org/my-operator-e2e",
			wantOperator: "my-operator",
			wantTag:      "",
		},
		{
			name:         "simple image",
			image:        "simple:v1",
			wantOperator: "simple",
			wantTag:      "v1",
		},
		{
			name:         "empty image",
			image:        "",
			wantOperator: "",
			wantTag:      "",
		},
		{
			name:         "no suffix to strip",
			image:        "quay.io/org/osd-example-operator:v2.0.0",
			wantOperator: "osd-example-operator",
			wantTag:      "v2.0.0",
		},
		{
			name:         "harness suffix",
			image:        "quay.io/org/my-operator-harness:abc123",
			wantOperator: "my-operator",
			wantTag:      "abc123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOperator, gotTag := ExtractOperatorVersionFromImage(tt.image)
			if gotOperator != tt.wantOperator {
				t.Errorf("operator = %q, want %q", gotOperator, tt.wantOperator)
			}
			if gotTag != tt.wantTag {
				t.Errorf("tag = %q, want %q", gotTag, tt.wantTag)
			}
		})
	}
}

func TestExtractCommitFromTag(t *testing.T) {
	tests := []struct {
		name string
		tag  string
		want string
	}{
		{
			name: "git describe format",
			tag:  "v0.1.459-g3fa5c0d",
			want: "3fa5c0d",
		},
		{
			name: "longer commit hash",
			tag:  "v0.1.455-gc8fe2a1",
			want: "c8fe2a1",
		},
		{
			name: "no commit hash",
			tag:  "v1.2.3",
			want: "v1.2.3",
		},
		{
			name: "plain commit",
			tag:  "abc1234",
			want: "abc1234",
		},
		{
			name: "empty",
			tag:  "",
			want: "",
		},
		{
			name: "latest",
			tag:  "latest",
			want: "latest",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExtractCommitFromTag(tt.tag)
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestCSVMatchesVersion(t *testing.T) {
	tests := []struct {
		name        string
		csvName     string
		expectedTag string
		want        bool
	}{
		{
			name:        "exact tag match",
			csvName:     "rbac-permissions-operator.v0.1.459-g3fa5c0d",
			expectedTag: "v0.1.459-g3fa5c0d",
			want:        true,
		},
		{
			name:        "commit hash match",
			csvName:     "rbac-permissions-operator.v0.1.459-g3fa5c0d",
			expectedTag: "v0.1.459-g3fa5c0d",
			want:        true,
		},
		{
			name:        "stale version",
			csvName:     "rbac-permissions-operator.v0.1.455-gc8fe2a1",
			expectedTag: "v0.1.459-g3fa5c0d",
			want:        false,
		},
		{
			name:        "empty csv name",
			csvName:     "",
			expectedTag: "v0.1.459-g3fa5c0d",
			want:        false,
		},
		{
			name:        "empty expected tag",
			csvName:     "rbac-permissions-operator.v0.1.459-g3fa5c0d",
			expectedTag: "",
			want:        false,
		},
		{
			name:        "commit only match",
			csvName:     "my-operator.v1.0.0-gabcdef1",
			expectedTag: "v1.0.0-gabcdef1",
			want:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CSVMatchesVersion(tt.csvName, tt.expectedTag)
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}
