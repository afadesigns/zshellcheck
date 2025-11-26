package lexer

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/token"
)

func FuzzLexer(f *testing.F) {
	testcases := []string{
		"echo hello",
		"if [[ -f foo ]]; then echo bar; fi",
		"ls -la | grep foo",
		"${var:-default}",
		"$(command)",
		"`backticks`",
		"# comment",
		"var=value",
		"foo=bar baz=qux",
	}

	for _, tc := range testcases {
		f.Add(tc)
	}

	f.Fuzz(func(t *testing.T, input string) {
		l := New(input)
		for {
			tok := l.NextToken()
			if tok.Type == token.EOF {
				break
			}
		}
	})
}
