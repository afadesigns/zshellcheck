// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package parser

import "testing"

func TestParseTypesetArray(t *testing.T) {
	parseSourceClean(t, "typeset -a arr=(a b c)\n")
}

func TestParseLocalAssign(t *testing.T) {
	parseSourceClean(t, "local x=1\n")
}

func TestParseDeclareInteger(t *testing.T) {
	parseSourceClean(t, "declare -i n=42\n")
}

func TestParseReadonly(t *testing.T) {
	parseSourceClean(t, "readonly y=hello\n")
}

func TestParseTypesetAssoc(t *testing.T) {
	parseSourceClean(t, "typeset -A m=(k1 v1 k2 v2)\n")
}
