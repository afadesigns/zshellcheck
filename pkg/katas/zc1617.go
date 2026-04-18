package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1617",
		Title:    "Warn on `xargs -P 0` — unbounded parallelism risks CPU / fd / memory exhaustion",
		Severity: SeverityWarning,
		Description: "`xargs -P 0` tells xargs to spawn as many concurrent children as input " +
			"lines. On any non-trivial input that number can blow past `RLIMIT_NPROC`, " +
			"saturate the downstream tool's file-descriptor limit, or drive the host OOM. " +
			"Pick an explicit cap — `xargs -P $(nproc)` for CPU-bound work, `-P 4..8` for " +
			"I/O-bound — so the failure mode is bounded and predictable.",
		Check: checkZC1617,
	})
}

func checkZC1617(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "xargs" {
		return nil
	}

	for i, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-P0" {
			return zc1617Hit(cmd)
		}
		if v == "-P" && i+1 < len(cmd.Arguments) && cmd.Arguments[i+1].String() == "0" {
			return zc1617Hit(cmd)
		}
	}
	return nil
}

func zc1617Hit(cmd *ast.SimpleCommand) []Violation {
	return []Violation{{
		KataID: "ZC1617",
		Message: "`xargs -P 0` spawns one child per input line — CPU / FD / memory " +
			"exhaustion risk. Use `-P $(nproc)` or an explicit cap.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}
