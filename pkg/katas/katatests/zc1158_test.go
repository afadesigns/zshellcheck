package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1158(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid chown without -R",
			input:    `chown user:group file`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid chown -Rh",
			input:    `chown -Rh user:group dir`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid chown -R without -h",
			input: `chown -R user:group dir`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1158",
					Message: "Use `chown -Rh` or `chown -R --no-dereference` to prevent following symlinks during recursive ownership changes.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1158")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
