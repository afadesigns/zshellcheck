package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1663(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — tune2fs -c 30 (reduced cadence)",
			input:    `tune2fs -c 30 $DISK`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — tune2fs -l (listing)",
			input:    `tune2fs -l $DISK`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — tune2fs -c 0",
			input: `tune2fs -c 0 $DISK`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1663",
					Message: "`tune2fs -c 0` disables periodic fsck on the filesystem — lower the cadence (e.g. `-c 30` / `-i 3m`) instead of turning it off.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — tune2fs -i 0",
			input: `tune2fs -i 0 $DISK`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1663",
					Message: "`tune2fs -i 0` disables periodic fsck on the filesystem — lower the cadence (e.g. `-c 30` / `-i 3m`) instead of turning it off.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1663")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
