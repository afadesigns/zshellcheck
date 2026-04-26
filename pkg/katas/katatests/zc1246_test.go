// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1246(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid mysql with -p prompt",
			input:    `mysql -u root -p mydb`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid mysql with inline password",
			input: `mysql -u root -pMySecret mydb`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1246",
					Message: "Avoid passing passwords as command arguments — they appear in process lists. Use environment variables (e.g., `MYSQL_PWD`) or credential files instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1246")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
