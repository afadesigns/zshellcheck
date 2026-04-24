package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1949(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `rmmod nft_chain_nat` (no force)",
			input:    `rmmod nft_chain_nat`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `modprobe -r bluetooth`",
			input:    `modprobe -r bluetooth`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `rmmod -f nvidia`",
			input: `rmmod -f nvidia`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1949",
					Message: "`rmmod -f` tears down a module even when its refcount is non-zero — in-use drivers dangle, kernel oopses on the next callback. Stop holders first (`lsof`/`umount`/`ip link down`), then `rmmod` without `-f`.",
					Line:    1,
					Column:  7,
				},
			},
		},
		{
			name:  "invalid — `rmmod --force foo`",
			input: `rmmod --force foo`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1949",
					Message: "`rmmod --force` tears down a module even when its refcount is non-zero — in-use drivers dangle, kernel oopses on the next callback. Stop holders first (`lsof`/`umount`/`ip link down`), then `rmmod` without `-f`.",
					Line:    1,
					Column:  8,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1949")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
