package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1586(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — systemctl enable sshd",
			input:    `systemctl enable sshd`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — chkconfig sshd on",
			input: `chkconfig sshd on`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1586",
					Message: "`chkconfig` is a SysV-init relic. Use `systemctl enable|disable <unit>` directly.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — update-rc.d sshd defaults",
			input: `update-rc.d sshd defaults`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1586",
					Message: "`update-rc.d` is a SysV-init relic. Use `systemctl enable|disable <unit>` directly.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — insserv sshd",
			input: `insserv sshd`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1586",
					Message: "`insserv` is a SysV-init relic. Use `systemctl enable|disable <unit>` directly.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1586")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
