package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1872",
		Title:    "Error on `badblocks -w` — destructive write-mode pattern test wipes the device",
		Severity: SeverityError,
		Description: "`badblocks -w` (alias `--write-mode`) runs the write-mode bad-block check, " +
			"which overwrites every sector of the target device with a test pattern and " +
			"reads it back. On a fresh drive about to be formatted that is exactly what " +
			"you want; on an already-populated disk it is a silent data-wipe — the " +
			"command returns success even as it bulldozes the filesystem. If only " +
			"non-destructive checking is needed, use `badblocks -n` (read-test-restore) " +
			"or `badblocks` without any mode flag (read-only). When a true destructive " +
			"test is intended, gate the call behind a confirmation prompt and a freshly " +
			"partitioned device.",
		Check: checkZC1872,
	})
}

func checkZC1872(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "badblocks" {
		return nil
	}
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-w" || v == "--write-mode" {
			return zc1872Hit(cmd)
		}
		// Also catch combined short-option clusters like `-wsv`.
		if len(v) > 1 && v[0] == '-' && v[1] != '-' && strings.ContainsRune(v[1:], 'w') {
			return zc1872Hit(cmd)
		}
	}
	return nil
}

func zc1872Hit(cmd *ast.SimpleCommand) []Violation {
	return []Violation{{
		KataID: "ZC1872",
		Message: "`badblocks -w` overwrites every sector of the target device — " +
			"silent data wipe on a populated disk. Use `-n` (non-destructive) " +
			"or gate destructive runs behind a confirmation and a fresh partition.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}
