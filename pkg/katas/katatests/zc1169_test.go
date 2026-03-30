package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1169(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid install without -m",
			input:    `install -d /usr/local/bin`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid install -m",
			input: `install -m 755 script /usr/local/bin/script`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1169",
					Message: "Consider using `cp` + `chmod` instead of `install -m`. Separate commands are clearer in shell scripts.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1169")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
