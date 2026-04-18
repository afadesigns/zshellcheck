package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1736",
		Title:    "Error on `pulumi destroy --yes` / `up --yes` — silent infra mutation in CI",
		Severity: SeverityError,
		Description: "`pulumi destroy --yes` (or `-y`) skips the preview-and-confirm step that " +
			"normally surfaces every resource scheduled for deletion. A single misresolved " +
			"stack name or wrong AWS credential resolves to a one-shot wipe of cloud " +
			"infrastructure. `pulumi up --yes` and `pulumi refresh --yes` carry the same " +
			"footgun for resource creation/replacement. Pipe `pulumi preview` output into " +
			"a review step (manual approval, GitHub Actions environment protection rule) " +
			"before applying, and never combine `--yes` with the `destroy` verb in " +
			"automation.",
		Check: checkZC1736,
	})
}

func checkZC1736(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "pulumi" {
		return nil
	}
	if len(cmd.Arguments) < 2 {
		return nil
	}

	sub := cmd.Arguments[0].String()
	switch sub {
	case "destroy", "up", "refresh":
	default:
		return nil
	}

	for _, arg := range cmd.Arguments[1:] {
		v := arg.String()
		if v == "--yes" || v == "-y" {
			return []Violation{{
				KataID: "ZC1736",
				Message: "`pulumi " + sub + " " + v + "` skips the preview-and-confirm — " +
					"a misresolved stack or credential wipes / mutates infrastructure " +
					"with no review. Gate behind `pulumi preview` plus a manual " +
					"approval step.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}
	return nil
}
