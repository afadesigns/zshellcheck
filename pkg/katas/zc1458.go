package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1458",
		Title:    "Warn on explicit `docker run --user root` / `--user 0`",
		Severity: SeverityWarning,
		Description: "Running as UID 0 inside a container means a break-out bug leaves the " +
			"attacker as root on the host (absent user namespaces). Build images with a " +
			"non-root `USER` directive and avoid overriding to root at runtime.",
		Check: checkZC1458,
	})
}

func checkZC1458(node ast.Node) []Violation {
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

	var prevU bool
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if prevU {
			prevU = false
			if v == "root" || v == "0" || strings.HasPrefix(v, "0:") || strings.HasPrefix(v, "root:") {
				return []Violation{{
					KataID: "ZC1458",
					Message: "Explicit root UID inside a container lets container-escape bugs " +
						"become host root. Use a non-root USER in the image.",
					Line:   cmd.Token.Line,
					Column: cmd.Token.Column,
					Level:  SeverityWarning,
				}}
			}
		}
		if v == "-u" || v == "--user" {
			prevU = true
		}
	}

	return nil
}
