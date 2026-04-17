package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1593(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — different command",
			input:    `lsblk $DISK`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — blkdiscard $DISK",
			input: `blkdiscard $DISK`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1593",
					Message: "`blkdiscard` issues TRIM/DISCARD across the full device — data is unrecoverable once the controller acknowledges. Require operator confirmation before running.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — blkdiscard -z $DISK",
			input: `blkdiscard -z $DISK`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1593",
					Message: "`blkdiscard` issues TRIM/DISCARD across the full device — data is unrecoverable once the controller acknowledges. Require operator confirmation before running.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1593")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
