package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1441",
		Title:    "Warn on `docker system prune -af` / `-a --force` (or similar podman/k8s)",
		Severity: SeverityWarning,
		Description: "`docker system prune -af` deletes every unused image, container, network, " +
			"and (with `--volumes`) volume. On shared CI runners or build hosts this obliterates " +
			"cached layers and slows future builds. Scope prunes with `--filter \"until=168h\"` " +
			"or target one resource type at a time.",
		Check: checkZC1441,
	})
}

func checkZC1441(node ast.Node) []Violation {
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

	seenPrune := false
	seenA := false
	seenF := false
	for _, arg := range cmd.Arguments {
		v := arg.String()
		switch v {
		case "prune":
			seenPrune = true
		case "-a", "--all":
			seenA = true
		case "-f", "--force":
			seenF = true
		case "-af", "-fa":
			seenA = true
			seenF = true
		}
	}
	if seenPrune && seenA && seenF {
		return []Violation{{
			KataID: "ZC1441",
			Message: "`docker prune -af` / `-a --force` deletes all unused resources without " +
				"prompt. Scope with `--filter` or target one resource type.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityWarning,
		}}
	}

	return nil
}
