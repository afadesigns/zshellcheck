package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1522(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — ip route add 10.0.0.0/24 dev eth1",
			input:    `ip route add 10.0.0.0/24 dev eth1`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — ip route show default",
			input:    `ip route show default`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — ip route add default via 1.2.3.4",
			input: `ip route add default via 1.2.3.4`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1522",
					Message: "`ip route add default` silently reroutes every non-local packet through the new gateway — MITM surface or CI footgun. Use NetworkManager/networkd config.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — route add default gw 1.2.3.4",
			input: `route add default gw 1.2.3.4`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1522",
					Message: "`route add default` silently reroutes every non-local packet through the new gateway — MITM surface or CI footgun. Use NetworkManager/networkd config.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1522")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
