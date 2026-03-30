package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1045(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "safe local declaration",
			input:    `local var`,
			expected: []katas.Violation{},
		},
		{
			name:     "regular command",
			input:    `echo hello`,
			expected: []katas.Violation{},
		},
		{
			name:     "local with simple value",
			input:    `local var=hello`,
			expected: []katas.Violation{},
		},
		{
			name:  "local with command substitution",
			input: `local var=$(date)`,
			expected: []katas.Violation{
				{
					KataID: "ZC1045",
					Message: "Declare and assign separately to avoid masking return values. " +
						"`local var=$(cmd)` masks the exit code of `cmd`.",
					Line:   1,
					Column: 7,
				},
			},
		},
		{
			name:  "readonly with command substitution",
			input: `readonly var=$(whoami)`,
			expected: []katas.Violation{
				{
					KataID: "ZC1045",
					Message: "Declare and assign separately to avoid masking return values. " +
						"`readonly var=$(cmd)` masks the exit code of `cmd`.",
					Line:   1,
					Column: 10,
				},
			},
		},
		{
			name:  "declare with command substitution",
			input: `declare var=$(date)`,
			expected: []katas.Violation{
				{
					KataID: "ZC1045",
					Message: "Declare and assign separately to avoid masking return values. " +
						"`declare var=$(cmd)` masks the exit code of `cmd`.",
					Line:   1,
					Column: 1,
				},
			},
		},
		{
			name:     "echo is not local or readonly",
			input:    `echo $(date)`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1045")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
