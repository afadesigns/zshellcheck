// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1772(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `hdparm -I $DISK` (info only)",
			input:    `hdparm -I $DISK`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `hdparm -tT $DISK` (benchmark)",
			input:    `hdparm -tT $DISK`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `hdparm --security-erase PASS $DISK`",
			input: `hdparm --security-erase PASS $DISK`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1772",
					Message: "`hdparm --security-erase` issues an ATA-level operation that ignores filesystems and cannot be rolled back. Pin the disk by `/dev/disk/by-id/…`, keep it behind a runbook, keep the password out of argv.",
					Line:    1,
					Column:  9,
				},
			},
		},
		{
			name:  "invalid — `hdparm --trim-sector-ranges 0:1 $DISK`",
			input: `hdparm --trim-sector-ranges 0:1 $DISK`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1772",
					Message: "`hdparm --trim-sector-ranges` issues an ATA-level operation that ignores filesystems and cannot be rolled back. Pin the disk by `/dev/disk/by-id/…`, keep it behind a runbook, keep the password out of argv.",
					Line:    1,
					Column:  9,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1772")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
