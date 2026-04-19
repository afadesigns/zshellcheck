package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1806(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `zmv -n '*.JPG' '*.jpg'` (dry-run)",
			input:    `zmv -n '*.JPG' '*.jpg'`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `zmv -i '(*).txt' '$1.md'` (interactive)",
			input:    `zmv -i '(*).txt' '$1.md'`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `zmv` alone (help)",
			input:    `zmv`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `zmv '*.txt' '*.md'`",
			input: `zmv '*.txt' '*.md'`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1806",
					Message: "`zmv` without `-n` (dry-run) or `-i` (interactive) renames every matched file in one shot — a pattern typo can collide names. Preview with `zmv -n`, then re-run once the list looks right.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `zmv -W '(*).jpg' 'archive/$1.jpg'`",
			input: `zmv -W '(*).jpg' 'archive/$1.jpg'`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1806",
					Message: "`zmv` without `-n` (dry-run) or `-i` (interactive) renames every matched file in one shot — a pattern typo can collide names. Preview with `zmv -n`, then re-run once the list looks right.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1806")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
