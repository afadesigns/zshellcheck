package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1900",
		Title:    "Warn on `curl --location-trusted` — Authorization/cookies forwarded across redirects",
		Severity: SeverityWarning,
		Description: "`curl --location-trusted` (alias of `curl -L --location-trusted`) tells " +
			"curl to replay the `Authorization` header, cookies, and `-u user:pass` credential " +
			"on every redirect hop, even across hosts. A 302 to an attacker-controlled origin " +
			"(or a compromised CDN edge) then receives the bearer token verbatim. Drop " +
			"`--location-trusted`; if cross-origin auth is truly required, scope a short-lived " +
			"token per destination and verify the final hostname before sending secrets.",
		Check: checkZC1900,
	})
}

func checkZC1900(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	// Parser caveat: `curl --location-trusted …` mangles the command name
	// to `location-trusted`.
	switch ident.Value {
	case "location-trusted":
		return zc1900Hit(cmd)
	case "curl":
		for _, arg := range cmd.Arguments {
			if arg.String() == "--location-trusted" {
				return zc1900Hit(cmd)
			}
		}
	}
	return nil
}

func zc1900Hit(cmd *ast.SimpleCommand) []Violation {
	return []Violation{{
		KataID: "ZC1900",
		Message: "`curl --location-trusted` replays `Authorization`, cookies, and " +
			"`-u user:pass` on every redirect — a 302 to attacker-controlled host " +
			"leaks the token. Drop the flag; verify final hostname before sending secrets.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}
