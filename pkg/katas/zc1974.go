package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1974",
		Title:    "Error on `ipset flush` / `ipset destroy` — nukes named sets referenced by iptables/nft rules",
		Severity: SeverityError,
		Description: "`ipset flush` empties every entry from a named IP set; `ipset destroy` " +
			"(no args) removes every set on the host. iptables/nft rules of the form " +
			"`-m set --match-set $NAME src` then reference a set that is either empty " +
			"or gone, so block-lists disappear instantly and allow-lists stop " +
			"whitelisting — the ruleset falls through to its default policy. Target a " +
			"specific set by name (`ipset destroy $NAME` after confirming no rule " +
			"references it), or add new entries with `ipset add` instead of rebuilding " +
			"from scratch. Reload atomically with `ipset restore -! < snapshot` if a " +
			"full replace is genuinely needed.",
		Check: checkZC1974,
	})
}

func checkZC1974(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "ipset" {
		return nil
	}
	if len(cmd.Arguments) == 0 {
		return nil
	}
	sub := cmd.Arguments[0].String()
	switch sub {
	case "flush", "-F":
		return zc1974Hit(cmd, "ipset flush")
	case "destroy", "-X":
		if len(cmd.Arguments) == 1 {
			return zc1974Hit(cmd, "ipset destroy")
		}
	}
	return nil
}

func zc1974Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1974",
		Message: "`" + form + "` drops named IP sets wholesale — iptables/nft rules " +
			"that reference them fall through to the default policy (block-list " +
			"empty, allow-list gone). Target by name; reload atomically via " +
			"`ipset restore -! < snapshot`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}
