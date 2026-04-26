// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1804(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `aws ec2 describe-instances`",
			input:    `aws ec2 describe-instances`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `aws ec2 terminate-instances --instance-ids i-1 --dry-run`",
			input:    `aws ec2 terminate-instances --instance-ids i-1 --dry-run`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `aws ec2 terminate-instances --instance-ids i-1 i-2`",
			input: `aws ec2 terminate-instances --instance-ids i-1 i-2`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1804",
					Message: "`aws ec2 terminate-instances` tears down EC2 instance(s) and their instance-store volumes with no automatic backup. Review with `aws ec2 describe-…`, add `--dry-run` to verify the target, and pin IDs through `--cli-input-json`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `aws ec2 delete-snapshot --snapshot-id snap-abc`",
			input: `aws ec2 delete-snapshot --snapshot-id snap-abc`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1804",
					Message: "`aws ec2 delete-snapshot` deletes the EBS / RDS snapshot with no automatic backup. Review with `aws ec2 describe-…`, add `--dry-run` to verify the target, and pin IDs through `--cli-input-json`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1804")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
