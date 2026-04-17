package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1499(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — docker pull nginx:1.27",
			input:    `docker pull nginx:1.27`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — docker pull nginx@sha256:abc",
			input:    `docker pull nginx@sha256:abcdef`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — docker pull nginx (no tag)",
			input: `docker pull nginx`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1499",
					Message: "`nginx` is unpinned (implicit `:latest`). Pin to a specific tag or an immutable `@sha256:` digest for reproducibility.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — docker pull nginx:latest",
			input: `docker pull nginx:latest`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1499",
					Message: "`nginx:latest` is unpinned (implicit `:latest`). Pin to a specific tag or an immutable `@sha256:` digest for reproducibility.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1499")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
