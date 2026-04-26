// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1253(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid docker build --no-cache",
			input:    `docker build --no-cache -t myapp .`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid docker run (not build)",
			input:    `docker run --rm alpine`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid docker build without --no-cache",
			input: `docker build -t myapp .`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1253",
					Message: "Consider `docker build --no-cache` in CI for reproducible builds. Layer caching can mask changed dependencies.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1253")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
