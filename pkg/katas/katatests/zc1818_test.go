// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1818(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `rsync -avn --delete src/ dst/` (dry-run short)",
			input:    `rsync -avn --delete src/ dst/`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `rsync -av --delete --dry-run src/ dst/`",
			input:    `rsync -av --delete --dry-run src/ dst/`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `rsync -av src/ dst/` (no delete)",
			input:    `rsync -av src/ dst/`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `rsync -av --delete src/ dst/`",
			input: `rsync -av --delete src/ dst/`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1818",
					Message: "`rsync --delete` without `--dry-run` removes anything in DST that isn't in SRC. Preview with `rsync -av --delete --dry-run SRC DST`, and pin `--max-delete=N` so an accidentally empty SRC can't cascade.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1818")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
