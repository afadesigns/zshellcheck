package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1589",
		Title:    "Warn on `trap 'set -x' ERR/RETURN/EXIT/ZERR` — trace hook leaks env to stderr",
		Severity: SeverityWarning,
		Description: "Installing a trap that enables `set -x` (or `set -o xtrace` / `set -v`) " +
			"causes every subsequent expanded command to hit stderr. Expansions embed " +
			"environment variables — API tokens, passwords, signed URLs — directly into " +
			"the trace. In CI, that stderr lands in build logs and gets shipped to long-term " +
			"log retention. Scope `set -x` to a `set -x ... set +x` block around the suspect " +
			"code, or replace the trap with `trap 'safe_dump' ERR` that prints only non-" +
			"sensitive state.",
		Check: checkZC1589,
	})
}

func checkZC1589(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "trap" {
		return nil
	}
	if len(cmd.Arguments) < 2 {
		return nil
	}

	action := strings.Trim(cmd.Arguments[0].String(), "'\"")
	if !strings.Contains(action, "set -x") &&
		!strings.Contains(action, "set -o xtrace") &&
		!strings.Contains(action, "set -v") &&
		!strings.Contains(action, "set -o verbose") {
		return nil
	}

	sig := cmd.Arguments[1].String()
	switch sig {
	case "ERR", "RETURN", "EXIT", "ZERR":
	default:
		return nil
	}

	return []Violation{{
		KataID: "ZC1589",
		Message: "`trap 'set -x' " + sig + "` enables shell trace from a trap — expansions " +
			"leak env vars (tokens, passwords) to stderr / CI logs. Use a scoped `set -x " +
			"... set +x`, not a trap.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}
