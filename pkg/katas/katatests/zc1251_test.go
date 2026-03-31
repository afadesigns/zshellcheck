package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1251(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid mount with -o",
			input:    `mount -o noexec,nosuid /dev/sdb1 /mnt`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid mount without device",
			input:    `mount -a`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid mount device without -o",
			input: `mount /dev/sdb1 /mnt`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1251",
					Message: "Use `mount -o noexec,nosuid,nodev` when mounting external media. Without restrictions, mounted filesystems can contain executable exploits.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1251")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
