package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1497(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — useradd alice",
			input:    `useradd alice`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — useradd -u 1001 alice",
			input:    `useradd -u 1001 alice`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — useradd -u 0 svc",
			input: `useradd -u 0 svc`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1497",
					Message: "Creating a user with UID 0 produces a second root account — classic persistence technique. Use sudo rules tied to a non-0 UID instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — usermod -u0 backup",
			input: `usermod -u0 backup`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1497",
					Message: "Creating a user with UID 0 produces a second root account — classic persistence technique. Use sudo rules tied to a non-0 UID instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1497")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
