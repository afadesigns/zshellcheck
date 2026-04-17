package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1565(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — command -v cmd",
			input:    `command -v git`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — whereis git",
			input: `whereis git`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1565",
					Message: "`whereis` is index-based and stale-prone. Use `command -v <cmd>` for runtime existence checks.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — locate foo",
			input: `locate foo`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1565",
					Message: "`locate` is index-based and stale-prone. Use `command -v <cmd>` for runtime existence checks.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1565")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
