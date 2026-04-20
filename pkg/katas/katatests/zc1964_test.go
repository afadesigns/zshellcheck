package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1964(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `uvx ruff@0.5.7 check`",
			input:    `uvx ruff@0.5.7 check`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `pipx run 'black==24.8.0' --version`",
			input:    `pipx run 'black==24.8.0' --version`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `uvx ruff check`",
			input: `uvx ruff check`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1964",
					Message: "`uvx ruff` resolves to the PyPI `latest` release — a squatted name or compromised maintainer lands untested code. Pin `pkg==X.Y.Z` (or `pkg@X.Y.Z` for uv).",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `pipx run black .`",
			input: `pipx run black .`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1964",
					Message: "`pipx run black` resolves to the PyPI `latest` release — a squatted name or compromised maintainer lands untested code. Pin `pkg==X.Y.Z` (or `pkg@X.Y.Z` for uv).",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1964")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
