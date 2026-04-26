// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1801(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `fwupdmgr get-devices` (read only)",
			input:    `fwupdmgr get-devices`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `fwupdmgr refresh` (metadata, not flash)",
			input:    `fwupdmgr refresh`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `fwupdmgr update` (all devices)",
			input: `fwupdmgr update`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1801",
					Message: "`fwupdmgr update` flashes firmware — a mid-write interruption can brick BIOS, SSD, Thunderbolt, or NIC microcontrollers. Inhibit reboot triggers (`systemd-inhibit`) and ensure battery / UPS before running.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `fwupdmgr install firmware.cab`",
			input: `fwupdmgr install firmware.cab`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1801",
					Message: "`fwupdmgr install` flashes firmware — a mid-write interruption can brick BIOS, SSD, Thunderbolt, or NIC microcontrollers. Inhibit reboot triggers (`systemd-inhibit`) and ensure battery / UPS before running.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1801")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
