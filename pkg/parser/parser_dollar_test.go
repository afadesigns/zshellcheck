// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package parser

import "testing"

func TestParseBareDollarBeforeSemicolon(t *testing.T) {
	parseSourceClean(t, "echo $;\n")
}

func TestParseBareDollarBeforePipe(t *testing.T) {
	parseSourceClean(t, "echo $ | wc\n")
}

func TestParseBareDollarBeforeRparen(t *testing.T) {
	parseSourceClean(t, "x=( $ )\n")
}

func TestParseDollarBracketArithmetic(t *testing.T) {
	parseSourceClean(t, "echo $[1+2]\n")
}

func TestParseDollarHashLength(t *testing.T) {
	parseSourceClean(t, "echo $#name\n")
}

func TestParseDollarInteger(t *testing.T) {
	parseSourceClean(t, "echo $1\n")
}

func TestParseDollarBang(t *testing.T) {
	parseSourceClean(t, "echo $!\n")
}

func TestParseDollarMinus(t *testing.T) {
	parseSourceClean(t, "echo $-\n")
}

func TestParseDollarStar(t *testing.T) {
	parseSourceClean(t, "echo $*\n")
}

func TestParseDollarPlusName(t *testing.T) {
	parseSourceClean(t, "(( $+name ))\n")
}

func TestParseDollarPlusNameSubscript(t *testing.T) {
	parseSourceClean(t, "(( $+name[key] ))\n")
}
