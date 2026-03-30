package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1182(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:  "invalid nc usage",
			input: `nc localhost 8080`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1182",
					Message: "Avoid `nc` for network operations in scripts. Use `curl` for HTTP or `zmodload zsh/net/tcp` for raw TCP connections with TLS support.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid netcat usage",
			input: `netcat -l 9090`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1182",
					Message: "Avoid `netcat` for network operations in scripts. Use `curl` for HTTP or `zmodload zsh/net/tcp` for raw TCP connections with TLS support.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1182")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
