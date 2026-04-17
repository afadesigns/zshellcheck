package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1478(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — mktemp (default)",
			input:    `mktemp`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — mktemp -d",
			input:    `mktemp -d`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — mktemp -u",
			input: `mktemp -u`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1478",
					Message: "`mktemp -u` returns a unique name but does not create the file — TOCTOU race. Let `mktemp` create the file (or use `-d` for a directory).",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — mktemp -u -t foo.XXXX",
			input: `mktemp -u -t foo.XXXX`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1478",
					Message: "`mktemp -u` returns a unique name but does not create the file — TOCTOU race. Let `mktemp` create the file (or use `-d` for a directory).",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1478")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
