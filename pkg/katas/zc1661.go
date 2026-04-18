package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1661",
		Title:    "Error on `curl --cacert /dev/null` — empty trust store, any cert passes",
		Severity: SeverityError,
		Description: "Pointing `--cacert` (or `--capath`) at `/dev/null` hands curl an empty " +
			"trust anchor set. Counter-intuitively, curl treats the peer certificate as " +
			"valid when no issuers are configured for the selected TLS backend (OpenSSL, " +
			"wolfSSL, Schannel all accept any cert chain against an empty CA bundle). This is " +
			"the TLS equivalent of `--insecure` with one more keystroke of plausible " +
			"deniability. Use a real bundle (`/etc/ssl/certs/ca-certificates.crt`) or " +
			"`--pinnedpubkey sha256//…` for known endpoints.",
		Check: checkZC1661,
	})
}

func checkZC1661(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	switch ident.Value {
	case "cacert", "capath":
		if len(cmd.Arguments) > 0 && cmd.Arguments[0].String() == "/dev/null" {
			return zc1661Hit(cmd)
		}
		return nil
	case "curl":
	default:
		return nil
	}

	for i, arg := range cmd.Arguments {
		v := arg.String()
		if v != "--cacert" && v != "--capath" {
			continue
		}
		if i+1 >= len(cmd.Arguments) {
			continue
		}
		if cmd.Arguments[i+1].String() == "/dev/null" {
			return zc1661Hit(cmd)
		}
	}
	return nil
}

func zc1661Hit(cmd *ast.SimpleCommand) []Violation {
	return []Violation{{
		KataID: "ZC1661",
		Message: "`curl --cacert /dev/null` feeds curl an empty trust store — most TLS " +
			"backends then accept any peer cert. Use a real bundle or `--pinnedpubkey`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}
