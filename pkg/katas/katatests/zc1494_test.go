// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1494(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — tcpdump without -w (stdout)",
			input:    `tcpdump -i eth0 port 443`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — tcpdump -w capture.pcap -Z tcpdump",
			input:    `tcpdump -i eth0 -w capture.pcap -Z tcpdump`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — tcpdump -i eth0 -w capture.pcap",
			input: `tcpdump -i eth0 -w capture.pcap`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1494",
					Message: "`tcpdump -w` without `-Z <user>` writes the pcap as root and never drops privileges. Add `-Z tcpdump` (or a dedicated capture user).",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1494")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
