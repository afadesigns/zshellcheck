package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1648(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — logrotate",
			input:    `logrotate -f /etc/logrotate.d/app`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — cp /dev/null to app tmp",
			input:    `cp /dev/null /tmp/marker`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — cp /dev/null /var/log/auth.log",
			input: `cp /dev/null /var/log/auth.log`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1648",
					Message: "`cp /dev/null /var/log/auth.log` wipes an audit log — use `logrotate -f` or `journalctl --vacuum-time=...` instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — truncate -s 0 /var/log/wtmp",
			input: `truncate -s 0 /var/log/wtmp`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1648",
					Message: "`truncate -s 0 /var/log/wtmp` wipes an audit log — use `logrotate -f` or `journalctl --vacuum-time=...` instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1648")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
