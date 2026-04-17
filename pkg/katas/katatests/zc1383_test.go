package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1383(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — echo $TIMEFMT (Zsh)",
			input:    `echo $TIMEFMT`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — echo $TIMEFORMAT",
			input: `echo $TIMEFORMAT`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1383",
					Message: "`$TIMEFORMAT` is Bash-only. Zsh reads `$TIMEFMT` (shorter name) for the `time` builtin's output format.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1383")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
