// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1238(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid docker exec without -it",
			input:    `docker exec mycontainer ls`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid docker exec -it",
			input: `docker exec -it mycontainer bash`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1238",
					Message: "Avoid `docker exec -it` in scripts — TTY allocation hangs without a terminal. Use `docker exec` without `-it` for non-interactive commands.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1238")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
