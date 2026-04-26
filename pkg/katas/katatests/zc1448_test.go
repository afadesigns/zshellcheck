// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1448(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — apt-get install -y",
			input:    `apt-get install -y curl`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — apt-get update (not install)",
			input:    `apt-get update`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — apt-get install without -y",
			input: `apt-get install curl`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1448",
					Message: "`apt-get install`/`apt install` without `-y` hangs on the interactive prompt in scripts. Add `-y` and set DEBIAN_FRONTEND=noninteractive.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1448")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
