package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1757(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `gh auth refresh --scopes repo`",
			input:    `gh auth refresh --scopes repo`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `gh auth login --scopes workflow,read:org`",
			input:    `gh auth login --scopes workflow,read:org`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `gh auth refresh --scopes delete_repo`",
			input: `gh auth refresh --scopes delete_repo`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1757",
					Message: "`gh auth refresh --scopes delete_repo` escalates the token to destructive privileges that outlast the script. Request the minimum scope (`repo`, `workflow`) and rotate the token when done.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `gh auth refresh -s admin:org`",
			input: `gh auth refresh -s admin:org`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1757",
					Message: "`gh auth refresh --scopes admin:org` escalates the token to destructive privileges that outlast the script. Request the minimum scope (`repo`, `workflow`) and rotate the token when done.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `gh auth login --scopes=repo,delete_repo`",
			input: `gh auth login --scopes=repo,delete_repo`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1757",
					Message: "`gh auth login --scopes delete_repo` escalates the token to destructive privileges that outlast the script. Request the minimum scope (`repo`, `workflow`) and rotate the token when done.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1757")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
