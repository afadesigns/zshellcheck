package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1970(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `losetup -r $LOOP $IMG` (readonly, no partscan)",
			input:    `losetup -r $LOOP $IMG`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `sfdisk --dump $IMG` (offline parser)",
			input:    `sfdisk --dump $IMG`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `losetup -P $LOOP $IMG`",
			input: `losetup -P $LOOP $IMG`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1970",
					Message: "`losetup -P` asks the kernel to parse the partition table of the image — attacker-controlled bytes have tripped kernel CVEs. Use `fdisk -l`/`sfdisk --dump` offline first, scan only known-good images.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `kpartx -av $IMG`",
			input: `kpartx -av $IMG`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1970",
					Message: "`kpartx -a` asks the kernel to parse the partition table of the image — attacker-controlled bytes have tripped kernel CVEs. Use `fdisk -l`/`sfdisk --dump` offline first, scan only known-good images.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `partprobe $LOOP`",
			input: `partprobe $LOOP`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1970",
					Message: "`partprobe` asks the kernel to parse the partition table of the image — attacker-controlled bytes have tripped kernel CVEs. Use `fdisk -l`/`sfdisk --dump` offline first, scan only known-good images.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1970")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
