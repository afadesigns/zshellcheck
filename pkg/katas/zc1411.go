package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1411",
		Title:    "Use Zsh `disable` instead of Bash `enable -n` to hide builtins",
		Severity: SeverityStyle,
		Description: "Bash's `enable -n name` disables a builtin so that the external of the same " +
			"name is used. Zsh provides a dedicated `disable` builtin: `disable name` achieves " +
			"the same in one verb. Re-enable later with `enable name`.",
		Check: checkZC1411,
		Fix:   fixZC1411,
	})
}

// fixZC1411 collapses `enable -n NAME` into `disable NAME`. The span
// covers the `enable` command name and the `-n` flag in a single edit;
// trailing builtin name(s) stay in place.
func fixZC1411(node ast.Node, v Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "enable" {
		return nil
	}
	var dashN ast.Expression
	for _, arg := range cmd.Arguments {
		if arg.String() == "-n" {
			dashN = arg
			break
		}
	}
	if dashN == nil {
		return nil
	}
	nameOff := LineColToByteOffset(source, v.Line, v.Column)
	if nameOff < 0 || nameOff+len("enable") > len(source) {
		return nil
	}
	if string(source[nameOff:nameOff+len("enable")]) != "enable" {
		return nil
	}
	dashTok := dashN.TokenLiteralNode()
	dashOff := LineColToByteOffset(source, dashTok.Line, dashTok.Column)
	if dashOff < 0 || dashOff+2 > len(source) {
		return nil
	}
	if string(source[dashOff:dashOff+2]) != "-n" {
		return nil
	}
	return []FixEdit{{
		Line:    v.Line,
		Column:  v.Column,
		Length:  dashOff + 2 - nameOff,
		Replace: "disable",
	}}
}

func checkZC1411(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "enable" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		if arg.String() == "-n" {
			return []Violation{{
				KataID: "ZC1411",
				Message: "Use Zsh `disable name` instead of `enable -n name`. Zsh has a " +
					"dedicated `disable` builtin that reads clearer.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityStyle,
			}}
		}
	}

	return nil
}
