package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1077(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid expansion",
			input:    `upper=${var:u}`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid tr upper",
			input: `upper=$(echo $var | tr 'a-z' 'A-Z')`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1077",
					Message: "Use `${var:u}` instead of `tr` for uppercase conversion. It is faster and built-in.",
					Line:    1,
					Column:  7,
				},
			},
		},
		{
			name:  "invalid tr lower",
			input: `lower=$(echo $var | tr 'A-Z' 'a-z')`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1077",
					Message: "Use `${var:l}` instead of `tr` for lowercase conversion. It is faster and built-in.",
					Line:    1,
					Column:  7,
				},
			},
		},
		{
			name:  "invalid tr upper posix",
			input: `upper=$(echo $var | tr '[:lower:]' '[:upper:]')`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1077",
					Message: "Use `${var:u}` instead of `tr` for uppercase conversion. It is faster and built-in.",
					Line:    1,
					Column:  7,
				},
			},
		},
		{
			name:     "other tr usage",
			input:    `clean=$(echo $var | tr -d '\n')`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1077")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
