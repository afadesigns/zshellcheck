// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1508(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — objdump -p",
			input:    `objdump -p /bin/ls`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — readelf -d",
			input:    `readelf -d /bin/ls`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — ldd /bin/ls",
			input: `ldd /bin/ls`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1508",
					Message: "`ldd` on glibc can execute the target binary. Use `objdump -p` or `readelf -d` to inspect ELF dependencies safely.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — ldd /tmp/downloaded.bin",
			input: `ldd /tmp/downloaded.bin`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1508",
					Message: "`ldd` on glibc can execute the target binary. Use `objdump -p` or `readelf -d` to inspect ELF dependencies safely.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1508")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
