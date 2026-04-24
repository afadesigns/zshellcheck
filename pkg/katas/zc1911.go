package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1911",
		Title:    "Warn on `umount -l` / `--lazy` — detach now, leaves open fds pointing at a ghost mount",
		Severity: SeverityWarning,
		Description: "`umount -l` (lazy unmount) detaches the filesystem from the directory tree " +
			"immediately but defers the real cleanup until every open file descriptor on it is " +
			"closed. Any process still holding an fd keeps reading/writing into a mount that " +
			"`mount | grep` no longer lists — cron jobs drop logs into a phantom directory, a " +
			"re-mount of the same path stacks invisibly, and `lsof`/`fuser` often miss the " +
			"stale handles. Find and stop the holder (`lsof`/`fuser`/`systemd-cgls`) first, " +
			"then do a normal `umount`; reserve `-l` for break-glass recovery, not scripts.",
		Check: checkZC1911,
	})
}

var zc1911LazyFlags = map[string]bool{
	"-l":     true,
	"--lazy": true,
}

func checkZC1911(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "umount" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-l" || v == "--lazy" {
			line, col := FlagArgPosition(cmd, zc1911LazyFlags)
			return zc1911Hit(v, line, col)
		}
		if strings.HasPrefix(v, "-") && !strings.HasPrefix(v, "--") {
			// Clustered short flags, e.g. `-fl` / `-lf`.
			for i := 1; i < len(v); i++ {
				if v[i] == 'l' {
					tok := arg.TokenLiteralNode()
					return zc1911Hit("-l", tok.Line, tok.Column)
				}
			}
		}
	}
	return nil
}

func zc1911Hit(flag string, line, col int) []Violation {
	return []Violation{{
		KataID: "ZC1911",
		Message: "`umount " + flag + "` detaches the mount but leaves any open fd pointing at " +
			"a ghost filesystem — writers keep writing, re-mounts stack invisibly. Stop the " +
			"fd holder first (`lsof`/`fuser`), then do a normal `umount`.",
		Line:   line,
		Column: col,
		Level:  SeverityWarning,
	}}
}
