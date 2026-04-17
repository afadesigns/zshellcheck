package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1581",
		Title:    "Warn on `ssh -o PubkeyAuthentication=no` / `-o PasswordAuthentication=yes`",
		Severity: SeverityWarning,
		Description: "Forcing password authentication on a connection that has a working key " +
			"turns a strong (challenge-response, no password leaves the client) into a weak " +
			"(password-in-the-clear-on-disk-or-prompt) authentication path. Similarly " +
			"disabling pubkey skips the good path entirely. Leave the defaults, let the " +
			"server's `PubkeyAuthentication yes` pick the key, and document any exception.",
		Check: checkZC1581,
	})
}

func checkZC1581(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "ssh" && ident.Value != "scp" && ident.Value != "sftp" {
		return nil
	}

	var prevO bool
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if prevO {
			prevO = false
			s := strings.TrimSpace(strings.ToLower(v))
			if s == "pubkeyauthentication=no" || s == "passwordauthentication=yes" ||
				s == "preferredauthentications=password" {
				return []Violation{{
					KataID: "ZC1581",
					Message: "`" + ident.Value + " -o " + v + "` forces password auth — " +
						"weaker than key auth. Let the default preference pick.",
					Line:   cmd.Token.Line,
					Column: cmd.Token.Column,
					Level:  SeverityWarning,
				}}
			}
		}
		if v == "-o" {
			prevO = true
		}
	}
	return nil
}
