package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1195(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid umask 022",
			input:    `umask 022`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid umask 000",
			input: `umask 000`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1195",
					Message: "Avoid `umask 000` — it creates world-writable files. Use `umask 022` or `umask 077` for secure default permissions.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1195")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
