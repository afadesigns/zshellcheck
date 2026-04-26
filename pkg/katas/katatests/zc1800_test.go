// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1800(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `pg_ctl stop -m fast`",
			input:    `pg_ctl stop -m fast`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `pg_ctl start` (no stop)",
			input:    `pg_ctl start`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `pg_ctl stop -m immediate -D /var/lib/pg`",
			input: `pg_ctl stop -m immediate -D /var/lib/pg`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1800",
					Message: "`pg_ctl stop -m immediate` kills the postmaster without a shutdown checkpoint — WAL replay on restart can lose committed transactions if WAL is corrupt. Use `-m smart` or `-m fast` for routine shutdowns.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `pg_ctl restart --mode=immediate`",
			input: `pg_ctl restart --mode=immediate`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1800",
					Message: "`pg_ctl stop -m immediate` kills the postmaster without a shutdown checkpoint — WAL replay on restart can lose committed transactions if WAL is corrupt. Use `-m smart` or `-m fast` for routine shutdowns.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1800")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
