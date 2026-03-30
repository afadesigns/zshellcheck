package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1144",
		Title:    "Avoid `trap` with signal numbers — use names",
		Severity: SeverityInfo,
		Description: "Signal numbers vary across platforms. Use signal names like " +
			"`SIGTERM`, `SIGINT`, `EXIT` instead of numeric values for portability.",
		Check: checkZC1144,
	})
}

func checkZC1144(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "trap" {
		return nil
	}

	if len(cmd.Arguments) < 2 {
		return nil
	}

	// Check last arguments for numeric signal values
	for i := 1; i < len(cmd.Arguments); i++ {
		val := cmd.Arguments[i].String()
		// Numeric signals: 1-31
		isNumeric := len(val) > 0
		for _, ch := range val {
			if ch < '0' || ch > '9' {
				isNumeric = false
				break
			}
		}
		if isNumeric && val != "0" {
			return []Violation{{
				KataID: "ZC1144",
				Message: "Use signal names (`SIGTERM`, `SIGINT`, `EXIT`) instead of numbers in `trap`. " +
					"Signal numbers vary across platforms.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityInfo,
			}}
		}
	}

	return nil
}
