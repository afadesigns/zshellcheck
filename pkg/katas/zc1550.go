package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1550",
		Title:    "Warn on `apt-mark hold <pkg>` — pins a package, blocks security updates",
		Severity: SeverityWarning,
		Description: "`apt-mark hold` tells apt to leave the package at its current version on " +
			"`apt upgrade` and `unattended-upgrades`. That is occasionally correct (pinning a " +
			"kernel variant for a driver, or a broken-upstream version) but silently keeps the " +
			"package vulnerable to every subsequent CVE. Document the reason in a comment, " +
			"schedule a review, and prefer `apt-mark unhold` + `apt upgrade <pkg>` over leaving " +
			"the pin in place indefinitely.",
		Check: checkZC1550,
	})
}

func checkZC1550(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "apt-mark" && ident.Value != "dpkg" {
		return nil
	}

	if ident.Value == "apt-mark" && len(cmd.Arguments) >= 2 &&
		cmd.Arguments[0].String() == "hold" {
		return []Violation{{
			KataID: "ZC1550",
			Message: "`apt-mark hold` pins the package — blocks future CVE fixes. Document " +
				"the reason and schedule an unhold review.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityWarning,
		}}
	}

	// `echo "<pkg> hold" | dpkg --set-selections` is the legacy equivalent; flag when
	// dpkg is called with --set-selections.
	if ident.Value == "dpkg" {
		for _, arg := range cmd.Arguments {
			if arg.String() == "--set-selections" {
				return []Violation{{
					KataID: "ZC1550",
					Message: "`dpkg --set-selections` with a `hold` entry pins a package — " +
						"blocks future CVE fixes. Use apt-mark hold and document.",
					Line:   cmd.Token.Line,
					Column: cmd.Token.Column,
					Level:  SeverityWarning,
				}}
			}
		}
	}
	return nil
}
