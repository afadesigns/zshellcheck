package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1652",
		Title:    "Warn on `ssh -Y` — trusted X11 forwarding grants full X-server access to remote clients",
		Severity: SeverityWarning,
		Description: "`ssh -Y` enables trusted X11 forwarding. Remote X clients can read every " +
			"keystroke on the local display, take screenshots, inject synthetic events, and " +
			"otherwise drive the local session with no sandbox. `ssh -X` enables the " +
			"untrusted variant, which routes X traffic through the X SECURITY extension so " +
			"those capabilities are limited (some GUI features break, which is why people " +
			"reach for `-Y` — usually at far higher risk than they realised). Prefer `-X` " +
			"when X11 forwarding is genuinely needed; better yet drop it for Wayland tools " +
			"or VNC-over-SSH with its own auth.",
		Check: checkZC1652,
	})
}

func checkZC1652(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "ssh" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		if arg.String() == "-Y" {
			return []Violation{{
				KataID: "ZC1652",
				Message: "`ssh -Y` enables trusted X11 forwarding — remote clients get full " +
					"access to the local X server. Use `-X` (untrusted) or drop X11 " +
					"forwarding entirely.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}
