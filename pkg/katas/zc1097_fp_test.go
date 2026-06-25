// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import "testing"

// TestZC1097IsPositionalName covers every branch of the positional-name
// predicate, including the empty-string guard the loop-variable path
// never reaches (a for-loop always names a variable).
func TestZC1097IsPositionalName(t *testing.T) {
	cases := map[string]bool{
		"":   false,
		"1":  true,
		"12": true,
		"x":  false,
		"1a": false,
		"a1": false,
	}
	for name, want := range cases {
		if got := zc1097IsPositionalName(name); got != want {
			t.Errorf("zc1097IsPositionalName(%q) = %v, want %v", name, got, want)
		}
	}
}
