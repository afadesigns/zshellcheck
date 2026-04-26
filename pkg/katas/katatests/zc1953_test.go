// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1953(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `mount --make-private /sys`",
			input:    `mount --make-private /sys`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `mount -t tmpfs tmpfs /run`",
			input:    `mount -t tmpfs tmpfs /run`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `mount --make-shared /data`",
			input: `mount --make-shared /data`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1953",
					Message: "`mount --make-shared` puts the mount in a shared-subtree group — later bind-mounts propagate to every peer, including containers. Classic escape stepping stone. Use `--make-private` on sensitive paths.",
					Line:    1,
					Column:  8,
				},
			},
		},
		{
			name:  "invalid — `mount /srv --make-rshared`",
			input: `mount /srv --make-rshared`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1953",
					Message: "`mount --make-rshared` puts the mount in a shared-subtree group — later bind-mounts propagate to every peer, including containers. Classic escape stepping stone. Use `--make-private` on sensitive paths.",
					Line:    1,
					Column:  13,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1953")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
