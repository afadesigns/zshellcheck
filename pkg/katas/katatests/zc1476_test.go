package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1476(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — apt-key list (read-only)",
			input:    `apt-key list`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — apt-key export",
			input:    `apt-key export ABCD`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — apt-key add -",
			input: `apt-key add -`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1476",
					Message: "`apt-key add` adds to a global keyring that signs every repo. Use `/etc/apt/keyrings/<vendor>.gpg` + `signed-by=` instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — apt-key adv --recv-keys",
			input: `apt-key adv --keyserver keyserver.ubuntu.com --recv-keys ABCD`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1476",
					Message: "`apt-key adv` adds to a global keyring that signs every repo. Use `/etc/apt/keyrings/<vendor>.gpg` + `signed-by=` instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1476")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
