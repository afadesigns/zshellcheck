package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1802(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `dnf history list`",
			input:    `dnf history list`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `dnf history info 5`",
			input:    `dnf history info 5`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `dnf history undo 5`",
			input: `dnf history undo 5`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1802",
					Message: "`dnf history undo` reverses the past transaction — deps downgrade, security patches can get reverted. Review with `dnf history info`, or restore a filesystem snapshot.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `yum history rollback 3`",
			input: `yum history rollback 3`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1802",
					Message: "`yum history rollback` reverses the past transaction — deps downgrade, security patches can get reverted. Review with `dnf history info`, or restore a filesystem snapshot.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1802")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
