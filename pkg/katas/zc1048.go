package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1048",
		Title: "Avoid `source` with relative paths",
		Description: "Sourcing a file with a relative path (e.g. `source ./lib.zsh`) depends on the current " +
			"working directory. Use `${0:a:h}/lib.zsh` to source relative to the script location.",
		Check: checkZC1048,
	})
}

func checkZC1048(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	// Check if command is source or .
	name := cmd.Name.String()
	if name != "source" && name != "." {
		return nil
	}

	if len(cmd.Arguments) == 0 {
		return nil
	}

	arg := cmd.Arguments[0]

	// Check if arg is a StringLiteral or ConcatenatedExpression starting with "./" or "../"
	val := getStringValueZC1048(arg)

	// Remove quoting for check manually to avoid tool call escaping issues
	if len(val) > 0 && (val[0] == '"' || val[0] == '\'') {
		val = val[1:]
	}
	if len(val) > 0 && (val[len(val)-1] == '"' || val[len(val)-1] == '\'') {
		val = val[:len(val)-1]
	}

	if strings.HasPrefix(val, "./") || strings.HasPrefix(val, "../") {
		return []Violation{{
			KataID:  "ZC1048",
			Message: "Avoid `source` with relative paths. Use `${0:a:h}/...` to resolve relative to the script.",
			Line:    arg.TokenLiteralNode().Line,
			Column:  arg.TokenLiteralNode().Column,
		}}
	}

	return nil
}

func getStringValueZC1048(node ast.Node) string {
	switch n := node.(type) {
	case *ast.StringLiteral:
		return n.Value
	case *ast.ConcatenatedExpression:
		var sb strings.Builder
		for _, p := range n.Parts {
			sb.WriteString(getStringValueZC1048(p))
		}
		return sb.String()
	}
	return ""
}
