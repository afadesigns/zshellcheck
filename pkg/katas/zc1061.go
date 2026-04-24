package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:          "ZC1061",
		Title:       "Prefer `{start..end}` over `seq`",
		Description: "Using `seq` creates an external process. Zsh supports integer range expansion natively: `{1..10}`.",
		Severity:    SeverityStyle,
		Check:       checkZC1061,
		Fix:         fixZC1061,
	})
}

// fixZC1061 rewrites `seq M` / `seq M N` / `seq M S N` with integer
// literal arguments into Zsh's brace range expansion `{M..N}` or
// `{M..N..S}`. Forms with `-s sep` separator, variable arguments,
// or floats are left alone because the rewrite semantics differ.
func fixZC1061(node ast.Node, v Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "seq" {
		return nil
	}
	if len(cmd.Arguments) == 0 || len(cmd.Arguments) > 3 {
		return nil
	}
	nums := make([]string, 0, len(cmd.Arguments))
	for _, arg := range cmd.Arguments {
		s := arg.String()
		if !isAllDigits(s) {
			return nil
		}
		nums = append(nums, s)
	}
	var rng string
	switch len(nums) {
	case 1:
		rng = "{1.." + nums[0] + "}"
	case 2:
		rng = "{" + nums[0] + ".." + nums[1] + "}"
	case 3:
		rng = "{" + nums[0] + ".." + nums[2] + ".." + nums[1] + "}"
	}
	// Replace from the `seq` command name through the last argument
	// so the whole invocation becomes a single brace expansion.
	start := LineColToByteOffset(source, v.Line, v.Column)
	if start < 0 {
		return nil
	}
	lastArg := cmd.Arguments[len(cmd.Arguments)-1]
	lastTok := lastArg.TokenLiteralNode()
	lastOff := LineColToByteOffset(source, lastTok.Line, lastTok.Column)
	if lastOff < 0 {
		return nil
	}
	end := lastOff + len(lastTok.Literal)
	if end <= start {
		return nil
	}
	return []FixEdit{{
		Line:    v.Line,
		Column:  v.Column,
		Length:  end - start,
		Replace: rng,
	}}
}

func isAllDigits(s string) bool {
	if len(s) == 0 {
		return false
	}
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	return true
}

func checkZC1061(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	// Check if command is seq
	if name, ok := cmd.Name.(*ast.Identifier); ok && name.Value == "seq" {
		return []Violation{{
			KataID:  "ZC1061",
			Message: "Prefer `{start..end}` range expansion over `seq`. It is built-in and faster.",
			Line:    name.Token.Line,
			Column:  name.Token.Column,
			Level:   SeverityStyle,
		}}
	}

	return nil
}
