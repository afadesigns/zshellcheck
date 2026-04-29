// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import "testing"

// TestOffsetLineColHelpers exercises every offsetLineColZCxxxx duplicate
// over the multi-line newline branch. The helpers are unexported and
// shaped identically; this test ensures coverage reports record both
// the offset-out-of-range guard and the newline-step branch for every
// kata-specific copy.
func TestOffsetLineColHelpers(t *testing.T) {
	src := []byte("a\nb\nc")
	for _, fn := range []func([]byte, int) (int, int){
		offsetLineColZC1016,
		offsetLineColZC1040,
		offsetLineColZC1051,
		offsetLineColZC1053,
		offsetLineColZC1076,
		offsetLineColZC1078,
		offsetLineColZC1079,
		offsetLineColZC1084,
		offsetLineColZC1085,
		offsetLineColZC1086,
		offsetLineColZC1091,
		offsetLineColZC1126,
		offsetLineColZC1146,
		offsetLineColZC1147,
		offsetLineColZC1163,
		offsetLineColZC1170,
		offsetLineColZC1190,
		offsetLineColZC1209,
		offsetLineColZC1210,
		offsetLineColZC1213,
		offsetLineColZC1226,
		offsetLineColZC1227,
		offsetLineColZC1230,
		offsetLineColZC1231,
		offsetLineColZC1234,
		offsetLineColZC1238,
		offsetLineColZC1241,
		offsetLineColZC1253,
		offsetLineColZC1255,
		offsetLineColZC1257,
		offsetLineColZC1265,
		offsetLineColZC1267,
		offsetLineColZC1268,
		offsetLineColZC1273,
		offsetLineColZC1293,
		offsetLineColZC1377,
		offsetLineColZC1378,
		offsetLineColZC1380,
		offsetLineColZC1381,
		offsetLineColZC1382,
		offsetLineColZC1383,
		offsetLineColZC1394,
		offsetLineColZC1403,
		offsetLineColZC1404,
		offsetLineColZC1448,
		offsetLineColZC1502,
		offsetLineColZC1643,
		offsetLineColZC1717,
		offsetLineColZC1773,
	} {
		// out-of-range guard
		if l, c := fn(nil, -1); l != -1 || c != -1 {
			t.Errorf("expected (-1,-1) for negative offset, got (%d,%d)", l, c)
		}
		if l, c := fn(src, len(src)+1); l != -1 || c != -1 {
			t.Errorf("expected (-1,-1) for over-range, got (%d,%d)", l, c)
		}
		// newline-step branch (`a\nb\nc` at offset 4 → line 3, col 1)
		if l, c := fn(src, 4); l != 3 || c != 1 {
			t.Errorf("expected (3,1) at offset 4, got (%d,%d)", l, c)
		}
	}
}
