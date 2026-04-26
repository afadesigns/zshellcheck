// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1431(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — crontab -l (list)",
			input:    `crontab -l`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — crontab -r",
			input: `crontab -r`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1431",
					Message: "`crontab -r` removes all cron jobs with no backup. Save first (`crontab -l > cron.bak`) and use `crontab -ir` for interactive confirmation.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1431")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
