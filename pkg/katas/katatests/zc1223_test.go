package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1223(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid ip -br addr",
			input:    `ip -br addr`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid ip addr show",
			input: `ip addr show`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1223",
					Message: "Use `ip -br addr` for machine-readable output instead of parsing `ip addr show` with grep or awk.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1223")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
