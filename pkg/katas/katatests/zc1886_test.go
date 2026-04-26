// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1886(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `cp /tmp/app.tar /opt/app/`",
			input:    `cp /tmp/app.tar /opt/app/`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `tee /var/log/install.log`",
			input:    `tee /var/log/install.log`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `tee /etc/profile`",
			input: `tee /etc/profile`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1886",
					Message: "`tee ... /etc/profile` writes a shell-init file sourced by every interactive shell — persistent foothold for the next root login. Stage a temp file, validate, and `install -m 644` atomically.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `cp new.sh /etc/profile.d/custom.sh`",
			input: `cp new.sh /etc/profile.d/custom.sh`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1886",
					Message: "`cp ... /etc/profile.d/custom.sh` writes a shell-init file sourced by every interactive shell — persistent foothold for the next root login. Stage a temp file, validate, and `install -m 644` atomically.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1886")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
