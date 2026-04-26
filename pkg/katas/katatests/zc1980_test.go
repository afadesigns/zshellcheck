// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1980(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `udevadm control --reload`",
			input:    `udevadm control --reload`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `udevadm trigger --action=change`",
			input:    `udevadm trigger --action=change`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `udevadm trigger --action=remove`",
			input: `udevadm trigger --action=remove`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1980",
					Message: "`udevadm trigger --action=remove` replays `remove` uevents across `/sys` — SATA/NIC/GPU nodes detach on a live host. Reload rules with `udevadm control --reload`; scope with `--subsystem-match=`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `udevadm trigger -c remove`",
			input: `udevadm trigger -c remove`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1980",
					Message: "`udevadm trigger --action=remove` replays `remove` uevents across `/sys` — SATA/NIC/GPU nodes detach on a live host. Reload rules with `udevadm control --reload`; scope with `--subsystem-match=`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1980")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
