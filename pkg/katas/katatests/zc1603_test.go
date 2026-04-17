package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1603(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — gdb on core file",
			input:    `gdb /usr/bin/app /var/lib/cores/app.core`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — coredumpctl",
			input:    `coredumpctl debug myapp`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — gdb -p 1234",
			input: `gdb -p 1234`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1603",
					Message: "`gdb -p PID` attaches via ptrace — memory, registers, env, and stack of the target are readable. Use `coredumpctl` on a captured core, not a live attach from a script.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — ltrace -p $PID",
			input: `ltrace -p $PID`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1603",
					Message: "`ltrace -p PID` attaches via ptrace — memory, registers, env, and stack of the target are readable. Use `coredumpctl` on a captured core, not a live attach from a script.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1603")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
