package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1605",
		Title:    "Error on `debugfs -w DEV` — write-mode filesystem debugger bypasses journal",
		Severity: SeverityError,
		Description: "`debugfs -w` opens the filesystem in write mode. It sidesteps the kernel's " +
			"normal write path — the journal doesn't see the changes, filesystem locks are " +
			"ignored, and inodes / blocks can be edited directly. On a mounted filesystem this " +
			"corrupts state silently; even on an unmounted one, the operator can repoint a " +
			"directory entry at an arbitrary inode. Scripts should never need this — keep " +
			"`debugfs -w` as an interactive last-resort from a rescue environment.",
		Check: checkZC1605,
	})
}

func checkZC1605(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "debugfs" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		if arg.String() == "-w" {
			return []Violation{{
				KataID: "ZC1605",
				Message: "`debugfs -w` writes to the filesystem outside the kernel's normal " +
					"path — journal bypassed, locks ignored. Keep it as an interactive " +
					"rescue tool, not a script path.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}
	return nil
}
