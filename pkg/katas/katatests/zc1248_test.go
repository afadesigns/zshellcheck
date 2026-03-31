package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1248(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid ufw",
			input:    `ufw allow 22`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid iptables",
			input: `iptables -A INPUT -p tcp --dport 22 -j ACCEPT`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1248",
					Message: "Prefer `ufw` or `firewalld` over raw `iptables`. Firewall frontends provide persistent, manageable rules.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1248")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
