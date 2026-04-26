// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1976(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `exportfs -ra` (re-sync)",
			input:    `exportfs -ra`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `exportfs -f` (flush cache after edit)",
			input:    `exportfs -f`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `exportfs -au`",
			input: `exportfs -au`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1976",
					Message: "`exportfs -au` unexports live NFS shares — mounted clients see `ESTALE` on every open fd. Use `exportfs -f` after editing `/etc/exports`; reserve `-u`/`-au` for coordinated shutdowns.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `exportfs -u $HOST:$PATH`",
			input: `exportfs -u $HOST:$PATH`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1976",
					Message: "`exportfs -u` unexports live NFS shares — mounted clients see `ESTALE` on every open fd. Use `exportfs -f` after editing `/etc/exports`; reserve `-u`/`-au` for coordinated shutdowns.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1976")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
