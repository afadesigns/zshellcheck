package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1442(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — kubectl get pods",
			input:    `kubectl get pods`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — kubectl delete specific pod",
			input:    `kubectl delete pod myapp-abc123`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — kubectl delete --all pods",
			input: `kubectl delete pods --all`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1442",
					Message: "`kubectl delete --all` (or `-A`) deletes resources cluster-wide. Dry-run with `--dry-run=client -o yaml` first, and scope with `-n` namespace.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1442")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
