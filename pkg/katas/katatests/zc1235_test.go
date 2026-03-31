package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1235(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid git push",
			input:    `git push origin main`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid git push -f",
			input: `git push -f origin main`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1235",
					Message: "Use `git push --force-with-lease` instead of `-f`/`--force`. It prevents overwriting remote changes made by others.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1235")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
