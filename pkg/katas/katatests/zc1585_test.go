package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1585(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — ufw allow from 10.0.0.0/8 to any port 22",
			input:    `ufw allow from 10.0.0.0/8 to any port 22`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — ufw status",
			input:    `ufw status`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — ufw allow from any to any port 22",
			input: `ufw allow from any to any port 22`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1585",
					Message: "`ufw allow from any …` opens the port to the whole internet. Scope to a specific source CIDR.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1585")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
