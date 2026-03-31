package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1233(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid npx",
			input:    `npx create-react-app myapp`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid npm install local",
			input:    `npm install express`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid npm install -g",
			input: `npm install -g typescript`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1233",
					Message: "Avoid `npm install -g`. Use `npx` for one-off tool execution or `npm install --save-dev` for project-scoped dependencies.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1233")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
