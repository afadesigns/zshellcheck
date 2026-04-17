package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1438(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — systemctl disable",
			input:    `systemctl disable some.service`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — systemctl mask",
			input: `systemctl mask some.service`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1438",
					Message: "`systemctl mask` permanently blocks service start. If this is a policy choice, document the `unmask` path. For a softer block, use `disable`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1438")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
