// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package parser

import "testing"

// Regression tests for parser gaps surfaced by the pinned integration
// corpora. Each input is a minimal form of a construct that previously
// produced a spurious parser error.

func TestParseCaseSubjectConcatenation(t *testing.T) {
	// A case subject is a shell word: it can concatenate expansions and
	// literals with no separating space. Previously the parser stopped
	// at the first `/` or `:` and reported the tail as unexpected.
	cases := []string{
		"case $state/$line[1] in\n  a) echo x ;;\nesac\n",
		"case ${variant}:${${service#ping}:-4} in\n  4) echo v4 ;;\nesac\n",
		"case $variant:$OSTYPE in\n  *) echo o ;;\nesac\n",
		"case $x in\n  a) echo a ;;\nesac\n",
	}
	for _, src := range cases {
		parseSourceClean(t, src)
	}
}

func TestParseArithmeticCharCodeOperator(t *testing.T) {
	// Inside `((…))`, `#name` / `##c` is the character-code prefix
	// operator. A bare `#` keeps its positional-count meaning.
	cases := []string{
		"(( #disk_name ))\n",
		"if (( #exts != 0 )); then echo y; fi\n",
		"(( # > 0 ))\n",
		"(( ! # ))\n",
	}
	for _, src := range cases {
		parseSourceClean(t, src)
	}
}

func TestParseFunctionNameWithPositional(t *testing.T) {
	// A `function` name can glue in a positional parameter, e.g.
	// `function _$0_fmt() { … }` (the lexer emits `$0` as DOLLAR + INT).
	parseSourceClean(t, "function _$0_fmt() {\n  echo hi\n}\n")
}
