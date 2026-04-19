package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1926",
		Title:    "Warn on `telinit 0/1/6` / `init 0/1/6` — SysV runlevel change halts, reboots, or isolates the host",
		Severity: SeverityWarning,
		Description: "`init 0`, `init 6`, `init 1`, and their `telinit` aliases ask systemd (or " +
			"SysV) to switch runlevel: `0` → `poweroff.target`, `6` → `reboot.target`, " +
			"`1`/`S` → `rescue.target`. From a script the side effect is a remote SSH " +
			"disconnect, an immediate service teardown for every other session on the host, " +
			"and — in the `1`/`S` case — dropping to single-user mode without a console to " +
			"recover. Use `systemctl poweroff`/`reboot`/`rescue` (which are clearer in " +
			"reviews) or schedule via `shutdown -h +N` so the operator has a cancel window.",
		Check: checkZC1926,
	})
}

func checkZC1926(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "init" && ident.Value != "telinit" {
		return nil
	}
	if len(cmd.Arguments) == 0 {
		return nil
	}
	lvl := cmd.Arguments[0].String()
	switch lvl {
	case "0", "1", "6", "S", "s":
		return []Violation{{
			KataID: "ZC1926",
			Message: "`" + ident.Value + " " + lvl + "` changes runlevel — `0` halts, `6` " +
				"reboots, `1`/`S` drops to single-user. Use `systemctl poweroff`/`reboot`/" +
				"`rescue` or `shutdown -h +N` so reviewers can read the intent.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityWarning,
		}}
	}
	return nil
}
