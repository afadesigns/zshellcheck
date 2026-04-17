package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1450(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — pacman -Ss (search)",
			input:    `pacman -Ss vim`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — pacman -S without --noconfirm",
			input: `pacman -S vim`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1450",
					Message: "`pacman -S` without `--noconfirm` hangs in scripts.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — zypper install without -n",
			input: `zypper install vim`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1450",
					Message: "`zypper install` without `--non-interactive` (`-n`) hangs in scripts.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1450")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
