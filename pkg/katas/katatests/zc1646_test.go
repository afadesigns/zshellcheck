package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1646(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — btrfs scrub",
			input:    `btrfs scrub start /mnt`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — btrfs check (read-only)",
			input:    `btrfs check $DEV`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — xfs_repair without -L",
			input:    `xfs_repair $DEV`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — btrfs check --repair",
			input: `btrfs check --repair $DEV`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1646",
					Message: "`btrfs check --repair` may worsen damage — try `btrfs scrub` and read-only `btrfs check` first, and snapshot the block device before running.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — xfs_repair -L",
			input: `xfs_repair -L $DEV`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1646",
					Message: "`xfs_repair -L` zeroes the log — uncommitted transactions are lost. Snapshot the block device first; mount read-only and read the log if possible.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1646")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
