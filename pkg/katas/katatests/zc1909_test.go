package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1909(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `kexec --help`",
			input:    `kexec --help`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `kexec -h`",
			input:    `kexec -h`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `kexec -l $KERN`",
			input: `kexec -l $KERN`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1909",
					Message: "`kexec -l` stages or jumps to a kernel without firmware / bootloader verification — Secure Boot never checks the signature. Gate behind `sudo` + audit and prefer `systemctl kexec` or a real reboot.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `kexec -e` (execute)",
			input: `kexec -e`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1909",
					Message: "`kexec -e` stages or jumps to a kernel without firmware / bootloader verification — Secure Boot never checks the signature. Gate behind `sudo` + audit and prefer `systemctl kexec` or a real reboot.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1909")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
