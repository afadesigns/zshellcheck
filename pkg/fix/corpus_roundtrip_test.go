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

// parseErrorCount lexes and parses source and reports how many parser
// errors it raises.
func parseErrorCount(source string) int {
	p := parser.New(lexer.New(source))
	p.ParseProgram()
	return len(p.Errors())
}

// fixOnce applies every registry fix edit to source in a single pass and
// returns the rewrite plus the parser-error count of that rewrite. It
// mirrors the collect-and-apply path used by the CLI and the per-kata
// integration tests.
func fixOnce(source string) (string, int, error) {
	program := parser.New(lexer.New(source)).ParseProgram()
	var edits []katas.FixEdit
	ast.Walk(program, func(n ast.Node) bool {
		_, es := katas.Registry.CheckAndFix(n, nil, []byte(source))
		edits = append(edits, es...)
		return true
	})
	out, err := Apply(source, edits)
	if err != nil {
		return "", 0, err
	}
	return out, parseErrorCount(out), nil
}

// roundTripCorpus stresses the fixer on inputs that historically broke
// it — overlapping subscript and arithmetic rewrites (ZC1001 vs ZC1073),
// glob-qualifier stacking (ZC1040), the echo cluster (ZC1030/1037/1092),
// nested backticks (ZC1002 vs ZC1005) — plus a broad sample of fixable
// constructs. Each entry parses cleanly, so a parser error after fixing
// can only mean the fix corrupted the source.
var roundTripCorpus = []string{
	// Subscript wrap vs redundant-$ removal, the v1.3.5 collision class.
	"echo $arr[1]\n",
	"(( total = $counts[key] + 1 ))\n",
	"(( $arr[i] ))\n",
	"x=$arr[1]\necho ${arr[2]}\n",
	// Glob-qualifier stacking — must converge, not append (N) forever.
	"for f in dir/*; do echo $f; done\n",
	"for f in *.log; do print $f; done\n",
	// Echo cluster — three katas, one fix.
	"echo hello world\n",
	"echo -n prompt\n",
	"echo -e \"a\\tb\"\n",
	"echo \"$var\"\n",
	// Backticks, including nesting an external-command rewrite.
	"result=`ls -la`\n",
	"result=`which git`\n",
	"x=`echo hi`\n",
	// let / arithmetic rewrites.
	"let x=5\n",
	"let i+=1\n",
	"let counter=counter+1\n",
	// test / [ ] to [[ ]].
	"[ $x -eq 1 ]\n",
	"[ -f /tmp/foo ] && echo yes\n",
	"[ \"$a\" = \"$b\" ]\n",
	// External commands with Zsh-native alternatives.
	"which git\n",
	"x=$(seq -s , 1 5)\n",
	"cat file | wc -l\n",
	// Read, sensitive read, print forms.
	"read name\n",
	"read -p \"pw: \" secret\n",
	// $(cat f) collapse.
	"x=$(cat /etc/hostname)\n",
	// readonly / typeset.
	"readonly NAME=value\n",
	// Mixed multi-line program exercising several fixes at once.
	"result=`which git`\nfor f in dir/*; do echo $f; done\n[ $x -eq 1 ] && echo $arr[1]\nlet n=n+1\n",
	// Quoted and already-idiomatic forms must stay no-ops.
	"result=$(ls -la)\n",
	"echo ${arr[1]}\n",
	"[[ $x -eq 1 ]]\n",
	"print -r -- hello\n",
}

// maxFixPasses bounds the fixed-point search. The CLI converges its
// multi-pass fixer (applyFixesUntilStable) within a handful of passes;
// a corpus entry that has not stabilised well before this cap is stuck
// in a non-convergent rewrite loop — the v1.3.5 glob-qualifier stacking
// class — and fails the gate.
const maxFixPasses = 10

// fixToFixedPoint applies fixOnce repeatedly until the source stops
// changing or maxFixPasses is reached. It returns the final source,
// whether a fixed point was reached, and the highest parser-error count
// observed on any pass. Iterating to a fixed point mirrors the CLI,
// whose fixer is deliberately multi-pass: one fix can enable another
// (for example `[ $x -eq 1 ]` becomes `(( $x == 1 ))`, which then loses
// its redundant `$`). The interesting invariants are therefore that the
// rewrite converges and never corrupts, not that a single pass is a
// no-op.
func fixToFixedPoint(source string) (string, bool, int, error) {
	cur := source
	worstErrs := parseErrorCount(source)
	for i := 0; i < maxFixPasses; i++ {
		next, errs, err := fixOnce(cur)
		if err != nil {
			return cur, false, worstErrs, err
		}
		if errs > worstErrs {
			worstErrs = errs
		}
		if next == cur {
			return cur, true, worstErrs, nil
		}
		cur = next
	}
	return cur, false, worstErrs, nil
}

// TestFixRoundTrip_NoCorruption_Converges is the always-on autofix
// safety gate. For every corpus entry it asserts the fixer (1) never
// raises the parser-error count on any pass — a fix that turns parseable
// source into unparseable source is a destructive bug — and (2)
// converges to a fixed point, so a stable `-fix` result exists and the
// rewrite never loops. The corpus sweep script
// (scripts/fix-corpus-sweep.sh) extends the same invariants to the full
// pinned upstream corpora in CI.
func TestFixRoundTrip_NoCorruption_Converges(t *testing.T) {
	for _, src := range roundTripCorpus {
		origErrs := parseErrorCount(src)

		final, converged, worstErrs, err := fixToFixedPoint(src)
		if err != nil {
			t.Errorf("apply failed for %q: %v", src, err)
			continue
		}
		if worstErrs > origErrs {
			t.Errorf("fix corrupted source\ninput: %q\nparser errors %d -> %d\nfinal: %q",
				src, origErrs, worstErrs, final)
		}
		if !converged {
			t.Errorf("fix did not converge within %d passes (non-convergent rewrite)\ninput: %q\nfinal: %q",
				maxFixPasses, src, final)
		}
	}
}
