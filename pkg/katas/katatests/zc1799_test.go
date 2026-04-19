package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1799(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `rclone sync --dry-run src dst`",
			input:    `rclone sync --dry-run src dst`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `rclone copy src dst` (copy, not sync)",
			input:    `rclone copy src dst`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `rclone sync src dst`",
			input: `rclone sync src dst`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1799",
					Message: "`rclone sync` deletes anything in DST that's not in SRC — empty / wrong SRC silently wipes DST. Preview with `rclone sync --dry-run`, and pin guards like `--backup-dir`, `--max-delete`, or `--min-age`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `rclone sync local: remote:bucket`",
			input: `rclone sync local: remote:bucket`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1799",
					Message: "`rclone sync` deletes anything in DST that's not in SRC — empty / wrong SRC silently wipes DST. Preview with `rclone sync --dry-run`, and pin guards like `--backup-dir`, `--max-delete`, or `--min-age`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1799")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
