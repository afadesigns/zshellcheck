package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1226",
		Title:    "Use `dmesg -T` or `--time-format=iso` for readable timestamps",
		Severity: SeverityStyle,
		Description: "`dmesg` without `-T` shows raw kernel timestamps in seconds since boot. " +
			"Use `-T` for human-readable timestamps or `--time-format=iso` for ISO 8601.",
		Check: checkZC1226,
	})
}

func checkZC1226(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "dmesg" {
		return nil
	}

	hasTimeFlag := false
	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-T" || val == "-t" || val == "--ctime" || val == "--reltime" {
			hasTimeFlag = true
		}
	}

	if !hasTimeFlag && len(cmd.Arguments) > 0 {
		return []Violation{{
			KataID: "ZC1226",
			Message: "Use `dmesg -T` for human-readable timestamps instead of raw " +
				"kernel boot-seconds. Or use `--time-format=iso` for ISO 8601.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityStyle,
		}}
	}

	return nil
}
