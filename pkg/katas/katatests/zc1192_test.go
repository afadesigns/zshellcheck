package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1192(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid sleep 1",
			input:    `sleep 1`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid sleep 0",
			input: `sleep 0`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1192",
					Message: "Remove `sleep 0` — it spawns a process that does nothing. Use `:` if an explicit no-op is needed.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1192")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
