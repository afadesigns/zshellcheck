package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1784(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `git config core.hooksPath .githooks` (repo-relative)",
			input:    `git config core.hooksPath .githooks`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `git config user.email me@example.com`",
			input:    `git config user.email me@example.com`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `git config core.hooksPath /tmp/hooks`",
			input: `git config core.hooksPath /tmp/hooks`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1784",
					Message: "`git config core.hooksPath /tmp/hooks` runs hooks from a mutable path — supply-chain primitive. Keep hooks in the repo's `.git/hooks/` (or a tracked `.githooks/`) and point `core.hooksPath` at repo-owned paths only.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `git config --global core.hooksPath /home/attacker/hooks`",
			input: `git config --global core.hooksPath /home/attacker/hooks`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1784",
					Message: "`git config core.hooksPath /home/attacker/hooks` runs hooks from a mutable path — supply-chain primitive. Keep hooks in the repo's `.git/hooks/` (or a tracked `.githooks/`) and point `core.hooksPath` at repo-owned paths only.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1784")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
