package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1713",
		Title:    "Error on `consul kv delete -recurse /` — wipes the entire Consul KV store",
		Severity: SeverityError,
		Description: "`consul kv delete -recurse PREFIX` removes every key under PREFIX. With " +
			"PREFIX `/` (or an empty string) the command nukes the whole KV store, " +
			"including service-discovery payloads, ACL bootstrap tokens, and any " +
			"application-level config the cluster relies on. Scope the prefix to the app " +
			"namespace (`consul kv delete -recurse /app/staging/`), confirm the keys you " +
			"are about to lose with `consul kv get -recurse -keys`, and snapshot the " +
			"datacenter (`consul snapshot save snap.bin`) before any large delete.",
		Check: checkZC1713,
	})
}

func checkZC1713(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "consul" {
		return nil
	}
	if len(cmd.Arguments) < 2 {
		return nil
	}
	if cmd.Arguments[0].String() != "kv" || cmd.Arguments[1].String() != "delete" {
		return nil
	}

	hasRecurse := false
	rootPrefix := false
	for _, arg := range cmd.Arguments[2:] {
		v := arg.String()
		switch v {
		case "-recurse", "--recurse":
			hasRecurse = true
		case "/", "", `""`, "''":
			rootPrefix = true
		}
	}
	if !hasRecurse || !rootPrefix {
		return nil
	}

	return []Violation{{
		KataID: "ZC1713",
		Message: "`consul kv delete -recurse /` removes the entire KV store — service " +
			"discovery, ACL bootstrap, app config. Scope the prefix and snapshot " +
			"(`consul snapshot save`) first.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}
