package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1676(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — helm rollback without --force",
			input:    `helm rollback myapp 2`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — helm history",
			input:    `helm history myapp`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — helm rollback --force",
			input: `helm rollback myapp 2 --force`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1676",
					Message: "`helm rollback --force` deletes and recreates unpatched resources — loses in-flight traffic and bypasses PodDisruptionBudget. Drop `--force` and gate the rollback via change review.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1676")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
