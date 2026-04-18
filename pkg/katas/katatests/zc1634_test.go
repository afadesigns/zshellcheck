package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1634(t *testing.T) {
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
			name:     "valid — umask 002 (group-write collab)",
			input:    `umask 002`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — umask 111",
			input: `umask 111`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1634",
					Message: "`umask 111` leaves world-write on new files — the \"other\" digit must be `2`/`3`/`6`/`7` to mask the write bit. Use `022` for public, `077` for secrets.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — umask 115 (last digit 5 leaves world-write)",
			input: `umask 115`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1634",
					Message: "`umask 115` leaves world-write on new files — the \"other\" digit must be `2`/`3`/`6`/`7` to mask the write bit. Use `022` for public, `077` for secrets.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1634")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
