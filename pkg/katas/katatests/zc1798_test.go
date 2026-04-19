package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1798(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `ufw status numbered`",
			input:    `ufw status numbered`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `ufw delete 3`",
			input:    `ufw delete 3`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `ufw reset`",
			input: `ufw reset`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1798",
					Message: "`ufw reset` drops every user-defined firewall rule. Snapshot (`ufw status numbered > /tmp/ufw.bak`) first, and prefer `ufw delete <num>` for targeted removals.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `ufw reset --force`",
			input: `ufw reset --force`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1798",
					Message: "`ufw reset` drops every user-defined firewall rule. Snapshot (`ufw status numbered > /tmp/ufw.bak`) first, and prefer `ufw delete <num>` for targeted removals.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1798")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
