// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1187(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:  "invalid notify-send",
			input: `notify-send "Build complete"`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1187",
					Message: "Wrap `notify-send` with an `$OSTYPE` check or `command -v` guard. It is Linux-only and will fail silently on macOS.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1187")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
