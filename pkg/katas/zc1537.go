package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1537",
		Title:    "Error on `lvremove -f` / `vgremove -f` / `pvremove -f` — force-destroys LVM metadata",
		Severity: SeverityError,
		Description: "The `-f`/`--force` flag on the LVM destructive commands skips the " +
			"confirmation prompt that protects against a typo in the volume name. If the " +
			"target variable resolves to the wrong VG/LV/PV (empty, unset, different host), " +
			"a single line destroys every filesystem on top of that LVM stack. Leave the " +
			"prompt in and pipe `yes` to it only when you have explicitly confirmed the " +
			"target immediately beforehand.",
		Check: checkZC1537,
	})
}

func checkZC1537(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "lvremove" && ident.Value != "vgremove" && ident.Value != "pvremove" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-f" || v == "-ff" || v == "--force" {
			return []Violation{{
				KataID: "ZC1537",
				Message: "`" + ident.Value + " " + v + "` skips the confirmation — a typo in " +
					"the volume name destroys every filesystem on top of it.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}
	return nil
}
