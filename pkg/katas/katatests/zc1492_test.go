package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1492(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — at -l (list jobs)",
			input:    `at -l`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — at -r 3 (remove job)",
			input:    `at -r 3`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — at now + 1 hour",
			input: `at now + 1 hour`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1492",
					Message: "`at` schedules via atd with no unit file — harder to audit. Prefer `systemd-run --on-calendar=` or a `.timer` unit.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — at -f script.sh midnight",
			input: `at -f script.sh midnight`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1492",
					Message: "`at` schedules via atd with no unit file — harder to audit. Prefer `systemd-run --on-calendar=` or a `.timer` unit.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1492")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
