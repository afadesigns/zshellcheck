package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1457",
		Title:    "Warn on bind-mount of `/var/run/docker.sock` — container escape vector",
		Severity: SeverityWarning,
		Description: "Mounting `/var/run/docker.sock` into a container lets the container start " +
			"any privileged container, mount host filesystems, and effectively gain root on the " +
			"host. Reserve this for trusted CI/tooling images; for general workloads use " +
			"rootless containers or a dedicated orchestrator API.",
		Check: checkZC1457,
	})
}

func checkZC1457(node ast.Node) []Violation {
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

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if strings.Contains(v, "docker.sock") || strings.Contains(v, "/var/run/docker") {
			return []Violation{{
				KataID: "ZC1457",
				Message: "Mounting `/var/run/docker.sock` gives the container effective root on " +
					"the host. Reserve for trusted tooling.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}

	return nil
}
