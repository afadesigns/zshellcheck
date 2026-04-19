package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1812(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `aws ssm put-parameter --type String --value plain --name /app/region`",
			input:    `aws ssm put-parameter --type String --value plain --name /app/region`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `aws ssm put-parameter --type SecureString --value file://secret --name /app/token`",
			input:    `aws ssm put-parameter --type SecureString --value file://secret --name /app/token`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `aws ssm put-parameter --type SecureString --value hunter2 --name /app/token`",
			input: `aws ssm put-parameter --type SecureString --value hunter2 --name /app/token`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1812",
					Message: "`aws ssm put-parameter --type SecureString --value …` puts the plaintext in argv — `ps` / `/proc/PID/cmdline` / history / CLI debug logs can read it. Use `--cli-input-json file://…` (mode 0600) or the `file://` form for `--value`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `aws ssm put-parameter --type=SecureString --value=hunter2 --name /app/token`",
			input: `aws ssm put-parameter --type=SecureString --value=hunter2 --name /app/token`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1812",
					Message: "`aws ssm put-parameter --type SecureString --value …` puts the plaintext in argv — `ps` / `/proc/PID/cmdline` / history / CLI debug logs can read it. Use `--cli-input-json file://…` (mode 0600) or the `file://` form for `--value`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1812")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
