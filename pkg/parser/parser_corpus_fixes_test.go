// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package parser

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/lexer"
)

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

// Exercise every operand form the character-code prefix operator accepts,
// plus the bare-`#` fallback when no operand is glued on.
func TestParseArithmeticCharCodeOperandForms(t *testing.T) {
	cases := []string{
		"(( #name ))\n", // IDENT operand
		"(( #$var ))\n", // VARIABLE operand
		"(( ##c ))\n",   // nested HASH operand
		"(( #${x} ))\n", // ${…} operand
		"(( #0 ))\n",    // INT operand
		"(( #>0 ))\n",   // no operand: bare `#` then `>` (fallback)
		"(( # > 0 ))\n", // bare `#` with spacing
	}
	for _, src := range cases {
		parseSourceClean(t, src)
	}
}

// Malformed C-style for headers exercise the error-return paths of
// parseArithForHeader (a missing `;` or closing `))`). The parser must
// record an error and not panic; the program value is irrelevant here.
func TestParseArithmeticForLoopMalformed(t *testing.T) {
	cases := []string{
		"for ((i=0 i<3; i++)); do :; done\n", // missing first `;`
		"for ((i=0; i<3 i++)); do :; done\n", // missing second `;`
		"for ((i=0; i<3; i++ do :; done\n",   // missing closing `))`
		"for ((&&; i<3; i++)); do :; done\n", // init slot opens on a non-operand
		"for ((i=0; &&; i++)); do :; done\n", // cond slot opens on a non-operand
		"for ((i=0; i<3; &&)); do :; done\n", // post slot opens on a non-operand
	}
	for _, src := range cases {
		p := New(lexer.New(src))
		p.ParseProgram()
		if len(p.Errors()) == 0 {
			t.Fatalf("expected a parser error for malformed for-header %q", src)
		}
	}
}
