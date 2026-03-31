package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1230(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid ping -c",
			input:    `ping -c 3 localhost`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid ping without -c",
			input: `ping localhost`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1230",
					Message: "Use `ping -c N` in scripts. Without `-c`, ping runs indefinitely on Linux and will hang the script.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1230")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
