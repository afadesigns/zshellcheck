// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1703(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — sysctl net.ipv4.conf.all.rp_filter=1 (strict)",
			input:    `sysctl -w net.ipv4.conf.all.rp_filter=1`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — sysctl unrelated knob",
			input:    `sysctl -w net.ipv4.tcp_syncookies=1`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — sysctl rp_filter=0",
			input: `sysctl -w net.ipv4.conf.all.rp_filter=0`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1703",
					Message: "`sysctl net.ipv4.conf.all.rp_filter=0` disables reverse-path filtering (anti-spoofing) — classic layer-3 attacks (spoofing / smurf / redirect tamper) reopen. Keep the protective default.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — sysctl accept_source_route=1",
			input: `sysctl -w net.ipv4.conf.all.accept_source_route=1`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1703",
					Message: "`sysctl net.ipv4.conf.all.accept_source_route=1` disables source-routed packet acceptance — classic layer-3 attacks (spoofing / smurf / redirect tamper) reopen. Keep the protective default.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1703")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
