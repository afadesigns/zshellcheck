package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1829(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `tailscale status`",
			input:    `tailscale status`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `nmcli connection show`",
			input:    `nmcli connection show`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `tailscale down`",
			input: `tailscale down`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1829",
					Message: "`tailscale down` tears down the VPN — if the SSH session rides on it, the script cuts itself off with no rollback. Schedule recovery via `systemd-run --on-active=30s`, or run from console / out-of-band.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `wg-quick down wg0`",
			input: `wg-quick down wg0`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1829",
					Message: "`wg-quick down` tears down the VPN — if the SSH session rides on it, the script cuts itself off with no rollback. Schedule recovery via `systemd-run --on-active=30s`, or run from console / out-of-band.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1829")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
