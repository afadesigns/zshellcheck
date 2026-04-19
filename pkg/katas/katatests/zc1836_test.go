package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1836(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `helm uninstall mychart`",
			input:    `helm uninstall mychart`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `helm uninstall mychart --keep-history` (unrelated)",
			input:    `helm uninstall mychart --keep-history`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `helm uninstall mychart --no-hooks`",
			input: `helm uninstall mychart --no-hooks`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1836",
					Message: "`helm uninstall --no-hooks` skips pre/post-delete cleanup hooks — orphaned locks, DNS, missed PVC backups. Drop the flag; fix stuck hooks via `helm.sh/hook-delete-policy`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `helm delete mychart --no-hooks` (Helm v2 spelling)",
			input: `helm delete mychart --no-hooks`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1836",
					Message: "`helm delete --no-hooks` skips pre/post-delete cleanup hooks — orphaned locks, DNS, missed PVC backups. Drop the flag; fix stuck hooks via `helm.sh/hook-delete-policy`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1836")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
