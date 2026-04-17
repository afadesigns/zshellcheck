package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1514(t *testing.T) {
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
			name:  "invalid — useradd -p hash alice",
			input: `useradd -p $6$salt$hashhashhash alice`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1514",
					Message: "`useradd -p <hash>` puts the hashed password in ps / /proc / history. Use `chpasswd --crypt-method=SHA512` from stdin.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — usermod -p hash bob",
			input: `usermod -p $6$salt$hashhash bob`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1514",
					Message: "`usermod -p <hash>` puts the hashed password in ps / /proc / history. Use `chpasswd --crypt-method=SHA512` from stdin.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1514")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
