package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1701",
		Title:    "Info: `dpkg -i FILE.deb` installs without automatic signature verification",
		Severity: SeverityInfo,
		Description: "Unlike `apt install`, which verifies package signatures against the apt " +
			"repository's `Release.gpg`, plain `dpkg -i FILE.deb` applies the package with " +
			"no integrity check beyond Debian's own `.deb` format. In a provisioning " +
			"pipeline that downloaded the file over HTTPS from a vendor, that is usually " +
			"fine — the TLS chain vouches for the bytes. In scripts that pick the file up " +
			"from `/tmp`, `/var/tmp`, `/dev/shm`, or a mutable cache, a local user could " +
			"swap the file between download and install. Verify with `sha256sum -c`, " +
			"`debsig-verify`, or `dpkg-sig --verify` before invoking `dpkg -i`.",
		Check: checkZC1701,
	})
}

func checkZC1701(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "dpkg" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-i" {
			return []Violation{{
				KataID: "ZC1701",
				Message: "`dpkg -i FILE.deb` runs the package without signature verification — " +
					"`sha256sum -c` or `debsig-verify` the file first, or install via `apt " +
					"install` from a signed repo.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityInfo,
			}}
		}
	}
	return nil
}
