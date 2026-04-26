// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package parser

import "testing"

func TestParseIndexExpressionFlagR(t *testing.T) {
	parseSourceClean(t, "echo ${arr[(R)x]}\n")
}

func TestParseIndexExpressionFlagLowerR(t *testing.T) {
	parseSourceClean(t, "echo ${arr[(r)pat]}\n")
}

func TestParseIndexExpressionFlagI(t *testing.T) {
	parseSourceClean(t, "echo ${arr[(I)i]}\n")
}

func TestParseIndexExpressionMultiFlag(t *testing.T) {
	parseSourceClean(t, "echo ${arr[(ri)pat]}\n")
}

func TestParseIndexExpressionPlain(t *testing.T) {
	parseSourceClean(t, "echo ${arr[1]}\n")
}
