package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1054(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "no range pattern",
			input:    `grep foo bar`,
			expected: []katas.Violation{},
		},
		{
			name:     "command with no args",
			input:    `ls`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1054")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
