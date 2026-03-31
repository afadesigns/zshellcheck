package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1216(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid dig",
			input:    `dig example.com`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid nslookup",
			input: `nslookup example.com`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1216",
					Message: "Avoid `nslookup` — it is deprecated on many systems. Use `dig` for detailed DNS queries or `host` for simple lookups.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1216")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
