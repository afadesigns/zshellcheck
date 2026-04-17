package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1401(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — awk for other field",
			input:    `awk '{print $1}' file`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — cut on MACHTYPE",
			input: `cut -d- -f2 $MACHTYPE`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1401",
					Message: "Use Zsh `$VENDOR` for vendor field instead of splitting `$MACHTYPE` with `cut`/`awk`/`sed`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1401")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
