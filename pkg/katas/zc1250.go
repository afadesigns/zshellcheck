package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1250",
		Title:    "Use `gpg --batch` in scripts for non-interactive operation",
		Severity: SeverityWarning,
		Description: "`gpg` without `--batch` may prompt for passphrases or confirmations. " +
			"Use `--batch` and `--yes` for fully non-interactive GPG operations in scripts.",
		Check: checkZC1250,
	})
}

func checkZC1250(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "gpg" {
		return nil
	}

	hasBatch := false
	hasOperation := false

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-b" {
			hasBatch = true
		}
		if val == "-d" || val == "-s" || val == "-e" || val == "-c" {
			hasOperation = true
		}
	}

	if hasOperation && !hasBatch {
		return []Violation{{
			KataID: "ZC1250",
			Message: "Use `gpg --batch` in scripts for non-interactive operation. " +
				"Without `--batch`, gpg may prompt for passphrases or confirmations.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityWarning,
		}}
	}

	return nil
}
