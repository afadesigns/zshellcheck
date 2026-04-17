package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1516(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — umask 022",
			input:    `umask 022`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — umask 077",
			input:    `umask 077`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — umask 000 (parser normalizes to 0)",
			input: `umask 000`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1516",
					Message: "`umask 0` leaves new files world-readable and world-writable. Use `022` for public software, `077` for secrets handling.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — umask 0",
			input: `umask 0`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1516",
					Message: "`umask 0` leaves new files world-readable and world-writable. Use `022` for public software, `077` for secrets handling.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1516")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
