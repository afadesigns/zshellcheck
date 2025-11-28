package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1088(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid subshell",
			input:    `( ls )`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid subshell with return",
			input:    `( return 1 )`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid subshell with exit",
			input:    `( exit 1 )`,
			expected: []katas.Violation{},
		},
		/*
			{
				name:     "valid subshell checked exit status",
				input:    `( cd /tmp ) || exit`,
				expected: []katas.Violation{},
			},
		*/
		{
			name:     "valid subshell used in condition",
			input:    `if ( cd /tmp ); then :; fi`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid subshell side effect",
			input: `( cd /tmp )`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1088",
					Message: "Subshell `( ... )` isolates state changes. The changes (e.g. `cd`, variable assignment) will be lost. Use `{ ... }` to preserve them, or add commands that use the changed state.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid subshell variable assignment",
			input: `( var=1 )`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1088",
					Message: "Subshell `( ... )` isolates state changes. The changes (e.g. `cd`, variable assignment) will be lost. Use `{ ... }` to preserve them, or add commands that use the changed state.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:     "valid subshell output capture",
			input:    `out=$( ( cd /tmp; pwd ) )`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1088")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
