package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1770",
		Title:    "Warn on `gpg --always-trust` / `--trust-model always` — bypasses Web-of-Trust",
		Severity: SeverityWarning,
		Description: "`gpg --always-trust` (equivalent to `--trust-model always`) accepts every key " +
			"in the keyring as fully trusted, regardless of signatures from the owner or any " +
			"introducer. A signature made by an attacker-controlled key pair that was imported " +
			"with no further vetting will verify cleanly. In automation this turns signature " +
			"verification into a presence check — any key bundled with the payload satisfies " +
			"`gpg --verify`. Remove the flag and build a proper trust path: either mark the " +
			"expected signer key trusted once (`gpg --edit-key KEYID trust`), or pin the " +
			"expected fingerprint and match it against the signer after `gpg --verify --status-fd 1`.",
		Check: checkZC1770,
	})
}

func checkZC1770(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	if ident.Value == "always-trust" {
		return zc1770Hit(cmd, "--always-trust")
	}
	if ident.Value == "trust-model" {
		if len(cmd.Arguments) > 0 && cmd.Arguments[0].String() == "always" {
			return zc1770Hit(cmd, "--trust-model always")
		}
		return nil
	}

	if ident.Value != "gpg" && ident.Value != "gpg2" {
		return nil
	}

	for i, arg := range cmd.Arguments {
		v := arg.String()
		if v == "--always-trust" {
			return zc1770Hit(cmd, "--always-trust")
		}
		if v == "--trust-model" && i+1 < len(cmd.Arguments) && cmd.Arguments[i+1].String() == "always" {
			return zc1770Hit(cmd, "--trust-model always")
		}
		if v == "--trust-model=always" {
			return zc1770Hit(cmd, "--trust-model=always")
		}
	}
	return nil
}

func zc1770Hit(cmd *ast.SimpleCommand, flag string) []Violation {
	return []Violation{{
		KataID: "ZC1770",
		Message: "`gpg " + flag + "` marks every imported key as fully trusted — a " +
			"signature from an attacker-supplied key verifies cleanly. Drop the flag " +
			"and pin the expected fingerprint, or assign trust via `gpg --edit-key KEYID trust`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}
