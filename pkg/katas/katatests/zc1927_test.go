package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1927(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `xfreerdp /u:alice /v:host.example` (no password)",
			input:    `xfreerdp /u:alice /v:host.example`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `rdesktop -u alice host.example` (prompts)",
			input:    `rdesktop -u alice host.example`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `xfreerdp /u:alice /p:$PASS /v:host`",
			input: `xfreerdp /u:alice /p:$PASS /v:host`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1927",
					Message: "`xfreerdp /p:$PASS` puts the RDP password in argv — visible in `ps`, `/proc`, and shell history. Pipe via `/from-stdin`, read from a protected `.rdp` file, or use NLA with a cached credential.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `rdesktop -u alice -p hunter2 host.example`",
			input: `rdesktop -u alice -p hunter2 host.example`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1927",
					Message: "`rdesktop -p hunter2` puts the RDP password in argv — visible in `ps`, `/proc`, and shell history. Pipe via `/from-stdin`, read from a protected `.rdp` file, or use NLA with a cached credential.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1927")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
