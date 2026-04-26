// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1446",
		Title:    "Dangerous: `aws s3 rm --recursive` / `s3 rb --force` — bulk S3 deletion",
		Severity: SeverityError,
		Description: "`aws s3 rm s3://bucket/prefix --recursive` deletes every key under the " +
			"prefix. `aws s3 rb --force` deletes the bucket along with its contents. Combine " +
			"with a wrong prefix or bucket name and data loss is total. Enable versioning on " +
			"production buckets and use `aws s3api list-object-versions` before bulk removals.",
		Check: checkZC1446,
	})
}

func checkZC1446(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "aws" {
		return nil
	}

	var seenS3 bool
	var seenRm bool
	var seenRb bool
	var seenRecursive bool
	var seenForce bool
	for _, arg := range cmd.Arguments {
		v := arg.String()
		switch v {
		case "s3":
			seenS3 = true
		case "rm":
			seenRm = true
		case "rb":
			seenRb = true
		case "--recursive":
			seenRecursive = true
		case "--force":
			seenForce = true
		}
	}
	if seenS3 && ((seenRm && seenRecursive) || (seenRb && seenForce)) {
		return []Violation{{
			KataID: "ZC1446",
			Message: "`aws s3 rm --recursive` / `s3 rb --force` mass-deletes objects/buckets. " +
				"Enable versioning and dry-run with `--dryrun`.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityError,
		}}
	}

	return nil
}
