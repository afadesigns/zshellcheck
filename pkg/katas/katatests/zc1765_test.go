package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1765(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `snap remove mysnap` (snapshot kept)",
			input:    `snap remove mysnap`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `snap list`",
			input:    `snap list`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `snap remove --purge mysnap`",
			input: `snap remove --purge mysnap`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1765",
					Message: "`snap remove --purge` skips the pre-remove data snapshot — the snap's files are gone with no rollback. Drop `--purge` or capture a `snap save` set ID first.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1765")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
