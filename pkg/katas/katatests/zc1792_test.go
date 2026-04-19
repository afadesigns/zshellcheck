package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1792(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `btrfs subvolume list /` (read only)",
			input:    `btrfs subvolume list /`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `btrfs device usage /` (read only)",
			input:    `btrfs device usage /`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `btrfs subvolume delete /snapshots/2025-01-01`",
			input: `btrfs subvolume delete /snapshots/2025-01-01`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1792",
					Message: "`btrfs subvolume delete` drops btrfs state with no automatic rollback — snapshots vanish on `subvolume delete`, and `device remove` can leave the filesystem degraded. Confirm the target explicitly.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `btrfs device remove $DEV /mnt/pool`",
			input: `btrfs device remove $DEV /mnt/pool`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1792",
					Message: "`btrfs device remove` drops btrfs state with no automatic rollback — snapshots vanish on `subvolume delete`, and `device remove` can leave the filesystem degraded. Confirm the target explicitly.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1792")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
