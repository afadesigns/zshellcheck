package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1926(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `init 3` (multi-user, legacy but non-destructive)",
			input:    `init 3`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `systemctl reboot`",
			input:    `systemctl reboot`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `init 0` (halt)",
			input: `init 0`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1926",
					Message: "`init 0` changes runlevel — `0` halts, `6` reboots, `1`/`S` drops to single-user. Use `systemctl poweroff`/`reboot`/`rescue` or `shutdown -h +N` so reviewers can read the intent.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `telinit 6` (reboot)",
			input: `telinit 6`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1926",
					Message: "`telinit 6` changes runlevel — `0` halts, `6` reboots, `1`/`S` drops to single-user. Use `systemctl poweroff`/`reboot`/`rescue` or `shutdown -h +N` so reviewers can read the intent.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1926")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
