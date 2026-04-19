package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1742",
		Title:    "Error on `mc alias set NAME URL ACCESS_KEY SECRET_KEY` — S3 keys in process list",
		Severity: SeverityError,
		Description: "MinIO's `mc alias set NAME URL ACCESS_KEY SECRET_KEY` (also `mc config " +
			"host add ALIAS URL ACCESS SECRET` on legacy versions) accepts the S3 access " +
			"and secret keys as positional arguments. Both land in argv — visible in " +
			"`ps`, `/proc/<pid>/cmdline`, shell history, and CI logs. Drop the trailing " +
			"keys and let `mc alias set NAME URL` prompt for them, or use the `MC_HOST_<" +
			"alias>=https://ACCESS:SECRET@host` env-var form scoped to a single command " +
			"and unset immediately after.",
		Check: checkZC1742,
	})
}

func checkZC1742(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "mc" && ident.Value != "mcli" {
		return nil
	}
	if len(cmd.Arguments) < 4 {
		return nil
	}

	sub0 := cmd.Arguments[0].String()
	sub1 := cmd.Arguments[1].String()
	rest := cmd.Arguments[2:]

	if sub0 == "alias" && sub1 == "set" {
		// Need NAME URL ACCESS_KEY SECRET_KEY (4 positionals after `alias set`).
		if zc1742PositionalCount(rest) >= 4 {
			return zc1742Hit(cmd, "mc alias set ... ACCESS_KEY SECRET_KEY")
		}
	}
	if sub0 == "config" && sub1 == "host" && len(cmd.Arguments) >= 5 {
		// Legacy: `mc config host add ALIAS URL ACCESS SECRET` (5 args).
		if cmd.Arguments[2].String() == "add" && zc1742PositionalCount(cmd.Arguments[3:]) >= 4 {
			return zc1742Hit(cmd, "mc config host add ... ACCESS SECRET")
		}
	}
	return nil
}

func zc1742PositionalCount(args []ast.Expression) int {
	count := 0
	for _, arg := range args {
		v := arg.String()
		if v == "" || v[0] == '-' {
			continue
		}
		count++
	}
	return count
}

func zc1742Hit(cmd *ast.SimpleCommand, what string) []Violation {
	return []Violation{{
		KataID: "ZC1742",
		Message: "`" + what + "` puts S3 access and secret keys in argv — visible in " +
			"`ps`, `/proc`, history. Drop the trailing keys (mc prompts) or use " +
			"`MC_HOST_<alias>=URL` env-var form scoped to one command.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}
