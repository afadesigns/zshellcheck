package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1647",
		Title:    "Warn on `kubectl apply -f URL` — remote manifest applied without digest verification",
		Severity: SeverityWarning,
		Description: "`kubectl apply -f https://...` fetches the manifest over the network and " +
			"applies it to the cluster. TLS (when present) verifies transport but not " +
			"authorship — if the URL is compromised or the content changes between reviews, " +
			"the cluster picks up the new definition. Pin the content: download to disk, " +
			"verify a known SHA256, then `kubectl apply -f local.yaml`. For plain HTTP the " +
			"attacker controls the response directly — never acceptable.",
		Check: checkZC1647,
	})
}

func checkZC1647(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "kubectl" {
		return nil
	}
	if len(cmd.Arguments) < 2 {
		return nil
	}
	sub := cmd.Arguments[0].String()
	if sub != "apply" && sub != "create" && sub != "replace" {
		return nil
	}

	for i, arg := range cmd.Arguments[1:] {
		if arg.String() != "-f" {
			continue
		}
		idx := i + 2
		if idx >= len(cmd.Arguments) {
			continue
		}
		target := cmd.Arguments[idx].String()
		if strings.HasPrefix(target, "http://") || strings.HasPrefix(target, "https://") {
			return []Violation{{
				KataID: "ZC1647",
				Message: "`kubectl " + sub + " -f " + target + "` applies a remote " +
					"manifest — verify digest first. Download, check SHA256, then apply " +
					"the local file.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}
