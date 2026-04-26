// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1776(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `psql postgresql://user@host/db` (no password)",
			input:    `psql postgresql://user@host/db`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `psql postgresql://host:5432/db` (port, not password)",
			input:    `psql postgresql://host:5432/db`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `psql $PG_URL`",
			input:    `psql $PG_URL`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `psql postgresql://user:hunter2@host/db`",
			input: `psql postgresql://user:hunter2@host/db`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1776",
					Message: "`postgresql://user:SECRET@…` in argv puts the password in `ps` / `/proc/PID/cmdline` / history. Use a password file (`~/.pgpass`, `~/.my.cnf`), `PGPASSWORD` / `REDISCLI_AUTH` env var, or build the URI from a secret variable.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `mongosh mongodb+srv://u:p@cluster/db`",
			input: `mongosh "mongodb+srv://u:p@cluster/db"`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1776",
					Message: "`mongodb+srv://user:SECRET@…` in argv puts the password in `ps` / `/proc/PID/cmdline` / history. Use a password file (`~/.pgpass`, `~/.my.cnf`), `PGPASSWORD` / `REDISCLI_AUTH` env var, or build the URI from a secret variable.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1776")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
