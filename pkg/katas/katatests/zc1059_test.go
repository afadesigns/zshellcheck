package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1059(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "rm with literal path",
			input:    `rm /tmp/foo`,
			expected: []katas.Violation{},
		},
		{
			name:     "not rm command",
			input:    `echo hello`,
			expected: []katas.Violation{},
		},
		{
			name:     "rm with no arguments",
			input:    `rm`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1059")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
