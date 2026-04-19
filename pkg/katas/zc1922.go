package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1922",
		Title:    "Error on `rpm --import http://…` / `rpmkeys --import ftp://…` — plaintext GPG key fetch",
		Severity: SeverityError,
		Description: "`rpm --import` (and `rpmkeys --import`) add the supplied ASCII-armoured " +
			"key to the system RPM trust store. When the source is a plain `http://` / `ftp://` " +
			"URL an on-path attacker swaps the key, and every subsequent package they sign " +
			"installs cleanly. Serve keys over HTTPS from a TLS-authenticated origin, pin the " +
			"key's SHA-256 before import, or stage an offline copy verified out of band " +
			"(`gpg --verify` against a known-good fingerprint).",
		Check: checkZC1922,
	})
}

func checkZC1922(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	// Parser caveat: `rpm --import URL` / `rpmkeys --import URL` both mangle
	// the command name to `import`.
	if ident.Value == "import" && len(cmd.Arguments) >= 1 {
		url := cmd.Arguments[0].String()
		if zc1922IsPlaintextURL(url) {
			return zc1922Hit(cmd, url)
		}
		return nil
	}
	if ident.Value != "rpm" && ident.Value != "rpmkeys" {
		return nil
	}

	for i, arg := range cmd.Arguments {
		if arg.String() == "--import" && i+1 < len(cmd.Arguments) {
			url := cmd.Arguments[i+1].String()
			if zc1922IsPlaintextURL(url) {
				return zc1922Hit(cmd, url)
			}
		}
	}
	return nil
}

func zc1922IsPlaintextURL(s string) bool {
	return strings.HasPrefix(s, "http://") || strings.HasPrefix(s, "ftp://")
}

func zc1922Hit(cmd *ast.SimpleCommand, url string) []Violation {
	return []Violation{{
		KataID: "ZC1922",
		Message: "`rpm --import " + url + "` fetches a GPG key over plaintext — on-path " +
			"attackers swap it, every future signed package installs. Use `https://` from " +
			"a pinned origin, or `gpg --verify` against a known fingerprint.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}
