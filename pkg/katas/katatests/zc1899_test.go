// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1899(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `mokutil --list-enrolled`",
			input:    `mokutil --list-enrolled`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `mokutil --import /root/MOK.der`",
			input:    `mokutil --import /root/MOK.der`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `mokutil --disable-validation now` (leading)",
			input: `mokutil --disable-validation now`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1899",
					Message: "`mokutil --disable-validation` stops the shim from validating kernel/modules against enrolled keys — Secure Boot becomes advisory. Leave validation on; enrol specific keys with `mokutil --import`.",
					Line:    1,
					Column:  10,
				},
			},
		},
		{
			name:  "invalid — `mokutil -l --disable-validation` (trailing)",
			input: `mokutil -l --disable-validation`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1899",
					Message: "`mokutil --disable-validation` stops the shim from validating kernel/modules against enrolled keys — Secure Boot becomes advisory. Leave validation on; enrol specific keys with `mokutil --import`.",
					Line:    1,
					Column:  13,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1899")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
