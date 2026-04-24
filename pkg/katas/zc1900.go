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

var zc1900LocationFlags = map[string]bool{
	"--location-trusted": true,
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
	if ident.Value != "curl" {
		return nil
	}
	for _, arg := range cmd.Arguments {
		if zc1900LocationFlags[arg.String()] {
			line, col := FlagArgPosition(cmd, zc1900LocationFlags)
			return []Violation{{
				KataID: "ZC1900",
				Message: "`curl --location-trusted` replays `Authorization`, cookies, and " +
					"`-u user:pass` on every redirect — a 302 to attacker-controlled host " +
					"leaks the token. Drop the flag; verify final hostname before sending secrets.",
				Line:   line,
				Column: col,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}
