// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import "testing"

func TestZC1075HasWidthFlag(t *testing.T) {
	cases := map[string]bool{
		"l:5:": true,  // left-pad
		"r:3:": true,  // right-pad
		"r":    false, // reverse, no width arg
		"":     false, // no flags
		"s:,:": false, // split delimiter, not a width flag
	}
	for flags, want := range cases {
		if got := zc1075HasWidthFlag(flags); got != want {
			t.Errorf("zc1075HasWidthFlag(%q) = %v, want %v", flags, got, want)
		}
	}
}
