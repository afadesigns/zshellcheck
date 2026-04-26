// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1635(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — mysql -p (prompts)",
			input:    `mysql -u root -p -h db.example.com`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — mysql --login-path",
			input:    `mysql --login-path=prod mydb`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — mysql -psecret",
			input: `mysql -u root -psecret -h db.example.com`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1635",
					Message: "`mysql -psecret` puts the MySQL password in argv. Use `-p` with no arg (prompt), `--login-path`, or a 0600 `~/.my.cnf`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — mysqldump -p$PW",
			input: `mysqldump -u root -p$PW mydb`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1635",
					Message: "`mysqldump -p$PW` puts the MySQL password in argv. Use `-p` with no arg (prompt), `--login-path`, or a 0600 `~/.my.cnf`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1635")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
