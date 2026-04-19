package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1919(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `ss -tunlp` (read-only socket list)",
			input:    `ss -tunlp`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `ss state established` (preview)",
			input:    `ss state established`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `ss -K state close-wait`",
			input: `ss -K state close-wait`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1919",
					Message: "`ss -K` terminates every socket the filter matches — broad filters (`state established`, `dport 22`) kill the running SSH session. Preview with the same filter minus `-K`, and pin to a specific dst/port/state tuple.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `ss --kill state established`",
			input: `ss --kill state established`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1919",
					Message: "`ss -K` terminates every socket the filter matches — broad filters (`state established`, `dport 22`) kill the running SSH session. Preview with the same filter minus `-K`, and pin to a specific dst/port/state tuple.",
					Line:    1,
					Column:  6,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1919")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
