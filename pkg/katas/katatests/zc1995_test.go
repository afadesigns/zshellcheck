package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1995(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `setopt BG_NICE` (keeps default on)",
			input:    `setopt BG_NICE`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `unsetopt NO_BG_NICE`",
			input:    `unsetopt NO_BG_NICE`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `unsetopt BG_NICE`",
			input: `unsetopt BG_NICE`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1995",
					Message: "`unsetopt BG_NICE` drops the `nice +5` that bg jobs get by default — a CPU-bound `cmd &` now competes with SSH/editor work. Wrap specific jobs with `nice -n 0` or a systemd `Nice=` unit instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `setopt NO_BG_NICE`",
			input: `setopt NO_BG_NICE`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1995",
					Message: "`setopt NO_BG_NICE` drops the `nice +5` that bg jobs get by default — a CPU-bound `cmd &` now competes with SSH/editor work. Wrap specific jobs with `nice -n 0` or a systemd `Nice=` unit instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1995")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
