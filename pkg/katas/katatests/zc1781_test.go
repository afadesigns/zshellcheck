package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1781(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `git clone https://github.com/owner/repo.git`",
			input:    `git clone https://github.com/owner/repo.git`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `git clone git@github.com:owner/repo.git` (SSH)",
			input:    `git clone git@github.com:owner/repo.git`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `git clone https://token@github.com/owner/repo.git` (no password segment)",
			input:    `git clone https://token@github.com/owner/repo.git`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `git clone https://user:ghp_xxx@github.com/owner/repo.git`",
			input: `git clone https://user:ghp_xxx@github.com/owner/repo.git`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1781",
					Message: "`git clone https://user:ghp_xxx@github.com/owner/repo.git` puts the token in argv and `.git/config`. Use a credential helper, `GIT_ASKPASS`, or `http.extraHeader` with an env-sourced bearer.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `git fetch https://u:p@host/repo`",
			input: `git fetch https://u:p@host/repo`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1781",
					Message: "`git fetch https://u:p@host/repo` puts the token in argv and `.git/config`. Use a credential helper, `GIT_ASKPASS`, or `http.extraHeader` with an env-sourced bearer.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1781")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
