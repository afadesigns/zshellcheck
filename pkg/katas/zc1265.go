package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1265",
		Title:    "Use `systemctl enable --now` to enable and start together",
		Severity: SeverityStyle,
		Description: "`systemctl enable` without `--now` only enables on next boot. " +
			"Use `--now` to enable and immediately start the service.",
		Check: checkZC1265,
		Fix:   fixZC1265,
	})
}

// fixZC1265 inserts ` --now` after the `enable` subcommand in a
// `systemctl enable …` invocation. Same subcommand-level insertion
// pattern as ZC1234's docker-run --rm.
func fixZC1265(node ast.Node, _ Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "systemctl" {
		return nil
	}
	var enableArg ast.Expression
	for _, arg := range cmd.Arguments {
		if arg.String() == "enable" {
			enableArg = arg
			break
		}
	}
	if enableArg == nil {
		return nil
	}
	tok := enableArg.TokenLiteralNode()
	off := LineColToByteOffset(source, tok.Line, tok.Column)
	if off < 0 || off+len("enable") > len(source) {
		return nil
	}
	if string(source[off:off+len("enable")]) != "enable" {
		return nil
	}
	insertAt := off + len("enable")
	insLine, insCol := offsetLineColZC1265(source, insertAt)
	if insLine < 0 {
		return nil
	}
	return []FixEdit{{
		Line:    insLine,
		Column:  insCol,
		Length:  0,
		Replace: " --now",
	}}
}

func offsetLineColZC1265(source []byte, offset int) (int, int) {
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

func checkZC1265(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "systemctl" {
		return nil
	}

	hasEnable := false
	hasNow := false

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "enable" {
			hasEnable = true
		}
		if val == "--now" {
			hasNow = true
		}
	}

	if hasEnable && !hasNow {
		return []Violation{{
			KataID: "ZC1265",
			Message: "Use `systemctl enable --now` to enable and start the service immediately. " +
				"Without `--now`, the service only starts on next boot.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityStyle,
		}}
	}

	return nil
}
