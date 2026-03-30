package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1065(t *testing.T) {
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
			name:     "properly spaced brackets",
			input:    `[ -f /tmp/foo ]`,
			expected: []katas.Violation{},
		},
		{
			name:     "properly spaced double brackets",
			input:    `[[ -f /tmp/foo ]]`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1065")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
