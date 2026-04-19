package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1808(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `kubectl apply -f deploy.yaml`",
			input:    `kubectl apply -f deploy.yaml`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `kubectl replace -f deploy.yaml` (no --force)",
			input:    `kubectl replace -f deploy.yaml`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `kubectl replace --force -f deploy.yaml`",
			input: `kubectl replace --force -f deploy.yaml`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1808",
					Message: "`kubectl replace --force` is delete + create — pods die, PDBs are ignored, in-flight requests drop. Prefer `kubectl apply -f FILE` and reserve `replace --force` for schema changes `apply` cannot patch, after draining traffic.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1808")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
