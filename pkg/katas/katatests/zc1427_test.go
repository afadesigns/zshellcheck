// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1427(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — nc -l listener without -e",
			input:    `nc -l 1234`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — nc -e",
			input: `nc -e sh 127.0.0.1 1234`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1427",
					Message: "`nc -e` spawns an arbitrary command for each connection — reverse-shell territory. Remove from scripts unless audited and restricted.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1427")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
