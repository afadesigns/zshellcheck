// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1659(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — fuser -k port target",
			input:    `fuser -k 8080/tcp`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — fuser without -k",
			input:    `fuser /var/log/syslog`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — fuser -k /mnt",
			input: `fuser -k /mnt`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1659",
					Message: "`fuser -k /mnt` signals every process with a file open anywhere under the path — use PID / port targets or `systemctl stop` for services.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — fuser -kim /",
			input: `fuser -kim /`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1659",
					Message: "`fuser -k /` signals every process with a file open anywhere under the path — use PID / port targets or `systemctl stop` for services.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1659")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
