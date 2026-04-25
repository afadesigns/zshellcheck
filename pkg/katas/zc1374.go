package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1374",
		Title:    "Avoid `$FUNCNEST` — Zsh uses `$FUNCNEST` as a limit, not a depth indicator",
		Severity: SeverityWarning,
		Description: "Bash's `$FUNCNEST` is both a writable limit and (implicitly) the current " +
			"depth-query vehicle. Zsh's `$FUNCNEST` is only the limit — to read the current depth " +
			"use `${#funcstack}`. Reading `$FUNCNEST` expecting depth returns the limit, not " +
			"the current depth.",
		Check: checkZC1374,
		Fix:   fixZC1374,
	})
}

// fixZC1374 rewrites `$FUNCNEST` / `${FUNCNEST}` arguments to
// `${#funcstack}` inside echo / print / printf calls. One edit per
// matching arg. Idempotent — a re-run sees `${#funcstack}`, which
// the detector's exact-match guard won't match.
func fixZC1374(node ast.Node, _ Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "echo" && ident.Value != "print" && ident.Value != "printf" {
		return nil
	}
	var edits []FixEdit
	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val != "$FUNCNEST" && val != "${FUNCNEST}" {
			continue
		}
		tok := arg.TokenLiteralNode()
		off := LineColToByteOffset(source, tok.Line, tok.Column)
		if off < 0 || off+len(val) > len(source) {
			continue
		}
		if string(source[off:off+len(val)]) != val {
			continue
		}
		edits = append(edits, FixEdit{
			Line:    tok.Line,
			Column:  tok.Column,
			Length:  len(val),
			Replace: "${#funcstack}",
		})
	}
	return edits
}

func checkZC1374(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "echo" && ident.Value != "print" && ident.Value != "printf" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "$FUNCNEST" || v == "${FUNCNEST}" {
			return []Violation{{
				KataID: "ZC1374",
				Message: "In Zsh, `$FUNCNEST` is the configured limit, not the current depth. " +
					"Use `${#funcstack}` for current function nesting depth.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}

	return nil
}
