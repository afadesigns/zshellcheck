package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

var zc1784HooksMutablePrefixes = []string{
	"/tmp/",
	"/var/tmp/",
	"/dev/shm/",
	"/home/",
	"/root/",
	"/opt/",
	"/srv/",
	"/mnt/",
	"/media/",
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1784",
		Title:    "Warn on `git config core.hooksPath /tmp/...` — hook execution from a mutable path",
		Severity: SeverityWarning,
		Description: "`core.hooksPath` tells git which directory to run repository hooks from. " +
			"Any file named `pre-commit`, `post-checkout`, `post-merge`, etc. under that " +
			"directory becomes executable code invoked by routine git operations. Pointing " +
			"`core.hooksPath` at `/tmp`, `/var/tmp`, `/dev/shm`, `/home/<other>`, `/opt`, " +
			"`/srv`, or `/mnt` hands the git CLI an execution primitive from a path that a " +
			"non-root (or another) user can write at will — a classic supply-chain entry " +
			"point on shared hosts and CI runners. Keep hooks inside the repo's `.git/hooks/` " +
			"(or a repo-owned `.githooks/` directory) and configure `core.hooksPath` only to " +
			"paths that share the repo's owner and permissions.",
		Check: checkZC1784,
	})
}

func checkZC1784(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "git" {
		return nil
	}
	if len(cmd.Arguments) < 3 {
		return nil
	}
	if cmd.Arguments[0].String() != "config" {
		return nil
	}

	for i, arg := range cmd.Arguments[1:] {
		v := arg.String()
		if v != "core.hooksPath" {
			continue
		}
		// Path follows.
		if 1+i+1 >= len(cmd.Arguments) {
			return nil
		}
		path := cmd.Arguments[1+i+1].String()
		path = strings.Trim(path, "\"'")
		if !strings.HasPrefix(path, "/") {
			return nil
		}
		for _, prefix := range zc1784HooksMutablePrefixes {
			if strings.HasPrefix(path, prefix) {
				return []Violation{{
					KataID: "ZC1784",
					Message: "`git config core.hooksPath " + path + "` runs hooks from " +
						"a mutable path — supply-chain primitive. Keep hooks in the " +
						"repo's `.git/hooks/` (or a tracked `.githooks/`) and point " +
						"`core.hooksPath` at repo-owned paths only.",
					Line:   cmd.Token.Line,
					Column: cmd.Token.Column,
					Level:  SeverityWarning,
				}}
			}
		}
	}
	return nil
}
