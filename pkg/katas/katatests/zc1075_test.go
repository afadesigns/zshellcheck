package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1075(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "quoted variable",
			input:    `rm "$var"`,
			expected: []katas.Violation{},
		},
		{
			name:     "unquoted variable",
			input:    `rm $var`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1075",
					Message: "Unquoted variable expansion '$var' is subject to globbing. Quote it: \"$var\".",
					Line:    1,
					Column:  4,
				},
			},
		},
		{
			name:     "unquoted array access",
			input:    `ls ${files[1]}`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1075",
					Message: "Unquoted array access is subject to globbing. Quote it.",
					Line:    1,
					Column:  4,
				},
			},
		},
		{
			name:     "quoted array access",
			input:    `ls "${files[1]}"`,
			expected: []katas.Violation{},
		},
		{
			name:     "unquoted concatenated",
			input:    `cp $src/file dest`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1075",
					Message: "Unquoted variable expansion '$src/file' is subject to globbing. Quote it: \"$src/file\".",
					Line:    1,
					Column:  4,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1075")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
