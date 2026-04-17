package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1484(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — npm config set cafile /path",
			input:    `npm config set cafile /etc/ssl/ca.pem`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — npm config set strict-ssl true",
			input:    `npm config set strict-ssl true`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — npm config set strict-ssl false",
			input: `npm config set strict-ssl false`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1484",
					Message: "`strict-ssl=false` disables npm/yarn/pnpm registry TLS verification — any MITM swaps packages. Point `cafile` at the right CA bundle instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — yarn config set --global strict-ssl false",
			input: `yarn config set --global strict-ssl false`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1484",
					Message: "`strict-ssl=false` disables npm/yarn/pnpm registry TLS verification — any MITM swaps packages. Point `cafile` at the right CA bundle instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — npm install --strict-ssl=false",
			input: `npm install foo --strict-ssl=false`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1484",
					Message: "`strict-ssl=false` disables npm/yarn/pnpm registry TLS verification — any MITM swaps packages. Point `cafile` at the right CA bundle instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1484")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
