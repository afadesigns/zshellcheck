package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1658(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — curl -O without -J",
			input:    `curl -O https://example.com/file`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — curl -o with fixed name",
			input:    `curl -o out.bin https://example.com/file`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — curl -OJ combined",
			input: `curl -OJ https://example.com/file`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1658",
					Message: "`curl -OJ` saves the response under the name the server picks in `Content-Disposition` — path traversal is blocked but arbitrary same-dir overwrites are not. Pass `-o NAME` with a filename you control.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — curl -O -J split",
			input: `curl -O -J https://example.com/file`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1658",
					Message: "`curl -OJ` saves the response under the name the server picks in `Content-Disposition` — path traversal is blocked but arbitrary same-dir overwrites are not. Pass `-o NAME` with a filename you control.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1658")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
