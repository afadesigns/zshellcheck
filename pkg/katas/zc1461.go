package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1461",
		Title:    "Avoid `docker run --pid=host` — shares host PID namespace with the container",
		Severity: SeverityWarning,
		Description: "`--pid=host` lets the container see every host process and send signals to " +
			"them, including sending SIGKILL to init-managed daemons or attaching a debugger to " +
			"host-side processes. Use only for diagnostic tools (e.g. strace/perf containers) and " +
			"never for general workloads.",
		Check: checkZC1461,
	})
}

func checkZC1461(node ast.Node) []Violation {
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

	var prevPid bool
	for _, arg := range cmd.Arguments {
		v := arg.String()

		if v == "--pid=host" {
			return violateZC1461(cmd)
		}
		if prevPid {
			prevPid = false
			if v == "host" {
				return violateZC1461(cmd)
			}
		}
		if v == "--pid" {
			prevPid = true
		}
	}

	return nil
}

func violateZC1461(cmd *ast.SimpleCommand) []Violation {
	return []Violation{{
		KataID: "ZC1461",
		Message: "`--pid=host` shares the host PID namespace — container can signal and " +
			"inspect every host process. Avoid outside debug tools.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}
