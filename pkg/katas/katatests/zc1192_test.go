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
			name:     "valid grep -c",
			input:    `grep -c pattern file`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid grep | wc -l",
			input: `grep pattern file | wc -l`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1192",
					Message: "Use `grep -c pattern` instead of `grep pattern | wc -l`. The `-c` flag counts matches without a pipeline.",
					Line:    1,
					Column:  19,
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
