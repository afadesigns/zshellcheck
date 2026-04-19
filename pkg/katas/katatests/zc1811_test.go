package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1811(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `chown -R user:group /srv/app`",
			input:    `chown -R user:group /srv/app`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `chmod -R 0750 /srv/app`",
			input:    `chmod -R 0750 /srv/app`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `chown -R --no-preserve-root user /target`",
			input: `chown -R --no-preserve-root user /target`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1811",
					Message: "`chown --no-preserve-root` disables the GNU safeguard against recursing into `/`. Remove the flag; list explicit paths instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `chmod -R --no-preserve-root 0755 /target`",
			input: `chmod -R --no-preserve-root 0755 /target`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1811",
					Message: "`chmod --no-preserve-root` disables the GNU safeguard against recursing into `/`. Remove the flag; list explicit paths instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1811")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
