package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1619(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — nfs with nosuid,nodev",
			input:    `mount -t nfs -o rw,nosuid,nodev host:/export /mnt`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — local ext4",
			input:    `mount -t ext4 /dev/nvme0n1p1 /data`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — nfs without nosuid/nodev",
			input: `mount -t nfs -o rw host:/export /mnt`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1619",
					Message: "`mount -t nfs` without nosuid,nodev — a hostile server can plant setuid binaries or device nodes that the client kernel honors. Add `nosuid,nodev` to the `-o` options.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — cifs with only nosuid",
			input: `mount -t cifs -o username=foo,nosuid //host/share /mnt`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1619",
					Message: "`mount -t cifs` without nodev — a hostile server can plant setuid binaries or device nodes that the client kernel honors. Add `nosuid,nodev` to the `-o` options.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1619")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
