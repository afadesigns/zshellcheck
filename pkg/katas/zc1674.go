package katas

import (
	"strconv"
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1674",
		Title:    "Warn on `docker/podman run --oom-kill-disable` or `--oom-score-adj <= -500`",
		Severity: SeverityWarning,
		Description: "`--oom-kill-disable` tells the kernel OOM killer to never touch the " +
			"container's memory cgroup — a leak inside then drives the whole host into OOM " +
			"reclaim until `sshd`, `systemd-journald`, or the init daemon itself gets " +
			"killed. `--oom-score-adj <= -500` stops short of full immunity but still " +
			"preferentially kills unrelated host processes under pressure. If the workload " +
			"genuinely needs resilience, cap memory with `--memory=<limit>` and accept the " +
			"container being killed on overrun; shift the heavy workload to a dedicated " +
			"node instead of rigging OOM scores.",
		Check: checkZC1674,
	})
}

func checkZC1674(node ast.Node) []Violation {
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
	if len(cmd.Arguments) == 0 || cmd.Arguments[0].String() != "run" {
		return nil
	}

	for i, arg := range cmd.Arguments[1:] {
		v := arg.String()
		if v == "--oom-kill-disable" {
			return zc1674Hit(cmd, v)
		}
		if strings.HasPrefix(v, "--oom-score-adj=") {
			adj := strings.TrimPrefix(v, "--oom-score-adj=")
			if zc1674Harsh(adj) {
				return zc1674Hit(cmd, v)
			}
			continue
		}
		if v == "--oom-score-adj" {
			idx := i + 2
			if idx >= len(cmd.Arguments) {
				continue
			}
			adj := cmd.Arguments[idx].String()
			if zc1674Harsh(adj) {
				return zc1674Hit(cmd, v+" "+adj)
			}
		}
	}
	return nil
}

func zc1674Harsh(adj string) bool {
	n, err := strconv.Atoi(adj)
	if err != nil {
		return false
	}
	return n <= -500
}

func zc1674Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1674",
		Message: "`" + form + "` shifts OOM pressure onto the rest of the host — cap " +
			"memory with `--memory=<limit>` instead of rigging the OOM score.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}
