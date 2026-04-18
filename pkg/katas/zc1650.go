package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1650",
		Title:    "Warn on `setopt RM_STAR_SILENT` / `unsetopt RM_STAR_WAIT` — removes `rm *` prompt",
		Severity: SeverityWarning,
		Description: "Zsh's default behaviour on an interactive `rm *` (or `rm /path/*`) is to " +
			"pause for 10 seconds and ask \"do you really want to delete N files?\" — the " +
			"`RM_STAR_WAIT` option. `setopt RM_STAR_SILENT` or `unsetopt RM_STAR_WAIT` both " +
			"disable the prompt. In a profile / dot file the option leaks to every future " +
			"interactive shell and removes a safety net that has saved countless home " +
			"directories.",
		Check: checkZC1650,
	})
}

func checkZC1650(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	switch ident.Value {
	case "setopt":
		for _, arg := range cmd.Arguments {
			norm := strings.ToLower(strings.ReplaceAll(arg.String(), "_", ""))
			if norm == "rmstarsilent" {
				return zc1650Hit(cmd, "setopt "+arg.String())
			}
		}
	case "unsetopt":
		for _, arg := range cmd.Arguments {
			norm := strings.ToLower(strings.ReplaceAll(arg.String(), "_", ""))
			if norm == "rmstarwait" {
				return zc1650Hit(cmd, "unsetopt "+arg.String())
			}
		}
	}
	return nil
}

func zc1650Hit(cmd *ast.SimpleCommand, desc string) []Violation {
	return []Violation{{
		KataID: "ZC1650",
		Message: "`" + desc + "` removes the `rm *` confirmation prompt — keep the default " +
			"`RM_STAR_WAIT` so accidental deletions pause before they happen.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}
