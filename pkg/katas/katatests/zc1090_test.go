package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1090(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid unquoted regex",
			input:    `[[ $v =~ ^foo ]]`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid mixed regex",
			input:    `[[ $v =~ "user_"* ]]`,
			expected: []katas.Violation{},
		},
		{
			name:     "invalid quoted start anchor",
			input:    `[[ $v =~ "^foo" ]]`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1090",
					Message: "Quoted regex pattern matches literally. Remove quotes to enable regex matching.",
					Line:    1,
					Column:  10, // Points to string
				},
			},
		},
		{
			name:     "invalid quoted wildcard",
			input:    `[[ $v =~ "foo.*" ]]`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1090",
					Message: "Quoted regex pattern matches literally. Remove quotes to enable regex matching.",
					Line:    1,
					Column:  10,
				},
			},
		},
		{
			name:     "invalid quoted alternation",
			input:    `[[ $v =~ "a|b" ]]`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1090",
					Message: "Quoted regex pattern matches literally. Remove quotes to enable regex matching.",
					Line:    1,
					Column:  10,
				},
			},
		},
		{
			name:     "valid quoted literal",
			input:    `[[ $v =~ "foo" ]]`, // No metachars, arguably valid literal match (though == is better)
			expected: []katas.Violation{},
		},
		{
			name:     "valid quoted variable",
			input:    `[[ $v =~ "$pat" ]]`, // Treating $pat content literally.
			// If $pat contains regex, it WON'T work.
			// But strictly "$pat" contains `$` which I excluded from check.
			// So this should PASS (silently allowed or handled as literal).
			expected: []katas.Violation{},
		},
		{
			name:     "invalid quoted variable with meta",
			input:    `[[ $v =~ "^$pat" ]]`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1090",
					Message: "Quoted regex pattern matches literally. Remove quotes to enable regex matching.",
					Line:    1,
					Column:  10,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1090")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
