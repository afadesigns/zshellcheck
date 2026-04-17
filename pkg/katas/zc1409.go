package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1409",
		Title:    "Avoid `[ -N file ]` / `test -N file` — Bash-only, use Zsh `zstat` for mtime comparison",
		Severity: SeverityInfo,
		Description: "`[ -N file ]` and `test -N file` test whether a file has been modified since " +
			"last read (Bash extension). Zsh does not implement `-N`. Use the `zsh/stat` module " +
			"to compare `atime` and `mtime` explicitly: `zstat -H s file; (( s[mtime] > s[atime] ))`.",
		Check: checkZC1409,
	})
}

func checkZC1409(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "test" && ident.Value != "[" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		if arg.String() == "-N" {
			return []Violation{{
				KataID: "ZC1409",
				Message: "`test -N file` (modified-since-read) is a Bash extension. In Zsh use " +
					"`zmodload zsh/stat; zstat -H s file; (( s[mtime] > s[atime] ))`.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityInfo,
			}}
		}
	}

	return nil
}
