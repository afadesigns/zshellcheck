package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1352(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — xargs without -I",
			input:    `xargs -n 1 echo`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — xargs -I{}",
			input: `xargs -I{} echo hi`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1352",
					Message: "Avoid `xargs -I{}` — iterate with `for x in ${(f)\"$(cmd)\"}; do ...; done` in Zsh. A for loop avoids one-subprocess-per-item and keeps variables in scope.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — xargs -Ixx custom replace-string",
			input: `xargs -Ifile cp file /tmp`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1352",
					Message: "Avoid `xargs -I{}` — iterate with `for x in ${(f)\"$(cmd)\"}; do ...; done` in Zsh. A for loop avoids one-subprocess-per-item and keeps variables in scope.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1352")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
