package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1571",
		Title:    "Style: `ntpdate` is deprecated — use `chronyc makestep` / `systemd-timesyncd`",
		Severity: SeverityStyle,
		Description: "`ntpdate` was retired by the ntp.org project around 4.2.6. Distros " +
			"increasingly ship without it, and packaging it breaks the invariant that only " +
			"one program writes the clock at a time (if `chrony` or `timesyncd` is also " +
			"running the two fight). Use `chronyc makestep` (if chrony is active) or " +
			"`systemctl restart systemd-timesyncd` (if timesyncd is active) for a one-shot " +
			"step, and leave the daemon to keep it synchronised.",
		Check: checkZC1571,
	})
}

func checkZC1571(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "ntpdate" && ident.Value != "sntp" {
		return nil
	}

	return []Violation{{
		KataID: "ZC1571",
		Message: "`" + ident.Value + "` is deprecated and races any running chrony/timesyncd. " +
			"Use `chronyc makestep` or `systemctl restart systemd-timesyncd`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityStyle,
	}}
}
