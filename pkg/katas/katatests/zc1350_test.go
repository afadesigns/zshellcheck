package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1350(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — expr arithmetic (not substr)",
			input:    `expr 2 + 3`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — expr substr",
			input: `expr substr "$s" 1 5`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1350",
					Message: "Use `${str:pos:len}` instead of `expr substr` for substring extraction. Parameter expansion avoids spawning an external process.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1350")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
