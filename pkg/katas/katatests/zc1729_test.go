package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1729(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `ip route flush dev eth1`",
			input:    `ip route flush dev eth1`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `ip route add default via 192.168.1.1`",
			input:    `ip route add default via 192.168.1.1`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `ip route replace default via 192.168.1.1`",
			input:    `ip route replace default via 192.168.1.1`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `ip route flush all`",
			input: `ip route flush all`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1729",
					Message: "`ip route flush all` removes the default gateway — the SSH session that just ran it loses connectivity. Scope the flush (`flush dev <iface>`) or use `ip route replace default via <gw>` so the new route is in place first.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `ip route del default`",
			input: `ip route del default`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1729",
					Message: "`ip route del default` removes the default gateway — the SSH session that just ran it loses connectivity. Scope the flush (`flush dev <iface>`) or use `ip route replace default via <gw>` so the new route is in place first.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `ip -6 route flush all`",
			input: `ip -6 route flush all`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1729",
					Message: "`ip route flush all` removes the default gateway — the SSH session that just ran it loses connectivity. Scope the flush (`flush dev <iface>`) or use `ip route replace default via <gw>` so the new route is in place first.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1729")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
