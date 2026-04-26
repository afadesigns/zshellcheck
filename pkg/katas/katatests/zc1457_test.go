// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1457(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — docker run normal mount",
			input:    `docker run -v data:/app alpine`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — docker run with docker.sock mount",
			input: `docker run -v docker.sock:/var/run/docker.sock alpine`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1457",
					Message: "Mounting `/var/run/docker.sock` gives the container effective root on the host. Reserve for trusted tooling.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1457")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
