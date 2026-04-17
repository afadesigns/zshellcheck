package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1422(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — sudo -u user cmd",
			input:    `sudo -u alice whoami`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — sudo -S cmd",
			input: `sudo -S apt update`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1422",
					Message: "`sudo -S` enables password-via-stdin. Avoid piping plaintext credentials. Use `sudo -A` (askpass), `NOPASSWD:` in sudoers, or `pkexec`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1422")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
