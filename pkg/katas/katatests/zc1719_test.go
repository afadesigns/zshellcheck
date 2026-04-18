package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1719(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `git filter-repo`",
			input:    `git filter-repo --path secret.txt --invert-paths`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `git rebase`",
			input:    `git rebase main`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `git filter-branch --tree-filter`",
			input: `git filter-branch --tree-filter rm secret.txt HEAD`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1719",
					Message: "`git filter-branch` is deprecated (Git 2.24+) and its manpage redirects to `git filter-repo`. Use that instead — faster, safer defaults, no orphaned objects.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — bare `git filter-branch`",
			input: `git filter-branch`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1719",
					Message: "`git filter-branch` is deprecated (Git 2.24+) and its manpage redirects to `git filter-repo`. Use that instead — faster, safer defaults, no orphaned objects.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1719")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
