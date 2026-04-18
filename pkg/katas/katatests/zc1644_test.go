package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1644(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — unzip without -P (prompts)",
			input:    `unzip archive.zip`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — zip without password",
			input:    `zip archive.zip files/`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — unzip -P secret",
			input: `unzip -P s3cret archive.zip`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1644",
					Message: "`unzip -P` places the archive password in argv — visible via `ps`. Drop `-P` for interactive prompt, or switch to `7z -p` (reads from stdin) / `age` / `gpg` with keys in a protected file.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — zip -Psecret",
			input: `zip -Ps3cret archive.zip files/`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1644",
					Message: "`zip -P` places the archive password in argv — visible via `ps`. Drop `-P` for interactive prompt, or switch to `7z -p` (reads from stdin) / `age` / `gpg` with keys in a protected file.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1644")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
