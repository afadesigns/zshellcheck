package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1346(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — find without -user",
			input:    `find . -type f`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — find -user",
			input: `find . -user alice`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1346",
					Message: "Use Zsh `*(u:name:)` / `*(u+uid)` / `*(U)` glob qualifiers instead of `find -user`/`-uid`/`-nouser`. Ownership predicates live entirely in the shell.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — find -uid",
			input: `find / -uid 1000`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1346",
					Message: "Use Zsh `*(u:name:)` / `*(u+uid)` / `*(U)` glob qualifiers instead of `find -user`/`-uid`/`-nouser`. Ownership predicates live entirely in the shell.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — find -nouser",
			input: `find / -nouser`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1346",
					Message: "Use Zsh `*(u:name:)` / `*(u+uid)` / `*(U)` glob qualifiers instead of `find -user`/`-uid`/`-nouser`. Ownership predicates live entirely in the shell.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1346")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
