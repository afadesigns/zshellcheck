package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1529(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — fsck -n $disk (dry run)",
			input:    `fsck -n $DISK`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — fsck -p $disk (preen)",
			input:    `fsck -p $DISK`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — fsck -y $disk",
			input: `fsck -y $DISK`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1529",
					Message: "`fsck -y` answers yes to every repair prompt — can destroy salvageable data. Prefer `-n` (dry-run) or `-p` (preen).",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — fsck.ext4 -y $disk",
			input: `fsck.ext4 -y $DISK`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1529",
					Message: "`fsck.ext4 -y` answers yes to every repair prompt — can destroy salvageable data. Prefer `-n` (dry-run) or `-p` (preen).",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1529")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
