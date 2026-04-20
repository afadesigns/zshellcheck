package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC2001(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `setopt EVAL_LINENO` (default on)",
			input:    `setopt EVAL_LINENO`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `unsetopt NO_EVAL_LINENO`",
			input:    `unsetopt NO_EVAL_LINENO`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `unsetopt EVAL_LINENO`",
			input: `unsetopt EVAL_LINENO`,
			expected: []katas.Violation{
				{
					KataID:  "ZC2001",
					Message: "`unsetopt EVAL_LINENO` reverts `$LINENO` inside `eval` to the outer line — errors in generated configs collapse to a single source line and stack frames past `eval` vanish. Keep on; scope via `emulate -LR sh`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `setopt NO_EVAL_LINENO`",
			input: `setopt NO_EVAL_LINENO`,
			expected: []katas.Violation{
				{
					KataID:  "ZC2001",
					Message: "`setopt NO_EVAL_LINENO` reverts `$LINENO` inside `eval` to the outer line — errors in generated configs collapse to a single source line and stack frames past `eval` vanish. Keep on; scope via `emulate -LR sh`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC2001")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
