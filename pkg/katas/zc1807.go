package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1807",
		Title:    "Warn on `gh api -X DELETE` — raw GitHub DELETE bypasses `gh` command confirmations",
		Severity: SeverityWarning,
		Description: "`gh api -X DELETE /repos/OWNER/REPO` (and `--method=DELETE` variants) " +
			"sends a raw GitHub API request with the caller's token. There is no confirmation " +
			"prompt, no `--yes` guard, and no friendly dry-run — a script that builds the " +
			"path from a variable can wipe repos, releases, deploy keys, workflow runs, " +
			"issue comments, or whole organisations in one call. Use the high-level `gh` " +
			"subcommand for the target (`gh repo delete`, `gh release delete`, `gh workflow " +
			"disable`) which still at least requires `--yes`, or wrap the raw call with a " +
			"preflight `gh api -X GET /path` and an explicit confirmation in the script.",
		Check: checkZC1807,
	})
}

func checkZC1807(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "gh" {
		return nil
	}
	if len(cmd.Arguments) == 0 || cmd.Arguments[0].String() != "api" {
		return nil
	}

	for i, arg := range cmd.Arguments[1:] {
		v := arg.String()
		switch {
		case v == "-X" || v == "--method":
			if 1+i+1 < len(cmd.Arguments) {
				next := cmd.Arguments[1+i+1].String()
				if strings.EqualFold(next, "DELETE") {
					return zc1807Hit(cmd)
				}
			}
		case strings.EqualFold(v, "-XDELETE"):
			return zc1807Hit(cmd)
		case strings.EqualFold(v, "--method=DELETE"):
			return zc1807Hit(cmd)
		}
	}
	return nil
}

func zc1807Hit(cmd *ast.SimpleCommand) []Violation {
	return []Violation{{
		KataID: "ZC1807",
		Message: "`gh api -X DELETE` sends a raw DELETE to the GitHub API with the " +
			"caller's token — no `--yes` guard, no dry-run. Use the high-level `gh` " +
			"subcommand for the target, or wrap with a preflight GET + explicit " +
			"confirmation.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}
