// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1472",
		Title:    "Error on `aws s3 --acl public-read` / `public-read-write` (public bucket)",
		Severity: SeverityError,
		Description: "Using the `public-read` or `public-read-write` canned ACL when uploading, " +
			"syncing, or setting a bucket policy makes the object (and often the bucket) readable " +
			"by anyone on the internet. Prefer bucket policies scoped to specific principals, or " +
			"CloudFront with Origin Access Identity if you truly need public read.",
		Check: checkZC1472,
	})
}

func checkZC1472(node ast.Node) []Violation {
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

	// Must see `s3` or `s3api` service argument anywhere before `--acl`.
	var sawService bool
	var prevAcl bool
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "s3" || v == "s3api" {
			sawService = true
		}
		if !sawService {
			continue
		}
		if prevAcl {
			prevAcl = false
			if v == "public-read" || v == "public-read-write" {
				return zc1472Violation(cmd, v)
			}
		}
		if v == "--acl" {
			prevAcl = true
		}
		if v == "--acl=public-read" || v == "--acl=public-read-write" {
			return zc1472Violation(cmd, v[len("--acl="):])
		}
	}
	return nil
}

func zc1472Violation(cmd *ast.SimpleCommand, acl string) []Violation {
	return []Violation{{
		KataID:  "ZC1472",
		Message: "Canned ACL `" + acl + "` makes the object (often the bucket) world-readable. Use a scoped bucket policy instead.",
		Line:    cmd.Token.Line,
		Column:  cmd.Token.Column,
		Level:   SeverityError,
	}}
}
