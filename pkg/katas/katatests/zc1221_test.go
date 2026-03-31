package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1221(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid fdisk -l",
			input:    `fdisk -l /dev/sda`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid fdisk interactive",
			input: `fdisk /dev/sda`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1221",
					Message: "Avoid `fdisk` in scripts — it is interactive. Use `parted -s` or `sfdisk` for scriptable disk partitioning.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1221")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
