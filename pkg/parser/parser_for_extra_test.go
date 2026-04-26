// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package parser

import "testing"

func TestParseForLoopArithmeticConditionOnly(t *testing.T) {
	parseSourceClean(t, "for ((i=0; i<3; )) do echo $i; done\n")
}

func TestParseForLoopArithmeticInitOnly(t *testing.T) {
	parseSourceClean(t, "for ((i=0; ; )) do break; done\n")
}

func TestParseForLoopMultiVariable(t *testing.T) {
	parseSourceClean(t, "for k v in a 1 b 2; do echo $k $v; done\n")
}

func TestParseForLoopShortForm(t *testing.T) {
	parseSourceClean(t, "for f (a b c) echo $f\n")
}

func TestParseForLoopImplicitList(t *testing.T) {
	parseSourceClean(t, "for k v w; do echo $k; done\n")
}

func TestParseForLoopNumericName(t *testing.T) {
	parseSourceClean(t, "for 1 in a b c; do echo $1; done\n")
}

func TestParseIfStatementInline(t *testing.T) {
	parseSourceClean(t, "if [ -f f ]; then echo y; fi\n")
}

func TestParseIfStatementWithCommandSubstitution(t *testing.T) {
	parseSourceClean(t, "if [[ $(date +%H) -lt 12 ]]; then echo morning; fi\n")
}

func TestParseSubshellStatement(t *testing.T) {
	parseSourceClean(t, "(cd /tmp && rm -rf foo)\n")
}

func TestParseSubshellWithPipeline(t *testing.T) {
	parseSourceClean(t, "(echo a; echo b) | sort\n")
}

func TestParseDeclarationValueArray(t *testing.T) {
	parseSourceClean(t, "x=(1 2 3 4)\n")
}

func TestParseDeclarationValueAssoc(t *testing.T) {
	parseSourceClean(t, "typeset -A m=(k1 v1 k2 v2)\n")
}
