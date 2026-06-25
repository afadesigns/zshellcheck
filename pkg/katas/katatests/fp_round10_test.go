// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

// TestZC1075GluedCommandNameNotFlagged pins the glued-command-name false
// positive. A `$var` glued with no separating space to a literal prefix
// (the command-name form `_clear$fsuf`, common in gitstatus) is part of a
// concatenation: the literal prefix keeps the word non-empty, so an empty
// value cannot elide it. The parser splits `name$var` in command position
// into the name plus this glued argument, which previously looked bare.
func TestZC1075GluedCommandNameNotFlagged(t *testing.T) {
	for _, src := range []string{
		"_gitstatus_clear$fsuf",
		"gitstatus_stop$fsuf",
		"prefix$suffix arg",
	} {
		if n := len(testutil.Check(src, "ZC1075")); n != 0 {
			t.Errorf("ZC1075 should not flag a glued-suffix expansion: %q (got %d)", src, n)
		}
	}
	// The glued suffix is skipped, but a following space-separated bare
	// expansion in the same command is a real elision hazard and fires.
	if n := len(testutil.Check("gitstatus_stop$fsuf $name", "ZC1075")); n != 1 {
		t.Errorf("ZC1075 should flag the space-separated `$name` (not the glued `$fsuf`): got %d", n)
	}
	// A space-separated bare expansion still elides and must fire.
	for _, src := range []string{"rm $file", "cd $dir"} {
		if n := len(testutil.Check(src, "ZC1075")); n == 0 {
			t.Errorf("ZC1075 should still flag a bare space-separated expansion: %q", src)
		}
	}
}

// TestZC1075QuoteFlagNotFlagged pins the `(q)`-flag false positive. A
// quoting flag (`q`, `qq`, `q-`, `q+`, or a combination such as `Vq-`)
// renders an empty value as a quoted empty string — a non-empty word — so
// the expansion never elides.
func TestZC1075QuoteFlagNotFlagged(t *testing.T) {
	for _, src := range []string{
		"f ${(q)result}",
		"f ${(qq)x}",
		"f ${(Vq-)result[3]}",
		"print -lr ${(qqqq)input[@]}",
	} {
		if n := len(testutil.Check(src, "ZC1075")); n != 0 {
			t.Errorf("ZC1075 should not flag a (q)-quoted expansion: %q (got %d)", src, n)
		}
	}
}

// TestZC1071SingleSelfReferenceNotFlagged pins that a single-element
// self-reference is an identity reassignment, not an append, so `+=` does
// not apply. The parser drops the `:#`/`//` modifier, so a filter rebuild
// (`arr=(${arr[@]:#pat})`) is indistinguishable from `arr=(${arr[@]})` at
// the AST level; the element-count guard separates a real append (whole
// array plus at least one new element) from both.
func TestZC1071SingleSelfReferenceNotFlagged(t *testing.T) {
	for _, src := range []string{
		"arr=($arr)",
		"del_list=(${del_list[@]:#pattern})",
		"arr=(${arr[@]//a/b})",
	} {
		if n := len(testutil.Check(src, "ZC1071")); n != 0 {
			t.Errorf("ZC1071 should not flag a single-element rebuild: %q (got %d)", src, n)
		}
	}
	// A whole-array reference followed by appended elements is a real
	// append and must still fire.
	for _, src := range []string{"arr=($arr a b)", "arr=(${arr[@]} x)"} {
		if n := len(testutil.Check(src, "ZC1071")); n == 0 {
			t.Errorf("ZC1071 should still flag a true append: %q", src)
		}
	}
}

// TestZC1149StderrNotFlagged pins the stderr false positive. A print/echo
// already routed to stderr (`print -u2`, `print -u 2`, or a `>&2` / `2>`
// redirection) has its error message on the right stream, so the kata has
// nothing to recommend.
func TestZC1149StderrNotFlagged(t *testing.T) {
	for _, src := range []string{
		`print -u2 "Error: boom"`,
		`print -u 2 "Error: boom"`,
		`echo "Error: boom" >&2`,
	} {
		if n := len(testutil.Check(src, "ZC1149")); n != 0 {
			t.Errorf("ZC1149 should not flag an error message already on stderr: %q (got %d)", src, n)
		}
	}
	// An error message on stdout still warrants the redirect advice.
	for _, src := range []string{`echo "Error: boom"`, `print "Error: boom"`} {
		if n := len(testutil.Check(src, "ZC1149")); n == 0 {
			t.Errorf("ZC1149 should still flag an error message on stdout: %q", src)
		}
	}
}

// TestLogicalChainCompoundBodyLinted is the integration proof that the
// parser walk fix makes katas see code inside a `&&` / `||` compound body.
// A dangerous `rm -rf $unquoted` inside `cmd && { … }` was previously
// invisible; it now draws its findings.
func TestLogicalChainCompoundBodyLinted(t *testing.T) {
	for _, src := range []string{
		"(( x > 0 )) && { rm -rf $target }",
		"[[ -n $x ]] && { cp $a $b }",
		"ok && if true; then rm -rf $target; fi",
	} {
		if n := len(testutil.Check(src, "ZC1075")); n == 0 {
			t.Errorf("a body inside a logical-chain compound should be linted: %q drew no ZC1075", src)
		}
	}
}
