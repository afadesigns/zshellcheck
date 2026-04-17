package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1440(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — usermod -aG (append)",
			input:    `usermod -aG docker alice`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — usermod -L (lock)",
			input:    `usermod -L alice`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — usermod -G (replace)",
			input: `usermod -G docker alice`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1440",
					Message: "`usermod -G` without `-a` overwrites supplementary groups. Use `-aG` to append — existing memberships are preserved.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1440")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
