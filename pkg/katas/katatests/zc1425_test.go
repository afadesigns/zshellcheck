package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1425(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — unrelated command",
			input:    `echo goodbye`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — shutdown now",
			input: `shutdown now`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1425",
					Message: "`shutdown` takes down the system. In scripts, confirm the caller really wants this (interactive prompt, feature flag, or CI guard).",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — reboot",
			input: `reboot now`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1425",
					Message: "`reboot` takes down the system. In scripts, confirm the caller really wants this (interactive prompt, feature flag, or CI guard).",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1425")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
