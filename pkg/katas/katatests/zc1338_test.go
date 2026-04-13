package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1338(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid seq without -s",
			input:    `seq 10`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid seq -s usage",
			input: `seq -s , 10`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1338",
					Message: "Avoid `seq -s` in Zsh — use `${(j:sep:)array}` with brace expansion for joined sequences.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1338")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
