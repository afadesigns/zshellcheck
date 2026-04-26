// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1896(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `docker run -v /etc/app:/app/etc ubuntu`",
			input:    `docker run -v /etc/app:/app/etc ubuntu`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `docker run -v /home/user:/work ubuntu`",
			input:    `docker run -v /home/user:/work ubuntu`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `docker run -v /proc:/host/proc:ro ubuntu`",
			input: `docker run -v /proc:/host/proc:ro ubuntu`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1896",
					Message: "`docker ... -v /proc:/host/proc:ro` bind-mounts host /proc into the container — every process's `environ`/`cmdline` and `/proc/1/ns/` breakout handles become readable. Use `--cap-add=SYS_PTRACE` or host-side monitoring instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `podman run --volume=/sys:/host/sys alpine`",
			input: `podman run --volume=/sys:/host/sys alpine`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1896",
					Message: "`podman ... -v /sys:/host/sys` bind-mounts host /sys into the container — every process's `environ`/`cmdline` and `/proc/1/ns/` breakout handles become readable. Use `--cap-add=SYS_PTRACE` or host-side monitoring instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1896")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
