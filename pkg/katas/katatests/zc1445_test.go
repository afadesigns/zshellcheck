// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1445(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — mysqladmin ping",
			input:    `mysqladmin ping`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — dropdb",
			input: `dropdb mydb`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1445",
					Message: "`dropdb` removes a PostgreSQL database. Verify target and backup first (`pg_dump`).",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — mysqladmin drop",
			input: `mysqladmin drop mydb`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1445",
					Message: "`mysqladmin drop` removes a MySQL database. Verify target and backup first (`mysqldump`).",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1445")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
