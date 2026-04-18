package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1726(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `gcloud projects delete PROJECT_ID` (no --quiet)",
			input:    `gcloud projects delete PROJECT_ID`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `gcloud projects list --quiet` (not delete)",
			input:    `gcloud projects list --quiet`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `gcloud projects delete PROJECT_ID --quiet`",
			input: `gcloud projects delete PROJECT_ID --quiet`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1726",
					Message: "`gcloud ... delete --quiet` skips confirmation — a wrong argument wipes the resource (compute disks, secrets, BigQuery tables have no soft-delete). Drop `--quiet` or destroy through a Terraform plan with review.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `gcloud sql instances delete INSTANCE -q`",
			input: `gcloud sql instances delete INSTANCE -q`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1726",
					Message: "`gcloud ... delete --quiet` skips confirmation — a wrong argument wipes the resource (compute disks, secrets, BigQuery tables have no soft-delete). Drop `--quiet` or destroy through a Terraform plan with review.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1726")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
