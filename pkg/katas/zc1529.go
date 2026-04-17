package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1529",
		Title:    "Warn on `fsck -y` / `fsck.<fs> -y` — auto-answer yes can corrupt",
		Severity: SeverityWarning,
		Description: "`fsck -y` answers `yes` to every repair prompt. For the happy case it is a " +
			"timesaver, but on a filesystem with unusual corruption (bad sector storm, mangled " +
			"journal after power loss) the automatic answer can turn salvageable data into " +
			"`lost+found` entries or zero it outright. In scripts, prefer `fsck -n` for a " +
			"dry-run and let a human adjudicate a real repair, or run with `-p` (preen: only " +
			"safe automatic fixes).",
		Check: checkZC1529,
	})
}

func checkZC1529(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "fsck" && !strings.HasPrefix(ident.Value, "fsck.") &&
		ident.Value != "e2fsck" && ident.Value != "xfs_repair" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		if arg.String() == "-y" {
			return []Violation{{
				KataID: "ZC1529",
				Message: "`" + ident.Value + " -y` answers yes to every repair prompt — can " +
					"destroy salvageable data. Prefer `-n` (dry-run) or `-p` (preen).",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}
