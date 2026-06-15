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

// `until` is a Zsh reserved word with the same grammar as `while`. It
// previously lexed as a plain command, so its `(( … ))` condition was
// parsed as an argument and a grouped sub-expression (`( a ) != b`)
// errored. The zsh-z plugin uses this form; sharing the WHILE token
// routes `until` through the loop parser. Each input is `zsh -n` clean.
func TestParseUntilLoop(t *testing.T) {
	cases := []string{
		"until true; do break; done\n",
		"until [[ -f x ]]; do sleep 1; done\n",
		"until (( ( a ) != b )); do :; done\n",
		"until (( ( ${#cd:h} - ${#${${cd:h}//${~q}/}} ) != q_chars )); do :; done\n",
		"until (( ( ${#cd:h} - ${#${${${cd:h}:l}//${~${q:l}}/}} ) != q )); do :; done\n",
		"until cmd; do echo loop; done\n",
	}
	for _, src := range cases {
		parseSourceClean(t, src)
	}
}

// A `[` glued to a path-glob word (a `/` before it) is a glob bracket
// class, not an array subscript. Inside a `for … in` list the INDEX
// infix used to treat the path as the array name and swallow the rest
// of the loop, leaving its `do`/`done` orphaned. The zsh-z and
// powerlevel10k plugins use `$dir/[^[:space:]]##(/N)`. A `$var[idx]`
// with no `/` stays a real subscript.
func TestParseGlobBracketAfterPath(t *testing.T) {
	clean := []string{
		"for x in a/[[:space:]]##; do :; done\n",
		"for x in /[^[:space:]]##; do :; done\n",
		"for plugin in $root/plugins/[^[:space:]]##(/N); do :; done\n",
		"echo $path[1]\n",
		"echo ${arr[1]}\n",
		"echo $arr[$i/2]\n",
	}
	for _, src := range clean {
		parseSourceClean(t, src)
	}
}

// Zsh's `$=name` forces word-splitting on the expansion (the bare-`$`
// counterpart of `${=name}`). The split flag is a single `=`
// (token.ASSIGN), not `==`; the dollar-flag dispatch checked the wrong
// token, so `$=line` errored with "expected IDENT". Sibling forms
// `$^name` and `$~name` already worked. p10k uses `local w=($=line)`.
func TestParseDollarForcedSplitFlag(t *testing.T) {
	clean := []string{
		"echo $=var\n",
		"local words=($=line)\n",
		"local header=($=lines[1])\n",
		"for a b in $=x[1]; do :; done\n",
		"x=( $=foo )\n",
	}
	for _, src := range clean {
		parseSourceClean(t, src)
	}
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

// A `( … )` subshell as the final operand of an `if`/`while` condition
// (`if [[ a ]] && ( [[ b ]] )`) followed by `then`/`do` on a NEW line
// must parse as a single compound statement. The subshell leaves
// curToken on its `)`; the condition block's RPAREN terminator (for the
// `if ( cond ) cmd` shortcut) used to break there, orphaning the
// `then`/body/`fi` into separate top-level statements (no error, but a
// wrong AST). The prezto terminal module uses this form. Issue:
// multi-line-condition leak.
func TestParseMultilineCondTrailingSubshell(t *testing.T) {
	cases := []string{
		"if [[ a ]] && ( [[ b ]] )\nthen\n  :\nfi\n",
		"if [[ a ]] && ( ! [[ b ]] )\nthen\n  :\nfi\n",
		"if [[ $T == X ]] \\\n  && ( ! [[ -n \"$S\" || -n \"$M\" ]] )\nthen\n  echo hi\nfi\n",
		"while [[ a ]] && ( [[ b ]] )\ndo\n  :\ndone\n",
	}
	for _, src := range cases {
		p := New(lexer.New(src))
		prog := p.ParseProgram()
		if len(p.Errors()) != 0 {
			t.Errorf("unexpected errors for %q: %v", src, p.Errors())
		}
		if len(prog.Statements) != 1 {
			t.Errorf("want 1 compound statement for %q, got %d (the condition's subshell orphaned the body)", src, len(prog.Statements))
		}
	}
}

// A redirection trailing a brace-group statement (`{ … } 2>/dev/null`)
// was left unconsumed on the statement-dispatch path, orphaning the
// redirect into a bogus `(2 > /dev/null)` statement that swallowed the
// next command. zinit's install path uses the `} 2>/dev/null | while`
// form. Each input is `zsh -n` clean and is one statement plus the
// trailing `print`.
func TestParseBraceGroupTrailingRedirection(t *testing.T) {
	cases := map[string]int{
		"{\n  print x\n} 2>/dev/null\nprint after\n":                                                2,
		"{\n  print x\n} >/dev/null\nprint after\n":                                                 2,
		"{\n  print x\n} >>log\nprint after\n":                                                      2,
		"foo() {\n  {\n    print\n  } 2>/dev/null | while read -r l; do\n    print $l\n  done\n}\n": 1,
	}
	for src, want := range cases {
		p := New(lexer.New(src))
		prog := p.ParseProgram()
		if len(p.Errors()) != 0 {
			t.Errorf("unexpected errors for %q: %v", src, p.Errors())
		}
		if len(prog.Statements) != want {
			t.Errorf("want %d statements for %q, got %d (the trailing redirect orphaned)", want, src, len(prog.Statements))
		}
	}
}

// A reserved closer word (`done`, `fi`, `esac`, `then`, `do`, `else`,
// `elif`) is a literal argument in Zsh when it is not in command
// position. `echo done` used to stop argument gathering at `done` and
// orphan it into a second bogus statement; it now parses as one command
// with the closer captured as an argument.
func TestParseReservedWordAsArgument(t *testing.T) {
	cases := []string{
		"echo done\n",
		"echo fi esac then\n",
		"print -l function do done\n",
		"local x=done\n",
		"args=(do done fi)\n",
	}
	for _, src := range cases {
		p := New(lexer.New(src))
		prog := p.ParseProgram()
		if len(p.Errors()) != 0 {
			t.Errorf("unexpected errors for %q: %v", src, p.Errors())
		}
		if len(prog.Statements) != 1 {
			t.Errorf("want 1 statement for %q, got %d (the reserved word orphaned)", src, len(prog.Statements))
		}
	}
}

// A process substitution as the first argument (`diff <(a) <(b)`) fell
// to the expression path, which parsed only the bare command name and
// orphaned each `<(…)` into its own bogus top-level statement. Every
// input is a single SimpleCommand and `zsh -n` clean.
func TestParseProcessSubstitutionFirstArg(t *testing.T) {
	cases := []string{
		"diff <(echo a) <(echo b)\n",
		"cat <(gen)\n",
		"source <(zstyle -L | awk '{print $2}')\n",
		"tee >(logger)\n",
	}
	for _, src := range cases {
		p := New(lexer.New(src))
		prog := p.ParseProgram()
		if len(p.Errors()) != 0 {
			t.Errorf("unexpected errors for %q: %v", src, p.Errors())
		}
		if len(prog.Statements) != 1 {
			t.Errorf("want 1 statement for %q, got %d (the process substitution orphaned)", src, len(prog.Statements))
		}
	}
}

// A bare `X=` at end of line is a standalone empty assignment. Zsh does
// not continue an assignment RHS across an unescaped newline, so the
// following line is a separate statement — the `=` infix used to swallow
// the next line's command (`X = print`) and orphan its arguments,
// leaving them unlinted. Each input is `zsh -n` clean and two statements.
func TestParseEmptyAssignmentNewline(t *testing.T) {
	cases := []string{
		"X=\nprint $A\n",
		"RPROMPT=\nprint $A\n",
		"local x=\nprint $y\n",
		// The RHS still attaches on the same line and across a `\` join.
		"X=value\nprint $A\n",
		"X=\\\nvalue\nprint $A\n",
	}
	for _, src := range cases {
		p := New(lexer.New(src))
		prog := p.ParseProgram()
		if len(p.Errors()) != 0 {
			t.Errorf("unexpected errors for %q: %v", src, p.Errors())
		}
		if len(prog.Statements) != 2 {
			t.Errorf("want 2 statements for %q, got %d (the empty assignment swallowed the next line)", src, len(prog.Statements))
		}
	}
}

// The Zsh-only loops `repeat <count> …` and `foreach <name> (<list>) …
// end` had no grammar, so `repeat`/`foreach` lexed as plain commands and
// the `do`/`done`/`end`/`}` body markers orphaned. Each case is one
// statement (the loop) plus any trailing command, and `zsh -n` clean.
func TestParseRepeatAndForeachLoops(t *testing.T) {
	cases := map[string]int{
		"repeat 3; do print a; done\n":                                       1,
		"repeat 3 do print a; done\n":                                        1,
		"repeat 3 { print a }\n":                                             1,
		"repeat $n print a\n":                                                1,
		"repeat 2; do print a; done\nprint after\n":                          2,
		"foreach f (a b)\n  rm $f\nend\n":                                    1,
		"foreach x (1 2 3); print $x; end\n":                                 1,
		"foreach f (*.txt)\n  rm $f\nend\nprint after\n":                     2,
		"foreach x (1 2)\n  if [[ -n $x ]]; then\n    print $x\n  fi\nend\n": 1,
		// The `in <list>` spelling of foreach.
		"foreach x in a b c; print $x; end\n": 1,
		"foreach x in $items\n  rm $x\nend\n": 1,
		// `end` stays a valid identifier outside foreach.
		"end() { print x }\nend\n": 2,
		"local end=5\n":            1,
	}
	for src, want := range cases {
		p := New(lexer.New(src))
		prog := p.ParseProgram()
		if len(p.Errors()) != 0 {
			t.Errorf("unexpected errors for %q: %v", src, p.Errors())
		}
		if len(prog.Statements) != want {
			t.Errorf("want %d statements for %q, got %d", want, src, len(prog.Statements))
		}
	}
}
