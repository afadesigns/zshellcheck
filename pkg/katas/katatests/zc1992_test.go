package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1992(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `sudo $CMD` (targeted sudoers drop-in)",
			input:    `sudo $CMD`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `pkexec $CMD arg`",
			input: `pkexec $CMD arg`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1992",
					Message: "`pkexec` elevates via PolicyKit — no agent to prompt in a script, poor CVE history (pwnkit), split audit trail. Use `sudo` with a targeted `sudoers.d` drop-in or a systemd unit with `User=`/`AmbientCapabilities=`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1992")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
