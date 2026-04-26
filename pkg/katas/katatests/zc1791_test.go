// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1791(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `curl https://example.com`",
			input:    `curl https://example.com`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `curl --unix-socket /run/user/1000/bus http://localhost/` (dbus)",
			input:    `curl --unix-socket /run/user/1000/bus http://localhost/`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `curl --unix-socket /var/run/docker.sock http://localhost/containers/json`",
			input: `curl --unix-socket /var/run/docker.sock http://localhost/containers/json`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1791",
					Message: "`curl --unix-socket /var/run/docker.sock` speaks the container-daemon API — a `POST /containers/create` with `Privileged=true` is a host-root primitive. Use the CLI (`docker`/`podman`) instead.",
					Line:    1,
					Column:  7,
				},
			},
		},
		{
			name:  "invalid — `curl http://localhost/v1/services --unix-socket /run/containerd/containerd.sock` (trailing form)",
			input: `curl http://localhost/v1/services --unix-socket /run/containerd/containerd.sock`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1791",
					Message: "`curl --unix-socket /run/containerd/containerd.sock` speaks the container-daemon API — a `POST /containers/create` with `Privileged=true` is a host-root primitive. Use the CLI (`docker`/`podman`) instead.",
					Line:    1,
					Column:  36,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1791")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
