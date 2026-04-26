// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

// ZC1024 was retired as a duplicate of ZC1022 (see issue #345). It is
// kept as a no-op stub so legacy `disabled_katas` lists that reference
// it keep parsing; the canonical `let` → `$((...))` guidance fires
// under ZC1022 now.

func TestCheckZC1024(t *testing.T) {
	tests := []struct {
		input    string
		expected []katas.Violation
	}{
		{input: "let x=1+1", expected: []katas.Violation{}},
		{input: "x=$((1+1))", expected: []katas.Violation{}},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1024")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
