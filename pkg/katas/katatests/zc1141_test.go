package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1141(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid curl to file",
			input:    `curl -o file.tar.gz https://example.com/file.tar.gz`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid curl -sSL",
			input: `curl -sSL https://example.com/install.sh`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1141",
					Message: "Avoid `curl -s URL | sh`. Download the script first, verify its integrity, then execute. Piping directly from the internet is a supply-chain risk.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1141")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
