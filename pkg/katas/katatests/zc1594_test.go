// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1594(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — plain docker run",
			input:    `docker run --rm alpine`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — podman run with different security-opt",
			input:    `podman run --security-opt=no-new-privileges alpine`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — docker run --security-opt=systempaths=unconfined",
			input: `docker run --security-opt=systempaths=unconfined alpine`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1594",
					Message: "`docker run --security-opt=systempaths=unconfined` unhides `/proc/sys`, `/proc/sysrq-trigger`, and other kernel knobs. A compromise in the container can then panic or re-tune the host.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — podman create systempaths=unconfined",
			input: `podman create --security-opt=systempaths=unconfined alpine`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1594",
					Message: "`podman run --security-opt=systempaths=unconfined` unhides `/proc/sys`, `/proc/sysrq-trigger`, and other kernel knobs. A compromise in the container can then panic or re-tune the host.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1594")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
