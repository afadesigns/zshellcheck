package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1124",
		Title: "Use `: > file` instead of `cat /dev/null > file` to truncate",
		Description: "Truncating a file with `cat /dev/null > file` spawns an unnecessary process. " +
			"Use `: > file` or simply `> file` in Zsh.",
		Severity: SeverityStyle,
		Check:    checkZC1124,
		Fix:      fixZC1124,
	})
}

// fixZC1124 replaces the `cat /dev/null` prefix with the `:` builtin.
// The redirection and anything following stays in place:
// `cat /dev/null > file` becomes `: > file`. Only fires when
// `/dev/null` is the first argument — the detector already requires
// that shape.
func fixZC1124(node ast.Node, v Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	if len(cmd.Arguments) == 0 {
		return nil
	}
	var devNull ast.Expression
	for _, arg := range cmd.Arguments {
		if arg.String() == "/dev/null" {
			devNull = arg
			break
		}
	}
	if devNull == nil {
		return nil
	}
	nameOff := LineColToByteOffset(source, v.Line, v.Column)
	if nameOff < 0 || nameOff+3 > len(source) {
		return nil
	}
	if string(source[nameOff:nameOff+3]) != "cat" {
		return nil
	}
	argTok := devNull.TokenLiteralNode()
	argOff := LineColToByteOffset(source, argTok.Line, argTok.Column)
	if argOff < 0 {
		return nil
	}
	end := argOff + len("/dev/null")
	if end > len(source) || string(source[argOff:end]) != "/dev/null" {
		return nil
	}
	return []FixEdit{{
		Line:    v.Line,
		Column:  v.Column,
		Length:  end - nameOff,
		Replace: ":",
	}}
}

func checkZC1124(node ast.Node) []Violation {
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
		if val == "/dev/null" {
			return []Violation{{
				KataID: "ZC1124",
				Message: "Use `: > file` instead of `cat /dev/null > file` to truncate. " +
					"The `:` builtin avoids spawning cat.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityStyle,
			}}
		}
	}

	return nil
}
