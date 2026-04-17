package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1512(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — systemctl restart sshd",
			input:    `systemctl restart sshd`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — service with unrecognized verb",
			input:    `service --help`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — service sshd restart",
			input: `service sshd restart`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1512",
					Message: "`service sshd restart` — prefer `systemctl restart sshd` for consistency with other systemd commands.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — service nginx reload",
			input: `service nginx reload`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1512",
					Message: "`service nginx reload` — prefer `systemctl reload nginx` for consistency with other systemd commands.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1512")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
