// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1713(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — consul kv delete scoped",
			input:    `consul kv delete -recurse /app/staging/`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — consul kv delete single key",
			input:    `consul kv delete /key`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — consul kv get (read-only)",
			input:    `consul kv get -recurse /`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — consul kv delete -recurse /",
			input: `consul kv delete -recurse /`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1713",
					Message: "`consul kv delete -recurse /` removes the entire KV store — service discovery, ACL bootstrap, app config. Scope the prefix and snapshot (`consul snapshot save`) first.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1713")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
