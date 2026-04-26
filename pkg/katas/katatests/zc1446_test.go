// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1446(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — aws s3 ls",
			input:    `aws s3 ls`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — aws s3 cp",
			input:    `aws s3 cp local remote`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — aws s3 rm --recursive",
			input: `aws s3 rm prefix --recursive`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1446",
					Message: "`aws s3 rm --recursive` / `s3 rb --force` mass-deletes objects/buckets. Enable versioning and dry-run with `--dryrun`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1446")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
