package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1706(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — lvresize grow without -r",
			input:    `lvresize -L +2G vg/lv`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — lvresize shrink with -r",
			input:    `lvresize -L -2G -r vg/lv`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — lvextend (always grows)",
			input:    `lvextend -L +2G vg/lv`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — lvresize shrink without -r",
			input: `lvresize -L -2G vg/lv`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1706",
					Message: "`lvresize` shrinks the LV without `-r` / `--resizefs` — the filesystem on top is not shrunk first and writes past the new boundary corrupt metadata. Add `-r` (or shrink the FS manually first).",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — lvreduce without -r",
			input: `lvreduce -L 1G vg/lv`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1706",
					Message: "`lvreduce` shrinks the LV without `-r` / `--resizefs` — the filesystem on top is not shrunk first and writes past the new boundary corrupt metadata. Add `-r` (or shrink the FS manually first).",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1706")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
