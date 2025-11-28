package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1082(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid expansion",
			input:    `new=${var//foo/bar}`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid sed s///",
			input: `new=$(echo $var | sed 's/foo/bar/')`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1082",
					Message: "Use `${var//old/new}` for string replacement. Pipeline to `sed` is inefficient.",
					Line:    1,
					Column:  17,
				},
			},
		},
		{
			name:  "invalid sed s///g",
			input: `new=$(echo $var | sed "s/foo/bar/g")`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1082",
					Message: "Use `${var//old/new}` for string replacement. Pipeline to `sed` is inefficient.",
					Line:    1,
					Column:  17,
				},
			},
		},
		{
			name:  "invalid sed different separator",
			input: `new=$(print $var | sed 's|foo|bar|')`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1082",
					Message: "Use `${var//old/new}` for string replacement. Pipeline to `sed` is inefficient.",
					Line:    1,
					Column:  18,
				},
			},
		},
		{
			name:     "valid sed other usage",
			input:    `echo $var | sed -n '/p/p'`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1082")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
