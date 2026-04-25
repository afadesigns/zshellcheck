package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1380",
		Title:    "Avoid `$HISTIGNORE` — use Zsh `$HISTORY_IGNORE`",
		Severity: SeverityWarning,
		Description: "Bash filters history entries matching `$HISTIGNORE` patterns. Zsh uses a " +
			"parameter named `$HISTORY_IGNORE` (underscore in the middle). Setting `HISTIGNORE` " +
			"in Zsh is a no-op.",
		Check: checkZC1380,
		Fix:   fixZC1380,
	})
}

// fixZC1380 rewrites the Bash `HISTIGNORE` parameter name to the Zsh
// `HISTORY_IGNORE` spelling. The detector ignores args that already
// contain `HISTORY_IGNORE`, so the rewrite is idempotent on a re-run.
// Span covers only the bare name occurrences inside the argument
// string; surrounding `=value` / quoting stays byte-identical.
func fixZC1380(node ast.Node, _ Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "echo" && ident.Value != "print" && ident.Value != "printf" && ident.Value != "export" {
		return nil
	}
	var edits []FixEdit
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if !strings.Contains(v, "HISTIGNORE") || strings.Contains(v, "HISTORY_IGNORE") {
			continue
		}
		tok := arg.TokenLiteralNode()
		argOff := LineColToByteOffset(source, tok.Line, tok.Column)
		if argOff < 0 {
			continue
		}
		// Verify the argument literal sits at this offset before
		// scanning. Some quote shapes alter the start byte; bail on
		// mismatch rather than risk a misaligned splice.
		if argOff+len(v) > len(source) || string(source[argOff:argOff+len(v)]) != v {
			continue
		}
		// Replace every occurrence of the bare name inside the arg.
		idx := 0
		for {
			rel := strings.Index(v[idx:], "HISTIGNORE")
			if rel < 0 {
				break
			}
			absStart := argOff + idx + rel
			line, col := offsetLineColZC1380(source, absStart)
			if line > 0 {
				edits = append(edits, FixEdit{
					Line:    line,
					Column:  col,
					Length:  len("HISTIGNORE"),
					Replace: "HISTORY_IGNORE",
				})
			}
			idx += rel + len("HISTIGNORE")
		}
	}
	if len(edits) == 0 {
		return nil
	}
	return edits
}

func offsetLineColZC1380(source []byte, offset int) (int, int) {
	if offset < 0 || offset > len(source) {
		return -1, -1
	}
	line := 1
	col := 1
	for i := 0; i < offset; i++ {
		if source[i] == '\n' {
			line++
			col = 1
			continue
		}
		col++
	}
	return line, col
}

func checkZC1380(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "echo" && ident.Value != "print" && ident.Value != "printf" && ident.Value != "export" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if strings.Contains(v, "HISTIGNORE") && !strings.Contains(v, "HISTORY_IGNORE") {
			return []Violation{{
				KataID: "ZC1380",
				Message: "`$HISTIGNORE` is Bash-only. In Zsh use `$HISTORY_IGNORE` (underscored) " +
					"for the same history-pattern filter.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}

	return nil
}
