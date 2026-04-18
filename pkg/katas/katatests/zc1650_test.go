package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1650(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — setopt NO_NOMATCH",
			input:    `setopt NO_NOMATCH`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — unsetopt BEEP",
			input:    `unsetopt BEEP`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — setopt RM_STAR_SILENT",
			input: `setopt RM_STAR_SILENT`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1650",
					Message: "`setopt RM_STAR_SILENT` removes the `rm *` confirmation prompt — keep the default `RM_STAR_WAIT` so accidental deletions pause before they happen.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — unsetopt rmstarwait (lowercase, no underscore)",
			input: `unsetopt rmstarwait`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1650",
					Message: "`unsetopt rmstarwait` removes the `rm *` confirmation prompt — keep the default `RM_STAR_WAIT` so accidental deletions pause before they happen.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1650")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
