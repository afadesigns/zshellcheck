// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package fix

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/ast"
	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/lexer"
	"github.com/afadesigns/zshellcheck/pkg/parser"
)

// isolatedFix applies only the named kata's fix to source, ignoring
// every other kata. Several fixable katas overlap on the same construct
// (the `test` rewrites ZC1006/ZC1020/ZC1036, the `let` rewrites
// ZC1008/ZC1022); collecting every edit through runFix would assert the
// cross-kata cascade instead of the one fix under test. Filtering the
// violations to a single KataID pins each kata's own rewrite.
func isolatedFix(t *testing.T, source, kataID string) string {
	t.Helper()
	program := parser.New(lexer.New(source)).ParseProgram()
	var edits []katas.FixEdit
	ast.Walk(program, func(n ast.Node) bool {
		for _, v := range katas.Registry.Check(n, nil) {
			if v.KataID == kataID {
				edits = append(edits, katas.Registry.FixesFor(n, v, []byte(source))...)
			}
		}
		return true
	})
	out, err := Apply(source, edits)
	if err != nil {
		t.Fatalf("apply %s: %v", kataID, err)
	}
	return out
}

// TestFixIsolated_GapKatas pins the exact rewrite of seven fixable katas
// that previously had no golden fix-output test, only detection and
// panic-safety coverage. Each entry asserts the kata's own fix in
// isolation so an overlapping kata cannot mask a regression.
func TestFixIsolated_GapKatas(t *testing.T) {
	cases := []struct {
		name string
		kata string
		src  string
		want string
	}{
		{"ZC1006_test_numeric", "ZC1006", "test 1 -eq 1\n", "[[ 1 -eq 1 ]]\n"},
		// The spaced `let x = 1` form carries its surrounding whitespace
		// into the rewrite, so the operator keeps padded spaces. The
		// result is correct and idempotent; tightening the spacing is a
		// cosmetic polish tracked separately.
		{"ZC1008_let_arith", "ZC1008", "let x = 1\n", "(( x  =  1 ))\n"},
		{"ZC1020_test_to_bracket", "ZC1020", "test 1 -eq 1\n", "[[ 1 -eq 1 ]]\n"},
		{"ZC1022_let_expr", "ZC1022", "let x=1+1\n", "(( x = 1+1 ))\n"},
		{"ZC1036_test_file", "ZC1036", "test -f file.txt\n", "[[ -f file.txt ]]\n"},
		{"ZC1037_echo_var", "ZC1037", "echo $foo\n", "print -r -- $foo\n"},
		{"ZC1217_service", "ZC1217", "service nginx start\n", "systemctl start nginx\n"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := isolatedFix(t, tc.src, tc.kata)
			if got != tc.want {
				t.Errorf("%s fix\ninput: %q\ngot:   %q\nwant:  %q", tc.kata, tc.src, got, tc.want)
			}
			// The rewrite must itself parse cleanly.
			if n := parseErrorCount(got); n != 0 {
				t.Errorf("%s fix produced %d parser error(s): %q", tc.kata, n, got)
			}
		})
	}
}
