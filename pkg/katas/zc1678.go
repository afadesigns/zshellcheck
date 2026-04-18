package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1678",
		Title:    "Error on `borg init --encryption=none` — unencrypted backup repository",
		Severity: SeverityError,
		Description: "`borg init --encryption=none REPO` creates a backup repository without " +
			"client-side encryption or authentication. Anyone with read access to the repo " +
			"gets every file in every archive, and no one can detect silent tampering — " +
			"borg will happily extract a modified chunk. Even for local-only repos the cost " +
			"of authenticated-encryption is tiny; use `--encryption=repokey-blake2` (or " +
			"`--encryption=keyfile-blake2` when you want the key off the server), and store " +
			"the passphrase in `BORG_PASSPHRASE_FILE` pointing at a mode-0400 file.",
		Check: checkZC1678,
	})
}

func checkZC1678(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "borg" {
		return nil
	}
	if len(cmd.Arguments) == 0 || cmd.Arguments[0].String() != "init" {
		return nil
	}

	for i, arg := range cmd.Arguments[1:] {
		v := arg.String()
		if v == "--encryption=none" {
			return zc1678Hit(cmd)
		}
		if (v == "--encryption" || v == "-e") && i+2 < len(cmd.Arguments) &&
			cmd.Arguments[i+2].String() == "none" {
			return zc1678Hit(cmd)
		}
	}
	return nil
}

func zc1678Hit(cmd *ast.SimpleCommand) []Violation {
	return []Violation{{
		KataID: "ZC1678",
		Message: "`borg init --encryption=none` leaves archives unauthenticated and " +
			"readable — use `--encryption=repokey-blake2` (or `keyfile-blake2`) and store " +
			"the passphrase in `BORG_PASSPHRASE_FILE`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}
