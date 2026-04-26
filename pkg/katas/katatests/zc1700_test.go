// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1700(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — ldapsearch -W (prompt)",
			input:    `ldapsearch -x -D cn=admin -W -b dc=example`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — ldapsearch -y FILE",
			input:    `ldapsearch -x -D cn=admin -y /etc/ldap.password`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — ldapsearch -w SECRET",
			input: `ldapsearch -x -D cn=admin -w SECRET -b dc=example`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1700",
					Message: "`ldapsearch -w PASSWORD` leaks the LDAP bind password into `ps` / `/proc/PID/cmdline` — use `-W` to prompt, `-y FILE` for a mode-0400 secret file, or SASL (`-Y GSSAPI`).",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — ldapmodify -w SECRET",
			input: `ldapmodify -w SECRET -f change.ldif`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1700",
					Message: "`ldapmodify -w PASSWORD` leaks the LDAP bind password into `ps` / `/proc/PID/cmdline` — use `-W` to prompt, `-y FILE` for a mode-0400 secret file, or SASL (`-Y GSSAPI`).",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1700")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
