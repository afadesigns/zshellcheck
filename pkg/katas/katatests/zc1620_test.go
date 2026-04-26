// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1620(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — tee on log file",
			input:    `tee -a /var/log/app.log`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — tee to tmp staging",
			input:    `tee /tmp/sudoers.new`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — tee /etc/sudoers",
			input: `tee /etc/sudoers`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1620",
					Message: "`tee /etc/sudoers` writes without syntax validation — a typo locks everyone out of sudo. Pipe through `visudo -cf /dev/stdin` or stage in a temp file and `visudo -cf` before `mv`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — tee -a /etc/sudoers.d/custom",
			input: `tee -a /etc/sudoers.d/custom`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1620",
					Message: "`tee /etc/sudoers.d/custom` writes without syntax validation — a typo locks everyone out of sudo. Pipe through `visudo -cf /dev/stdin` or stage in a temp file and `visudo -cf` before `mv`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1620")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
