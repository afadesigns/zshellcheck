package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.IdentifierNode, Kata{
		ID:       "ZC1301",
		Title:    "Avoid `$PIPESTATUS` — use `$pipestatus` (lowercase) in Zsh",
		Severity: SeverityWarning,
		Description: "`$PIPESTATUS` is a Bash array containing exit statuses from the last " +
			"pipeline. Zsh uses `$pipestatus` (lowercase) for the same purpose. " +
			"The uppercase form is undefined in Zsh.",
		Check: checkZC1301,
		Fix:   fixZC1301,
	})
}

// fixZC1301 rewrites the uppercase Bash `$PIPESTATUS` / `PIPESTATUS`
// identifier to the lowercase Zsh `$pipestatus` / `pipestatus`
// form. Span covers only the name itself — subscripts and surrounding
// context stay in place.
func fixZC1301(node ast.Node, v Violation, source []byte) []FixEdit {
	ident, ok := node.(*ast.Identifier)
	if !ok {
		return nil
	}
	switch ident.Value {
	case "$PIPESTATUS":
		return []FixEdit{{
			Line:    v.Line,
			Column:  v.Column,
			Length:  len("$PIPESTATUS"),
			Replace: "$pipestatus",
		}}
	case "PIPESTATUS":
		return []FixEdit{{
			Line:    v.Line,
			Column:  v.Column,
			Length:  len("PIPESTATUS"),
			Replace: "pipestatus",
		}}
	}
	return nil
}

func checkZC1301(node ast.Node) []Violation {
	ident, ok := node.(*ast.Identifier)
	if !ok {
		return nil
	}

	if ident.Value != "$PIPESTATUS" && ident.Value != "PIPESTATUS" {
		return nil
	}

	return []Violation{{
		KataID:  "ZC1301",
		Message: "Avoid `$PIPESTATUS` in Zsh — use `$pipestatus` (lowercase) instead. The uppercase form is Bash-specific.",
		Line:    ident.Token.Line,
		Column:  ident.Token.Column,
		Level:   SeverityWarning,
	}}
}
