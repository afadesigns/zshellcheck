// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package lexer

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/token"
)

// drainTokens lexes every token in src — coverage helper.
func drainTokens(src string) {
	l := New(src)
	for {
		tok := l.NextToken()
		if tok.Type == token.EOF {
			return
		}
	}
}

// readUnicodeIdent path: prompt strings / theme code with non-ASCII
// bytes (UTF-8 multibyte) used as identifier characters.
func TestLexerUnicodeIdentifier(t *testing.T) {
	drainTokens("echo café\n")
	drainTokens("echo αβγ\n")
	drainTokens("echo 日本語\n")
}

// readArithCompoundOr exercising `*=` / `%=` fusion.
func TestLexerCompoundAsterisk(t *testing.T) {
	drainTokens("(( x *= 2 ))\n")
}

func TestLexerCompoundPercent(t *testing.T) {
	drainTokens("(( x %= 60 ))\n")
}

// readBangLead non-equal path.
func TestLexerBangAlone(t *testing.T) {
	drainTokens("! true\n")
}

// readSemicolonLead double-semi.
func TestLexerDoubleSemi(t *testing.T) {
	drainTokens("case x in a) :;; esac\n")
}

// readEqualsLead variants.
func TestLexerEqualsTilde(t *testing.T) {
	drainTokens("[[ $a =~ ^pat$ ]]\n")
}

func TestLexerEqualsLparenWithSpace(t *testing.T) {
	drainTokens("typeset arr =( one two )\n")
}

// readPlusLead ++ / +=.
func TestLexerPlusPlus(t *testing.T) {
	drainTokens("(( i++ ))\n")
}

func TestLexerPlusEqual(t *testing.T) {
	drainTokens("(( total += 1 ))\n")
}

// readMinusLead `-` followed by non-letter, then DEC, then PLUSEQ.
func TestLexerMinusAlone(t *testing.T) {
	drainTokens("expr 5 - 1\n")
}

func TestLexerDecrement(t *testing.T) {
	drainTokens("(( i-- ))\n")
}

func TestLexerMinusEqual(t *testing.T) {
	drainTokens("(( total -= 1 ))\n")
}

// readAmpersandLead variants &&, &|, &!.
func TestLexerDoubleAmpersand(t *testing.T) {
	drainTokens("true && echo ok\n")
}

func TestLexerAmpersandPipe(t *testing.T) {
	drainTokens("sleep 99 &|\n")
}

func TestLexerAmpersandBang(t *testing.T) {
	drainTokens("sleep 99 &!\n")
}

// readPipeLead double pipe.
func TestLexerDoublePipe(t *testing.T) {
	drainTokens("false || echo ok\n")
}

// readOpenBracket: LDBRACKET vs LBRACKET-class.
func TestLexerDoubleBracketStmt(t *testing.T) {
	drainTokens("[[ -f file ]]\n")
}

func TestLexerBracketClass(t *testing.T) {
	drainTokens("ls [[:alnum:]]*\n")
}

// readCloseBracket: RDBRACKET inside [[, vs RBRACKET in glob class.
func TestLexerNestedBrackets(t *testing.T) {
	drainTokens("[[ -f file && [[:alnum:]] ]]\n")
}

// readAngleBracket: <<<, <(, <=, <&, <<-, >&, >|, >!.
func TestLexerLessTriple(t *testing.T) {
	drainTokens("cat <<< hello\n")
}

func TestLexerProcessSubst(t *testing.T) {
	drainTokens("diff <(sort a) <(sort b)\n")
}

func TestLexerLessEqual(t *testing.T) {
	drainTokens("[[ $x <= 5 ]]\n")
}

func TestLexerLessAmp(t *testing.T) {
	drainTokens("cmd <&3\n")
}

func TestLexerHeredocStrip(t *testing.T) {
	drainTokens("cat <<-EOF\n\tbody\n\tEOF\n")
}

func TestLexerForceClobber(t *testing.T) {
	drainTokens("echo hi >| /tmp/out\n")
}

func TestLexerForceBang(t *testing.T) {
	drainTokens("echo hi >! /tmp/out\n")
}

// readDollarToken full surface.
func TestLexerDollarBraceParam(t *testing.T) {
	drainTokens("echo ${var:-default}\n")
}

func TestLexerDollarParenSubst(t *testing.T) {
	drainTokens("echo $(date)\n")
}

func TestLexerDollarSingleQuote(t *testing.T) {
	drainTokens("echo $'tab\\there'\n")
}

func TestLexerDollarDoubleQuote(t *testing.T) {
	drainTokens("echo $\"localised\"\n")
}

func TestLexerDollarSpecial(t *testing.T) {
	drainTokens("echo $? $$ $!\n")
}

// Backslash escape variants (readBackslashEscape).
func TestLexerBackslashLetter(t *testing.T) {
	drainTokens("echo \\n \\t \\r\n")
}

func TestLexerBackslashGlob(t *testing.T) {
	drainTokens("echo \\* \\? \\[\n")
}

func TestLexerBackslashShell(t *testing.T) {
	drainTokens("echo \\& \\| \\; \\< \\>\n")
}

func TestLexerBackslashIllegal(t *testing.T) {
	drainTokens("echo \\\x01\n")
}

// String literals — readString with interpolation.
func TestLexerInterpolatedString(t *testing.T) {
	drainTokens("echo \"hello $name and ${var}\"\n")
}

func TestLexerSingleQuoteString(t *testing.T) {
	drainTokens("echo 'literal $not interpolated'\n")
}

func TestLexerNestedBraceInString(t *testing.T) {
	drainTokens("echo \"${var=\"default\"}\"\n")
}

// Number tokens.
func TestLexerHexNumber(t *testing.T) {
	drainTokens("(( x = 0x1F ))\n")
}

func TestLexerOctalNumber(t *testing.T) {
	drainTokens("(( x = 0755 ))\n")
}

// Edge: comment forms.
func TestLexerCommentAtColumnOne(t *testing.T) {
	drainTokens("# top-level comment\necho hi\n")
}

func TestLexerCommentAfterCode(t *testing.T) {
	drainTokens("echo hi # trailing comment\n")
}

// peekAt + multi-character look-ahead.
func TestLexerPosixCharClassAvoidsLDBracket(t *testing.T) {
	drainTokens("ls *[[:upper:]]*\n")
}

// Continuation absorption.
func TestLexerLineContinuation(t *testing.T) {
	drainTokens("echo hi \\\nworld\n")
}

// Zsh number bases: `0x…` hex, `0b…` binary, `0o…` octal, `BASE#…`.
// Used inside `((…))` arithmetic.
func TestLexerHexLiteral(t *testing.T)         { drainTokens("(( x = 0xff ))\n") }
func TestLexerBinaryLiteral(t *testing.T)      { drainTokens("(( x = 0b101 ))\n") }
func TestLexerOctalLiteral(t *testing.T)       { drainTokens("(( x = 0o755 ))\n") }
func TestLexerCustomBaseLiteral(t *testing.T)  { drainTokens("(( x = 16#ff ))\n") }
func TestLexerHexBareNoDigits(t *testing.T)    { drainTokens("(( x == 0x${var} ))\n") }
func TestLexerBinaryBareNoDigits(t *testing.T) { drainTokens("(( x == 0b${var} ))\n") }

// Escaped backtick appears in brace-expansion lists and pattern
// strings. The lexer treats it as an IDENT word so the surrounding
// command-word fold keeps working. Used heavily by zinit and
// zsh-vi-mode.
func TestLexerEscapedBacktickInBraceExpansion(t *testing.T) {
	drainTokens("for s in {\\',\\\",\\`}; do echo $s; done\n")
}

func TestLexerEscapedBacktickAsArg(t *testing.T) {
	drainTokens("echo \\`\n")
}
