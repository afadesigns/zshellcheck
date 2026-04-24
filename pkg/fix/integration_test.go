package fix

import (
	"strings"
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/ast"
	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/lexer"
	"github.com/afadesigns/zshellcheck/pkg/parser"
)

// runFix walks source through lex + parse, collects fix edits from the
// shipped registry, applies them, and returns the rewritten source.
// It is used by each kata-level integration test below; keeping it
// centralised means the tests stay small and identical in shape.
func runFix(t *testing.T, source string) string {
	t.Helper()
	l := lexer.New(source)
	p := parser.New(l)
	program := p.ParseProgram()
	if errs := p.Errors(); len(errs) != 0 {
		t.Fatalf("parse errors: %v", errs)
	}
	var edits []katas.FixEdit
	ast.Walk(program, func(n ast.Node) bool {
		_, es := katas.Registry.CheckAndFix(n, nil, []byte(source))
		edits = append(edits, es...)
		return true
	})
	out, err := Apply(source, edits)
	if err != nil {
		t.Fatalf("apply: %v", err)
	}
	return out
}

func TestFixIntegration_ZC1002_Backticks(t *testing.T) {
	src := "result=`ls -la`\n"
	want := "result=$(ls -la)\n"
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFixIntegration_ZC1005_Which(t *testing.T) {
	src := "which git\n"
	want := "whence git\n"
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFixIntegration_ZC1092_Echo(t *testing.T) {
	src := "echo hello world\n"
	got := runFix(t, src)
	// ZC1092 replaces just the command name; arguments stay intact.
	if !strings.HasPrefix(got, "print -r --") {
		t.Errorf("expected fix to replace echo, got %q", got)
	}
	if !strings.Contains(got, "hello world") {
		t.Errorf("fix mangled arguments: %q", got)
	}
}

func TestFixIntegration_EchoFlag_LeftAlone(t *testing.T) {
	// ZC1092 Fix skips flagged forms — `echo -n` translates to a
	// different print invocation than `print -r --`.
	src := "echo -n keep\n"
	if got := runFix(t, src); got != src {
		t.Errorf("flagged echo should not be auto-fixed, got %q", got)
	}
}

func TestFixIntegration_NestedKatas_OuterWins(t *testing.T) {
	// Outer backtick-span fix (ZC1002) wraps an inner `which` that
	// ZC1005 would rewrite. The outer edit wins on the first pass —
	// the output preserves `which` inside the new parens. A second
	// `-fix` pass would then convert `which` -> `whence`.
	src := "result=`which git`\n"
	want := "result=$(which git)\n"
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFixIntegration_ZC1010_TestToDoubleBracket(t *testing.T) {
	src := `if [ -f /tmp/foo ]; then
  :
fi
[ "$x" = "y" ] && :
`
	want := `if [[ -f /tmp/foo ]]; then
  :
fi
[[ "$x" = "y" ]] && :
`
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFixIntegration_ZC1013_LetToArith(t *testing.T) {
	src := `let x=5
let y=$((x + 1))
let counter=counter+1
`
	want := `(( x = 5 ))
(( y = $((x + 1)) ))
(( counter = counter+1 ))
`
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFixIntegration_ZC1004_ExitToReturn(t *testing.T) {
	src := `foo() {
  if [[ -z "$1" ]]; then
    exit 1
  fi
  exit
}
`
	want := `foo() {
  if [[ -z "$1" ]]; then
    return 1
  fi
  return
}
`
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFixIntegration_ZC1001_BraceArrayAccess(t *testing.T) {
	src := `x=$arr[1]
y=$other[2]
z=$pair[a,b]
`
	want := `x=${arr[1]}
y=${other[2]}
z=${pair[a,b]}
`
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFixIntegration_ZC1003_TestToArith(t *testing.T) {
	src := `if [ $count -eq 0 ]; then
  :
fi
[ "$n" -lt 10 ] && :
[ "$a" -ge "$b" ] || :
`
	want := `if (( $count == 0 )); then
  :
fi
(( "$n" < 10 )) && :
(( "$a" >= "$b" )) || :
`
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFixIntegration_ZC1012_ReadAddR(t *testing.T) {
	src := "read VAR\n"
	want := "read -r VAR\n"
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFixIntegration_ZC1012_ReadWithFlagsPreserved(t *testing.T) {
	src := "read -p prompt VAR\n"
	want := "read -r -p prompt VAR\n"
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFixIntegration_ZC1012_ReadDashRLeftAlone(t *testing.T) {
	src := "read -r VAR\n"
	if got := runFix(t, src); got != src {
		t.Errorf("already-fixed input should be idempotent, got %q", got)
	}
}

func TestFixIntegration_ZC1031_ShebangEnv(t *testing.T) {
	src := "#!/bin/zsh\nx=1\n"
	want := "#!/usr/bin/env zsh\nx=1\n"
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFixIntegration_ZC1031_AlreadyEnv(t *testing.T) {
	src := "#!/usr/bin/env zsh\nx=1\n"
	if got := runFix(t, src); got != src {
		t.Errorf("already-portable shebang should be left alone, got %q", got)
	}
}

func TestFixIntegration_ZC1055_EmptyCheckToDashZ(t *testing.T) {
	src := `[[ $x == "" ]]` + "\n"
	want := `[[ -z $x ]]` + "\n"
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFixIntegration_ZC1055_NonEmptyCheckToDashN(t *testing.T) {
	src := `[[ $x != "" ]]` + "\n"
	want := `[[ -n $x ]]` + "\n"
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFixIntegration_ZC1055_LeftEmpty(t *testing.T) {
	src := `[[ "" == $x ]]` + "\n"
	want := `[[ -z $x ]]` + "\n"
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFixIntegration_ZC1017_PrintAddR(t *testing.T) {
	src := `print "hello"` + "\n"
	want := `print -r "hello"` + "\n"
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFixIntegration_ZC1017_PrintRLeftAlone(t *testing.T) {
	src := `print -r "hello"` + "\n"
	if got := runFix(t, src); got != src {
		t.Errorf("already-fixed input should be idempotent, got %q", got)
	}
}

func TestFixIntegration_ZC1051_RmUnquotedVar(t *testing.T) {
	src := "rm $FILE\n"
	want := `rm "$FILE"` + "\n"
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFixIntegration_ZC1051_RmQuotedStaysIdempotent(t *testing.T) {
	src := `rm "$FILE"` + "\n"
	if got := runFix(t, src); got != src {
		t.Errorf("quoted input should be idempotent, got %q", got)
	}
}

func TestFixIntegration_ZC1062_EgrepToGrepE(t *testing.T) {
	src := "egrep pattern file\n"
	want := "grep -E pattern file\n"
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFixIntegration_ZC1063_FgrepToGrepF(t *testing.T) {
	src := "fgrep pattern file\n"
	want := "grep -F pattern file\n"
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFixIntegration_ZC1086_FunctionKeywordBareBody(t *testing.T) {
	src := "function foo { body }\n"
	want := "foo() { body }\n"
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFixIntegration_ZC1086_FunctionKeywordWithParens(t *testing.T) {
	src := "function foo() { body }\n"
	want := "foo() { body }\n"
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFixIntegration_ZC1086_NoFunctionKeywordUnchanged(t *testing.T) {
	src := "foo() { body }\n"
	if got := runFix(t, src); got != src {
		t.Errorf("plain form should be idempotent, got %q", got)
	}
}

func TestFixIntegration_ZC1073_DropDollarInArith(t *testing.T) {
	src := "(( $x > 0 ))\n"
	want := "(( x > 0 ))\n"
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFixIntegration_ZC1073_MultipleDollarsAllDropped(t *testing.T) {
	src := "(( $a + $b ))\n"
	want := "(( a + b ))\n"
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFixIntegration_ZC1064_TypeToCommandV(t *testing.T) {
	src := "type git\n"
	want := "command -v git\n"
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFixIntegration_ZC1061_SeqSingleArg(t *testing.T) {
	// ZC1085 also fires and wraps the for-loop item in quotes; the
	// combined output shows both rewrites applied in one pass.
	src := "for i in $(seq 5); do :; done\n"
	want := `for i in "$({1..5})"; do :; done` + "\n"
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFixIntegration_ZC1061_SeqTwoArgs(t *testing.T) {
	src := "for i in $(seq 3 8); do :; done\n"
	want := `for i in "$({3..8})"; do :; done` + "\n"
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFixIntegration_ZC1061_SeqThreeArgs(t *testing.T) {
	src := "for i in $(seq 1 2 10); do :; done\n"
	want := `for i in "$({1..10..2})"; do :; done` + "\n"
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFixIntegration_ZC1061_SeqVariableArgsSkipped(t *testing.T) {
	src := "seq $N\n"
	if got := runFix(t, src); got != src {
		t.Errorf("variable arg should be left alone, got %q", got)
	}
}

func TestFixIntegration_ZC1079_QuoteRhsInBrackets(t *testing.T) {
	src := `[[ $x == $y ]]` + "\n"
	want := `[[ $x == "$y" ]]` + "\n"
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFixIntegration_ZC1079_AlreadyQuotedRhsUnchanged(t *testing.T) {
	src := `[[ $x == "$y" ]]` + "\n"
	if got := runFix(t, src); got != src {
		t.Errorf("quoted RHS should be idempotent, got %q", got)
	}
}

func TestFixIntegration_ZC1085_QuoteForLoopExpansion(t *testing.T) {
	src := "for f in $files; do :; done\n"
	want := `for f in "$files"; do :; done` + "\n"
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFixIntegration_ZC1085_ArrayExpansionQuoted(t *testing.T) {
	src := "for f in ${files[@]}; do :; done\n"
	want := `for f in "${files[@]}"; do :; done` + "\n"
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFixIntegration_ZC1085_AlreadyQuotedUnchanged(t *testing.T) {
	src := `for f in "$files"; do :; done` + "\n"
	if got := runFix(t, src); got != src {
		t.Errorf("quoted input should be idempotent, got %q", got)
	}
}

func TestFixIntegration_SecondPass_ResolvesInner(t *testing.T) {
	src := "result=`which git`\n"
	first := runFix(t, src)
	final := runFix(t, first)
	want := "result=$(whence git)\n"
	if final != want {
		t.Errorf("got %q, want %q", final, want)
	}
}
