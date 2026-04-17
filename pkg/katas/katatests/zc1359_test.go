package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1359(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — id without group flag",
			input:    `id -u`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — id -Gn",
			input: `id -Gn`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1359",
					Message: "Avoid `id -Gn`/`-G`/`-gn`/`-g` — use Zsh `$groups` (names→gids assoc array) or `$GID` for the primary group after `zmodload zsh/parameter`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — id -g",
			input: `id -g`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1359",
					Message: "Avoid `id -Gn`/`-G`/`-gn`/`-g` — use Zsh `$groups` (names→gids assoc array) or `$GID` for the primary group after `zmodload zsh/parameter`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1359")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
