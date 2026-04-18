package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1694(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — ssh without forwarding",
			input:    `ssh user@host`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — ssh -J jump (ProxyJump)",
			input:    `ssh -J bastion user@target`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — ssh -A",
			input: `ssh -A user@host`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1694",
					Message: "`ssh -A` forwards the caller's `SSH_AUTH_SOCK` into the remote — any root on that host can reuse the keys. Use `ssh -J jumphost` instead, or a scoped key for the remote task.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — ssh -o ForwardAgent=yes",
			input: `ssh -o ForwardAgent=yes user@host`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1694",
					Message: "`ssh -o ForwardAgent=yes` forwards the caller's `SSH_AUTH_SOCK` into the remote — any root on that host can reuse the keys. Use `ssh -J jumphost` instead, or a scoped key for the remote task.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1694")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
