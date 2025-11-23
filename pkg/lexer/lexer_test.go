package lexer

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/token"
)

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
		{token.STRING, "foobar"},
		{token.STRING, "foo bar"},
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
