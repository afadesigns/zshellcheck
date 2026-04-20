package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1956(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `tailscale status`",
			input:    `tailscale status`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `tailscale up --auth-key=file:/etc/ts.key` (file source)",
			input:    `tailscale up --auth-key=file:/etc/ts.key`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `tailscale up host --auth-key=tskey-auth-abc123`",
			input: `tailscale up host --auth-key=tskey-auth-abc123`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1956",
					Message: "`tailscale --auth-key=tskey-auth-abc123` puts the pre-auth key in argv — visible in `ps`/`/proc`/history/crash dumps. Use `--auth-key=file:/etc/ts.key` (mode 0400) or `--authkey-env=TS_AUTHKEY`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `tailscale up host --authkey tskey-ZZZ`",
			input: `tailscale up host --authkey tskey-ZZZ`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1956",
					Message: "`tailscale --authkey tskey-ZZZ` puts the pre-auth key in argv — visible in `ps`/`/proc`/history/crash dumps. Use `--auth-key=file:/etc/ts.key` (mode 0400) or `--authkey-env=TS_AUTHKEY`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1956")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
