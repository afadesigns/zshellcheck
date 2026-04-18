package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1739(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `git submodule update --init --recursive` (pinned commits)",
			input:    `git submodule update --init --recursive`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `git submodule add URL path`",
			input:    `git submodule add https://example.com/repo path`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `git submodule update --remote`",
			input: `git submodule update --remote`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1739",
					Message: "`git submodule update --remote` ignores the pinned commits in the parent repo and pulls each submodule's branch HEAD — non-reproducible builds, supply-chain risk. Use `--init --recursive` and bump pins via reviewed PRs.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `git submodule update --remote --recursive`",
			input: `git submodule update --remote --recursive`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1739",
					Message: "`git submodule update --remote` ignores the pinned commits in the parent repo and pulls each submodule's branch HEAD — non-reproducible builds, supply-chain risk. Use `--init --recursive` and bump pins via reviewed PRs.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1739")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
