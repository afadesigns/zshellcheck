package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1536(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — iptables -L (list)",
			input:    `iptables -L`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — iptables -A INPUT -j DROP",
			input:    `iptables -A INPUT -j DROP`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — iptables -I PREROUTING ... -j DNAT",
			input: `iptables -t nat -I PREROUTING -p tcp -j DNAT`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1536",
					Message: "`iptables -j DNAT` rewrites packet destination — silent redirect surface. Use declarative nftables/firewalld config.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — iptables -A OUTPUT ... -j REDIRECT",
			input: `iptables -t nat -A OUTPUT -p tcp -j REDIRECT`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1536",
					Message: "`iptables -j REDIRECT` rewrites packet destination — silent redirect surface. Use declarative nftables/firewalld config.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1536")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
