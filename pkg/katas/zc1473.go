package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1473",
		Title:    "Warn on `openssl req ... -nodes` / `genrsa` without passphrase — unencrypted private key",
		Severity: SeverityWarning,
		Description: "`-nodes` tells OpenSSL not to encrypt the private key that is written to " +
			"disk. The file ends up at whatever filesystem permissions the umask dictates, and " +
			"any subsequent backup / container image / rsync picks up a usable key with no " +
			"passphrase. Use `-aes256` / `-aes-256-cbc` and keep the passphrase in a secrets " +
			"store, or rely on a hardware-backed key via PKCS#11 / TPM.",
		Check: checkZC1473,
	})
}

func checkZC1473(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "openssl" {
		return nil
	}

	if len(cmd.Arguments) == 0 {
		return nil
	}
	sub := cmd.Arguments[0].String()
	// Only flag on subcommands that actually produce a private key file.
	if sub != "req" && sub != "genrsa" && sub != "genpkey" && sub != "ecparam" &&
		sub != "dsaparam" && sub != "pkcs12" {
		return nil
	}

	for _, arg := range cmd.Arguments[1:] {
		v := arg.String()
		if v == "-nodes" || v == "-noenc" {
			return []Violation{{
				KataID: "ZC1473",
				Message: "`" + v + "` writes the private key to disk unencrypted. Use `-aes256` " +
					"(or an HSM/TPM) and keep the passphrase in a secrets store.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}
