// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1902(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `ln -s /opt/app/current /opt/app/live` (app release symlink)",
			input:    `ln -s /opt/app/current /opt/app/live`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `ln -s /dev/null /tmp/scratch` (non-sensitive target)",
			input:    `ln -s /dev/null /tmp/scratch`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `ln -sf /dev/null /var/log/auth.log`",
			input: `ln -sf /dev/null /var/log/auth.log`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1902",
					Message: "`ln -s /dev/null /var/log/auth.log` redirects every write to the bit-bucket — audit / history entries vanish silently. If the log must stop, disable the writer or rotate with `logrotate`, never symlink to `/dev/null`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `ln -s /dev/null $HOME/.bash_history`",
			input: `ln -s /dev/null $HOME/.bash_history`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1902",
					Message: "`ln -s /dev/null $HOME/.bash_history` redirects every write to the bit-bucket — audit / history entries vanish silently. If the log must stop, disable the writer or rotate with `logrotate`, never symlink to `/dev/null`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1902")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
