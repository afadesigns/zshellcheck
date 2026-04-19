package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1753",
		Title:    "Error on `rclone purge REMOTE:PATH` — bulk delete of every object under the remote path",
		Severity: SeverityError,
		Description: "`rclone purge REMOTE:PATH` removes every object and empty directory under " +
			"PATH on the remote — no dry-run gate, no confirmation, no soft-delete unless " +
			"the backend happens to version. A typo'd path or a stale variable turns one " +
			"line into a bucket-wide wipe (S3, GCS, Azure, Swift all honour the same API " +
			"call). Preview with `rclone lsf REMOTE:PATH` or `rclone delete --dry-run`, " +
			"then use `rclone delete` scoped narrower; enable object versioning on the " +
			"backend so a bad run can roll back.",
		Check: checkZC1753,
	})
}

func checkZC1753(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "rclone" {
		return nil
	}
	if len(cmd.Arguments) == 0 {
		return nil
	}
	if cmd.Arguments[0].String() != "purge" {
		return nil
	}

	return []Violation{{
		KataID: "ZC1753",
		Message: "`rclone purge` removes every object under the remote path with no dry-run " +
			"or soft-delete. Preview with `rclone lsf` / `rclone delete --dry-run` and " +
			"prefer narrower `rclone delete`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}
