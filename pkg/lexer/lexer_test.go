package lexer

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/token"
)

func TestDollarQuotedStrings(t *testing.T) {
	// Zsh ANSI-C quoting `$'…'` and gettext quoting `$"…"` must lex
	// as single STRING tokens. Previously the `$` was emitted as
	// DOLLAR and the quote contents as STRING, producing a pair the
	// parser could not consume on an assignment RHS. Widely used in
	// oh-my-zsh themes (darkblood and many more) to embed escape
	// sequences inside prompt strings.
	cases := []struct {
		in   string
		want string
	}{
		{`$'hello'`, `$'hello'`},
		{`$'a\nb'`, `$'a\nb'`},
		{`$"gettext"`, `$"gettext"`},
	}
	for _, c := range cases {
		l := New(c.in)
		tok := l.NextToken()
		if tok.Type != token.STRING {
			t.Errorf("%q: want STRING, got %s", c.in, tok.Type)
		}
		if tok.Literal != c.want {
			t.Errorf("%q: literal = %q, want %q", c.in, tok.Literal, c.want)
		}
	}
}

func TestNextToken(t *testing.T) {
	input := `let five = 5;
let ten = 10;

let add = function(x, y) {
  x + y;
};

let result = add(five, ten);
!-/+*5;5 == 5; 5 != 10;
"foobar"
"foo bar"
if (5 < 10) { return true; } else { return false; }
10 == 10;
10 != 9;
$#&|`

	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.LET, "let"},
		{token.IDENT, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "ten"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "add"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "function"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.IDENT, "x"},
		{token.PLUS, "+"},
		{token.IDENT, "y"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "result"},
		{token.ASSIGN, "="},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.IDENT, "five"},
		{token.COMMA, ","},
		{token.IDENT, "ten"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.BANG, "!"},
		{token.MINUS, "-"},
		{token.SLASH, "/"},
		{token.PLUS, "+"},
		{token.ASTERISK, "*"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.INT, "5"},
		{token.EQ, "=="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.INT, "5"},
		{token.NotEq, "!="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.STRING, "\"foobar\""},
		{token.STRING, "\"foo bar\""},
		{token.If, "if"},
		{token.LPAREN, "("},
		{token.INT, "5"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.TRUE, "true"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.ELSE, "else"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.FALSE, "false"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.INT, "10"},
		{token.EQ, "=="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.INT, "10"},
		{token.NotEq, "!="},
		{token.INT, "9"},
		{token.SEMICOLON, ";"},
		{token.DOLLAR, "$"},
		{token.HASH, "#"},
		{token.AMPERSAND, "&"},
		{token.PIPE, "|"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestNextToken_ZshSpecific(t *testing.T) {
	input := `
	local var=(1 2 3)
	array=(one two three)
	${(f)foo}
	${(z)bar}
	`

	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.IDENT, "local"},
		{token.IDENT, "var"},
		{token.ASSIGN, "="},
		{token.LPAREN, "("},
		{token.INT, "1"},
		{token.INT, "2"},
		{token.INT, "3"},
		{token.RPAREN, ")"},
		{token.IDENT, "array"},
		{token.ASSIGN, "="},
		{token.LPAREN, "("},
		{token.IDENT, "one"},
		{token.IDENT, "two"},
		{token.IDENT, "three"},
		{token.RPAREN, ")"},
		{token.DollarLbrace, "${"},
		{token.LPAREN, "("},
		{token.IDENT, "f"},
		{token.RPAREN, ")"},
		{token.IDENT, "foo"},
		{token.RBRACE, "}"},
		{token.DollarLbrace, "${"},
		{token.LPAREN, "("},
		{token.IDENT, "z"},
		{token.RPAREN, ")"},
		{token.IDENT, "bar"},
		{token.RBRACE, "}"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestNextToken_ZshArrayAndCommandSubstitution(t *testing.T) {
	input := "`${my_array[1]}`"

	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.BACKTICK, "`"},
		{token.DollarLbrace, "${"},
		{token.IDENT, "my_array"},
		{token.LBRACKET, "["},
		{token.INT, "1"},
		{token.RBRACKET, "]"},
		{token.RBRACE, "}"},
		{token.BACKTICK, "`"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestNextToken_DoubleParens(t *testing.T) {
	input := "(( )) ;; ++"
	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.DoubleLparen, "(("},
		{token.DoubleRparen, "))"},
		{token.DSEMI, ";;"},
		{token.INC, "++"},
		{token.EOF, ""},
	}

	l := New(input)
	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestNextToken_DecrementAndPlusEqual(t *testing.T) {
	input := "-- +="
	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.DEC, "--"},
		{token.PLUSEQ, "+="},
		{token.EOF, ""},
	}

	l := New(input)
	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestNextToken_RedirectionOperators(t *testing.T) {
	input := ">> << >& <& <( >("
	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.GTGT, ">>"},
		{token.LTLT, "<<"},
		{token.GTAMP, ">&"},
		{token.LTAMP, "<&"},
		{token.LT_LPAREN, "<("},
		{token.GT_LPAREN, ">("},
		{token.EOF, ""},
	}

	l := New(input)
	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestNextToken_DoubleBrackets(t *testing.T) {
	input := "[[ ]]"
	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.LDBRACKET, "[["},
		{token.RDBRACKET, "]]"},
		{token.EOF, ""},
	}

	l := New(input)
	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestNextToken_SpecialTokens(t *testing.T) {
	input := "~ ^ % . : ?"
	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.TILDE, "~"},
		{token.CARET, "^"},
		{token.PERCENT, "%"},
		{token.DOT, "."},
		{token.COLON, ":"},
		{token.QUESTION, "?"},
		{token.EOF, ""},
	}

	l := New(input)
	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestNextToken_EqTilde(t *testing.T) {
	input := "=~"
	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.EQTILDE, "=~"},
		{token.EOF, ""},
	}

	l := New(input)
	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestNextToken_EqLparen(t *testing.T) {
	input := " =("
	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.EQ_LPAREN, "=("},
		{token.EOF, ""},
	}

	l := New(input)
	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}
	}
}

func TestNextToken_EqLparenNoSpace(t *testing.T) {
	input := "=("
	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.ASSIGN, "="},
		{token.LPAREN, "("},
		{token.EOF, ""},
	}

	l := New(input)
	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}
	}
}

func TestNextToken_LogicalOperators(t *testing.T) {
	input := "&& ||"
	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.AND, "&&"},
		{token.OR, "||"},
		{token.EOF, ""},
	}

	l := New(input)
	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestNextToken_DollarLparen(t *testing.T) {
	input := "$("
	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.DOLLAR_LPAREN, "$("},
		{token.EOF, ""},
	}

	l := New(input)
	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}
	}
}

func TestNextToken_DollarVariable(t *testing.T) {
	input := "$HOME"
	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.VARIABLE, "$HOME"},
		{token.EOF, ""},
	}

	l := New(input)
	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestNextToken_SingleQuotedString(t *testing.T) {
	input := "'hello world'"
	l := New(input)
	tok := l.NextToken()
	if tok.Type != token.STRING {
		t.Fatalf("expected STRING, got %q", tok.Type)
	}
	if tok.Literal != "'hello world'" {
		t.Fatalf("expected literal %q, got %q", "'hello world'", tok.Literal)
	}
}

func TestNextToken_UnterminatedString(t *testing.T) {
	input := `"hello`
	l := New(input)
	tok := l.NextToken()
	if tok.Type != token.STRING {
		t.Fatalf("expected STRING, got %q", tok.Type)
	}
	// Unterminated string should still produce a token
	if tok.Literal != `"hello` {
		t.Fatalf("expected literal %q, got %q", `"hello`, tok.Literal)
	}
}

func TestNextToken_EscapedString(t *testing.T) {
	input := `"hello \"world\""`
	l := New(input)
	tok := l.NextToken()
	if tok.Type != token.STRING {
		t.Fatalf("expected STRING, got %q", tok.Type)
	}
}

func TestNextToken_Shebang(t *testing.T) {
	input := "#!/bin/zsh\necho"
	l := New(input)
	tok := l.NextToken()
	if tok.Type != token.SHEBANG {
		t.Fatalf("expected SHEBANG, got %q", tok.Type)
	}
}

func TestNextToken_InlineComment(t *testing.T) {
	input := "echo # this is a comment\nfoo"
	l := New(input)
	tok1 := l.NextToken() // echo
	if tok1.Type != token.IDENT || tok1.Literal != "echo" {
		t.Fatalf("expected IDENT 'echo', got %q %q", tok1.Type, tok1.Literal)
	}
	tok2 := l.NextToken() // # (as HASH, not a comment at non-bol without space context)
	_ = tok2
}

func TestNextToken_IllegalCharacter(t *testing.T) {
	input := "\x01"
	l := New(input)
	tok := l.NextToken()
	if tok.Type != token.ILLEGAL {
		t.Fatalf("expected ILLEGAL, got %q", tok.Type)
	}
}

func TestNextToken_HasPrecedingSpace(t *testing.T) {
	input := "a b"
	l := New(input)

	tok1 := l.NextToken()
	if tok1.HasPrecedingSpace {
		t.Error("first token should not have preceding space")
	}

	tok2 := l.NextToken()
	if !tok2.HasPrecedingSpace {
		t.Error("second token should have preceding space")
	}
}

func TestNextToken_MinusKeyword(t *testing.T) {
	input := "-eq -ne -lt -le -gt -ge"
	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.EQ_NUM, "-eq"},
		{token.NE_NUM, "-ne"},
		{token.LT_NUM, "-lt"},
		{token.LE_NUM, "-le"},
		{token.GT_NUM, "-gt"},
		{token.GE_NUM, "-ge"},
		{token.EOF, ""},
	}

	l := New(input)
	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestNextToken_MinusNonKeyword(t *testing.T) {
	// -f is not a keyword, so it should be MINUS followed by IDENT
	input := "-f"
	l := New(input)
	tok := l.NextToken()
	if tok.Type != token.MINUS {
		t.Fatalf("expected MINUS, got %q %q", tok.Type, tok.Literal)
	}
}

func TestNextToken_LineAndColumnTracking(t *testing.T) {
	input := "a\nb"
	l := New(input)

	tok1 := l.NextToken()
	if tok1.Line != 1 {
		t.Errorf("first token line expected 1, got %d", tok1.Line)
	}

	tok2 := l.NextToken()
	if tok2.Line != 2 {
		t.Errorf("second token line expected 2, got %d", tok2.Line)
	}
}

// TestNextToken_KeywordFollowedByEquals exercises the regression fix for
// https://github.com/afadesigns/zshellcheck/issues/435: when an identifier
// that happens to match a Zsh keyword (`if`, `of`, `while`, `do`, etc.) is
// immediately followed by `=`, it is a flag-style assignment token (as in
// `dd if=foo of=bar`), not the keyword. The lexer must return IDENT so the
// parser treats the whole run as a single word.
func TestNextToken_KeywordFollowedByEquals(t *testing.T) {
	cases := []struct {
		input  string
		expect []struct {
			t token.Type
			l string
		}
	}{
		{
			input: `dd if=src of=dst`,
			expect: []struct {
				t token.Type
				l string
			}{
				{token.IDENT, "dd"},
				{token.IDENT, "if"},
				{token.ASSIGN, "="},
				{token.IDENT, "src"},
				{token.IDENT, "of"},
				{token.ASSIGN, "="},
				{token.IDENT, "dst"},
			},
		},
		{
			input: `while=10 do=go`,
			expect: []struct {
				t token.Type
				l string
			}{
				{token.IDENT, "while"},
				{token.ASSIGN, "="},
				{token.INT, "10"},
				{token.IDENT, "do"},
				{token.ASSIGN, "="},
				{token.IDENT, "go"},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.input, func(t *testing.T) {
			l := New(c.input)
			for i, exp := range c.expect {
				tok := l.NextToken()
				if tok.Type != exp.t || tok.Literal != exp.l {
					t.Fatalf("token[%d]: expected {%s %q}, got {%s %q}", i, exp.t, exp.l, tok.Type, tok.Literal)
				}
			}
		})
	}
}
