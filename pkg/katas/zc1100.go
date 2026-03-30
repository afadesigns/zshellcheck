package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1100",
		Title: "Use parameter expansion instead of `dirname`/`basename`",
		Description: "Zsh parameter expansion `${var%/*}` (dirname) and `${var##*/}` (basename) " +
			"avoid spawning external processes for simple path manipulation.",
		Severity: SeverityStyle,
		Check:    checkZC1100,
	})
}

func checkZC1100(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	name := ident.Value
	if name != "dirname" && name != "basename" {
		return nil
	}

	// Only flag simple single-argument calls
	// basename with -s or -a flags is more complex
	for _, arg := range cmd.Arguments {
		val := arg.String()
		if len(val) > 0 && val[0] == '-' {
			return nil
		}
	}

	if len(cmd.Arguments) != 1 {
		return nil
	}

	var msg string
	if name == "dirname" {
		msg = "Use `${var%/*}` instead of `dirname` to extract the directory path. " +
			"Parameter expansion avoids spawning an external process."
	} else {
		msg = "Use `${var##*/}` instead of `basename` to extract the filename. " +
			"Parameter expansion avoids spawning an external process."
	}

	return []Violation{{
		KataID:  "ZC1100",
		Message: msg,
		Line:    cmd.Token.Line,
		Column:  cmd.Token.Column,
		Level:   SeverityStyle,
	}}
}
