// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1684(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — redis-cli PING",
			input:    `redis-cli PING`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — redis-cli -h host",
			input:    `redis-cli -h cache.example.com PING`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — redis-cli -a SECRET PING",
			input: `redis-cli -a SECRET PING`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1684",
					Message: "`redis-cli -a PASSWORD` leaks the password into `ps` / `/proc/PID/cmdline` — use `REDISCLI_AUTH` env var or `-askpass` instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — redis-cli -aSECRET joined",
			input: `redis-cli -aSECRET PING`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1684",
					Message: "`redis-cli -a PASSWORD` leaks the password into `ps` / `/proc/PID/cmdline` — use `REDISCLI_AUTH` env var or `-askpass` instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1684")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
