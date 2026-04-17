package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1525",
		Title:    "Warn on `ping -f` — flood ping sends packets as fast as possible (DoS vector)",
		Severity: SeverityWarning,
		Description: "`ping -f` (flood mode) removes the one-per-second rate limit and sends " +
			"ICMP echo requests in a tight loop. It's a root-only builtin specifically because " +
			"it can saturate a slow link or overload a low-end host. Legitimate uses exist " +
			"(latency benchmarking, stress testing known-internal targets), but in a script " +
			"aimed at arbitrary hosts it is a DoS primitive. Scope tightly and document.",
		Check: checkZC1525,
	})
}

func checkZC1525(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "ping" && ident.Value != "ping6" && ident.Value != "fping" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		if arg.String() == "-f" {
			return []Violation{{
				KataID: "ZC1525",
				Message: "`" + ident.Value + " -f` (flood) bypasses the rate limit and can " +
					"saturate slow links — DoS primitive. Scope tightly and document.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}
