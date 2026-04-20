package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1992",
		Title:    "Warn on `pkexec cmd` — PolicyKit privilege elevation is historically bug-prone and hard to audit from scripts",
		Severity: SeverityWarning,
		Description: "`pkexec` lifts a command to the UID configured in a PolicyKit `.policy` " +
			"file — typically root — after consulting an authorisation agent. From a " +
			"non-interactive script the agent has no way to prompt, so the call " +
			"either depends on a pre-authorised `.policy` override or fails in a " +
			"confusing manner. The binary also has a poor CVE track record (CVE-2021-" +
			"4034 pwnkit, CVE-2017-16089, envvar handling bugs) and its audit trail is " +
			"split across journald and `/var/log/auth.log`. Use `sudo` with a targeted " +
			"`sudoers` drop-in for scripted privilege elevation, or run the script " +
			"under a systemd unit with `User=` / `AmbientCapabilities=` when specific " +
			"capabilities are needed.",
		Check: checkZC1992,
	})
}

func checkZC1992(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "pkexec" {
		return nil
	}
	if len(cmd.Arguments) == 0 {
		return nil
	}
	return []Violation{{
		KataID: "ZC1992",
		Message: "`pkexec` elevates via PolicyKit — no agent to prompt in a script, " +
			"poor CVE history (pwnkit), split audit trail. Use `sudo` with a " +
			"targeted `sudoers.d` drop-in or a systemd unit with " +
			"`User=`/`AmbientCapabilities=`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}
