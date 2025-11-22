package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestCheckZC1031(t *testing.T) {
	tests := []struct {
		input    string
		expected []katas.Violation
	}{
		{
			input: `#!/bin/zsh
echo "hello"`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1031",
					Message: "Use `#!/usr/bin/env zsh` for portability instead of `#!/bin/zsh`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			input: `#!/usr/bin/env zsh
echo "hello"`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1031")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
