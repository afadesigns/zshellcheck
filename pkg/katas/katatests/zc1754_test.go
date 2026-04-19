package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1754(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `gh auth status`",
			input:    `gh auth status`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `gh auth token`",
			input:    `gh auth token`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `gh auth status -t`",
			input: `gh auth status -t`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1754",
					Message: "`gh auth status -t` prints the OAuth token in the status output — CI logs and scrollback become a repo-admin leak. Use `gh auth token` in automation so the secret path is explicit.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `gh auth status --show-token`",
			input: `gh auth status --show-token`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1754",
					Message: "`gh auth status --show-token` prints the OAuth token in the status output — CI logs and scrollback become a repo-admin leak. Use `gh auth token` in automation so the secret path is explicit.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1754")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
