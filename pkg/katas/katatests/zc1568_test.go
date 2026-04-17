package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1568(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — useradd -u 1000 alice",
			input:    `useradd -u 1000 alice`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — useradd -o -u 1000 alice",
			input: `useradd -o -u 1000 alice`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1568",
					Message: "`useradd -o` assigns a non-unique UID — the two accounts share kernel identity, indistinguishable in audit. Use a fresh UID.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — usermod -o -u 500 bob",
			input: `usermod -o -u 500 bob`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1568",
					Message: "`usermod -o` assigns a non-unique UID — the two accounts share kernel identity, indistinguishable in audit. Use a fresh UID.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1568")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
