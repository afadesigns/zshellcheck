package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1547(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — kubectl apply -f manifests/",
			input:    `kubectl apply -f manifests/`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — kubectl apply --prune -l app=x -f manifests/",
			input:    `kubectl apply --prune -l app=x -f manifests/`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — kubectl apply --prune --all -f m/",
			input: `kubectl apply --prune --all -f m/`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1547",
					Message: "`kubectl apply --prune --all` deletes every matching resource not in the manifest — manifest typo wipes other teams' resources. Scope with a narrow `-l <selector>`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — kubectl apply --prune -A -f m/",
			input: `kubectl apply --prune -A -f m/`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1547",
					Message: "`kubectl apply --prune --all` deletes every matching resource not in the manifest — manifest typo wipes other teams' resources. Scope with a narrow `-l <selector>`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1547")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
