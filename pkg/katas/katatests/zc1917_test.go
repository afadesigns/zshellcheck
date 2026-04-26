// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1917(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `iw dev wlan0 link` (passive link info)",
			input:    `iw dev wlan0 link`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `iwlist wlan0 channel`",
			input:    `iwlist wlan0 channel`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `iw dev wlan0 scan`",
			input: `iw dev wlan0 scan`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1917",
					Message: "`iw dev <if> scan` runs an active probe-request sweep — interrupts the current association and broadcasts the host to every nearby AP. Use cached `iw dev $IF link` for passive queries.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `iwlist wlan0 scanning`",
			input: `iwlist wlan0 scanning`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1917",
					Message: "`iwlist <if> scan` runs an active probe-request sweep — interrupts the current association and broadcasts the host to every nearby AP. Use cached `iw dev $IF link` for passive queries.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1917")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
