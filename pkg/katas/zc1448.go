package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1448",
		Title:    "`apt-get install` / `apt install` without `-y` hangs in non-interactive scripts",
		Severity: SeverityWarning,
		Description: "In provisioning scripts, `apt-get install foo` (no `-y`) waits for " +
			"interactive confirmation and stalls CI/Dockerfiles indefinitely. Always pass `-y` " +
			"(or `--yes`), and for unattended upgrades also set " +
			"`DEBIAN_FRONTEND=noninteractive` in the environment.",
		Check: checkZC1448,
	})
}

func checkZC1448(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "apt-get" && ident.Value != "apt" {
		return nil
	}

	hasInstall := false
	hasYes := false
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "install" || v == "upgrade" || v == "dist-upgrade" || v == "full-upgrade" {
			hasInstall = true
		}
		if v == "-y" || v == "--yes" || v == "--assume-yes" {
			hasYes = true
		}
	}
	if hasInstall && !hasYes {
		return []Violation{{
			KataID: "ZC1448",
			Message: "`apt-get install`/`apt install` without `-y` hangs on the interactive " +
				"prompt in scripts. Add `-y` and set DEBIAN_FRONTEND=noninteractive.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityWarning,
		}}
	}

	return nil
}
