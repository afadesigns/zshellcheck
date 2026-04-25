package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1378",
		Title:    "Avoid uppercase `$DIRSTACK` — Zsh uses lowercase `$dirstack`",
		Severity: SeverityError,
		Description: "Bash's `$DIRSTACK` is the `pushd`/`popd` directory stack. Zsh exposes the " +
			"same stack as lowercase `$dirstack` (per zsh/parameter module). Using uppercase " +
			"`$DIRSTACK` in Zsh accesses an unrelated (and usually empty) variable.",
		Check: checkZC1378,
		Fix:   fixZC1378,
	})
}

// fixZC1378 lower-cases every `DIRSTACK` token inside an echo / print /
// printf argument to `dirstack`. Each occurrence becomes its own edit at
// the absolute source offset of that arg's token + the substring index;
// surrounding quoting and adjoining text stay byte-exact.
func fixZC1378(node ast.Node, _ Violation, source []byte) []FixEdit {
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
		if !strings.Contains(val, "DIRSTACK") {
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
		idx := 0
		for {
			pos := strings.Index(val[idx:], "DIRSTACK")
			if pos < 0 {
				break
			}
			abs := off + idx + pos
			line, col := offsetLineColZC1378(source, abs)
			if line < 0 {
				break
			}
			edits = append(edits, FixEdit{
				Line:    line,
				Column:  col,
				Length:  len("DIRSTACK"),
				Replace: "dirstack",
			})
			idx += pos + len("DIRSTACK")
		}
	}
	return edits
}

func offsetLineColZC1378(source []byte, offset int) (int, int) {
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

func checkZC1378(node ast.Node) []Violation {
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
		if strings.Contains(v, "DIRSTACK") {
			return []Violation{{
				KataID:  "ZC1378",
				Message: "Use lowercase `$dirstack` in Zsh — uppercase `$DIRSTACK` is Bash-only.",
				Line:    cmd.Token.Line,
				Column:  cmd.Token.Column,
				Level:   SeverityError,
			}}
		}
	}

	return nil
}
