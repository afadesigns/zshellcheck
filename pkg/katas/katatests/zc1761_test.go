package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1761(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `gh gist create secret.env` (secret by default)",
			input:    `gh gist create secret.env`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `gh gist list`",
			input:    `gh gist list`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `gh gist create --public secret.env`",
			input: `gh gist create --public secret.env`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1761",
					Message: "`gh gist create --public` publishes the file to the public discover feed — search engines crawl it within minutes. Drop the flag unless public exposure is the explicit goal.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `gh gist create -p note.md`",
			input: `gh gist create -p note.md`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1761",
					Message: "`gh gist create -p` publishes the file to the public discover feed — search engines crawl it within minutes. Drop the flag unless public exposure is the explicit goal.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1761")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
