// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1200(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid sftp",
			input:    `sftp user@host`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid ftp",
			input: `ftp server.example.com`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1200",
					Message: "Avoid `ftp` — it transmits credentials in plain text. Use `sftp`, `scp`, or `curl` with HTTPS for secure file transfers.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1200")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
