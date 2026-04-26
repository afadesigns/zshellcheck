// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1569(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — nvme list",
			input:    `nvme list`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — nvme id-ctrl $DISK",
			input:    `nvme id-ctrl $DISK`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — nvme format -s1 $DISK",
			input: `nvme format -s1 $DISK`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1569",
					Message: "`nvme format -s1` unrecoverably erases the namespace in seconds. Do not run from automation.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — nvme format -s2 $DISK",
			input: `nvme format -s2 $DISK`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1569",
					Message: "`nvme format -s2` unrecoverably erases the namespace in seconds. Do not run from automation.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — nvme sanitize -a 4 $DISK",
			input: `nvme sanitize -a 4 $DISK`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1569",
					Message: "`nvme sanitize -a` unrecoverably erases the namespace in seconds. Do not run from automation.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1569")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
