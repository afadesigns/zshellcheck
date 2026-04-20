package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1960(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `az vm list`",
			input:    `az vm list`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `aws ssm describe-instance-information`",
			input:    `aws ssm describe-instance-information`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `az vm run-command invoke -g rg -n vm --command-id RunShellScript --scripts $CMD`",
			input: `az vm run-command invoke -g rg -n vm --command-id RunShellScript --scripts $CMD`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1960",
					Message: "`az vm run-command invoke` runs arbitrary shell on the VM via the cloud control plane — operator-composed command strings become IAM-driven RCE. Pin to a reviewed asset, template-escape input, require MFA.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `aws ssm send-command --document-name AWS-RunShellScript --parameters commands=$CMD`",
			input: `aws ssm send-command --document-name AWS-RunShellScript --parameters commands=$CMD`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1960",
					Message: "`aws ssm send-command` runs arbitrary shell on the VM via the cloud control plane — operator-composed command strings become IAM-driven RCE. Pin to a reviewed asset, template-escape input, require MFA.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1960")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
