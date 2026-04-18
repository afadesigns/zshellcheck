package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1617(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — xargs -P 4",
			input:    `xargs -P 4 -n 1 echo`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — xargs without -P",
			input:    `xargs -n 10 echo`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — xargs -P 0",
			input: `xargs -P 0 -n 1 echo`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1617",
					Message: "`xargs -P 0` spawns one child per input line — CPU / FD / memory exhaustion risk. Use `-P $(nproc)` or an explicit cap.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — xargs -P0 (joined)",
			input: `xargs -P0 -n1 echo`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1617",
					Message: "`xargs -P 0` spawns one child per input line — CPU / FD / memory exhaustion risk. Use `-P $(nproc)` or an explicit cap.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1617")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
