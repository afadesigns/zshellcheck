// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package parser

import "testing"

func TestParseKeywordReturnInLogicalChain(t *testing.T) {
	parseSourceClean(t, "true || return 0\n")
}

func TestParseKeywordReturnAfterAnd(t *testing.T) {
	parseSourceClean(t, "test -f x && return 1\n")
}

func TestParseKeywordReturnBare(t *testing.T) {
	parseSourceClean(t, "false || return\n")
}
