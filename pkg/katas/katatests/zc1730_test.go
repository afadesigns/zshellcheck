package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1730(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `brew install foo` (stable release)",
			input:    `brew install foo`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `brew list --HEAD` (not install/upgrade)",
			input:    `brew list --HEAD`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `brew install --HEAD foo`",
			input: `brew install --HEAD foo`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1730",
					Message: "`brew install --HEAD` builds from upstream HEAD — every run pulls a different commit. Pin to a stable formula release or vendor a private tap with a fixed revision.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `brew reinstall --HEAD foo`",
			input: `brew reinstall --HEAD foo`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1730",
					Message: "`brew reinstall --HEAD` builds from upstream HEAD — every run pulls a different commit. Pin to a stable formula release or vendor a private tap with a fixed revision.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1730")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
