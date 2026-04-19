package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1946(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `setopt HUP` (explicit default)",
			input:    `setopt HUP`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `unsetopt NOMATCH` (unrelated)",
			input:    `unsetopt NOMATCH`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `unsetopt HUP`",
			input: `unsetopt HUP`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1946",
					Message: "`unsetopt HUP` stops the shell from `SIGHUP`-ing background jobs on exit — long pipelines and spawned daemons outlive the session, orphans accumulate. Use `disown` or `systemd-run --scope` on specific commands instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `setopt NO_HUP`",
			input: `setopt NO_HUP`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1946",
					Message: "`setopt NO_HUP` stops the shell from `SIGHUP`-ing background jobs on exit — long pipelines and spawned daemons outlive the session, orphans accumulate. Use `disown` or `systemd-run --scope` on specific commands instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1946")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
