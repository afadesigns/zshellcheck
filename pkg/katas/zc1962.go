package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1962",
		Title:    "Warn on `kustomize build --load-restrictor=LoadRestrictionsNone` — path-traversal in overlays",
		Severity: SeverityWarning,
		Description: "Kustomize's default `LoadRestrictionsRootOnly` limits every base, patch, " +
			"configMapGenerator, and secretGenerator to paths under the current kustomization " +
			"root. `kustomize build … --load-restrictor=LoadRestrictionsNone` (also the legacy " +
			"spelling `--load_restrictor none` / `--load-restrictor=LoadRestrictionsNone_WarnForAll`) " +
			"drops that guard, so an overlay from an untrusted remote base can reference " +
			"`../../secrets/prod.env` or absolute paths and pull them into the render. Keep " +
			"the default; if a legitimate overlay needs a sibling file, vendor it in.",
		Check: checkZC1962,
	})
}

func checkZC1962(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "kustomize" {
		return nil
	}
	if len(cmd.Arguments) == 0 || cmd.Arguments[0].String() != "build" {
		return nil
	}

	for i, arg := range cmd.Arguments[1:] {
		v := arg.String()
		if strings.HasPrefix(v, "--load-restrictor=") ||
			strings.HasPrefix(v, "--load_restrictor=") {
			val := v[strings.IndexByte(v, '=')+1:]
			if zc1962IsNoneVariant(val) {
				return zc1962Hit(cmd, v)
			}
		}
		if (v == "--load-restrictor" || v == "--load_restrictor") && i+2 <= len(cmd.Arguments)-1 {
			val := cmd.Arguments[i+2].String()
			if zc1962IsNoneVariant(val) {
				return zc1962Hit(cmd, v+" "+val)
			}
		}
	}
	return nil
}

func zc1962IsNoneVariant(val string) bool {
	switch val {
	case "none", "None", "LoadRestrictionsNone":
		return true
	}
	return false
}

func zc1962Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1962",
		Message: "`kustomize build " + form + "` drops path-root restriction — untrusted " +
			"overlays can reference `../../secrets/prod.env` and pull them into the render. " +
			"Keep the default; vendor sibling files into the overlay.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}
