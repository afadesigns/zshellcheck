package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1690(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — pip install normal package",
			input:    `pip install requests`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — pip install git+URL@commit-hash",
			input:    `pip install git+https://github.com/org/repo@abc1234`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — pip install git+URL@v1.2.3 tag",
			input:    `pip install git+https://github.com/org/repo@v1.2.3`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — pip install git+URL without ref",
			input: `pip install git+https://github.com/org/repo`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1690",
					Message: "`pip install git+https://github.com/org/repo` tracks a moving git ref — pin to a commit SHA (`@abc1234…`) or signed tag (`@v1.2.3`), or use the PyPI release.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — pip install git+URL@main (moving branch)",
			input: `pip install git+https://github.com/org/repo@main`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1690",
					Message: "`pip install git+https://github.com/org/repo@main` tracks a moving git ref — pin to a commit SHA (`@abc1234…`) or signed tag (`@v1.2.3`), or use the PyPI release.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1690")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
