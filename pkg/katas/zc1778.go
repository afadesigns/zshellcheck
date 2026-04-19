package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1778",
		Title:    "Warn on `systemctl link /path/to/unit` — persistence from a mutable source path",
		Severity: SeverityWarning,
		Description: "`systemctl link` symlinks the given unit file into `/etc/systemd/system/` " +
			"so it can be `enable`d and `start`ed by name, but the unit definition lives at " +
			"the original path forever. If that path is writable by any non-root user " +
			"(`/tmp/*`, `/var/tmp/*`, `/home/*`, `/opt/` with wide perms, a build output " +
			"directory), a later tamper of the source file silently changes what systemd " +
			"runs the next time the unit starts. Copy the unit into `/etc/systemd/system/` " +
			"with root-only permissions, or install a package that ships it under " +
			"`/lib/systemd/system/`, rather than linking from a mutable location.",
		Check: checkZC1778,
	})
}

var zc1778MutablePrefixes = []string{
	"/tmp/",
	"/var/tmp/",
	"/home/",
	"/root/",
	"/opt/",
	"/srv/",
	"/mnt/",
	"/media/",
	"/var/lib/",
}

func checkZC1778(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "systemctl" {
		return nil
	}
	if len(cmd.Arguments) < 2 {
		return nil
	}

	linkIdx := -1
	for i, arg := range cmd.Arguments {
		if arg.String() == "link" {
			linkIdx = i
			break
		}
	}
	if linkIdx == -1 {
		return nil
	}
	for _, arg := range cmd.Arguments[linkIdx+1:] {
		v := arg.String()
		if strings.HasPrefix(v, "-") {
			continue
		}
		if !strings.HasPrefix(v, "/") {
			continue
		}
		for _, prefix := range zc1778MutablePrefixes {
			if strings.HasPrefix(v, prefix) {
				return []Violation{{
					KataID: "ZC1778",
					Message: "`systemctl link " + v + "` keeps the unit in a mutable " +
						"path — a tamper later changes what systemd runs. Copy the " +
						"unit into `/etc/systemd/system/` with root-only perms.",
					Line:   cmd.Token.Line,
					Column: cmd.Token.Column,
					Level:  SeverityWarning,
				}}
			}
		}
	}
	return nil
}
