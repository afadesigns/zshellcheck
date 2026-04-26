// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1534(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — dmesg -T",
			input:    `dmesg -T`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — dmesg -c",
			input: `dmesg -c`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1534",
					Message: "`dmesg -c` wipes the kernel ring buffer — subsequent readers see no OOM/panic/audit messages. Read without clearing.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — dmesg -C",
			input: `dmesg -C`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1534",
					Message: "`dmesg -C` wipes the kernel ring buffer — subsequent readers see no OOM/panic/audit messages. Read without clearing.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1534")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
