// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1541(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — apk add curl",
			input:    `apk add curl`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — apk add --allow-untrusted local.apk",
			input: `apk add --allow-untrusted ./local.apk`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1541",
					Message: "`apk --allow-untrusted` skips signature verification on the package — MITM-to-root on Alpine. Sign and place key in /etc/apk/keys/.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1541")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
