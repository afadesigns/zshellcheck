// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1988(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `nsupdate -k $KEYFILE`",
			input:    `nsupdate -k $KEYFILE`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `nsupdate -v`",
			input:    `nsupdate -v`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `nsupdate -y HMAC-SHA256:name:c2VjcmV0`",
			input: `nsupdate -y HMAC-SHA256:name:c2VjcmV0`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1988",
					Message: "`nsupdate -y …` puts the TSIG key in argv — `ps`, `/proc/*/cmdline`, and `$HISTFILE` all capture it. Use `nsupdate -k $KEYFILE` with a `0600` keyfile instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1988")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
