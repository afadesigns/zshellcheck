// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1933(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `ipvsadm -L -n` (list)",
			input:    `ipvsadm -L -n`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `ipvsadm --save` (backup)",
			input:    `ipvsadm --save`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `ipvsadm -C`",
			input: `ipvsadm -C`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1933",
					Message: "`ipvsadm -C` wipes every IPVS virtual service and real-server binding — load balancing stops, clients see 5xx. Save via `ipvsadm --save`, drain specific services with `-D`, reserve `--clear` for break-glass.",
					Line:    1,
					Column:  9,
				},
			},
		},
		{
			name:  "invalid — `ipvsadm --clear now`",
			input: `ipvsadm --clear now`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1933",
					Message: "`ipvsadm --clear` wipes every IPVS virtual service and real-server binding — load balancing stops, clients see 5xx. Save via `ipvsadm --save`, drain specific services with `-D`, reserve `--clear` for break-glass.",
					Line:    1,
					Column:  10,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1933")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
