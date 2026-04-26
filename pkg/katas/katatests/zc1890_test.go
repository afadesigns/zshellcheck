// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1890(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `kadmin -p admin/admin -k -t /etc/krb5.keytab`",
			input:    `kadmin -p admin/admin -k -t /etc/krb5.keytab`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `kinit admin/admin`",
			input:    `kinit admin/admin`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `kadmin -p admin/admin -w hunter2`",
			input: `kadmin -p admin/admin -w hunter2`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1890",
					Message: "`kadmin -w hunter2` embeds the Kerberos admin password in argv — visible to `ps`, `/proc`, shell history. Use `-k -t /etc/krb5.keytab` (keytab root-only) or pipe the password on stdin.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `kadmin.local -w hunter2 addprinc user`",
			input: `kadmin.local -w hunter2 addprinc user`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1890",
					Message: "`kadmin.local -w hunter2` embeds the Kerberos admin password in argv — visible to `ps`, `/proc`, shell history. Use `-k -t /etc/krb5.keytab` (keytab root-only) or pipe the password on stdin.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1890")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
