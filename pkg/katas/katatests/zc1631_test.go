// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1631(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — passin env:VAR",
			input:    `openssl pkcs12 -in f.p12 -passin env:PASS`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — passin file:path",
			input:    `openssl pkcs12 -in f.p12 -passin file:/run/secrets/p`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — passin stdin",
			input:    `openssl req -passin stdin`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — passin pass:LITERAL",
			input: `openssl pkcs12 -in f.p12 -passin pass:hunter2 -nocerts`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1631",
					Message: "`openssl -passin pass:hunter2` puts the password in argv — visible via `ps`. Use `env:VARNAME`, `file:PATH`, `fd:N`, or `stdin`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — passout pass:X",
			input: `openssl genrsa -passout pass:X 2048`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1631",
					Message: "`openssl -passout pass:X` puts the password in argv — visible via `ps`. Use `env:VARNAME`, `file:PATH`, `fd:N`, or `stdin`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1631")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
