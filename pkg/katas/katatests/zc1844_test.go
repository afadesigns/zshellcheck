package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1844(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `logger -p auth.notice` (audit)",
			input:    `logger -p auth.notice -t scriptaudit "user added"`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `logger message` (default)",
			input:    `logger "hello"`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `logger -p local0.info`",
			input: `logger -p local0.info "audit: user added to wheel"`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1844",
					Message: "`logger -p local0.info` writes to a `local*` facility — stock `rsyslog`/`journald` rarely collects these. Use `auth.notice`/`authpriv.info` for audit events, or `user.notice` + `-t TAG` for app logs.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `logger msg --priority=local7.notice` (trailing)",
			input: `logger "site event" --priority=local7.notice`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1844",
					Message: "`logger -p local7.notice` writes to a `local*` facility — stock `rsyslog`/`journald` rarely collects these. Use `auth.notice`/`authpriv.info` for audit events, or `user.notice` + `-t TAG` for app logs.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1844")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
