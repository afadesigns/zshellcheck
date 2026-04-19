package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1939(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `systemctl reboot`",
			input:    `systemctl reboot`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `shutdown -r +5`",
			input:    `shutdown -r +5`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `reboot -f`",
			input: `reboot -f`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1939",
					Message: "`reboot -f` fires `reboot(2)` immediately — no `ExecStop=`, no filesystem sync, no clean unmount. Databases replay from last checkpoint. Use `systemctl reboot` / `shutdown -r +N`; reserve `-f` for wedged recovery.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `poweroff --force now` (mangled)",
			input: `poweroff --force now`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1939",
					Message: "`reboot/halt/poweroff --force` fires `reboot(2)` immediately — no `ExecStop=`, no filesystem sync, no clean unmount. Databases replay from last checkpoint. Use `systemctl reboot` / `shutdown -r +N`; reserve `-f` for wedged recovery.",
					Line:    1,
					Column:  12,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1939")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
