package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1667",
		Title:    "Warn on `openssl enc` without `-pbkdf2` — legacy MD5-based key derivation",
		Severity: SeverityWarning,
		Description: "Without `-pbkdf2`, `openssl enc` derives the symmetric key through " +
			"EVP_BytesToKey, which is a single MD5 round over `password || salt`. A modern " +
			"GPU cracks that at billions of guesses per second. Add `-pbkdf2 -iter 100000` " +
			"(OpenSSL 1.1.1+) to switch to PBKDF2-HMAC-SHA256 with a real iteration count. " +
			"Even better, stop using `openssl enc` for new code — it has no AEAD support and " +
			"`-aes-256-gcm` silently drops the auth tag — and reach for `age`, " +
			"`gpg --symmetric --cipher-algo AES256`, or `openssl smime` instead.",
		Check: checkZC1667,
	})
}

func checkZC1667(node ast.Node) []Violation {
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

	if len(cmd.Arguments) == 0 || cmd.Arguments[0].String() != "enc" {
		return nil
	}

	for _, arg := range cmd.Arguments[1:] {
		if arg.String() == "-pbkdf2" {
			return nil
		}
	}

	return []Violation{{
		KataID: "ZC1667",
		Message: "`openssl enc` without `-pbkdf2` uses single-round EVP_BytesToKey (MD5) — " +
			"add `-pbkdf2 -iter 100000`, or prefer `age` / `gpg --symmetric`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}
