// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1219(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid wget to file",
			input:    `wget -O file.tar.gz https://example.com/file`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid wget -qO-",
			input: `wget -qO- https://example.com/script.sh`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1219",
					Message: "Use `curl -fsSL` instead of `wget -O -` for piped downloads. `curl` fails on HTTP errors and is available on more platforms.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1219")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
