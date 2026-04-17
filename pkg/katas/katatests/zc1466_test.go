package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1466(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — ufw allow 22",
			input:    `ufw allow 22`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — systemctl start firewalld",
			input:    `systemctl start firewalld`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — ufw disable",
			input: `ufw disable`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1466",
					Message: "Host firewall disabled (ufw disable). Keep it on and open specific ports.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — systemctl stop firewalld",
			input: `systemctl stop firewalld`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1466",
					Message: "Host firewall disabled (systemctl stop firewalld). Keep it on and open specific ports.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — systemctl mask ufw.service",
			input: `systemctl mask ufw.service`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1466",
					Message: "Host firewall disabled (systemctl mask ufw.service). Keep it on and open specific ports.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1466")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
