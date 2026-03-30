package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1056(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "normal command",
			input:    `echo hello`,
			expected: []katas.Violation{},
		},
		{
			name:     "command substitution as statement",
			input:    `$(echo hello)`,
			expected: []katas.Violation{},
		},
		{
			name:     "arithmetic command",
			input:    `((x + 1))`,
			expected: []katas.Violation{},
		},
		{
			name:     "non SimpleCommand node",
			input:    `for i in a b; do echo $i; done`,
			expected: []katas.Violation{},
		},
		{
			name:     "simple command with identifier name",
			input:    `ls -la`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1056")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
