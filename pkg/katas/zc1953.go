package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1953",
		Title:    "Warn on `mount --make-shared` / `--make-rshared` — flips propagation, container-escape vector",
		Severity: SeverityWarning,
		Description: "`mount --make-shared /path` (and the recursive `--make-rshared`) turns the " +
			"mount point into a peer in a shared-subtree group. Any later bind-mount that " +
			"lands inside it propagates to every other peer, including containers and other " +
			"namespaces. Combined with `CAP_SYS_ADMIN` inside a pod, that is one of the " +
			"classic container-escape stepping stones — a hostile workload can mount into the " +
			"host's `/` via the propagation group. Use `--make-private` on sensitive paths and " +
			"mount containers with `--mount-propagation=private` / `slave` unless the app " +
			"genuinely requires `shared`.",
		Check: checkZC1953,
	})
}

func checkZC1953(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	// Parser caveat: `mount --make-shared PATH` mangles the command name to
	// `make-shared` (or `make-rshared`).
	switch ident.Value {
	case "make-shared", "make-rshared":
		return zc1953Hit(cmd, "mount --"+ident.Value)
	case "mount":
		for _, arg := range cmd.Arguments {
			v := arg.String()
			if v == "--make-shared" || v == "--make-rshared" {
				return zc1953Hit(cmd, "mount "+v)
			}
		}
	}
	return nil
}

func zc1953Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1953",
		Message: "`" + form + "` puts the mount in a shared-subtree group — later " +
			"bind-mounts propagate to every peer, including containers. Classic escape " +
			"stepping stone. Use `--make-private` on sensitive paths.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}
