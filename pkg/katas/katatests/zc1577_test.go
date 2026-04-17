package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1577(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — dig A example.com",
			input:    `dig A example.com`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — dig MX example.com",
			input:    `dig MX example.com`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — dig ANY example.com",
			input: `dig ANY example.com`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1577",
					Message: "`dig ... ANY` is RFC 8482-deprecated — filtered by recursors. Query specific types (A / MX / NS / …) and combine.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — dig example.com ANY",
			input: `dig example.com ANY`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1577",
					Message: "`dig ... ANY` is RFC 8482-deprecated — filtered by recursors. Query specific types (A / MX / NS / …) and combine.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1577")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
