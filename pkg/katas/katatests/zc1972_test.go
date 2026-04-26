// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1972(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `dmsetup ls`",
			input:    `dmsetup ls`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `dmsetup remove $NAME` (no force)",
			input:    `dmsetup remove $NAME`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `dmsetup remove_all`",
			input: `dmsetup remove_all`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1972",
					Message: "`dmsetup remove_all` drops LVM/LUKS/multipath mappings while still in use — in-flight I/O returns `ENXIO`, metadata needs a reboot. `umount` + `vgchange -an` / `cryptsetup close` first, then `dmsetup remove`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `dmsetup remove -f $NAME`",
			input: `dmsetup remove -f $NAME`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1972",
					Message: "`dmsetup remove -f` drops LVM/LUKS/multipath mappings while still in use — in-flight I/O returns `ENXIO`, metadata needs a reboot. `umount` + `vgchange -an` / `cryptsetup close` first, then `dmsetup remove`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1972")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
