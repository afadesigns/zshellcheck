package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1042(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "safe loop with $@",
			input:    `for i in "$@"; do echo $i; done`,
			expected: []katas.Violation{},
		},
		{
			name:     "loop over plain items",
			input:    `for i in a b c; do echo $i; done`,
			expected: []katas.Violation{},
		},
		{
			name:     "arithmetic loop",
			input:    `for (( i=0; i<10; i++ )); do echo $i; done`,
			expected: []katas.Violation{},
		},
		{
			name:     "loop with $* (parsed as separate tokens)",
			input:    `for i in $*; do echo $i; done`,
			expected: []katas.Violation{},
		},
		{
			name:     "arithmetic loop (no items)",
			input:    `for (( i=0; i<5; i++ )); do echo $i; done`,
			expected: []katas.Violation{},
		},
		{
			name:     "non ForLoopStatement node",
			input:    `echo hello`,
			expected: []katas.Violation{},
		},
		{
			name:     "for-each with string literal items",
			input:    `for i in "one" "two"; do echo $i; done`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1042")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
