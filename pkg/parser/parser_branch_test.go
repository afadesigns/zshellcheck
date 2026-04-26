// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package parser

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/lexer"
)

// drainProgram parses src and discards the result. Used to drive code
// paths for coverage without making assertions about the AST shape.
func drainProgram(src string) {
	l := lexer.New(src)
	p := New(l)
	_ = p.ParseProgram()
}

// parseIfStatement multi-elif + alternative chain.
func TestBranchIfElifElse(t *testing.T) {
	drainProgram("if true; then echo a; elif false; then echo b; elif test; then echo c; else echo d; fi\n")
}

// parsePipelineStartingWithExpression: $(...)|sort, `...`|sort, $var|wc, ${name}|head.
func TestBranchPipelineDollarParen(t *testing.T)  { drainProgram("$(date) | tee log\n") }
func TestBranchPipelineBacktick(t *testing.T)     { drainProgram("`date` | tee log\n") }
func TestBranchPipelineVariableHead(t *testing.T) { drainProgram("$cmd | sort | head -1\n") }
func TestBranchPipelineDollarBraceHead(t *testing.T) {
	drainProgram("${cmd:-default} | wc -l\n")
}

// parseSingleCommand head-prefix branch (LBRACKET-rooted prefix).
func TestBranchLBracketAsCommand(t *testing.T) { drainProgram("[ -f file ]\n") }

// parseFunctionLiteral with ${=X} name + multi-name keyword.
func TestBranchFnLitDollarBraceName(t *testing.T) {
	drainProgram("function ${=NAMES} { :; }\n")
}

func TestBranchFnLitMultiName(t *testing.T) {
	drainProgram("function alpha beta gamma { :; }\n")
}

func TestBranchFnLitNoParens(t *testing.T) {
	drainProgram("function name { echo hi; }\n")
}

func TestBranchFnLitCompositeName(t *testing.T) {
	drainProgram("function name${1:-}suffix() { :; }\n")
}

// parseForLoop short-form.
func TestBranchForLoopShortFormCmd(t *testing.T) { drainProgram("for f (a b c) echo $f\n") }

func TestBranchForLoopShortFormBlock(t *testing.T) {
	drainProgram("for f (a b c) { echo $f }\n")
}

func TestBranchForLoopShortFormDoDone(t *testing.T) {
	drainProgram("for f (a b c); do echo $f; done\n")
}

func TestBranchForLoopArithEmptyAll(t *testing.T) {
	drainProgram("for ((;;)) do break; done\n")
}

func TestBranchForLoopArithInitOnly(t *testing.T) {
	drainProgram("for ((i=0; ; )) do break; done\n")
}

func TestBranchForLoopArithCondOnly(t *testing.T) {
	drainProgram("for ((; i<3; )) do break; done\n")
}

func TestBranchForLoopArithPostOnly(t *testing.T) {
	drainProgram("for ((; ; i++)) do break; done\n")
}

func TestBranchForLoopBraceBody(t *testing.T) {
	drainProgram("for x in 1 2; { echo $x }\n")
}

func TestBranchForLoopShortLineBody(t *testing.T) {
	drainProgram("for x in 1 2 3\n  echo $x\n")
}

// parseCaseClause variants.
func TestBranchCaseLeadingParen(t *testing.T) {
	drainProgram("case x in (a) :;; (b) :;; esac\n")
}

func TestBranchCaseGlobAlt(t *testing.T) {
	drainProgram("case x in (a|b)*) :;; *) :;; esac\n")
}

func TestBranchCaseInlineGlob(t *testing.T) {
	drainProgram("case x in plugin::(a|b|c)) echo hit ;; esac\n")
}

// parseInvalidArrayAccessPrefix variants.
func TestBranchDollarBeforeSemicolon(t *testing.T) { drainProgram("echo $;\n") }
func TestBranchDollarBracketArith(t *testing.T)    { drainProgram("echo $[1+2]\n") }
func TestBranchDollarBracketNested(t *testing.T)   { drainProgram("echo $[[1+2]+3]\n") }
func TestBranchDollarHashLength(t *testing.T)      { drainProgram("echo $#name\n") }
func TestBranchDollarPlusSubscript(t *testing.T) {
	drainProgram("(( $+commands[ls] ))\n")
}

func TestBranchDollarPlusBareSubscript(t *testing.T) {
	drainProgram("(( $+commands[$cmd] ))\n")
}
func TestBranchDollarIntegerArg(t *testing.T) { drainProgram("echo $1\n") }
func TestBranchDollarStarArg(t *testing.T)    { drainProgram("echo $*\n") }
func TestBranchDollarMinusArg(t *testing.T)   { drainProgram("echo $-\n") }
func TestBranchDollarBangArg(t *testing.T)    { drainProgram("echo $!\n") }

// parseArithmeticCommand chained with logical chain.
func TestBranchArithLogicalChain(t *testing.T) {
	drainProgram("(( x > 0 )) && echo positive || echo non-positive\n")
}

// parseDoubleBracketExpression chained.
func TestBranchDoubleBracketLogicalChain(t *testing.T) {
	drainProgram("[[ -f file ]] && echo yes || echo no\n")
}

// parseDeclarationStatement composite name + tail variants.
func TestBranchDeclTypesetCompositeName(t *testing.T) {
	drainProgram("typeset -g \"$1\"=\"$2\"\n")
}

func TestBranchDeclLocalAppend(t *testing.T) {
	drainProgram("local PATH+=:bin\n")
}

func TestBranchDeclTypesetEmptyRHS(t *testing.T) {
	drainProgram("typeset -g VAR=\nfor x in 1 2; do :; done\n")
}

// parseIndexExpression flag tuples + slice forms.
func TestBranchIndexExprUpperRFlag(t *testing.T) { drainProgram("echo ${arr[(R)pat]}\n") }
func TestBranchIndexExprLowerRFlag(t *testing.T) { drainProgram("echo ${arr[(r)pat]}\n") }
func TestBranchIndexExprIFlag(t *testing.T)      { drainProgram("echo ${arr[(I)i]}\n") }
func TestBranchIndexExprMultiFlag(t *testing.T)  { drainProgram("echo ${arr[(ri)pat]}\n") }
func TestBranchIndexExprSlice(t *testing.T)      { drainProgram("echo ${arr[1,8]}\n") }
func TestBranchIndexExprNestedSubscript(t *testing.T) {
	drainProgram("echo ${arr[$colors[red]]}\n")
}

// parseInfixExpression vs glob-alt break inside [[ ]].
func TestBranchDoubleBracketGlobAlt(t *testing.T) {
	drainProgram("[[ $x = (foo|bar)*.zsh ]]\n")
}

func TestBranchDoubleBracketTrailingGlob(t *testing.T) {
	drainProgram("[[ $x = /* ]] || true\n")
}

// parseLetStatement modifier-prefix path.
func TestBranchLetWithLocal(t *testing.T) {
	drainProgram("let local elapsed=1\n")
}

// parseStatement bang variants.
func TestBranchBangPipeline(t *testing.T) { drainProgram("! grep foo file\n") }
func TestBranchBangSubshell(t *testing.T) { drainProgram("! (cd /tmp)\n") }
func TestBranchBangDoubleBracket(t *testing.T) {
	drainProgram("! [[ -f file ]]\n")
}

// parseStatement for misc tokens.
func TestBranchStmtColon(t *testing.T)     { drainProgram(": noop\n") }
func TestBranchStmtDot(t *testing.T)       { drainProgram(". /etc/profile\n") }
func TestBranchStmtAmpersand(t *testing.T) { drainProgram("& cmd\n") }
