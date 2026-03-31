package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1266",
		Title:    "Use `nproc` instead of parsing `/proc/cpuinfo`",
		Severity: SeverityStyle,
		Description: "Parsing `/proc/cpuinfo` for CPU count is fragile and platform-specific. " +
			"`nproc` is a portable, dedicated tool for this purpose.",
		Check: checkZC1266,
	})
}

func checkZC1266(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "cat" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "/proc/cpuinfo" {
			return []Violation{{
				KataID: "ZC1266",
				Message: "Use `nproc` instead of parsing `/proc/cpuinfo` for CPU count. " +
					"`nproc` is portable and available on Linux and macOS (via coreutils).",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityStyle,
			}}
		}
	}

	return nil
}
