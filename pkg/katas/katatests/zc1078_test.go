package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1078(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "quoted arguments",
			input:    `cmd "$@"`,
			expected: []katas.Violation{},
		},
		{
			name:     "quoted star",
			input:    `cmd "$*"`,
			expected: []katas.Violation{},
		},
		{
			name:     "unquoted arguments",
			input:    `cmd $@`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1078",
					Message: "Unquoted $@ splits arguments. Use \"$@\" to preserve structure.",
					Line:    1,
					Column:  5,
				},
			},
		},
		{
			name:     "unquoted star",
			input:    `cmd $*`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1078",
					Message: "Unquoted $* splits arguments. Use \"$*\" to preserve structure.",
					Line:    1,
					Column:  5,
				},
			},
		},
		{
			name:     "mixed",
			input:    `cmd arg1 $@ arg2`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1078",
					Message: "Unquoted $@ splits arguments. Use \"$@\" to preserve structure.",
					Line:    1,
					Column:  10,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1078")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
