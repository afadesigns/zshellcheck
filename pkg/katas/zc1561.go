package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1561",
		Title:    "Error on `systemctl isolate rescue.target` / `emergency.target` from a script",
		Severity: SeverityError,
		Description: "`systemctl isolate rescue.target` drops the host into single-user rescue " +
			"mode; `emergency.target` goes even further, leaving only the root shell on the " +
			"console. Both terminate networking, SSH sessions, and most services. On a remote " +
			"host the script loses its own connection mid-run, and anyone relying on the box " +
			"is cut off without warning. Reserve these for console recovery, not script flow.",
		Check: checkZC1561,
	})
}

func checkZC1561(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "systemctl" {
		return nil
	}

	if len(cmd.Arguments) < 2 {
		return nil
	}
	if cmd.Arguments[0].String() != "isolate" {
		return nil
	}
	target := cmd.Arguments[1].String()
	switch target {
	case "rescue.target", "emergency.target", "rescue", "emergency":
		return []Violation{{
			KataID: "ZC1561",
			Message: "`systemctl isolate " + target + "` kills SSH and most services — " +
				"console-only recovery. Do not run from a script.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityError,
		}}
	}
	return nil
}
