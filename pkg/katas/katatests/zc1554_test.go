// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1554(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — unzip without -o",
			input:    `unzip file.zip`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — tar xf without --overwrite",
			input:    `tar xf foo.tar`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — unzip -o file.zip",
			input: `unzip -o file.zip`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1554",
					Message: "`unzip -o` overwrites existing files without prompting. Extract to a staging directory, diff, then move.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — tar xf foo.tar --overwrite",
			input: `tar xf foo.tar --overwrite`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1554",
					Message: "`tar --overwrite` discards existing files during extract. Use a staging directory and diff before rolling forward.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1554")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
