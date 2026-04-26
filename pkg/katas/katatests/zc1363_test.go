// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1363(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — find without -newer",
			input:    `find . -type f`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — find -newer",
			input: `find . -newer ref.txt`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1363",
					Message: "Use Zsh `*(e:'[[ $REPLY -nt REF ]]':)` eval glob qualifier instead of `find -newer`/`-anewer`/`-cnewer`/`-newerXY`. `$REPLY` holds the current match.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1363")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
