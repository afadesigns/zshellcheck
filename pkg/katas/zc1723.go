package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

var zc1723DeleteFlags = map[string]bool{
	"--delete-secret-keys":            true,
	"--delete-secret-and-public-keys": true,
	"--delete-keys":                   true,
	"--delete-key":                    true,
}

// Parser caveat: a leading `--long-flag` mangles the GPG name so the next
// SimpleCommand becomes the bare flag without the leading dashes.
var zc1723MangledNames = map[string]bool{
	"delete-secret-keys":            true,
	"delete-secret-and-public-keys": true,
	"delete-keys":                   true,
	"delete-key":                    true,
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1723",
		Title:    "Error on `gpg --delete-secret-keys` / `--delete-key` — irreversible key destruction",
		Severity: SeverityError,
		Description: "GPG key deletion is permanent. Once `--delete-secret-keys`, " +
			"`--delete-secret-and-public-keys`, `--delete-keys`, or `--delete-key` removes " +
			"the keyring entry there is no recovery short of a separate backup or off-card " +
			"reimport. Combined with `--batch --yes`, the confirmation prompt is bypassed " +
			"and a single accidental KEYID resolves to a one-shot wipe. Export the key " +
			"first (`gpg --export-secret-keys --armor KEYID > backup.asc`, store offline) " +
			"and never pair the delete flag with `--batch --yes` in automation.",
		Check: checkZC1723,
	})
}

func checkZC1723(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	if zc1723MangledNames[ident.Value] {
		return zc1723Hit(cmd, "--"+ident.Value)
	}

	if ident.Value != "gpg" && ident.Value != "gpg2" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if zc1723DeleteFlags[v] {
			return zc1723Hit(cmd, v)
		}
	}
	return nil
}

func zc1723Hit(cmd *ast.SimpleCommand, flag string) []Violation {
	return []Violation{{
		KataID: "ZC1723",
		Message: "`gpg " + flag + "` permanently destroys keyring entries — no recovery " +
			"without a separate backup. Export with `gpg --export-secret-keys --armor " +
			"KEYID` first; never pair this flag with `--batch --yes`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}
