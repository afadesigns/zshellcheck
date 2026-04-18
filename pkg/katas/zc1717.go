package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1717",
		Title:    "Warn on `docker pull/push --disable-content-trust` — bypasses image signature checks",
		Severity: SeverityWarning,
		Description: "When `DOCKER_CONTENT_TRUST=1` is enforced on a host (or set via `/etc/docker/" +
			"daemon.json`), Docker rejects unsigned image pulls and signs every push. The " +
			"`--disable-content-trust` flag overrides that per command: a `pull` accepts a " +
			"replaced or unsigned image into local storage, a `push` lands an unsigned tag in " +
			"the registry where downstream pulls cannot verify provenance. Drop the flag and " +
			"sign the artifact (`docker trust sign IMAGE:TAG`) instead, or scope the bypass " +
			"with a tight Notary signer policy.",
		Check: checkZC1717,
	})
}

func checkZC1717(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "docker" {
		return nil
	}

	var sub string
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if sub == "" {
			switch v {
			case "pull", "push", "build", "create", "run":
				sub = v
				continue
			}
		}
		if sub != "" && v == "--disable-content-trust" {
			return []Violation{{
				KataID: "ZC1717",
				Message: "`docker " + sub + " --disable-content-trust` overrides " +
					"`DOCKER_CONTENT_TRUST=1` — unsigned image moves into the registry " +
					"or local store. Sign the artifact (`docker trust sign`) instead.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}
