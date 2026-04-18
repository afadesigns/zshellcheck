package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1721(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `chmod 660 /dev/kvm` (group, not world)",
			input:    `chmod 660 /dev/kvm`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `chmod 644 /tmp/file` (not /dev/)",
			input:    `chmod 644 /tmp/file`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `chmod 644 /dev/null` (read-only, ignored)",
			input:    `chmod 644 /dev/null`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `chmod 666 /dev/kvm`",
			input: `chmod 666 /dev/kvm`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1721",
					Message: "`chmod 666 /dev/kvm` opens a kernel device node to every local user — privilege-escalation surface. Use a udev rule that grants the specific group access instead of world-write.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `chmod 0666 /dev/uinput` (leading-zero octal)",
			input: `chmod 0666 /dev/uinput`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1721",
					Message: "`chmod 438 /dev/uinput` opens a kernel device node to every local user — privilege-escalation surface. Use a udev rule that grants the specific group access instead of world-write.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `chmod 777 /dev/dri/card0`",
			input: `chmod 777 /dev/dri/card0`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1721",
					Message: "`chmod 777 /dev/dri/card0` opens a kernel device node to every local user — privilege-escalation surface. Use a udev rule that grants the specific group access instead of world-write.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1721")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
