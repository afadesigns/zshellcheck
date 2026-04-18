package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1664",
		Title:    "Error on `systemctl set-default rescue.target|emergency.target` — persistent single-user boot",
		Severity: SeverityError,
		Description: "`systemctl set-default` rewrites `/etc/systemd/system/default.target` as a " +
			"symlink to the named target. Pointing it at `rescue.target` or " +
			"`emergency.target` means every subsequent boot drops to single-user mode " +
			"before networking, sshd, or any normal unit starts — you lose remote access to " +
			"the box unless you have serial console / out-of-band management. Unlike " +
			"`systemctl isolate` (one-shot, caught by ZC1561) this persists across reboots. " +
			"Revert with `systemctl set-default multi-user.target` (servers) or `graphical." +
			"target` (desktops).",
		Check: checkZC1664,
	})
}

func checkZC1664(node ast.Node) []Violation {
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
	if cmd.Arguments[0].String() != "set-default" {
		return nil
	}
	target := cmd.Arguments[1].String()
	if target != "rescue.target" && target != "emergency.target" {
		return nil
	}

	return []Violation{{
		KataID: "ZC1664",
		Message: "`systemctl set-default " + target + "` makes every subsequent boot land " +
			"in single-user mode — revert with `set-default multi-user.target` or " +
			"`graphical.target`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}
