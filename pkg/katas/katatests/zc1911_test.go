package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1911(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `umount /mnt/scratch`",
			input:    `umount /mnt/scratch`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `umount -f /mnt/stuck` (force, not lazy)",
			input:    `umount -f /mnt/stuck`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `umount -l /mnt/scratch`",
			input: `umount -l /mnt/scratch`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1911",
					Message: "`umount -l` detaches the mount but leaves any open fd pointing at a ghost filesystem — writers keep writing, re-mounts stack invisibly. Stop the fd holder first (`lsof`/`fuser`), then do a normal `umount`.",
					Line:    1,
					Column:  8,
				},
			},
		},
		{
			name:  "invalid — `umount --lazy /mnt/scratch`",
			input: `umount --lazy /mnt/scratch`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1911",
					Message: "`umount --lazy` detaches the mount but leaves any open fd pointing at a ghost filesystem — writers keep writing, re-mounts stack invisibly. Stop the fd holder first (`lsof`/`fuser`), then do a normal `umount`.",
					Line:    1,
					Column:  9,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1911")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
