package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1444",
		Title:    "Dangerous: `redis-cli FLUSHALL` / `FLUSHDB` — wipes Redis data",
		Severity: SeverityError,
		Description: "`FLUSHALL` deletes every key in every database; `FLUSHDB` clears the current " +
			"DB. Running against production is usually catastrophic. Either rename the command " +
			"in `redis.conf` (`rename-command FLUSHALL \"\"`) or require an explicit " +
			"confirmation in scripts.",
		Check: checkZC1444,
	})
}

func checkZC1444(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "redis-cli" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := strings.ToUpper(arg.String())
		if v == "FLUSHALL" || v == "FLUSHDB" {
			return []Violation{{
				KataID: "ZC1444",
				Message: "`redis-cli FLUSHALL`/`FLUSHDB` wipes Redis data. Disable via " +
					"`rename-command` in redis.conf on production, or require explicit confirmation.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}

	return nil
}
