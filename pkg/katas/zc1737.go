package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1737",
		Title:    "Error on `wpa_passphrase SSID PASSWORD` — Wi-Fi passphrase in process list",
		Severity: SeverityError,
		Description: "`wpa_passphrase SSID PASSPHRASE` generates `wpa_supplicant.conf` content " +
			"on stdout. Putting PASSPHRASE on the command line lands it in `ps`, `/proc/<" +
			"pid>/cmdline`, shell history, and the audit log of every local user that can " +
			"list processes. Drop the second positional argument and let `wpa_passphrase " +
			"SSID < /run/secrets/wifi` (or piped via stdin from a secrets store) read the " +
			"passphrase from a file descriptor instead.",
		Check: checkZC1737,
	})
}

func checkZC1737(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "wpa_passphrase" {
		return nil
	}

	positionals := 0
	for _, arg := range cmd.Arguments {
		v := arg.String()
		// Skip redirection markers — they aren't true positional args.
		if v == "<" || v == ">" || v == ">>" || v == "<<" {
			break
		}
		if v == "" {
			continue
		}
		positionals++
	}

	if positionals < 2 {
		return nil
	}

	return []Violation{{
		KataID: "ZC1737",
		Message: "`wpa_passphrase SSID PASSWORD` puts the Wi-Fi passphrase in argv — " +
			"visible in `ps`, `/proc`, history. Drop the PASSWORD argument and pipe it " +
			"via stdin (`wpa_passphrase SSID < /run/secrets/wifi`).",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}
