package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1525(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — ping -c 4 host",
			input:    `ping -c 4 example.com`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — ping -f host",
			input: `ping -f example.com`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1525",
					Message: "`ping -f` (flood) bypasses the rate limit and can saturate slow links — DoS primitive. Scope tightly and document.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — ping6 -f host",
			input: `ping6 -f 2001:db8::1`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1525",
					Message: "`ping6 -f` (flood) bypasses the rate limit and can saturate slow links — DoS primitive. Scope tightly and document.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1525")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
