package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1141",
		Title: "Avoid `curl | sh` pattern",
		Description: "Piping curl output to sh/bash/zsh is a security risk. Download first, " +
			"verify integrity (checksum or signature), then execute.",
		Severity: SeverityWarning,
		Check:    checkZC1141,
	})
}

func checkZC1141(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "curl" {
		return nil
	}

	// Check for -s or -sSL flags which suggest piping intent
	hasSilent := false
	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-s" || val == "-sS" || val == "-sSL" || val == "-sL" {
			hasSilent = true
		}
	}

	if !hasSilent {
		return nil
	}

	return []Violation{{
		KataID: "ZC1141",
		Message: "Avoid `curl -s URL | sh`. Download the script first, verify its integrity, " +
			"then execute. Piping directly from the internet is a supply-chain risk.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}
