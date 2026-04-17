package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1527(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — crontab -l",
			input:    `crontab -l`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — crontab file",
			input:    `crontab /etc/cron.d/myfile`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — crontab -",
			input: `crontab -`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1527",
					Message: "`crontab -` overwrites the user's crontab from stdin — silently drops existing rows. Use /etc/cron.d/ files or a diff/merge workflow.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — crontab -u svc -",
			input: `crontab -u svc -`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1527",
					Message: "`crontab -` overwrites the user's crontab from stdin — silently drops existing rows. Use /etc/cron.d/ files or a diff/merge workflow.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1527")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
