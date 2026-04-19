package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1866",
		Title:    "Warn on `docker exec -u 0` — bypasses the image's non-root `USER` directive",
		Severity: SeverityWarning,
		Description: "A hardened image runs with a non-root `USER` set in its Dockerfile so " +
			"exploited processes inside the container are contained by the Linux " +
			"user-namespace mapping. `docker exec -u 0` (and `-u root`, `--user=0`, the " +
			"podman equivalent) overrides that choice on a per-exec basis and drops a " +
			"shell back into uid 0 — every subsequent file write, cap check, and namespace " +
			"test now runs as root inside the container, which on a default Docker setup " +
			"is also root on the host via the shared mount namespace. Keep exec sessions " +
			"as the container's configured user; if you genuinely need root for a one-off " +
			"fix, document it in the ticket and consider rebuilding the image with the " +
			"capability baked in so `-u 0` is never required.",
		Check: checkZC1866,
	})
}

func checkZC1866(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "docker" && ident.Value != "podman" {
		return nil
	}
	args := cmd.Arguments
	if len(args) == 0 || args[0].String() != "exec" {
		return nil
	}
	for i := 1; i < len(args); i++ {
		v := args[i].String()
		var user string
		switch {
		case (v == "-u" || v == "--user") && i+1 < len(args):
			user = args[i+1].String()
		case strings.HasPrefix(v, "-u") && v != "-u":
			user = strings.TrimPrefix(v, "-u")
		case strings.HasPrefix(v, "--user="):
			user = strings.TrimPrefix(v, "--user=")
		default:
			continue
		}
		user = strings.Trim(user, "\"'")
		if zc1866IsRoot(user) {
			return []Violation{{
				KataID: "ZC1866",
				Message: "`" + ident.Value + " exec -u " + user + "` drops a root " +
					"shell — bypasses the image's non-root `USER` and, without " +
					"userns remap, equals host root. Keep execs as the container " +
					"user.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}

func zc1866IsRoot(v string) bool {
	if v == "0" || v == "root" || v == "0:0" {
		return true
	}
	// `0:gid` or `0:groupname` — still uid 0.
	if strings.HasPrefix(v, "0:") {
		return true
	}
	if strings.HasPrefix(v, "root:") {
		return true
	}
	return false
}
