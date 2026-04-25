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
	// ZC1034 + ZC1271 also fire on this input and both rewrite to
	// `command -v`. Their edits arrive ahead of ZC1005's `whence` swap
	// in walk order (ExpressionStatement parent before SimpleCommand
	// child) so the conflict resolver keeps the `command -v` form. The
	// rewrite remains deterministic and idempotent on a re-run.
	src := "which git\n"
	want := "command -v git\n"
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

func TestFixIntegration_EchoFlag_HandledByZC1118(t *testing.T) {
	// ZC1092 skips flagged `echo` because `print -r --` is the wrong
	// rewrite for `-n`. ZC1118 picks it up and rewrites to the
	// matching `print -rn` form instead.
	src := "echo -n keep\n"
	want := "print -rn keep\n"
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
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

func TestFixIntegration_ZC1084_FindGlobQuoted(t *testing.T) {
	src := "find . -name *.txt\n"
	want := "find . -name '*.txt'\n"
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFixIntegration_ZC1084_FindAlreadyQuotedUnchanged(t *testing.T) {
	src := "find . -name '*.txt'\n"
	if got := runFix(t, src); got != src {
		t.Errorf("quoted pattern should be idempotent, got %q", got)
	}
}

func TestFixIntegration_ZC1078_QuoteDollarAt(t *testing.T) {
	src := "cmd $@\n"
	want := `cmd "$@"` + "\n"
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFixIntegration_ZC1078_QuoteDollarStar(t *testing.T) {
	src := "cmd $*\n"
	want := `cmd "$*"` + "\n"
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFixIntegration_ZC1078_AlreadyQuotedUnchanged(t *testing.T) {
	src := `cmd "$@"` + "\n"
	if got := runFix(t, src); got != src {
		t.Errorf("quoted input should be idempotent, got %q", got)
	}
}

func TestFixIntegration_ZC1076_AutoloadAddUZ(t *testing.T) {
	src := "autoload compinit\n"
	want := "autoload -Uz compinit\n"
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFixIntegration_ZC1076_AutoloadWithUZUnchanged(t *testing.T) {
	src := "autoload -Uz compinit\n"
	if got := runFix(t, src); got != src {
		t.Errorf("already-flagged autoload should be idempotent, got %q", got)
	}
}

func TestFixIntegration_ZC1040_AppendNullGlob(t *testing.T) {
	src := "for f in *.txt; do :; done\n"
	want := "for f in *.txt(N); do :; done\n"
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFixIntegration_ZC1040_AlreadyQualifiedUnchanged(t *testing.T) {
	src := "for f in *.txt(N); do :; done\n"
	if got := runFix(t, src); got != src {
		t.Errorf("already-qualified glob should be idempotent, got %q", got)
	}
}

func TestFixIntegration_ZC1091_BracketCmpToArith(t *testing.T) {
	src := "[[ x -lt 10 ]]\n"
	want := "(( x < 10 ))\n"
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFixIntegration_ZC1091_GeBracketCmp(t *testing.T) {
	src := "[[ a -ge b ]]\n"
	want := "(( a >= b ))\n"
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFixIntegration_ZC1091_MultipleOpsLeftAlone(t *testing.T) {
	// Ambiguous: two comparison operators — fix yields to avoid
	// corrupting the expression.
	src := "[[ a -lt b && c -gt d ]]\n"
	if got := runFix(t, src); got != src {
		t.Errorf("multi-op bracket should be left alone, got %q", got)
	}
}

func TestFixIntegration_ZC1126_SortPipeUniq(t *testing.T) {
	src := "sort file | uniq\n"
	want := "sort -u file\n"
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFixIntegration_ZC1126_SortNoArgsPipeUniq(t *testing.T) {
	src := "sort | uniq\n"
	want := "sort -u\n"
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFixIntegration_ZC1126_UniqCountUnchanged(t *testing.T) {
	// ZC1126 detector skips `uniq -c` etc; fix also stays silent.
	src := "sort | uniq -c\n"
	if got := runFix(t, src); got != src {
		t.Errorf("uniq with flags should be left alone, got %q", got)
	}
}

func TestFixIntegration_ZC1118_EchoDashNToPrintRN(t *testing.T) {
	src := "echo -n hello\n"
	want := "print -rn hello\n"
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFixIntegration_ZC1124_CatDevNullTruncate(t *testing.T) {
	src := "cat /dev/null > file\n"
	want := ": > file\n"
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFixIntegration_ZC1128_TouchToEmptyRedirect(t *testing.T) {
	src := "touch file\n"
	want := "> file\n"
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFixIntegration_ZC1128_TouchWithFlagsUnchanged(t *testing.T) {
	src := "touch -t 202504240000 file\n"
	if got := runFix(t, src); got != src {
		t.Errorf("flagged touch should stay, got %q", got)
	}
}

func TestFixIntegration_ZC1147_MkdirAddParentFlag(t *testing.T) {
	src := "mkdir a/b/c\n"
	want := "mkdir -p a/b/c\n"
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFixIntegration_ZC1147_FlatPathUnchanged(t *testing.T) {
	// Detector only fires on nested paths; flat mkdir stays.
	src := "mkdir foo\n"
	if got := runFix(t, src); got != src {
		t.Errorf("flat mkdir should stay, got %q", got)
	}
}

func TestFixIntegration_ZC1140_HashToCommandV(t *testing.T) {
	src := "hash git\n"
	want := "command -v git\n"
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFixIntegration_ZC1144_TrapSignalNumberToName(t *testing.T) {
	src := "trap 'cleanup' 15\n"
	want := "trap 'cleanup' TERM\n"
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFixIntegration_ZC1144_MultipleSignalsRewritten(t *testing.T) {
	src := "trap 'cleanup' 2 15\n"
	want := "trap 'cleanup' INT TERM\n"
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFixIntegration_ZC1162_CpRecursiveToArchive(t *testing.T) {
	src := "cp -r src dest\n"
	want := "cp -a src dest\n"
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFixIntegration_ZC1162_CpCapitalRToArchive(t *testing.T) {
	src := "cp -R src dest\n"
	want := "cp -a src dest\n"
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFixIntegration_ZC1171_EchoEscapeToPrint(t *testing.T) {
	src := `echo -e "line1\nline2"` + "\n"
	want := `print "line1\nline2"` + "\n"
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFixIntegration_ZC1135_EnvStripInline(t *testing.T) {
	src := "env FOO=bar cmd\n"
	want := "FOO=bar cmd\n"
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFixIntegration_ZC1170_PushdAddQuiet(t *testing.T) {
	src := "pushd /tmp\n"
	want := "pushd -q /tmp\n"
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFixIntegration_ZC1170_PopdWithDirAddQuiet(t *testing.T) {
	// Bare `popd` parses as Identifier (not SimpleCommand) so the
	// detector doesn't fire on it; giving popd an argument routes
	// through SimpleCommand where the detector lives.
	src := "popd /tmp\n"
	want := "popd -q /tmp\n"
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFixIntegration_ZC1209_SystemctlNoPager(t *testing.T) {
	src := "systemctl status nginx\n"
	want := "systemctl --no-pager status nginx\n"
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFixIntegration_ZC1210_JournalctlNoPager(t *testing.T) {
	src := "journalctl -u nginx\n"
	want := "journalctl --no-pager -u nginx\n"
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFixIntegration_ZC1213_AptGetAddYes(t *testing.T) {
	src := "apt-get install curl\n"
	want := "apt-get -y install curl\n"
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFixIntegration_ZC1226_DmesgAddTime(t *testing.T) {
	src := "dmesg --level err\n"
	want := "dmesg -T --level err\n"
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFixIntegration_ZC1227_CurlAddFail(t *testing.T) {
	src := "curl https://example.com/data\n"
	// ZC1255 (curl -L for HTTP redirects) shares this fixture and applies
	// alongside ZC1227's -f insertion in a single pass.
	want := "curl -L -f https://example.com/data\n"
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFixIntegration_ZC1234_DockerRunAddRm(t *testing.T) {
	src := "docker run alpine ls\n"
	want := "docker run --rm alpine ls\n"
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFixIntegration_ZC1231_GitCloneShallow(t *testing.T) {
	src := "git clone https://github.com/x/y\n"
	want := "git clone --depth 1 https://github.com/x/y\n"
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFixIntegration_ZC1241_XargsAddNullSep(t *testing.T) {
	src := "xargs rm\n"
	// ZC1773 (xargs -r to skip the no-input run) shares this fixture and
	// applies alongside ZC1241's -0 insertion in a single pass.
	want := "xargs -r -0 rm\n"
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFixIntegration_ZC1263_AptToAptGet(t *testing.T) {
	// ZC1448 also fires (apt install without -y) and inserts ` -y` at
	// the byte just past the original `apt` name. ZC1263 then rewrites
	// `apt` -> `apt-get`. The two edits do not overlap, so the combined
	// pass produces the apt-get + non-interactive form in one shot.
	src := "apt install curl\n"
	want := "apt-get -y install curl\n"
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFixIntegration_ZC1264_YumToDnf(t *testing.T) {
	src := "yum install httpd\n"
	want := "dnf install httpd\n"
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFixIntegration_ZC1253_DockerBuildNoCache(t *testing.T) {
	src := "docker build -t app .\n"
	want := "docker build --no-cache -t app .\n"
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFixIntegration_ZC1267_DfAddPortable(t *testing.T) {
	src := "df -h /\n"
	want := "df -P -h /\n"
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFixIntegration_ZC1355_EchoDashERawToPrintR(t *testing.T) {
	src := `echo -E "literal\tslash"` + "\n"
	want := `print -r "literal\tslash"` + "\n"
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFixIntegration_ZC1356_ReadArrayFlagUppercase(t *testing.T) {
	// ZC1012 fires simultaneously (missing raw flag) so the combined
	// rewrite adds the raw flag as well.
	src := "read -a arr\n"
	want := "read -r -A arr\n"
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFixIntegration_ZC1333_TimeformatToTimefmt(t *testing.T) {
	src := "fmt=$TIMEFORMAT\n"
	want := "fmt=$TIMEFMT\n"
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFixIntegration_ZC1331_BashRematchToMatch(t *testing.T) {
	src := "m=$BASH_REMATCH\n"
	want := "m=$match\n"
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFixIntegration_ZC1300_BashVersionToZsh(t *testing.T) {
	src := "v=$BASH_VERSION\n"
	want := "v=$ZSH_VERSION\n"
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFixIntegration_ZC1313_BashAliasesToAliases(t *testing.T) {
	src := "a=$BASH_ALIASES\n"
	want := "a=$aliases\n"
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFixIntegration_ZC1318_BashCmdsToCommands(t *testing.T) {
	src := "c=$BASH_CMDS\n"
	want := "c=$commands\n"
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFixIntegration_ZC1305_CompWordsToWords(t *testing.T) {
	src := "w=$COMP_WORDS\n"
	want := "w=$words\n"
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFixIntegration_ZC1306_CompCwordToCurrent(t *testing.T) {
	src := "c=$COMP_CWORD\n"
	want := "c=$CURRENT\n"
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFixIntegration_ZC1308_CompLineToBuffer(t *testing.T) {
	src := "b=$COMP_LINE\n"
	want := "b=$BUFFER\n"
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFixIntegration_ZC1298_FuncnameToFuncstack(t *testing.T) {
	src := "name=$FUNCNAME\n"
	want := "name=$funcstack\n"
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFixIntegration_ZC1304_BashSubshellToZsh(t *testing.T) {
	src := "depth=$BASH_SUBSHELL\n"
	want := "depth=$ZSH_SUBSHELL\n"
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFixIntegration_ZC1307_DirstackLowercase(t *testing.T) {
	src := "top=$DIRSTACK\n"
	want := "top=$dirstack\n"
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFixIntegration_ZC1301_PipestatusLowercase(t *testing.T) {
	src := "status=$PIPESTATUS\n"
	want := "status=$pipestatus\n"
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFixIntegration_ZC1265_SystemctlEnableNow(t *testing.T) {
	src := "systemctl enable nginx\n"
	want := "systemctl enable --now nginx\n"
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFixIntegration_ZC1283_SetOToSetopt(t *testing.T) {
	src := "set -o pipefail\n"
	want := "setopt pipefail\n"
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFixIntegration_ZC1288_DeclareToTypeset(t *testing.T) {
	src := "declare -i counter=0\n"
	want := "typeset -i counter=0\n"
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFixIntegration_ZC1192_SleepZeroToColon(t *testing.T) {
	src := "sleep 0\n"
	want := ":\n"
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
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

func TestFixIntegration_ZC1015_BackticksAlias(t *testing.T) {
	// ZC1015 shares ZC1002's fix shape — backticks become $(...)$
	// regardless of which kata id surfaces first.
	src := "result=`ls -la`\n"
	want := "result=$(ls -la)\n"
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFixIntegration_ZC1276_SeqAlias(t *testing.T) {
	src := "for i in $(seq 5); do :; done\n"
	want := `for i in "$({1..5})"; do :; done` + "\n"
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFixIntegration_ZC1271_WhichToCommandV(t *testing.T) {
	// ZC1271 fires alongside ZC1005 / ZC1034. The conflict resolver
	// keeps the `command -v` rewrite (parent ExpressionStatement edit
	// wins on walk order) and the result is idempotent on a re-run.
	src := "which git\n"
	want := "command -v git\n"
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFixIntegration_ZC1191_ClearToPrintAnsi(t *testing.T) {
	src := "clear\n"
	want := "print -rn $'\\e[2J\\e[H'\n"
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFixIntegration_ZC1202_IfconfigToIpAddr(t *testing.T) {
	src := "ifconfig eth0\n"
	want := "ip addr eth0\n"
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFixIntegration_ZC1203_NetstatToSs(t *testing.T) {
	src := "netstat -tulpn\n"
	want := "ss -tulpn\n"
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFixIntegration_ZC1216_NslookupToHost(t *testing.T) {
	src := "nslookup example.com\n"
	want := "host example.com\n"
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFixIntegration_ZC1501_DockerComposeHyphenToSpace(t *testing.T) {
	src := "docker-compose up -d\n"
	want := "docker compose up -d\n"
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFixIntegration_ZC1565_WhereisToCommandV(t *testing.T) {
	src := "whereis bash\n"
	want := "command -v bash\n"
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFixIntegration_ZC1565_LocateToCommandV(t *testing.T) {
	src := "locate bash\n"
	want := "command -v bash\n"
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFixIntegration_ZC1155_WhichDashAToWhence(t *testing.T) {
	// ZC1271 also fires on the `which` head and wants to rewrite to
	// `command -v`; ZC1155 / ZC1005 want `whence`. The conflict
	// resolver picks the first non-overlapping edit per pass; the
	// rewritten output is idempotent on re-run.
	src := "which -a python\n"
	got := runFix(t, src)
	if got == "which -a python\n" {
		t.Errorf("expected rewrite, got identical input %q", got)
	}
	if got != "whence -a python\n" && got != "command -v -a python\n" {
		t.Errorf("got %q, want a deterministic rewrite of `which -a`", got)
	}
}

func TestFixIntegration_ZC1334_TypeDashPToWhence(t *testing.T) {
	src := "type -p python\n"
	want := "whence -p python\n"
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFixIntegration_ZC1334_TypeDashCapPNormalised(t *testing.T) {
	src := "type -P python\n"
	want := "whence -p python\n"
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFixIntegration_ZC1411_EnableDashNToDisable(t *testing.T) {
	src := "enable -n cd\n"
	want := "disable cd\n"
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFixIntegration_ZC1235_GitPushDashF(t *testing.T) {
	src := "git push -f origin main\n"
	want := "git push --force-with-lease origin main\n"
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFixIntegration_ZC1260_GitBranchCapDToD(t *testing.T) {
	src := "git branch -D feat/old\n"
	want := "git branch -d feat/old\n"
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFixIntegration_ZC1448_AptInstallAddYes(t *testing.T) {
	// `apt install` triggers ZC1263 (apt -> apt-get) AND ZC1448
	// (insert -y). Both edits are non-overlapping; the combined
	// rewrite gives the unattended apt-get form.
	src := "apt install curl\n"
	want := "apt-get -y install curl\n"
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFixIntegration_ZC1448_AptUpgradeAddYes(t *testing.T) {
	src := "apt upgrade\n"
	want := "apt-get -y upgrade\n"
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFixIntegration_ZC1219_WgetDashOToCurl(t *testing.T) {
	src := "wget -O- https://example.com\n"
	want := "curl -fsSL https://example.com\n"
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFixIntegration_ZC1219_WgetDashQOToCurl(t *testing.T) {
	src := "wget -qO- https://example.com\n"
	want := "curl -fsSL https://example.com\n"
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFixIntegration_ZC1297_BashSourceToZsh(t *testing.T) {
	src := "src=$BASH_SOURCE\n"
	want := "src=${(%):-%x}\n"
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFixIntegration_ZC1377_BashAliasesInEcho(t *testing.T) {
	// ZC1092 also fires on `echo "..."` and rewrites the head to
	// `print -r --`. Both fixes apply in one pass; the variable
	// substitution inside the string literal is what ZC1377 owns.
	src := `echo "BASH_ALIASES=$BASH_ALIASES"` + "\n"
	want := `print -r -- "aliases=$aliases"` + "\n"
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFixIntegration_ZC1378_DirstackInPrint(t *testing.T) {
	src := `print "$DIRSTACK"` + "\n"
	want := `print -r "$dirstack"` + "\n"
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFixIntegration_ZC1383_TimeformatInEcho(t *testing.T) {
	// ZC1092 rewrites the `echo` head to `print -r --`; ZC1383 owns
	// the variable rename inside the quoted argument.
	src := `echo "$TIMEFORMAT"` + "\n"
	want := `print -r -- "$TIMEFMT"` + "\n"
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFixIntegration_ZC1394_BashInPrintf(t *testing.T) {
	src := `printf "%s\n" "$BASH"` + "\n"
	want := `printf "%s\n" "$ZSH_NAME"` + "\n"
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFixIntegration_ZC1163_GrepHeadOne(t *testing.T) {
	src := "grep PAT file | head -1\n"
	want := "grep -m 1 PAT file\n"
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFixIntegration_ZC1163_GrepHeadDashN1(t *testing.T) {
	src := "grep PAT file | head -n1\n"
	want := "grep -m 1 PAT file\n"
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFixIntegration_ZC1146_CatPipeTool(t *testing.T) {
	// Detector whitelist on the right of the pipe is awk/sed/sort/head/tail
	// — grep is intentionally excluded (grep accepts file args natively but
	// has its own kata for the pattern). Use sed to exercise the rewrite.
	src := "cat data.txt | sed s/foo/bar/\n"
	want := "sed s/foo/bar/ data.txt\n"
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFixIntegration_ZC1190_DoubleGrepInverted(t *testing.T) {
	// Fix gates on each side carrying exactly one non-flag pattern arg
	// (zc1190SinglePattern). Adding a `file` arg on the left would push it
	// over the limit, so the test fixture omits it.
	src := "grep -v p1 | grep -v p2\n"
	want := "grep -v -e p1 -e p2\n"
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFixIntegration_ZC1230_PingAddCount(t *testing.T) {
	src := "ping example.com\n"
	want := "ping -c 4 example.com\n"
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFixIntegration_ZC1238_DockerExecStripIt(t *testing.T) {
	src := "docker exec -it container ls\n"
	want := "docker exec container ls\n"
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFixIntegration_ZC1239_KubectlExecStripIt(t *testing.T) {
	src := "kubectl exec -it pod -- ls\n"
	want := "kubectl exec pod -- ls\n"
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFixIntegration_ZC1257_DockerStopAddTimeout(t *testing.T) {
	src := "docker stop container\n"
	want := "docker stop -t 10 container\n"
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFixIntegration_ZC1268_DuAddDoubleDash(t *testing.T) {
	src := "du -sh *\n"
	want := "du -sh -- *\n"
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFixIntegration_ZC1319_BashArgcRename(t *testing.T) {
	// ZC1092 (echo → print -r --) also fires on this fixture; both fixes
	// land in the same pass per the registry's first-edit-wins policy.
	src := "echo $BASH_ARGC\n"
	want := "print -r -- $#\n"
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFixIntegration_ZC1320_BashArgvRename(t *testing.T) {
	// ZC1092 (echo → print -r --) also fires on this fixture; both fixes
	// land in the same pass per the registry's first-edit-wins policy.
	src := "echo $BASH_ARGV\n"
	want := "print -r -- $argv\n"
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFixIntegration_ZC1380_HistignoreRename(t *testing.T) {
	src := "export HISTIGNORE='ls:cd'\n"
	want := "export HISTORY_IGNORE='ls:cd'\n"
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
