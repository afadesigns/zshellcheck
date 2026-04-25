package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.IdentifierNode, Kata{
		ID:       "ZC1319",
		Title:    "Avoid `$BASH_ARGC` — use `$#` in Zsh",
		Severity: SeverityWarning,
		Description: "`$BASH_ARGC` is a Bash array tracking argument counts per stack frame. " +
			"Zsh uses `$#` for argument count and `$argv` for the argument array.",
		Check: checkZC1319,
		Fix:   fixZC1319,
	})
}

// fixZC1319 rewrites the Bash `$BASH_ARGC` / `BASH_ARGC` identifier to
// the Zsh `$#` form. Caveat: `$BASH_ARGC` is per-frame in Bash; `$#`
// is the current-frame argument count in Zsh. The rewrite is correct
// for the common single-value usage; multi-frame stack inspection is
// not portable and stays the user's responsibility.
func fixZC1319(node ast.Node, v Violation, source []byte) []FixEdit {
	ident, ok := node.(*ast.Identifier)
	if !ok || ident == nil {
		return nil
	}
	switch ident.Value {
	case "$BASH_ARGC":
		return []FixEdit{{
			Line:    v.Line,
			Column:  v.Column,
			Length:  len("$BASH_ARGC"),
			Replace: "$#",
		}}
	case "BASH_ARGC":
		return []FixEdit{{
			Line:    v.Line,
			Column:  v.Column,
			Length:  len("BASH_ARGC"),
			Replace: "#",
		}}
	}
	return nil
}

func checkZC1319(node ast.Node) []Violation {
	ident, ok := node.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident == nil {
		return nil
	}

	if ident.Value != "$BASH_ARGC" && ident.Value != "BASH_ARGC" {
		return nil
	}

	return []Violation{{
		KataID:  "ZC1319",
		Message: "Avoid `$BASH_ARGC` in Zsh — use `$#` for argument count. `BASH_ARGC` is Bash-specific.",
		Line:    ident.Token.Line,
		Column:  ident.Token.Column,
		Level:   SeverityWarning,
	}}
}
