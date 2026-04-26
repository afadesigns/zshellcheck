// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1981(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `exec $BIN`",
			input:    `exec $BIN`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `exec $BIN arg1 arg2`",
			input:    `exec $BIN arg1 arg2`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `exec -a login $BIN`",
			input: `exec -a login $BIN`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1981",
					Message: "`exec -a NAME` sets `argv[0]` to `NAME` — `ps`/`top`/audit rules see the alias, not the real binary. Keep out of production scripts unless the alias is documented.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `exec -a $ALIAS $BIN`",
			input: `exec -a $ALIAS $BIN`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1981",
					Message: "`exec -a NAME` sets `argv[0]` to `NAME` — `ps`/`top`/audit rules see the alias, not the real binary. Keep out of production scripts unless the alias is documented.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1981")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
