package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1840(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `openssl enc -pass env:MYPASS`",
			input:    `openssl enc -aes-256-cbc -pass env:MYPASS -in in.txt -out out.bin`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `openssl enc` without `-k`",
			input:    `openssl enc -aes-256-cbc -in in.txt -out out.bin`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `openssl enc -k SECRET`",
			input: `openssl enc -aes-256-cbc -k hunter2 -in in.txt -out out.bin`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1840",
					Message: "`openssl -k hunter2` embeds the password in argv — visible to `ps`, `/proc/<pid>/cmdline`, and shell history. Use `-pass env:VAR`, `-pass file:PATH`, or `-pass fd:N`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1840")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
