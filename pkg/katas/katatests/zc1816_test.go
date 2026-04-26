// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1816(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `docker build -t myimage .`",
			input:    `docker build -t myimage .`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `podman ps`",
			input:    `podman ps`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `docker commit mycontainer myimage:latest`",
			input: `docker commit mycontainer myimage:latest`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1816",
					Message: "`docker commit` snapshots a running container — no Dockerfile trail, runtime env / `/tmp` scratch / shell history get baked in, and the layer metadata does not record what was installed. Build from a `Dockerfile` instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `podman commit web web:snap`",
			input: `podman commit web web:snap`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1816",
					Message: "`podman commit` snapshots a running container — no Dockerfile trail, runtime env / `/tmp` scratch / shell history get baked in, and the layer metadata does not record what was installed. Build from a `Dockerfile` instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1816")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
