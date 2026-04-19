package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1929(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `cpio -o -H newc` (create, not extract)",
			input:    `cpio -o -H newc`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `cpio -i --no-absolute-filenames`",
			input:    `cpio -i --no-absolute-filenames`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `cpio -i -d`",
			input: `cpio -i -d`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1929",
					Message: "`cpio -i` extracts paths verbatim — absolute and `..` entries escape the target dir. Pass `--no-absolute-filenames` and stage into a scratch dir before `mv` into place.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `cpio -idmv` (clustered)",
			input: `cpio -idmv`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1929",
					Message: "`cpio -i` extracts paths verbatim — absolute and `..` entries escape the target dir. Pass `--no-absolute-filenames` and stage into a scratch dir before `mv` into place.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1929")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
