// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1268(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid du with specific path",
			input:    `du -sh /tmp`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid du -sh *",
			input: `du -sh *`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1268",
					Message: "Use `du -sh -- *` instead of `du -sh *`. The `--` prevents filenames starting with `-` from being interpreted as options.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1268")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
