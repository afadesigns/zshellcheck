package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1915(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `mdadm --detail $MD` (read only)",
			input:    `mdadm --detail $MD`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `mdadm --examine $DISK` (read only)",
			input:    `mdadm --examine $DISK`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `mdadm --zero-superblock $DISK` (mangled)",
			input: `mdadm --zero-superblock $DISK`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1915",
					Message: "`mdadm --zero-superblock` drops RAID metadata or halts a live array — mounted root or /boot panics the host; a stale superblock scrambles data on next `--create`. Snapshot `mdadm --detail --export` first and keep behind a runbook.",
					Line:    1,
					Column:  9,
				},
			},
		},
		{
			name:  "invalid — `mdadm -S $MD` (stop array)",
			input: `mdadm -S $MD`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1915",
					Message: "`mdadm -S` drops RAID metadata or halts a live array — mounted root or /boot panics the host; a stale superblock scrambles data on next `--create`. Snapshot `mdadm --detail --export` first and keep behind a runbook.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1915")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
