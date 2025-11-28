package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
	"github.com/afadesigns/zshellcheck/pkg/token"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1084",
		Title: "Quote globs in `find` commands",
		Description: "Unquoted globs in `find` commands are expanded by the shell before `find` runs. " +
			"If files match, `find` receives the list of files instead of the pattern. " +
			"Quote arguments to `-name`, `-path`, etc.",
		Check: checkZC1084,
	})
}

func checkZC1084(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	// Check if command is 'find'
	cmdName := cmd.Name.String()
	if cmdName != "find" && cmdName != "gfind" {
		return nil
	}

	violations := []Violation{}

	for i := 0; i < len(cmd.Arguments); i++ {
		arg := cmd.Arguments[i]

		// Handle merged -name[...] case
		if concat, ok := arg.(*ast.ConcatenatedExpression); ok {
			var prefixStr string
			foundMerged := false
			for _, part := range concat.Parts {
				if str, ok := part.(*ast.StringLiteral); ok {
					if str.Value == "[" && str.Token.Type != token.STRING {
						if isFindGlobFlag(prefixStr) {
							violations = append(violations, Violation{
								KataID:  "ZC1084",
								Message: "Quote globs in `find` commands. `" + cleanString(arg.String()) + "` contains unquoted brackets.",
								Line:    str.Token.Line,
								Column:  str.Token.Column,
							})
							foundMerged = true
							break
						}
					}
				} else if idx, ok := part.(*ast.IndexExpression); ok {
					// IndexExpression implies unquoted [
					// Check if prefix + left part forms a flag
					candidate := prefixStr + cleanString(idx.Left.String())
					if isFindGlobFlag(candidate) {
						violations = append(violations, Violation{
							KataID:  "ZC1084",
							Message: "Quote globs in `find` commands. `" + cleanString(arg.String()) + "` contains unquoted brackets.",
							Line:    idx.Token.Line,
							Column:  idx.Token.Column,
						})
						foundMerged = true
						break
					}
				}
				prefixStr += cleanString(part.String())
			}
			if foundMerged {
				continue
			}
		}

		// Check for flags that take a pattern
		flag := getFlagName(arg)
		if !isFindGlobFlag(flag) {
			continue
		}

		// Check next argument
		if i+1 >= len(cmd.Arguments) {
			break
		}
		patternArg := cmd.Arguments[i+1]
		i++ // Advance

		if isUnquotedGlob(patternArg) {
			violations = append(violations, Violation{
				KataID:  "ZC1084",
				Message: "Quote globs in `find` commands. `" + cleanString(patternArg.String()) + "` is subject to shell expansion.",
				Line:    patternArg.TokenLiteralNode().Line,
				Column:  patternArg.TokenLiteralNode().Column,
			})
		}
	}

	return violations
}

func isFindGlobFlag(f string) bool {
	switch f {
	case "-name", "-iname", "-path", "-ipath", "-wholename", "-iwholename", "-lname", "-ilname":
		return true
	}
	return false
}

func cleanString(s string) string {
	// Remove all outer parens added by AST String() methods
	for len(s) >= 2 && s[0] == '(' && s[len(s)-1] == ')' {
		s = s[1 : len(s)-1]
	}
	return s
}

func getFlagName(node ast.Node) string {
	// -name parsed as PrefixExpression (- name)
	if prefix, ok := node.(*ast.PrefixExpression); ok {
		if prefix.Operator == "-" {
			return "-" + prefix.Right.String()
		}
	}
	// -name parsed as StringLiteral (if quoted or simple)
	// or ConcatenatedExpression
	s := node.String()
	s = cleanString(s)
	return s
}

func isUnquotedGlob(node ast.Expression) bool {
	// Check for SimpleCommand (e.g. [a-z])
	if sc, ok := node.(*ast.SimpleCommand); ok {
		if sc.Name.String() == "[" {
			return true
		}
	}

	// Check for IndexExpression (e.g. file[a-z])
	if _, ok := node.(*ast.IndexExpression); ok {
		return true
	}

	// Check for StringLiteral with glob tokens
	if str, ok := node.(*ast.StringLiteral); ok {
		tok := str.TokenLiteralNode()
		if tok.Type == token.STRING {
			return false
		}
		return isGlobToken(tok)
	}

	// Check ConcatenatedExpression
	if concat, ok := node.(*ast.ConcatenatedExpression); ok {
		escaped := false
		for _, part := range concat.Parts {
			if str, ok := part.(*ast.StringLiteral); ok {
				tok := str.TokenLiteralNode()
				if escaped {
					escaped = false
					continue
				}
				if tok.Literal == "\\" { // Backslash
					escaped = true
					continue
				}
				if isGlobToken(tok) {
					return true
				}
			} else if sc, ok := part.(*ast.SimpleCommand); ok {
				if sc.Name.String() == "[" {
					return true
				}
			} else if prefix, ok := part.(*ast.PrefixExpression); ok {
				if escaped {
					escaped = false
					continue
				}
				if prefix.Operator == "*" || prefix.Operator == "?" {
					return true
				}
			} else {
				escaped = false
			}
		}
	}

	// Check PrefixExpression (e.g. *.txt, ?foo)
	if prefix, ok := node.(*ast.PrefixExpression); ok {
		if prefix.Operator == "*" || prefix.Operator == "?" {
			return true
		}
	}

	return false
}

func isGlobToken(tok token.Token) bool {
	if tok.Type == token.ASTERISK { // *
		return true
	}
	if (tok.Type == token.ILLEGAL && tok.Literal == "?") || tok.Type == token.QUESTION { // ?
		return true
	}
	if tok.Type == token.LBRACKET { // [
		return true
	}
	return false
}
