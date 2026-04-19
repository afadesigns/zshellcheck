package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1795(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `git remote add origin git@github.com:owner/repo.git`",
			input:    `git remote add origin git@github.com:owner/repo.git`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `git remote set-url origin https://github.com/owner/repo.git`",
			input:    `git remote set-url origin https://github.com/owner/repo.git`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `git remote add origin https://user:ghp_xxx@github.com/owner/repo.git`",
			input: `git remote add origin https://user:ghp_xxx@github.com/owner/repo.git`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1795",
					Message: "`git remote add … https://user:ghp_xxx@github.com/owner/repo.git` stores the token in `.git/config` and leaks it via argv at creation. Use a credential helper, `GIT_ASKPASS`, or an SSH deploy key instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `git remote set-url origin https://u:p@host/repo`",
			input: `git remote set-url origin https://u:p@host/repo`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1795",
					Message: "`git remote set-url … https://u:p@host/repo` stores the token in `.git/config` and leaks it via argv at creation. Use a credential helper, `GIT_ASKPASS`, or an SSH deploy key instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1795")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
