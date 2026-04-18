package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1707(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — gpg --keyserver hkps:// trailing",
			input:    `gpg ABCD --keyserver hkps://keys.openpgp.org --recv-keys`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — gpg --recv-keys (default keyserver)",
			input:    `gpg --recv-keys ABCD`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — gpg --keyserver hkp:// trailing",
			input: `gpg ABCD --keyserver hkp://keys.example.com --recv-keys`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1707",
					Message: "`gpg --keyserver hkp://…` is plaintext — a MITM swaps the key bytes. Use `hkps://keys.openpgp.org` or fetch over HTTPS and verify the fingerprint.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1707")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
