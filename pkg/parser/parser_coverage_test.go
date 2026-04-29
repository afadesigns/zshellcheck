// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package parser

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/lexer"
)

// parseClean is a low-noise helper: parse src, fail when the parser
// emits an error. Used by coverage-targeted tests that only need to
// drive a code path.
func parseClean(t *testing.T, src string) {
	t.Helper()
	p := New(lexer.New(src))
	_ = p.ParseProgram()
	if errs := p.Errors(); len(errs) != 0 {
		t.Fatalf("unexpected parser errors for %q: %v", src, errs)
	}
}

// parseSingleCommand head-prefix path: command starts with $(…),
// `…`, $name, or ${name}.
func TestParseCommandHeadDollarParenSubst(t *testing.T) { parseClean(t, "$(date) -h\n") }
func TestParseCommandHeadBacktickSubst(t *testing.T)    { parseClean(t, "`date` -h\n") }
func TestParseCommandHeadVariableArg(t *testing.T)      { parseClean(t, "$cmd arg\n") }
func TestParseCommandHeadDollarBraceWithModifier(t *testing.T) {
	parseClean(t, "${cmd:-default} arg\n")
}

func TestParseCommandHeadProcessSubst(t *testing.T) {
	parseClean(t, "diff <(sort a) <(sort b)\n")
}

// parseSingleCommand function-definition path.
func TestParseFuncDefBareName(t *testing.T)       { parseClean(t, "myfn() { :; }\n") }
func TestParseFuncDefSingleArgParen(t *testing.T) { parseClean(t, "myfn ( arg )\n") }

// Simple command argument shapes.
func TestParseCommandWithDoubleDash(t *testing.T)  { parseClean(t, "rm -- -file\n") }
func TestParseCommandWithIncDecArg(t *testing.T)   { parseClean(t, "cmd -- ++foo --bar\n") }
func TestParseCommandWithReservedArg(t *testing.T) { parseClean(t, "print -l function do done\n") }
func TestParseCommandWithBraceUpstream(t *testing.T) {
	parseClean(t, "git log @{upstream}\n")
}

func TestParseCommandWithBraceExpansion(t *testing.T) {
	parseClean(t, "echo {a,b,c}.txt\n")
}

func TestParseCommandWithGlobAlternation(t *testing.T) {
	parseClean(t, "ls *(.).log\n")
}

// parseRedirection / readAngleBracket variants exercised at
// statement level so the AST surface gets reached.
func TestParseRedirectStdoutAppendChain(t *testing.T) { parseClean(t, "echo hi >> log\n") }

func TestParseRedirectStdinHerestring(t *testing.T) {
	parseClean(t, "cat <<< hello\n")
}

func TestParseRedirectStderrToStdout(t *testing.T) {
	parseClean(t, "cmd 2>&1\n")
}

func TestParseRedirectForceClobber(t *testing.T) {
	parseClean(t, "echo hi >| /tmp/out\n")
}

func TestParseRedirectProcessSubstOut(t *testing.T) {
	parseClean(t, "tee >(cat) >(cat) < input\n")
}

// parseDollarIdent / parseInvalidArrayAccessPrefix branches.
func TestParseDollarIdentWithSubscript(t *testing.T) {
	parseClean(t, "echo $arr[1]\n")
}
func TestParseDollarPlusNameInArith(t *testing.T)     { parseClean(t, "(( $+name ))\n") }
func TestParseDollarPlusSubscript(t *testing.T)       { parseClean(t, "(( $+commands[ls] ))\n") }
func TestParseDollarHashLengthSubscript(t *testing.T) { parseClean(t, "echo $#arr\n") }
func TestParseDollarBracketArith(t *testing.T)        { parseClean(t, "echo $[1+2]\n") }
func TestParseBareDollarBeforeRBrace(t *testing.T) {
	parseClean(t, "echo $\n")
}

// parseFunctionLiteral brace forms + composite names.
func TestParseFunctionLiteralKeywordBracesOnly(t *testing.T) {
	parseClean(t, "function name { echo hi; }\n")
}

func TestParseFunctionLiteralCompositeNameInline(t *testing.T) {
	parseClean(t, "function n${1:-}suffix() { :; }\n")
}

func TestParseFunctionLiteralMultipleNamesShort(t *testing.T) {
	parseClean(t, "function a b c { :; }\n")
}

// parseForLoopStatement variant coverage.
func TestParseForLoopArithEmpty(t *testing.T) {
	parseClean(t, "for (( ; ; )) do break; done\n")
}

func TestParseForLoopShortFormBraceBodyInline(t *testing.T) {
	parseClean(t, "for x in 1 2 3; { echo $x }\n")
}

func TestParseForLoopShortFormParenItemsCmd(t *testing.T) {
	parseClean(t, "for x (1 2 3) echo $x\n")
}

func TestParseForLoopMultiVarPairs(t *testing.T) {
	parseClean(t, "for k v in a 1 b 2; do echo $k $v; done\n")
}

func TestParseForLoopNumericPositional(t *testing.T) {
	parseClean(t, "for 1 in foo bar; do :; done\n")
}

// parseStatement branch coverage.
func TestParseStatementBangPipeline(t *testing.T) {
	parseClean(t, "! grep foo file\n")
}

func TestParseStatementBangCommand(t *testing.T) {
	parseClean(t, "! [[ -f file ]]\n")
}

func TestParseStatementCoproc(t *testing.T) {
	parseClean(t, "coproc cat\n")
}

func TestParseStatementBraceGroupPipe(t *testing.T) {
	parseClean(t, "{ echo a; echo b; } | sort\n")
}

// parseDeclarationStatement branch coverage.
func TestParseDeclTypesetIntegerArray(t *testing.T) {
	parseClean(t, "typeset -ai nums=(1 2 3)\n")
}

func TestParseDeclLocalReadonly(t *testing.T) {
	parseClean(t, "local -r CONST=42\n")
}

func TestParseDeclDeclareAssoc(t *testing.T) {
	parseClean(t, "declare -A m=([k]=v)\n")
}

func TestParseDeclTypesetFlags(t *testing.T) {
	parseClean(t, "typeset -gx EXPORTED=1\n")
}

// peekStartsCommand / peekStartsArgPrefix branch coverage.
func TestParseIdentBangArg(t *testing.T)   { parseClean(t, "cmd !=value\n") }
func TestParseIdentLBraceArg(t *testing.T) { parseClean(t, "cmd {a,b}\n") }
func TestParseIdentTildeArg(t *testing.T)  { parseClean(t, "cmd ~/path\n") }
func TestParseIdentSlashArg(t *testing.T)  { parseClean(t, "cmd /path/file\n") }

// Heredoc bodies through the lexer.
func TestParseHeredocBodyInPipeline(t *testing.T) {
	parseClean(t, "cat <<EOF | wc\nbody\nEOF\n")
}

func TestParseHeredocStripTabs(t *testing.T) {
	parseClean(t, "cat <<-EOF\n\tbody\n\tEOF\n")
}

// Case statement variants.
func TestParseCaseGlobAlternationLabel(t *testing.T) {
	parseClean(t, "case $x in (darwin|freebsd)*) echo bsd ;; *) echo other ;; esac\n")
}

func TestParseCaseLeadingParenLabel(t *testing.T) {
	parseClean(t, "case $x in (a) echo a ;; (b) echo b ;; esac\n")
}

// Subshell variants.
func TestParseSubshellPipe(t *testing.T) {
	parseClean(t, "( echo a; echo b ) | wc\n")
}

func TestParseSubshellAnonymousFn(t *testing.T) {
	parseClean(t, "() { echo anon; }\n")
}

// parseDoubleParenExpression prefix path: `((…))` in expression slot.
func TestParseDoubleParenExpressionInLet(t *testing.T) {
	parseClean(t, "let x=(( 1 + 2 ))\n")
}

func TestParseDoubleParenExpressionRadix(t *testing.T) {
	parseClean(t, "(([#16] 0xff))\n")
}

func TestParseDoubleParenExpressionInChain(t *testing.T) {
	parseClean(t, "true && (( x++ ))\n")
}

// parseRedirection infix paths (`>>`, `<<<`, `>&`, `<&`).
func TestParseRedirectionAppendArg(t *testing.T)     { parseClean(t, "echo a >> file\n") }
func TestParseRedirectionHerestringArg(t *testing.T) { parseClean(t, "cat <<< 'inline'\n") }
func TestParseRedirectionFdMerge(t *testing.T)       { parseClean(t, "cmd 2>&1\n") }
func TestParseRedirectionFdInputDup(t *testing.T)    { parseClean(t, "exec 3<&0\n") }
func TestParseRedirectionChain(t *testing.T) {
	parseClean(t, "cmd >> out.log 2>&1\n")
}

// parseKeywordAsCommand: `return` as expression in a logical chain.
func TestParseReturnAsExprInChain(t *testing.T) {
	parseClean(t, "cond || return 1\n")
}

func TestParseReturnAsExprBare(t *testing.T) {
	parseClean(t, "func() { check || return; }\n")
}

func TestParseReturnAsExprMultiArg(t *testing.T) {
	parseClean(t, "guard && return 0\n")
}

// parseDollarIdent invalid-array-access path: `$name[idx]`.
func TestParseDollarIdentSubscript(t *testing.T) {
	parseClean(t, "echo $arr[1]\n")
}

func TestParseDollarIdentNestedSubscript(t *testing.T) {
	parseClean(t, "echo $arr[$idx]\n")
}

// finalizeInvalidArrayAccess drain path: subscript with a deeper
// bracket mismatch that the drainer must walk past.
func TestParseDollarIdentSubscriptDeep(t *testing.T) {
	parseClean(t, "echo $arr[$nested[2]]\n")
}

// drainSubscriptBody used by `${var[idx]}` modifier-tail walk.
func TestParseDollarBraceSubscriptModifier(t *testing.T) {
	parseClean(t, "echo ${arr[1]:-default}\n")
}

func TestParseDollarBraceSubscriptComplex(t *testing.T) {
	parseClean(t, "echo ${arr[$i+1]##prefix}\n")
}

// parseSingleCommand: trailing-redirection + arg-prefix variants.
func TestParseSingleCommandRedirToFD(t *testing.T) {
	parseClean(t, "cmd 1>&2\n")
}

func TestParseSingleCommandWithEnvPrefix(t *testing.T) {
	parseClean(t, "FOO=bar BAR=baz cmd arg\n")
}

func TestParseSingleCommandLeadingNewlines(t *testing.T) {
	parseClean(t, "\n\n\necho a\n")
}

// parsePipelineStartingWithExpression: head is `${…}` / `$(…)` /
// VARIABLE / BACKTICK.
func TestParsePipelineHeadFromBacktick(t *testing.T) {
	parseClean(t, "`echo cmd` arg | wc\n")
}

func TestParsePipelineHeadFromVariable(t *testing.T) {
	parseClean(t, "$prog arg1 arg2\n")
}

// peekStartsCommand variants (used by `! cmd`).
func TestParseBangBeforeBracket(t *testing.T) {
	parseClean(t, "! [[ -f file ]]\n")
}

func TestParseBangBeforeArith(t *testing.T) {
	parseClean(t, "! (( x > 0 ))\n")
}

// parseArithmeticSubscript: `arr[expr]` with non-trivial expr.
func TestParseArithmeticSubscriptArith(t *testing.T) {
	parseClean(t, "echo ${arr[i+1]}\n")
}

// parseProcessSubstitution: `>(…)` write side and nested forms.
func TestParseProcessSubstWrite(t *testing.T) {
	parseClean(t, "tee >(grep err) >(gzip > out.gz)\n")
}

// parseFlaggedSubscript: `${(flag)arr[idx]}`.
func TestParseFlaggedSubscriptKeyArr(t *testing.T) {
	parseClean(t, "echo ${(k)assoc}\n")
}

// parseGroupedExpression keyword-headed bodies (for/while/if inside
// subshell).
func TestParseGroupedKeywordFor(t *testing.T) {
	parseClean(t, "( for f in *.txt; do echo $f; done )\n")
}

func TestParseGroupedKeywordWhile(t *testing.T) {
	parseClean(t, "( while read l; do echo $l; done )\n")
}

// parseSingleCommand head-prefix coverage targeting the
// DOLLAR_LPAREN / BACKTICK / VARIABLE / ${} branch pre-arg-loop.
func TestParseSingleCommandDollarParenHead(t *testing.T) {
	parseClean(t, "$(which cmd) -n\n")
}

func TestParseSingleCommandBacktickHead(t *testing.T) {
	parseClean(t, "`which cmd` -n\n")
}

func TestParseSingleCommandVariableHead(t *testing.T) {
	parseClean(t, "$cmd -h --foo bar\n")
}

func TestParseSingleCommandDollarBraceHead(t *testing.T) {
	parseClean(t, "${cmd} arg1 arg2\n")
}

// parseSingleCommand `name (arg)` non-fn-def path: parens after a
// command name with content inside, NOT followed immediately by `)`.
func TestParseSingleCommandNameParenArg(t *testing.T) {
	parseClean(t, "ls (file)\n")
}

// parseDollarSpecialOp: `$?`, `$$`, `$@`, positional.
func TestParseDollarQuestionBare(t *testing.T) { parseClean(t, "echo $?\n") }
func TestParseDollarDollarBare(t *testing.T)   { parseClean(t, "echo $$\n") }
func TestParseDollarAtBare(t *testing.T)       { parseClean(t, "echo $@\n") }
func TestParseDollarPositionalChain(t *testing.T) {
	parseClean(t, "echo $0 $1 $2\n")
}

// parseDollarSpecialOp `$+` (zsh: subscript flag): `$+commands`.
func TestParseDollarPlusInArith(t *testing.T) {
	parseClean(t, "(( $+commands[ls] ))\n")
}

// parseArrayAccessSubject keyword-as-subject path.
func TestParseDollarBraceKeywordAsSubject(t *testing.T) {
	parseClean(t, "echo ${for}\n")
}

// drainSubscriptBody depth-tracking branches.
func TestParseDollarBraceNestedBracket(t *testing.T) {
	parseClean(t, "echo ${arr[$nested[1]]}\n")
}

// parseDollarParenExpression keyword-headed body forms.
func TestParseDollarParenForLoop(t *testing.T) {
	parseClean(t, "echo $(for f in *; do print $f; done)\n")
}

func TestParseDollarParenIfStatement(t *testing.T) {
	parseClean(t, "echo $(if [[ -f $1 ]]; then echo yes; fi)\n")
}

// parseArithmeticSubscript with operator chain.
func TestParseArithSubscriptChain(t *testing.T) {
	parseClean(t, "echo ${arr[i*2+1]}\n")
}

// parseFlaggedSubscript Zsh subscript-flag tuples.
func TestParseFlaggedSubscriptKey(t *testing.T) {
	parseClean(t, "echo ${(k)assoc}\n")
}

func TestParseFlaggedSubscriptValue(t *testing.T) {
	parseClean(t, "echo ${(v)assoc}\n")
}

func TestParseFlaggedSubscriptKeyValue(t *testing.T) {
	parseClean(t, "echo ${(kv)assoc}\n")
}

// parseProcessSubstitution write+read mix and bare path.
func TestParseProcessSubstReadAndWrite(t *testing.T) {
	parseClean(t, "diff <(sort a) >(sort b)\n")
}

// parseCommandWord with mixed quoting + concat.
func TestParseCommandWordConcatMix(t *testing.T) {
	parseClean(t, "echo prefix${var}suffix\n")
}

func TestParseCommandWordDoubleQuoteWithSub(t *testing.T) {
	parseClean(t, "echo \"value=$(cmd) end\"\n")
}

// parseExpression / parseEqualsForm: env-prefix assignment chain.
func TestParseEnvPrefixChain(t *testing.T) {
	parseClean(t, "X=1 Y=2 Z=3 cmd arg\n")
}

// parseStatement coverage for HASH (top-level comment).
func TestParseStatementHashCommentOnly(t *testing.T) {
	parseClean(t, "# top-level comment line\n")
}

// parsePipelineHeadStatement select / coproc.
func TestParseSelectStatementBranches(t *testing.T) {
	parseClean(t, "select x in a b c; do echo $x; break; done\n")
}

// parseFunctionLiteral with composite name + body shapes.
func TestParseFunctionCompositeNameKeyword(t *testing.T) {
	parseClean(t, "function ::my-fn { echo hi; }\n")
}

// parseDeclarationStatement assoc-array assignment edge cases.
func TestParseDeclAssocArrayMulti(t *testing.T) {
	parseClean(t, "typeset -A m=([k1]=v1 [k2]=v2)\n")
}

// parseCaseStatement with empty / fall-through clause body.
func TestParseCaseEmptyBody(t *testing.T) {
	parseClean(t, "case $x in a) ;; b) ;; *) ;; esac\n")
}

func TestParseCaseFallThroughSemiAmp(t *testing.T) {
	parseClean(t, "case $x in a) echo a;& b) echo b;; esac\n")
}

func TestParseCaseFallThroughSemiPipe(t *testing.T) {
	parseClean(t, "case $x in a) echo a;| b) echo b;; esac\n")
}

// parseCommandPipeline redirection loop after pipeline head — fires
// when a parens / brace-group / arith head is followed by redir.
func TestParsePipelineRedirAfterSubshell(t *testing.T) {
	parseClean(t, "( echo hi ) > /tmp/out\n")
}

func TestParsePipelineRedirAfterArith(t *testing.T) {
	parseClean(t, "(( x++ )) > /dev/null\n")
}

func TestParsePipelineRedirAfterBraceGroup(t *testing.T) {
	parseClean(t, "{ echo a; } 2> /dev/null\n")
}

func TestParsePipelineRedirChain(t *testing.T) {
	parseClean(t, "( cmd ) >> log 2>&1\n")
}

// parseSingleCommand function-definition body variants.
func TestParseFuncDefSimpleBody(t *testing.T) {
	parseClean(t, "myfn() echo body\n")
}

func TestParseFuncDefBraceBody(t *testing.T) {
	parseClean(t, "myfn() { echo body; }\n")
}

func TestParseFuncDefSubshellBody(t *testing.T) {
	parseClean(t, "myfn() ( echo body )\n")
}

// parseDollarIdent: bare DOLLAR followed by IDENT.
func TestParseDollarSpaceIdent(t *testing.T) {
	parseClean(t, "echo $ name\n")
}

// parsePipelineStartingWithExpression pipeline tail variants.
func TestParsePipelineExprPipeChain(t *testing.T) {
	parseClean(t, "$(date) | grep -v warn | head -1\n")
}

func TestParsePipelineExprAndChain(t *testing.T) {
	parseClean(t, "$cmd && echo done\n")
}

// parseExpression bit operations in arith.
func TestParseArithBitOr(t *testing.T) {
	parseClean(t, "(( x = a | b ))\n")
}

func TestParseArithBitAnd(t *testing.T) {
	parseClean(t, "(( x = a & b ))\n")
}

// parseStatement COPROC bodies.
func TestParseCoprocPipeline(t *testing.T) {
	parseClean(t, "coproc cat | sort\n")
}

func TestParseCoprocBlock(t *testing.T) {
	parseClean(t, "coproc { echo a; echo b; }\n")
}

// parseDeclarationStatement flag combinations.
func TestParseDeclTypesetExportArray(t *testing.T) {
	parseClean(t, "typeset -gx -a items=(a b c)\n")
}

// parseGroupedExpression with select keyword body.
func TestParseGroupedSelect(t *testing.T) {
	parseClean(t, "( select x in a b c; do echo $x; break; done )\n")
}

func TestParseGroupedTypesetDeclaration(t *testing.T) {
	parseClean(t, "( typeset -A m=([k]=v); echo done )\n")
}
