// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1862(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `ssh-keygen -t ed25519 -f id_host`",
			input:    `ssh-keygen -t ed25519 -f id_host`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `ssh-keygen -lf id_host.pub`",
			input:    `ssh-keygen -lf id_host.pub`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `ssh-keygen -R server.example`",
			input: `ssh-keygen -R server.example`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1862",
					Message: "`ssh-keygen -R server.example` deletes a known-hosts entry — the next `ssh` silently re-trusts whatever key the network returns. Fetch the new fingerprint out-of-band and verify before re-adding.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1862")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
