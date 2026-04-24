package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1076",
		Title: "Use `autoload -Uz` for lazy loading",
		Description: "When using `autoload`, prefer `-Uz` to ensure standard Zsh behavior (no alias expansion, zsh style). " +
			"`-U` prevents alias expansion, and `-z` ensures Zsh style autoloading.",
		Severity: SeverityStyle,
		Check:    checkZC1076,
		Fix:      fixZC1076,
	})
}

// fixZC1076 inserts ` -Uz` after the `autoload` command name. Only
// fires when neither `U` nor `z` are already present; the detector
// already gates on that. Idempotent on re-run once both flags exist.
func fixZC1076(node ast.Node, v Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	if cmd.Name.String() != "autoload" {
		return nil
	}
	nameOffset := LineColToByteOffset(source, v.Line, v.Column)
	if nameOffset < 0 {
		return nil
	}
	nameLen := IdentLenAt(source, nameOffset)
	if nameLen != len("autoload") {
		return nil
	}
	insertAt := nameOffset + nameLen
	insLine, insCol := offsetLineColZC1076(source, insertAt)
	if insLine < 0 {
		return nil
	}
	return []FixEdit{{
		Line:    insLine,
		Column:  insCol,
		Length:  0,
		Replace: " -Uz",
	}}
}

func offsetLineColZC1076(source []byte, offset int) (int, int) {
	if offset < 0 || offset > len(source) {
		return -1, -1
	}
	line := 1
	col := 1
	for i := 0; i < offset; i++ {
		if source[i] == '\n' {
			line++
			col = 1
			continue
		}
		col++
	}
	return line, col
}

func checkZC1076(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	if cmd.Name.String() != "autoload" {
		return nil
	}

	hasU := false
	hasZ := false

	for _, arg := range cmd.Arguments {
		argStr := arg.String()
		argStr = strings.Trim(argStr, "\"'")
		if strings.HasPrefix(argStr, "-") {
			if strings.Contains(argStr, "U") {
				hasU = true
			}
			if strings.Contains(argStr, "z") {
				hasZ = true
			}
		}
	}

	if !hasU || !hasZ {
		return []Violation{{
			KataID:  "ZC1076",
			Message: "Use `autoload -Uz` to ensure consistent and safe function loading.",
			Line:    cmd.TokenLiteralNode().Line,
			Column:  cmd.TokenLiteralNode().Column,
			Level:   SeverityStyle,
		}}
	}

	return nil
}
