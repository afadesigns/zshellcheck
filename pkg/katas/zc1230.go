package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1230",
		Title:    "Use `ping -c N` in scripts to limit ping count",
		Severity: SeverityWarning,
		Description: "`ping` without `-c` runs indefinitely on Linux, hanging scripts. " +
			"Always specify `-c N` to limit the number of packets.",
		Check: checkZC1230,
	})
}

func checkZC1230(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "ping" {
		return nil
	}

	hasCount := false
	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-c" || val == "-W" {
			hasCount = true
		}
	}

	if !hasCount {
		return []Violation{{
			KataID: "ZC1230",
			Message: "Use `ping -c N` in scripts. Without `-c`, ping runs " +
				"indefinitely on Linux and will hang the script.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityWarning,
		}}
	}

	return nil
}
