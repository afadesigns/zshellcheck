package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1361(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — awk with generic program",
			input:    `awk '{print $1}' file`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — awk NR==5",
			input: `awk 'NR==5' file`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1361",
					Message: "Avoid `awk 'NR==N'` — split with `${(f)\"$(<file)\"}` in Zsh and index: `lines=(${(f)\"$(<file)\"}); print $lines[N]`. No awk process needed.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1361")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
