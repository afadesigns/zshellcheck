package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1135",
		Title: "Avoid `env VAR=val cmd` — use inline assignment",
		Description: "Zsh supports inline environment variable assignment with `VAR=val cmd`. " +
			"Avoid spawning `env` for simple variable-prefixed command execution.",
		Severity: SeverityStyle,
		Check:    checkZC1135,
		Fix:      fixZC1135,
	})
}

// fixZC1135 strips the `env ` prefix from `env VAR=val cmd`. Detector
// already forbids `env` flags, so the remaining args form a valid
// inline-assignment command.
func fixZC1135(node ast.Node, v Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "env" {
		return nil
	}
	nameOff := LineColToByteOffset(source, v.Line, v.Column)
	if nameOff < 0 || nameOff+len("env") > len(source) {
		return nil
	}
	if string(source[nameOff:nameOff+len("env")]) != "env" {
		return nil
	}
	// Span covers `env` plus the whitespace that follows it.
	end := nameOff + len("env")
	for end < len(source) && (source[end] == ' ' || source[end] == '\t') {
		end++
	}
	return []FixEdit{{
		Line:    v.Line,
		Column:  v.Column,
		Length:  end - nameOff,
		Replace: "",
	}}
}

func checkZC1135(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "env" {
		return nil
	}

	// Only flag env with VAR=val patterns followed by a command
	// Skip env -i (clean environment), env -u (unset), env -S
	for _, arg := range cmd.Arguments {
		val := arg.String()
		if len(val) > 0 && val[0] == '-' {
			return nil
		}
	}

	// Check if any argument contains = (env var assignment)
	for _, arg := range cmd.Arguments {
		val := arg.String()
		if strings.Contains(val, "=") {
			return []Violation{{
				KataID: "ZC1135",
				Message: "Use inline `VAR=val cmd` instead of `env VAR=val cmd`. " +
					"Zsh supports inline env assignment without spawning env.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityStyle,
			}}
		}
	}

	return nil
}
