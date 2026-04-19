package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1807(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `gh api /repos/owner/repo`",
			input:    `gh api /repos/owner/repo`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `gh api -X GET /repos/owner/repo`",
			input:    `gh api -X GET /repos/owner/repo`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `gh api -X DELETE /repos/owner/repo`",
			input: `gh api -X DELETE /repos/owner/repo`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1807",
					Message: "`gh api -X DELETE` sends a raw DELETE to the GitHub API with the caller's token — no `--yes` guard, no dry-run. Use the high-level `gh` subcommand for the target, or wrap with a preflight GET + explicit confirmation.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `gh api --method=DELETE /repos/owner/repo/releases/123`",
			input: `gh api --method=DELETE /repos/owner/repo/releases/123`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1807",
					Message: "`gh api -X DELETE` sends a raw DELETE to the GitHub API with the caller's token — no `--yes` guard, no dry-run. Use the high-level `gh` subcommand for the target, or wrap with a preflight GET + explicit confirmation.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1807")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
