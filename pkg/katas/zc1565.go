package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1565",
		Title:    "Style: use `command -v` instead of `whereis` / `locate` for command existence",
		Severity: SeverityStyle,
		Description: "`whereis` searches a hard-coded list of binary/manual/source directories " +
			"and returns everything it finds, including stale paths on custom `$PATH` layouts. " +
			"`locate` relies on a cron-maintained index that may be hours or days stale. For " +
			"a scripted \"does this command exist?\" check, `command -v <cmd>` respects the " +
			"current `$PATH`, returns the selected resolution, and has no index-refresh " +
			"coupling.",
		Check: checkZC1565,
		Fix:   fixZC1565,
	})
}

// fixZC1565 rewrites a `whereis` / `locate` / `mlocate` / `plocate`
// command-name lookup into `command -v`. The detector restricts to the
// four index-based forms so the swap is safe; arguments stay untouched.
func fixZC1565(node ast.Node, v Violation, _ []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	switch ident.Value {
	case "whereis", "locate", "mlocate", "plocate":
		return []FixEdit{{
			Line:    v.Line,
			Column:  v.Column,
			Length:  len(ident.Value),
			Replace: "command -v",
		}}
	}
	return nil
}

func checkZC1565(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "whereis" && ident.Value != "locate" && ident.Value != "mlocate" &&
		ident.Value != "plocate" {
		return nil
	}

	return []Violation{{
		KataID: "ZC1565",
		Message: "`" + ident.Value + "` is index-based and stale-prone. Use `command -v " +
			"<cmd>` for runtime existence checks.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityStyle,
	}}
}
