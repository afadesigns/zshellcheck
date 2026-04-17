package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1454",
		Title:    "Avoid `docker/podman run --privileged` — disables most container isolation",
		Severity: SeverityError,
		Description: "`--privileged` disables the seccomp profile, grants all Linux capabilities, " +
			"and lets the container access all host devices. It is effectively equivalent to " +
			"running the process as host root. Add specific capabilities with `--cap-add` and " +
			"bind-mount specific devices with `--device` instead.",
		Check: checkZC1454,
	})
}

func checkZC1454(node ast.Node) []Violation {
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

	hasRun := false
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "run" || v == "exec" || v == "create" {
			hasRun = true
		}
		if hasRun && v == "--privileged" {
			return []Violation{{
				KataID: "ZC1454",
				Message: "`--privileged` disables container isolation — effectively host root. " +
					"Use `--cap-add` + `--device` for narrow permissions.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}

	return nil
}
