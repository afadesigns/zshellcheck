package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1154(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid find without exec",
			input:    `find . -name "*.go"`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid find -print0",
			input:    `find . -print0`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1154")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
