package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1963",
		Title:    "Warn on `npx pkg` / `pnpm dlx pkg` / `bunx pkg` without a version pin — runs latest registry code",
		Severity: SeverityWarning,
		Description: "`npx PKG`, `pnpm dlx PKG`, `bunx PKG`, and `bun x PKG` fetch the named " +
			"package from the npm registry and execute its `bin` entry. Without a version " +
			"pin (`pkg@1.2.3`), each run resolves to the registry's `latest` tag — a " +
			"compromised maintainer, squatted name, or even a mistyped package is enough to " +
			"land attacker code in the build. Pin the exact version (`npx pkg@1.2.3`), cache " +
			"the binary under `./node_modules/.bin/` via a regular `npm install`, or verify " +
			"the tarball signature before execution.",
		Check: checkZC1963,
	})
}

func checkZC1963(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	var form string
	var pkgs []ast.Expression
	switch ident.Value {
	case "npx":
		form = "npx"
		pkgs = cmd.Arguments
	case "bunx":
		form = "bunx"
		pkgs = cmd.Arguments
	case "pnpm", "bun":
		if len(cmd.Arguments) < 2 {
			return nil
		}
		sub := cmd.Arguments[0].String()
		if ident.Value == "pnpm" && sub != "dlx" {
			return nil
		}
		if ident.Value == "bun" && sub != "x" {
			return nil
		}
		form = ident.Value + " " + sub
		pkgs = cmd.Arguments[1:]
	default:
		return nil
	}

	for _, arg := range pkgs {
		v := arg.String()
		if strings.HasPrefix(v, "-") {
			continue
		}
		if strings.Contains(v, "@") && !strings.HasPrefix(v, "@") {
			// pkg@version — pinned.
			return nil
		}
		if strings.HasPrefix(v, "@") {
			// scoped name like @scope/pkg — check for second @ (version).
			rest := v[1:]
			if strings.Contains(rest, "@") {
				return nil
			}
		}
		if strings.HasPrefix(v, "$") {
			return nil
		}
		return []Violation{{
			KataID: "ZC1963",
			Message: "`" + form + " " + v + "` pulls the `latest` tag every run — " +
				"a squatted or compromised package lands attacker code. Pin the version " +
				"(`" + v + "@X.Y.Z`) or use a regular `npm install` + `./node_modules/.bin/`.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityWarning,
		}}
	}
	return nil
}
