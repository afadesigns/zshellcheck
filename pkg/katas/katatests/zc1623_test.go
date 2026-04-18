package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1623(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — kill -TERM",
			input:    `kill -TERM $PID`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — kill -CONT (resume)",
			input:    `kill -CONT $PID`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — kill -STOP $PID",
			input: `kill -STOP $PID`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1623",
					Message: "`kill -STOP` halts the target until SIGCONT arrives. Pair every STOP with `trap \"kill -CONT PID\" EXIT` so the resume fires even on failure.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — kill -s STOP $PID",
			input: `kill -s STOP $PID`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1623",
					Message: "`kill -STOP` halts the target until SIGCONT arrives. Pair every STOP with `trap \"kill -CONT PID\" EXIT` so the resume fires even on failure.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — pkill -19",
			input: `pkill -19 slowproc`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1623",
					Message: "`kill -STOP` halts the target until SIGCONT arrives. Pair every STOP with `trap \"kill -CONT PID\" EXIT` so the resume fires even on failure.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1623")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
