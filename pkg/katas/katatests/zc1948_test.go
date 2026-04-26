// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1948(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `ipmitool -I lan -H bmc -U admin -f /etc/ipmi.pass chassis status`",
			input:    `ipmitool -I lan -H bmc -U admin -f /etc/ipmi.pass chassis status`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `ipmitool -E -H bmc chassis status` (env password)",
			input:    `ipmitool -E -H bmc chassis status`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `ipmitool -H bmc -U admin -P hunter2 chassis status`",
			input: `ipmitool -H bmc -U admin -P hunter2 chassis status`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1948",
					Message: "`ipmitool -P hunter2` leaks the BMC password into argv — visible in `ps`/`/proc`/crash dumps. Use `-f <password_file>` (mode 0400) or `IPMI_PASSWORD=… ipmitool -E`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `ipmitool -Phunter2 -H bmc chassis power status` (joined)",
			input: `ipmitool -Phunter2 -H bmc chassis power status`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1948",
					Message: "`ipmitool -Phunter2` leaks the BMC password into argv — visible in `ps`/`/proc`/crash dumps. Use `-f <password_file>` (mode 0400) or `IPMI_PASSWORD=… ipmitool -E`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1948")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
