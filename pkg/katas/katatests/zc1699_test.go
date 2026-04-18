package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1699(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — kubectl drain with --ignore-daemonsets only",
			input:    `kubectl drain NODE --ignore-daemonsets`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — kubectl cordon (no drain)",
			input:    `kubectl cordon NODE`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — kubectl drain --delete-emptydir-data",
			input: `kubectl drain NODE --delete-emptydir-data --ignore-daemonsets`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1699",
					Message: "`kubectl drain --delete-emptydir-data` deletes `emptyDir` volumes along with the evicted pods — caches / WAL / scratch state are lost. Verify tolerance or migrate to a PersistentVolumeClaim first.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — kubectl drain --delete-local-data (deprecated alias)",
			input: `kubectl drain NODE --force --delete-local-data`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1699",
					Message: "`kubectl drain --delete-local-data` deletes `emptyDir` volumes along with the evicted pods — caches / WAL / scratch state are lost. Verify tolerance or migrate to a PersistentVolumeClaim first.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1699")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
