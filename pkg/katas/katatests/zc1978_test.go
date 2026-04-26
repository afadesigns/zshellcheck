// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1978(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `sftp $HOST`",
			input:    `sftp $HOST`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `curl -u user: https://$HOST/file`",
			input:    `curl -u user: https://$HOST/file`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `ftp $HOST` (owned by ZC1200)",
			input:    `ftp $HOST`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `tftp $HOST`",
			input: `tftp $HOST`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1978",
					Message: "`tftp` transfers over plaintext UDP/69 with no authentication — capture the payload, or push a crafted file under the expected name. Use a signed-payload `curl` over HTTPS and verify the signature before use.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1978")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
