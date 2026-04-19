package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1854",
		Title:    "Error on `yum-config-manager --add-repo http://…` / `zypper addrepo http://…` — plaintext repo allows MITM",
		Severity: SeverityError,
		Description: "Adding a package repository over plain HTTP (`yum-config-manager " +
			"--add-repo http://…`, `dnf config-manager --add-repo http://…`, `zypper " +
			"addrepo http://…`) tells the package manager to fetch metadata and RPMs " +
			"without TLS — any on-path attacker can substitute packages, and even GPG " +
			"signature checks do not help because the attacker can simply strip the " +
			"`repo_gpgcheck=1` line from the unsigned `.repo` file. Use the `https://` " +
			"mirror (every major distro now publishes one), or pin to a local mirror over " +
			"TLS and drop a `gpgkey=file:///etc/pki/...` entry in the same `.repo` so " +
			"signatures cannot be disabled mid-install.",
		Check: checkZC1854,
	})
}

func checkZC1854(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	switch ident.Value {
	case "yum-config-manager":
		for _, arg := range cmd.Arguments {
			if zc1854IsHTTPURL(arg.String()) {
				return zc1854Hit(cmd, "yum-config-manager --add-repo "+arg.String())
			}
		}
	case "add-repo":
		// Parser caveat: `yum-config-manager --add-repo URL` mangles the
		// command name to `add-repo`.
		for _, arg := range cmd.Arguments {
			if zc1854IsHTTPURL(arg.String()) {
				return zc1854Hit(cmd, "yum-config-manager --add-repo "+arg.String())
			}
		}
	case "dnf":
		if len(cmd.Arguments) >= 3 &&
			cmd.Arguments[0].String() == "config-manager" &&
			cmd.Arguments[1].String() == "--add-repo" {
			for _, arg := range cmd.Arguments[2:] {
				if zc1854IsHTTPURL(arg.String()) {
					return zc1854Hit(cmd, "dnf config-manager --add-repo "+arg.String())
				}
			}
		}
	case "zypper":
		if len(cmd.Arguments) >= 2 {
			sub := cmd.Arguments[0].String()
			if sub == "addrepo" || sub == "ar" {
				for _, arg := range cmd.Arguments[1:] {
					if zc1854IsHTTPURL(arg.String()) {
						return zc1854Hit(cmd, "zypper addrepo "+arg.String())
					}
				}
			}
		}
	}
	return nil
}

func zc1854IsHTTPURL(v string) bool {
	return strings.HasPrefix(v, "http://")
}

func zc1854Hit(cmd *ast.SimpleCommand, where string) []Violation {
	return []Violation{{
		KataID: "ZC1854",
		Message: "`" + where + "` registers a plaintext repo — on-path attacker can " +
			"substitute packages and strip GPG-check directives. Use `https://` and " +
			"pin `gpgkey=file://` in the `.repo`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}
