package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1898",
		Title:    "Error on `gpg --export-secret-keys` — private-key material leaks to stdout",
		Severity: SeverityError,
		Description: "`gpg --export-secret-keys KEYID` and `--export-secret-subkeys` write the " +
			"ASCII-armoured private key to stdout. In a script, that stream usually lands " +
			"in a file the operator plans to move off-box — and any misstep (wrong " +
			"`cd`, script-wide stdout captured by CI, tee to a world-readable log, " +
			"piped into a remote unencrypted channel) permanently leaks the key. Backup " +
			"the key interactively on an air-gapped machine; if automation is required, " +
			"write the output to a `umask 077`-protected path and immediately encrypt " +
			"with a second symmetric passphrase.",
		Check: checkZC1898,
	})
}

func checkZC1898(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	// Parser caveat: `gpg --export-secret-keys …` mangles the command name
	// to `export-secret-keys` (or `-subkeys`).
	switch ident.Value {
	case "export-secret-keys", "export-secret-subkeys":
		return zc1898Hit(cmd, "--"+ident.Value)
	case "gpg", "gpg2":
		for _, arg := range cmd.Arguments {
			v := arg.String()
			if v == "--export-secret-keys" || v == "--export-secret-subkeys" {
				return zc1898Hit(cmd, v)
			}
		}
	}
	return nil
}

func zc1898Hit(cmd *ast.SimpleCommand, flag string) []Violation {
	return []Violation{{
		KataID: "ZC1898",
		Message: "`gpg " + flag + "` writes the private key to stdout — one " +
			"CI-log or wrong-tty redirect leaks it. Back up interactively on an " +
			"air-gapped host, or write to a `umask 077` path and re-encrypt.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}
