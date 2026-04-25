package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1095",
		Title: "Use `repeat N` for simple repetition",
		Description: "Zsh provides `repeat N do ... done` for running a block a fixed number of times. " +
			"It is cleaner than `for i in {1..N}` or C-style for loops when the iterator variable is unused.",
		Severity: SeverityStyle,
		Check:    checkZC1095,
		// Reuse the seq → {start..end} rewrite from ZC1061. The detector
		// here fires on a single-numeric-arg `seq N`, which fixZC1061
		// rewrites to `{1..N}` — exactly the brace expansion this kata
		// suggests for `for i in {1..N}`.
		Fix: fixZC1061,
	})
}

func checkZC1095(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "seq" {
		return nil
	}

	// Flag bare `seq N` calls (often used in `for i in $(seq N)`)
	// Only flag if seq has exactly one numeric argument
	if len(cmd.Arguments) != 1 {
		return nil
	}

	arg := cmd.Arguments[0].String()
	for _, ch := range arg {
		if ch < '0' || ch > '9' {
			return nil
		}
	}

	return []Violation{{
		KataID: "ZC1095",
		Message: "Use `repeat N do ... done` or `for i in {1..N}` instead of `seq N`. " +
			"Zsh has built-in constructs for repetition that avoid spawning an external process.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityStyle,
	}}
}
