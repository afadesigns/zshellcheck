package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1488(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — ssh -R 2222:localhost:22 host (default bind)",
			input:    `ssh -R 2222:localhost:22 host`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — ssh -R 127.0.0.1:2222:localhost:22 host",
			input:    `ssh -R 127.0.0.1:2222:localhost:22 host`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — ssh -R 0.0.0.0:2222:localhost:22 host",
			input: `ssh -R 0.0.0.0:2222:localhost:22 host`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1488",
					Message: "SSH reverse tunnel bound to all interfaces (`0.0.0.0:2222:localhost:22`) — forwarded port reachable from any network. Bind to a specific IP.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — ssh -D 0.0.0.0:1080 host (dynamic SOCKS)",
			input: `ssh -D 0.0.0.0:1080 host`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1488",
					Message: "SSH reverse tunnel bound to all interfaces (`0.0.0.0:1080`) — forwarded port reachable from any network. Bind to a specific IP.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1488")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
