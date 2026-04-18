package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1741(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `mkpasswd -s` (read from stdin)",
			input:    `mkpasswd -s`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `mkpasswd -m sha-512 -s`",
			input:    `mkpasswd -m sha-512 -s`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `mkpasswd --stdin`",
			input:    `mkpasswd --stdin`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `mkpasswd hunter2`",
			input: `mkpasswd hunter2`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1741",
					Message: "`mkpasswd PASSWORD` puts the cleartext password in argv — visible in `ps`, `/proc`, history. Use `mkpasswd -s` and pipe the secret via stdin (`printf %s \"$PASSWORD\" | mkpasswd -s`).",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `mkpasswd -m sha-512 hunter2`",
			input: `mkpasswd -m sha-512 hunter2`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1741",
					Message: "`mkpasswd PASSWORD` puts the cleartext password in argv — visible in `ps`, `/proc`, history. Use `mkpasswd -s` and pipe the secret via stdin (`printf %s \"$PASSWORD\" | mkpasswd -s`).",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1741")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
