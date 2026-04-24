package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1128",
		Title: "Use `> file` instead of `touch file` for creation",
		Description: "If the goal is to create an empty file, `> file` does it without " +
			"spawning `touch`. Use `touch` only when you need to update timestamps.",
		Severity: SeverityStyle,
		Check:    checkZC1128,
		Fix:      fixZC1128,
	})
}

// fixZC1128 rewrites `touch file` into `> file`. Detector already
// guards against flagged forms (timestamp updates) and multi-arg
// invocations, so the fix covers only the single-file case.
func fixZC1128(node ast.Node, v Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	if len(cmd.Arguments) != 1 {
		return nil
	}
	nameOff := LineColToByteOffset(source, v.Line, v.Column)
	if nameOff < 0 || nameOff+len("touch") > len(source) {
		return nil
	}
	if string(source[nameOff:nameOff+len("touch")]) != "touch" {
		return nil
	}
	return []FixEdit{{
		Line:    v.Line,
		Column:  v.Column,
		Length:  len("touch"),
		Replace: ">",
	}}
}

func checkZC1128(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "touch" {
		return nil
	}

	// Skip touch with flags (timestamps: -t, -d, -r, -a, -m)
	for _, arg := range cmd.Arguments {
		val := arg.String()
		if len(val) > 0 && val[0] == '-' {
			return nil
		}
	}

	// Only flag touch with a single file argument
	if len(cmd.Arguments) != 1 {
		return nil
	}

	return []Violation{{
		KataID: "ZC1128",
		Message: "Use `> file` instead of `touch file` to create an empty file. " +
			"This avoids spawning an external process.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityStyle,
	}}
}
