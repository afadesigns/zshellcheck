package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1493(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — watch ls",
			input:    `watch ls`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — watch -n 1 df",
			input:    `watch -n 1 df`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — watch -n 0.5 df",
			input:    `watch -n 0.5 df`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — watch -n 0 df",
			input: `watch -n 0 df`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1493",
					Message: "`watch -n 0` pins a core at 100% and saturates the terminal. Use a realistic interval or an event-driven API (inotifywait / journalctl -f).",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — watch -n0 df (joined)",
			input: `watch -n0 df`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1493",
					Message: "`watch -n -n0` pins a core at 100% and saturates the terminal. Use a realistic interval or an event-driven API (inotifywait / journalctl -f).",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1493")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
