package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1724",
		Title:    "Warn on `pacman -Sy <pkg>` — partial upgrade, breaks dependency closure",
		Severity: SeverityWarning,
		Description: "Arch Linux is rolling-release on the invariant that the local package " +
			"database and the installed package set move together. `pacman -Sy <pkg>` " +
			"refreshes the database and installs ONE package against the new metadata while " +
			"every other installed package stays at its old version. The new package's " +
			"dependency closure pulls libraries newer than what the rest of the system has, " +
			"leaving a half-upgraded state that often manifests as `error while loading " +
			"shared libraries`. Run a full `pacman -Syu` first, then install (`pacman -S " +
			"<pkg>`); for CI use `pacman -Syu --noconfirm <pkg>` so the upgrade and install " +
			"are atomic.",
		Check: checkZC1724,
	})
}

func checkZC1724(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "pacman" {
		return nil
	}

	hasSyNoU := false
	hasPkg := false
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if strings.HasPrefix(v, "-") && len(v) >= 2 {
			letters := v[1:]
			// Must contain both 'S' and 'y' but not 'u' (would be -Syu).
			if strings.Contains(letters, "S") && strings.Contains(letters, "y") && !strings.Contains(letters, "u") {
				hasSyNoU = true
			}
			continue
		}
		if v != "" {
			hasPkg = true
		}
	}

	if !hasSyNoU || !hasPkg {
		return nil
	}

	return []Violation{{
		KataID: "ZC1724",
		Message: "`pacman -Sy <pkg>` is a partial-upgrade footgun — refresh the DB but " +
			"install only one package against the newer metadata. Use `pacman -Syu` " +
			"first, then `pacman -S <pkg>` (or `pacman -Syu --noconfirm <pkg>` to keep " +
			"it atomic).",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}
