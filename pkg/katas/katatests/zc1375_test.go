package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1375(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — shift 1",
			input:    `shift 1`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — shift $#",
			input: `shift $#`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1375",
					Message: "Use `set --` instead of `shift $#` to clear positional arguments. Clearer intent, no dependency on `$#` accuracy at evaluation time.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1375")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
