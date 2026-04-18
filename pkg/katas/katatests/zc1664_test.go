package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1664(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — systemctl set-default multi-user.target",
			input:    `systemctl set-default multi-user.target`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — systemctl set-default graphical.target",
			input:    `systemctl set-default graphical.target`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — systemctl set-default rescue.target",
			input: `systemctl set-default rescue.target`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1664",
					Message: "`systemctl set-default rescue.target` makes every subsequent boot land in single-user mode — revert with `set-default multi-user.target` or `graphical.target`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — systemctl set-default emergency.target",
			input: `systemctl set-default emergency.target`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1664",
					Message: "`systemctl set-default emergency.target` makes every subsequent boot land in single-user mode — revert with `set-default multi-user.target` or `graphical.target`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1664")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
