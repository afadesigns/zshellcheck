package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1919",
		Title:    "Warn on `ss -K` / `ss --kill` — terminates every socket that matches the filter",
		Severity: SeverityWarning,
		Description: "`ss -K` issues `SOCK_DESTROY` to every socket matching the filter (requires " +
			"`CAP_NET_ADMIN`). With a broad filter — `ss -K state established`, `ss -K dport 22` " +
			"— the command happily terminates the SSH session that is running it, along with " +
			"every backend keep-alive that happens to match. Spell the filter tightly " +
			"(`ss -K dst 10.0.0.5 dport 5432 state close-wait`), test it first without `-K` " +
			"to confirm only the target sockets appear, and wrap the call in a review step " +
			"rather than a scheduled job.",
		Check: checkZC1919,
	})
}

func checkZC1919(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	// Parser caveat: `ss --kill <filter>` mangles the command name to `kill`.
	if ident.Value == "kill" {
		for _, arg := range cmd.Arguments {
			if arg.String() == "state" || arg.String() == "dst" || arg.String() == "src" ||
				arg.String() == "dport" || arg.String() == "sport" {
				return zc1919Hit(cmd)
			}
		}
	}
	if ident.Value != "ss" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-K" || v == "--kill" {
			return zc1919Hit(cmd)
		}
		if len(v) >= 2 && v[0] == '-' && v[1] != '-' {
			for i := 1; i < len(v); i++ {
				if v[i] == 'K' {
					return []Violation{{
						KataID: "ZC1919",
						Message: "`ss -K` terminates every socket the filter matches — " +
							"broad filters (`state established`, `dport 22`) kill the running SSH " +
							"session. Preview with the same filter minus `-K`, and pin to a specific " +
							"dst/port/state tuple.",
						Line:   cmd.Token.Line,
						Column: cmd.Token.Column,
						Level:  SeverityWarning,
					}}
				}
			}
		}
	}
	return nil
}

func zc1919Hit(cmd *ast.SimpleCommand) []Violation {
	return []Violation{{
		KataID: "ZC1919",
		Message: "`ss -K` terminates every socket the filter matches — broad filters " +
			"(`state established`, `dport 22`) kill the running SSH session. Preview " +
			"with the same filter minus `-K`, and pin to a specific dst/port/state tuple.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}
