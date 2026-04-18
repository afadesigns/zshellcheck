package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1706",
		Title:    "Error on `lvresize -L -SIZE` without `-r` — shrink without filesystem resize corrupts data",
		Severity: SeverityError,
		Description: "`lvresize -L -SIZE` (or `--size -SIZE`) shrinks the logical volume by " +
			"SIZE bytes/extents. The filesystem on top still thinks it owns the original " +
			"range; reads beyond the new LV end now return zeros, and the next write " +
			"corrupts metadata. The `-r` (`--resizefs`) flag tells lvresize to call " +
			"`fsadm` (which calls `resize2fs` / `xfs_growfs` / etc.) so the filesystem " +
			"shrinks first. For ext4, always shrink the FS before the LV; for XFS, online " +
			"shrink is impossible — back up, recreate, restore.",
		Check: checkZC1706,
	})
}

func checkZC1706(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "lvresize" && ident.Value != "lvreduce" {
		return nil
	}

	hasResizefs := false
	shrinking := ident.Value == "lvreduce" // lvreduce always shrinks
	for i, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-r" || v == "--resizefs" {
			hasResizefs = true
		}
		if (v == "-L" || v == "--size") && i+1 < len(cmd.Arguments) {
			next := cmd.Arguments[i+1].String()
			if strings.HasPrefix(next, "-") {
				shrinking = true
			}
		}
	}

	if !shrinking || hasResizefs {
		return nil
	}

	return []Violation{{
		KataID: "ZC1706",
		Message: "`" + ident.Value + "` shrinks the LV without `-r` / `--resizefs` — the " +
			"filesystem on top is not shrunk first and writes past the new boundary " +
			"corrupt metadata. Add `-r` (or shrink the FS manually first).",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}
