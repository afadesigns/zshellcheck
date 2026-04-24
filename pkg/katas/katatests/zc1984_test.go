package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1984(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `sgdisk -p $DISK` (print)",
			input:    `sgdisk -p $DISK`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `sgdisk --backup=/root/disk.gpt $DISK`",
			input:    `sgdisk --backup=/root/disk.gpt $DISK`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `sgdisk -Z $DISK`",
			input: `sgdisk -Z $DISK`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1984",
					Message: "`sgdisk -Z` erases the GPT on the target device — a wrong `$DISK` detaches every partition/LVM/LUKS header and bricks boot. `lsblk`/`blkid` preflight, `--backup` the old table, and test with `-t`/`--pretend` first.",
					Line:    1,
					Column:  8,
				},
			},
		},
		{
			name:  "invalid — `sgdisk -o $DISK`",
			input: `sgdisk -o $DISK`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1984",
					Message: "`sgdisk -o` erases the GPT on the target device — a wrong `$DISK` detaches every partition/LVM/LUKS header and bricks boot. `lsblk`/`blkid` preflight, `--backup` the old table, and test with `-t`/`--pretend` first.",
					Line:    1,
					Column:  8,
				},
			},
		},
		{
			name:  "invalid — `sgdisk --zap-all $DISK`",
			input: `sgdisk --zap-all $DISK`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1984",
					Message: "`sgdisk --zap-all` erases the GPT on the target device — a wrong `$DISK` detaches every partition/LVM/LUKS header and bricks boot. `lsblk`/`blkid` preflight, `--backup` the old table, and test with `-t`/`--pretend` first.",
					Line:    1,
					Column:  9,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1984")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
