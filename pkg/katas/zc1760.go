package katas

import (
	"strconv"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1760",
		Title:    "Warn on `openssl rand -hex|-base64 N` with N < 16 — generated value too short",
		Severity: SeverityWarning,
		Description: "`openssl rand -hex N` (and `-base64 N`) outputs N random bytes encoded into " +
			"the requested form. N below 16 (128 bits) produces a value short enough that " +
			"an attacker with modest GPU resources can brute-force it offline — too weak " +
			"for passwords, API tokens, reset URLs, or any other secret that sits at rest. " +
			"Use `-hex 32` (256-bit) for secrets and long-lived tokens; `-hex 16` is " +
			"acceptable only for short-validity nonces paired with rate-limited consumers.",
		Check: checkZC1760,
	})
}

func checkZC1760(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "openssl" {
		return nil
	}
	if len(cmd.Arguments) < 2 {
		return nil
	}
	if cmd.Arguments[0].String() != "rand" {
		return nil
	}

	prevEnc := ""
	for _, arg := range cmd.Arguments[1:] {
		v := arg.String()
		if prevEnc != "" {
			if n, err := strconv.Atoi(v); err == nil && n > 0 && n < 16 {
				return zc1760Hit(cmd, prevEnc+" "+v)
			}
			prevEnc = ""
			continue
		}
		if v == "-hex" || v == "-base64" {
			prevEnc = v
		}
	}
	return nil
}

func zc1760Hit(cmd *ast.SimpleCommand, what string) []Violation {
	return []Violation{{
		KataID: "ZC1760",
		Message: "`openssl rand " + what + "` produces a sub-128-bit value — brute-forceable " +
			"offline. Use `-hex 32` for secrets / long-lived tokens.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}
