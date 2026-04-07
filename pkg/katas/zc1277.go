package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1277",
		Title:    "Use Zsh `${var:l}` / `${var:u}` instead of `tr` for case conversion",
		Severity: SeverityStyle,
		Description: "Zsh provides the `:l` (lowercase) and `:u` (uppercase) modifiers for parameter " +
			"expansion, avoiding the overhead of piping through `tr` for case conversion.",
		Check: checkZC1277,
	})
}

func checkZC1277(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "tr" {
		return nil
	}

	if len(cmd.Arguments) < 2 {
		return nil
	}

	first := strings.Trim(cmd.Arguments[0].String(), "'\"")
	second := strings.Trim(cmd.Arguments[1].String(), "'\"")

	if (first == "[:upper:]" && second == "[:lower:]") ||
		(first == "[A-Z]" && second == "[a-z]") ||
		(first == "A-Z" && second == "a-z") {
		return []Violation{{
			KataID:  "ZC1277",
			Message: "Use Zsh parameter expansion `${var:l}` for lowercase conversion instead of `tr`. The `:l` modifier is built-in.",
			Line:    cmd.Token.Line,
			Column:  cmd.Token.Column,
			Level:   SeverityStyle,
		}}
	}

	if (first == "[:lower:]" && second == "[:upper:]") ||
		(first == "[a-z]" && second == "[A-Z]") ||
		(first == "a-z" && second == "A-Z") {
		return []Violation{{
			KataID:  "ZC1277",
			Message: "Use Zsh parameter expansion `${var:u}` for uppercase conversion instead of `tr`. The `:u` modifier is built-in.",
			Line:    cmd.Token.Line,
			Column:  cmd.Token.Column,
			Level:   SeverityStyle,
		}}
	}

	return nil
}
