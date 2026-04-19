package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1860",
		Title:    "Warn on `hostnamectl set-hostname NEW` — caches and certs still reference the old name",
		Severity: SeverityWarning,
		Description: "`hostnamectl set-hostname NEW` (and the new-style `hostnamectl hostname NEW` " +
			"and `hostname NEW`) updates `/etc/hostname` and `kernel.hostname` atomically, " +
			"but every process that called `gethostname()` at startup keeps the old " +
			"value until it restarts: syslog tags, Prometheus scrape labels, Docker " +
			"daemons, and anything that populated a TLS `subjectAltName` with `$(hostname)` " +
			"still speak as the previous host. Change the hostname interactively, then " +
			"plan a restart window — in automation, prefer shipping the new hostname via " +
			"cloud-init / Ignition so every service starts with it from boot.",
		Check: checkZC1860,
	})
}

func checkZC1860(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value == "hostnamectl" {
		if len(cmd.Arguments) < 2 {
			return nil
		}
		sub := cmd.Arguments[0].String()
		if sub != "set-hostname" && sub != "hostname" {
			return nil
		}
		return zc1860Hit(cmd, "hostnamectl "+sub+" "+cmd.Arguments[1].String())
	}
	if ident.Value == "hostname" {
		if len(cmd.Arguments) != 1 {
			return nil
		}
		v := cmd.Arguments[0].String()
		// `hostname` with no args just prints; `hostname -f` is read-only.
		if len(v) == 0 || v[0] == '-' {
			return nil
		}
		return zc1860Hit(cmd, "hostname "+v)
	}
	return nil
}

func zc1860Hit(cmd *ast.SimpleCommand, where string) []Violation {
	return []Violation{{
		KataID: "ZC1860",
		Message: "`" + where + "` updates the kernel hostname live, but running " +
			"services keep the old `gethostname()` — syslog tags, Prometheus " +
			"labels, TLS SANs stay stale. Apply at provisioning or reboot.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}
