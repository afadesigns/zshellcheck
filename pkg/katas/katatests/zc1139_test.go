package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1139(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid source local file",
			input:    `source /usr/local/lib/utils.zsh`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid source URL",
			input: `source https://example.com/script.zsh`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1139",
					Message: "Avoid sourcing scripts from URLs. Download, verify integrity, then source from a local path to prevent supply-chain attacks.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1139")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
