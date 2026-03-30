package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1140",
		Title: "Use `command -v` instead of `hash` for command existence",
		Description: "`hash cmd` is a POSIX way to check command existence but provides " +
			"poor error messages. Use `command -v cmd` for cleaner checks in Zsh.",
		Severity: SeverityStyle,
		Check:    checkZC1140,
	})
}

func checkZC1140(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "hash" {
		return nil
	}

	// Only flag bare hash (command existence check)
	// hash -r (rehash) is a different valid use
	for _, arg := range cmd.Arguments {
		val := strings.Trim(arg.String(), "'\"")
		if len(val) > 0 && val[0] == '-' {
			return nil
		}
	}

	if len(cmd.Arguments) != 1 {
		return nil
	}

	return []Violation{{
		KataID: "ZC1140",
		Message: "Use `command -v cmd` instead of `hash cmd` for command existence checks. " +
			"`command -v` provides clearer semantics in Zsh.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityStyle,
	}}
}
