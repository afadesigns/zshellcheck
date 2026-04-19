package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1758(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `gh codespace delete -c mycodespace` (prompt kept)",
			input:    `gh codespace delete -c mycodespace`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `gh codespace list`",
			input:    `gh codespace list`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `gh codespace delete --force`",
			input: `gh codespace delete --force`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1758",
					Message: "`gh codespace delete --force` skips the prompt and drops uncommitted work along with the codespace. Let the prompt list what's about to go and verify `git status` inside first.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `gh codespace delete -f --all`",
			input: `gh codespace delete -f --all`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1758",
					Message: "`gh codespace delete -f` skips the prompt and drops uncommitted work along with the codespace. Let the prompt list what's about to go and verify `git status` inside first.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1758")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
