package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1588",
		Title:    "Error on `nsenter --target 1` — joins host init namespaces (container escape)",
		Severity: SeverityError,
		Description: "`nsenter -t 1` attaches to the namespaces of pid 1. Inside a privileged " +
			"container or one with `CAP_SYS_ADMIN`, pid 1 is the host init — joining its " +
			"mount / pid / net / uts / ipc namespaces is the canonical escape primitive. " +
			"From that new shell the caller sees and writes the host filesystem, kills host " +
			"processes, and hijacks host network. Legit debugging runs from the host, not from " +
			"inside the container. If you need to exec into a container, use `docker exec` / " +
			"`kubectl exec`.",
		Check: checkZC1588,
	})
}

func checkZC1588(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "nsenter" {
		return nil
	}

	var expectTarget, hit bool
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "--target=1" || v == "-t1" {
			hit = true
			break
		}
		if expectTarget && v == "1" {
			hit = true
			break
		}
		expectTarget = v == "-t" || v == "--target"
	}
	if !hit {
		return nil
	}
	return []Violation{{
		KataID: "ZC1588",
		Message: "`nsenter --target 1` joins the host init namespaces — classic container-" +
			"escape primitive. Use `docker exec` / `kubectl exec` for legitimate debugging.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}
