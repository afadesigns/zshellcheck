package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1511(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — nmcli con up myssid",
			input:    `nmcli con up myssid`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — nmcli --ask",
			input:    `nmcli con up myssid --ask`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — nmcli con mod myssid 802-11-wireless-security.psk mypassword",
			input: `nmcli con mod myssid 802-11-wireless-security.psk mypassword`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1511",
					Message: "`nmcli` passed `802-11-wireless-security.psk <secret>` on the command line — ends up in ps/history. Use `--ask` or a keyfile profile.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — nmcli con mod vpn vpn.secrets.password pw",
			input: `nmcli con mod myvpn vpn.secrets.password vpnpass`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1511",
					Message: "`nmcli` passed `vpn.secrets.password <secret>` on the command line — ends up in ps/history. Use `--ask` or a keyfile profile.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1511")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
