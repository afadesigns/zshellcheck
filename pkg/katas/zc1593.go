package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1593",
		Title:    "Error on `blkdiscard` — issues TRIM/DISCARD across the whole device (data loss)",
		Severity: SeverityError,
		Description: "`blkdiscard $DEV` tells the underlying SSD controller to invalidate every " +
			"block in the range. On most modern drives the data is unrecoverable the moment the " +
			"controller acknowledges — even forensic recovery cannot pull it back. Scripts that " +
			"reach this command from any codepath an attacker or typo can trigger destroy the " +
			"drive. Gate it behind interactive confirmation, not shell flow control.",
		Check: checkZC1593,
	})
}

func checkZC1593(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "blkdiscard" {
		return nil
	}
	if len(cmd.Arguments) == 0 {
		return nil
	}

	return []Violation{{
		KataID: "ZC1593",
		Message: "`blkdiscard` issues TRIM/DISCARD across the full device — data is " +
			"unrecoverable once the controller acknowledges. Require operator confirmation " +
			"before running.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}
