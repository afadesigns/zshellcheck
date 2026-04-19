package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1874(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `sshuttle -r user@host 10.0.0.0/8`",
			input:    `sshuttle -r user@host 10.0.0.0/8`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `sshuttle -r user@host 192.168.1.0/24`",
			input:    `sshuttle -r user@host 192.168.1.0/24`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `sshuttle -r user@host 0/0`",
			input: `sshuttle -r user@host 0/0`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1874",
					Message: "`sshuttle ... 0/0` routes every outbound packet through the jump host — a compromise of `user@host` sees the whole fleet's traffic. Scope to the subnets you actually need.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `sshuttle -r user@host 0.0.0.0/0`",
			input: `sshuttle -r user@host 0.0.0.0/0`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1874",
					Message: "`sshuttle ... 0.0.0.0/0` routes every outbound packet through the jump host — a compromise of `user@host` sees the whole fleet's traffic. Scope to the subnets you actually need.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1874")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
