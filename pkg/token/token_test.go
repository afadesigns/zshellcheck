// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package token

import "testing"

func TestLookupIdent_Keywords(t *testing.T) {
	tests := []struct {
		input    string
		expected Type
	}{
		{"function", FUNCTION},
		{"let", LET},
		{"true", TRUE},
		{"false", FALSE},
		{"if", If},
		{"else", ELSE},
		{"return", RETURN},
		{"then", THEN},
		{"fi", Fi},
		{"for", FOR},
		{"while", WHILE},
		{"do", DO},
		{"done", DONE},
		{"in", IN},
		{"case", CASE},
		{"esac", ESAC},
		{"elif", ELIF},
		{"select", SELECT},
		{"coproc", COPROC},
		{"typeset", TYPESET},
		{"declare", DECLARE},
		{"-eq", EQ_NUM},
		{"-ne", NE_NUM},
		{"-lt", LT_NUM},
		{"-le", LE_NUM},
		{"-gt", GT_NUM},
		{"-ge", GE_NUM},
		{"/", SLASH},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := LookupIdent(tt.input)
			if got != tt.expected {
				t.Errorf("LookupIdent(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestLookupIdent_Identifiers(t *testing.T) {
	tests := []string{
		"foo",
		"bar",
		"myVar",
		"some_function",
		"x",
	}

	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			got := LookupIdent(input)
			if got != IDENT {
				t.Errorf("LookupIdent(%q) = %q, want %q", input, got, IDENT)
			}
		})
	}
}

func TestTokenStruct(t *testing.T) {
	tok := Token{
		Type:              IDENT,
		Literal:           "hello",
		Line:              10,
		Column:            5,
		HasPrecedingSpace: true,
	}

	if tok.Type != IDENT {
		t.Errorf("expected Type=IDENT, got %q", tok.Type)
	}
	if tok.Literal != "hello" {
		t.Errorf("expected Literal=hello, got %q", tok.Literal)
	}
	if tok.Line != 10 {
		t.Errorf("expected Line=10, got %d", tok.Line)
	}
	if tok.Column != 5 {
		t.Errorf("expected Column=5, got %d", tok.Column)
	}
	if !tok.HasPrecedingSpace {
		t.Error("expected HasPrecedingSpace=true")
	}
}
