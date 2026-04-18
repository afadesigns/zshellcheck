package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1688",
		Title:    "Warn on `aws s3 sync --delete` — destination objects deleted when source diverges",
		Severity: SeverityWarning,
		Description: "`aws s3 sync SRC DST --delete` removes every object in DST that does not " +
			"exist under SRC. A misspelled SRC, an empty build directory, or a stale " +
			"`cd` turns the sync into a full-bucket wipe with no second confirmation and " +
			"no recovery unless the bucket had versioning enabled. Restrict deletion to " +
			"the prefix that really changed (`aws s3 sync ./build s3://bucket/app/ " +
			"--delete`), add `--dryrun` behind a gate, or enable versioning and MFA-delete " +
			"before running the command from a pipeline.",
		Check: checkZC1688,
	})
}

func checkZC1688(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "aws" {
		return nil
	}

	if len(cmd.Arguments) < 2 {
		return nil
	}
	if cmd.Arguments[0].String() != "s3" || cmd.Arguments[1].String() != "sync" {
		return nil
	}

	for _, arg := range cmd.Arguments[2:] {
		if arg.String() == "--delete" {
			return []Violation{{
				KataID: "ZC1688",
				Message: "`aws s3 sync --delete` wipes DST objects that are missing from " +
					"SRC — a mistyped SRC bulk-deletes the bucket. Scope to the prefix, " +
					"dry-run first, or enable versioning + MFA-delete.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}
