package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1230",
		Title:    "Use `ping -c N` in scripts to limit ping count",
		Severity: SeverityWarning,
		Description: "`ping` without `-c` runs indefinitely on Linux, hanging scripts. " +
			"Always specify `-c N` to limit the number of packets.",
		Check: checkZC1230,
		Fix:   fixZC1230,
	})
}

// fixZC1230 inserts ` -c 4` after the `ping` command name. Detector
// already guards against an existing `-c` / `-W` flag, so the
// insertion is safe and idempotent on a re-run.
func fixZC1230(node ast.Node, v Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "ping" {
		return nil
	}
	nameOff := LineColToByteOffset(source, v.Line, v.Column)
	if nameOff < 0 {
		return nil
	}
	if IdentLenAt(source, nameOff) != len("ping") {
		return nil
	}
	insertAt := nameOff + len("ping")
	insLine, insCol := offsetLineColZC1230(source, insertAt)
	if insLine < 0 {
		return nil
	}
	return []FixEdit{{
		Line:    insLine,
		Column:  insCol,
		Length:  0,
		Replace: " -c 4",
	}}
}

func offsetLineColZC1230(source []byte, offset int) (int, int) {
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

func checkZC1230(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "ping" {
		return nil
	}

	hasCount := false
	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-c" || val == "-W" {
			hasCount = true
		}
	}

	if !hasCount {
		return []Violation{{
			KataID: "ZC1230",
			Message: "Use `ping -c N` in scripts. Without `-c`, ping runs " +
				"indefinitely on Linux and will hang the script.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityWarning,
		}}
	}

	return nil
}
