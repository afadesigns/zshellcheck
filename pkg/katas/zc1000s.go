// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"reflect"
	"regexp"
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
	"github.com/afadesigns/zshellcheck/pkg/token"
)

func init() {
	RegisterKata(ast.IndexExpressionNode, Kata{
		ID:    "ZC1001",
		Title: "Use ${} for array element access",
		Description: "In Zsh, accessing array elements with `$my_array[1]` doesn't work as expected. " +
			"It tries to access an element from an array named `my_array[1]`. " +
			"The correct way to access an array element is to use `${my_array[1]}`.",
		Severity: SeverityStyle,
		Check:    checkZC1001,
		Fix:      fixZC1001,
	})
	RegisterKata(ast.InvalidArrayAccessNode, Kata{
		ID:    "ZC1001",
		Title: "Use ${} for array element access",
		Description: "In Zsh, accessing array elements with `$my_array[1]` doesn't work as expected. " +
			"It tries to access an element from an array named `my_array[1]`. " +
			"The correct way to access an array element is to use `${my_array[1]}`.",
		Severity: SeverityStyle,
		Check:    checkZC1001,
		Fix:      fixZC1001,
	})
}

// fixZC1001 rewrites `$arr[i]` to `${arr[i]}`. Two edits: insert
// `{` between the `$` and the identifier, then insert `}` after
// the closing `]`. Source positions are derived from the violation
// column (which points at the leading `$`) and a quote/brace-aware
// scan for the matching `]`.
func fixZC1001(node ast.Node, v Violation, source []byte) []FixEdit {
	dollarOff := LineColToByteOffset(source, v.Line, v.Column)
	if dollarOff < 0 || dollarOff >= len(source) || source[dollarOff] != '$' {
		return nil
	}
	// Find the `[` that opens the subscript starting from after `$name`.
	// Walk identifier chars then expect `[`.
	i := dollarOff + 1
	for i < len(source) && (isIdentByte(source[i])) {
		i++
	}
	if i >= len(source) || source[i] != '[' {
		return nil
	}
	closeOff := findSubscriptClose(source, i)
	if closeOff < 0 {
		return nil
	}
	return []FixEdit{
		// Insert `{` immediately after the `$`.
		{Line: v.Line, Column: v.Column + 1, Length: 0, Replace: "{"},
		// Insert `}` immediately after the closing `]`.
		offsetToEdit(source, closeOff+1, 0, "}"),
	}
}

// findSubscriptClose returns the byte offset of the `]` that closes
// the subscript opened at open. Tracks single / double quotes and
// nested `[` so subscripts containing patterns like `arr[$other[i]]`
// or `arr[(R)pat]` close at the right bracket. Returns -1 when no
// match is found before EOF.
func findSubscriptClose(source []byte, open int) int {
	inSingle := false
	inDouble := false
	depth := 1
	for i := open + 1; i < len(source); i++ {
		c := source[i]
		switch {
		case inSingle:
			if c == '\'' {
				inSingle = false
			}
		case inDouble:
			if c == '\\' && i+1 < len(source) {
				i++
				continue
			}
			if c == '"' {
				inDouble = false
			}
		default:
			switch c {
			case '\'':
				inSingle = true
			case '"':
				inDouble = true
			case '[':
				depth++
			case ']':
				depth--
				if depth == 0 {
					return i
				}
			case '\n':
				return -1
			}
		}
	}
	return -1
}

func checkZC1001(node ast.Node) []Violation {
	violations := []Violation{}

	if indexExp, ok := node.(*ast.IndexExpression); ok {
		if ident, ok := indexExp.Left.(*ast.Identifier); ok {
			if len(ident.Value) > 0 && ident.Value[0] == '$' {
				violations = append(violations, Violation{
					KataID: "ZC1001",
					Message: "Use ${} for array element access. " +
						"Accessing array elements with `" + ident.Value + "[...]` is not the correct syntax in Zsh.",
					Line:   ident.Token.Line,
					Column: ident.Token.Column,
					Level:  SeverityStyle,
				})
			}
		}
	} else if arrayAccess, ok := node.(*ast.InvalidArrayAccess); ok {
		violations = append(violations, Violation{
			KataID: "ZC1001",
			Message: "Use ${} for array element access. " +
				"Accessing array elements with `$my_array[1]` is not the correct syntax in Zsh.",
			Line:   arrayAccess.Token.Line,
			Column: arrayAccess.Token.Column,
			Level:  SeverityStyle,
		})
	}

	return violations
}

func init() {
	RegisterKata(ast.CommandSubstitutionNode, Kata{
		ID:    "ZC1002",
		Title: "Use $(...) instead of backticks",
		Description: "Backticks are the old-style command substitution. " +
			"$(...) is nesting-safe, easier to read, and generally preferred.",
		Severity: SeverityStyle,
		Check:    checkZC1002,
		Fix:      fixZC1002,
	})
}

// fixZC1002 rewrites “ `cmd` “ -> `$(cmd)`. The violation's Line and
// Column point at the opening backtick. Locate the matching closing
// backtick by scanning forward, respecting backslash escapes, and emit
// a single replacement edit spanning both delimiters. Unterminated
// backtick spans are skipped (the parser rejects them earlier; this is
// defensive).
func fixZC1002(node ast.Node, v Violation, source []byte) []FixEdit {
	cs, ok := node.(*ast.CommandSubstitution)
	if !ok {
		return nil
	}
	start := LineColToByteOffset(source, v.Line, v.Column)
	if start < 0 || start >= len(source) || source[start] != '`' {
		return nil
	}

	// Walk forward to find the matching closing backtick. Double-quoted
	// strings and escaped backticks inside the span are preserved.
	end := -1
	for i := start + 1; i < len(source); i++ {
		switch source[i] {
		case '\\':
			i++ // skip escaped char
		case '`':
			end = i
		}
		if end >= 0 {
			break
		}
	}
	if end < 0 {
		return nil
	}

	inner := string(source[start+1 : end])
	// Defensive: if the inner payload already contains `$(...)` shaped
	// parens we could still round-trip, but our parser produced
	// cs.Command so use its stringified form as a sanity cross-check.
	if cs.Command != nil && cs.Command.String() == "" {
		return nil
	}
	return []FixEdit{{
		Line:    v.Line,
		Column:  v.Column,
		Length:  end - start + 1,
		Replace: "$(" + inner + ")",
	}}
}

func checkZC1002(node ast.Node) []Violation {
	violations := []Violation{}

	if cs, ok := node.(*ast.CommandSubstitution); ok {
		violations = append(violations, Violation{
			KataID: "ZC1002",
			Message: "Use $(...) instead of backticks for command substitution. " +
				"The `$(...)` syntax is more readable and can be nested easily.",
			Line:   cs.Token.Line,
			Column: cs.Token.Column,
			Level:  SeverityStyle,
		})
	}

	return violations
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1003",
		Title: "Use `((...))` for arithmetic comparisons instead of `[` or `test`",
		Description: "Bash/Zsh have a dedicated arithmetic context `((...))` " +
			"which is cleaner and faster than `[` or `test` for numeric comparisons.",
		Severity: SeverityStyle,
		Check:    checkZC1003,
		Fix:      fixZC1003,
	})
}

// arithCmpReplacements maps the dash-prefixed test-comparison
// operators to their arithmetic equivalents. Used by fixZC1003 to
// rewrite `[ x -eq y ]` → `(( x == y ))`.
var arithCmpReplacements = map[string]string{
	"-eq": "==",
	"-ne": "!=",
	"-lt": "<",
	"-le": "<=",
	"-gt": ">",
	"-ge": ">=",
}

func checkZC1003(node ast.Node) []Violation {
	violations := []Violation{}

	if cmd, ok := node.(*ast.SimpleCommand); ok {
		name := cmd.Name.String()
		if name == "[" || name == "test" {
			for _, arg := range cmd.Arguments {
				val := arg.String()
				// Trim parens added by AST String() method for expressions
				val = strings.Trim(val, "()")

				if _, found := arithCmpReplacements[val]; found {
					violations = append(violations, Violation{
						KataID:  "ZC1003",
						Message: "Use `((...))` for arithmetic comparisons instead of `[` or `test`.",
						Line:    cmd.Token.Line,
						Column:  cmd.Token.Column,
						Level:   SeverityStyle,
					})
					return violations
				}
			}
		}
	}

	return violations
}

// fixZC1003 rewrites `[ x -eq y ]` to `(( x == y ))`. Only the
// `[` form is auto-fixable: replace the opening `[` with `((`,
// the matching `]` with `))`, and the `-eq`/etc. operator with
// its arithmetic equivalent. The `test` form has no closing
// terminator on the line and the surrounding context (pipelines,
// chains) makes a safe rewrite ambiguous; it stays detection-only.
//
// Bails when the command shape doesn't match: more than one
// comparison operator, no `[` byte at the violation column,
// missing close bracket, or no recognised operator.
func fixZC1003(node ast.Node, v Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok || cmd.Name == nil {
		return nil
	}
	if cmd.Name.String() != "[" {
		return nil
	}
	var op string
	var opIdx int
	for i, arg := range cmd.Arguments {
		val := strings.Trim(arg.String(), "()")
		if _, found := arithCmpReplacements[val]; found {
			if op != "" {
				return nil // multiple ops — not safe to auto-fix
			}
			op = val
			opIdx = i
		}
	}
	if op == "" {
		return nil
	}
	openOff := LineColToByteOffset(source, v.Line, v.Column)
	if openOff < 0 || openOff >= len(source) || source[openOff] != '[' {
		return nil
	}
	closeOff := findTestCloseBracket(source, openOff)
	if closeOff < 0 {
		return nil
	}
	opTok := cmd.Arguments[opIdx]
	opLine := opTok.TokenLiteralNode().Line
	opCol := opTok.TokenLiteralNode().Column
	return []FixEdit{
		{Line: v.Line, Column: v.Column, Length: 1, Replace: "(("},
		{Line: opLine, Column: opCol, Length: len(op), Replace: arithCmpReplacements[op]},
		offsetToEdit(source, closeOff, 1, "))"),
	}
}

func init() {
	kata := Kata{
		ID:    "ZC1004",
		Title: "Use `return` instead of `exit` in functions",
		Description: "Using `exit` in a function terminates the entire shell, which is often unintended " +
			"in interactive sessions or sourced scripts. Use `return` to exit the function.",
		Severity: SeverityWarning,
		Check:    checkZC1004,
		Fix:      fixZC1004,
	}
	RegisterKata(ast.FunctionDefinitionNode, kata)
	RegisterKata(ast.FunctionLiteralNode, kata)
}

// fixZC1004 rewrites `exit` to `return` at the command-name position
// inside a function body. Arguments (the exit/return code) stay
// unchanged — `exit 1` becomes `return 1` with the `1` byte-identical.
// The violation's Line/Column already point at the command name.
func fixZC1004(node ast.Node, v Violation, source []byte) []FixEdit {
	return []FixEdit{{
		Line:    v.Line,
		Column:  v.Column,
		Length:  len("exit"),
		Replace: "return",
	}}
}

func checkZC1004(node ast.Node) []Violation {
	var body ast.Statement

	switch n := node.(type) {
	case *ast.FunctionDefinition:
		body = n.Body
	case *ast.FunctionLiteral:
		body = n.Body
	default:
		return nil
	}

	violations := []Violation{}

	ast.Walk(body, func(n ast.Node) bool {
		// Stop traversal at subshell boundaries where exit is safe/scoped
		switch t := n.(type) {
		case *ast.GroupedExpression: // ( ... )
			return false
		case *ast.Subshell: // ( ... ) as subshell
			return false
		case *ast.CommandSubstitution: // ` ... `
			return false
		case *ast.DollarParenExpression: // $( ... )
			return false
		case *ast.BlockStatement:
			if t.Token.Type == token.LPAREN { // ( ... ) as a statement block
				return false
			}
		}

		// Match both SimpleCommand (`exit 1`) and bare Identifier
		// wrapped in an ExpressionStatement (`exit` with no args —
		// the parser folds zero-arg command invocations into a plain
		// identifier expression rather than a SimpleCommand).
		switch sn := n.(type) {
		case *ast.SimpleCommand:
			if sn.Name != nil && sn.Name.String() == "exit" {
				violations = append(violations, Violation{
					KataID:  "ZC1004",
					Message: "Use `return` instead of `exit` in functions to avoid killing the shell.",
					Line:    sn.Token.Line,
					Column:  sn.Token.Column,
					Level:   SeverityWarning,
				})
			}
			// Don't descend — the Name Identifier would otherwise
			// double-count as a bare-`exit` hit below.
			return false
		case *ast.Identifier:
			if sn.Value == "exit" {
				violations = append(violations, Violation{
					KataID:  "ZC1004",
					Message: "Use `return` instead of `exit` in functions to avoid killing the shell.",
					Line:    sn.Token.Line,
					Column:  sn.Token.Column,
					Level:   SeverityWarning,
				})
			}
		}
		return true
	})

	return violations
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1005",
		Title: "Use whence instead of which",
		Description: "The `which` command is an external command and may not be available on all systems. " +
			"The `whence` command is a built-in Zsh command that provides a more reliable and consistent " +
			"way to find the location of a command.",
		Severity: SeverityInfo,
		Check:    checkZC1005,
		Fix:      fixZC1005,
	})
}

// fixZC1005 rewrites `which` -> `whence` at the command name position.
// Arguments are unchanged; the two builtins share the identifier-query
// shape for the common case.
func fixZC1005(node ast.Node, v Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "which" {
		return nil
	}
	return []FixEdit{{
		Line:    v.Line,
		Column:  v.Column,
		Length:  len("which"),
		Replace: "whence",
	}}
}

func checkZC1005(node ast.Node) []Violation {
	violations := []Violation{}

	if cmd, ok := node.(*ast.SimpleCommand); ok {
		if ident, ok := cmd.Name.(*ast.Identifier); ok {
			if ident.Value == "which" {
				violations = append(violations, Violation{
					KataID: "ZC1005",
					Message: "Use whence instead of which. The `whence` command is a built-in Zsh command " +
						"that provides a more reliable and consistent way to find the location of a command.",
					Line:   ident.Token.Line,
					Column: ident.Token.Column,
					Level:  SeverityInfo,
				})
			}
		}
	}

	return violations
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1006",
		Title: "Prefer [[ over test for tests",
		Description: "The `test` command is an external command and may not be available on all systems. " +
			"The `[[...]]` construct is a Zsh keyword, offering safer and more powerful conditional " +
			"expressions than the traditional `test` command. It prevents word splitting and pathname " +
			"expansion, and supports advanced features like regex matching.",
		Severity: SeverityStyle,
		Check:    checkZC1006,
		// All three of ZC1006 / ZC1020 / ZC1036 fire on the same `test`
		// shape and want the same `[[ … ]]` rewrite that ZC1293 ships.
		// The conflict resolver dedupes overlapping edits.
		Fix: fixZC1293,
	})
}

func checkZC1006(node ast.Node) []Violation {
	violations := []Violation{}

	if cmd, ok := node.(*ast.SimpleCommand); ok {
		if ident, ok := cmd.Name.(*ast.Identifier); ok {
			if ident.Value == "test" {
				violations = append(violations, Violation{
					KataID: "ZC1006",
					Message: "Prefer [[ over test for tests. " +
						"[[ is a Zsh keyword that offers safer and more powerful conditional expressions.",
					Line:   ident.Token.Line,
					Column: ident.Token.Column,
					Level:  SeverityStyle,
				})
			}
		}
	}

	return violations
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1007",
		Title: "Avoid using `chmod 777`",
		Description: "Using `chmod 777` is a security risk as it gives read, write, and execute " +
			"permissions to everyone. It's better to use more restrictive permissions.",
		Severity: SeverityWarning,
		Check:    checkZC1007,
	})
}

func checkZC1007(node ast.Node) []Violation {
	violations := []Violation{}

	if cmd, ok := node.(*ast.SimpleCommand); ok {
		if ident, ok := cmd.Name.(*ast.Identifier); ok {
			if ident.Value == "chmod" {
				for _, arg := range cmd.Arguments {
					switch v := arg.(type) {
					case *ast.Identifier:
						if v.Value == "777" {
							violations = append(violations, Violation{
								KataID:  "ZC1007",
								Message: "Avoid using `chmod 777`. It is a security risk.",
								Line:    ident.Token.Line,
								Column:  ident.Token.Column,
								Level:   SeverityWarning,
							})
						}
					case *ast.IntegerLiteral:
						if v.Token.Literal == "777" {
							violations = append(violations, Violation{
								KataID:  "ZC1007",
								Message: "Avoid using `chmod 777`. It is a security risk.",
								Line:    ident.Token.Line,
								Column:  ident.Token.Column,
								Level:   SeverityWarning,
							})
						}
					}
				}
			}
		}
	}

	return violations
}

func init() {
	RegisterKata(ast.LetStatementNode, Kata{
		ID:    "ZC1008",
		Title: "Use `\\$(())` for arithmetic operations",
		Description: "The `let` command is a shell builtin, but the `\\$(())` syntax is more portable " +
			"and generally preferred for arithmetic operations in Zsh. It's also more powerful as it " +
			"can be used in more contexts.",
		Severity: SeverityStyle,
		Check:    checkZC1008,
		// Reuse ZC1013's `let NAME=EXPR` → `(( NAME = EXPR ))` rewrite.
		// All three of ZC1008, ZC1013, ZC1022 fire on the same `let`
		// shape and want the same arithmetic-command form; the
		// conflict resolver dedupes overlapping edits.
		Fix: fixZC1013,
	})
}

func checkZC1008(node ast.Node) []Violation {
	// Duplicate check for 'let' covered by ZC1013?
	// ZC1008 title says \$(()) which is expansion.
	// But check was for LetStatement.
	// Let's keep it as 'let' check for now to match original intent, maybe redundant.
	stmt, ok := node.(*ast.LetStatement)
	if !ok {
		return nil
	}

	return []Violation{{
		KataID:  "ZC1008",
		Message: "Use `\\$(())` for arithmetic operations instead of `let`.",
		Line:    stmt.TokenLiteralNode().Line,
		Column:  stmt.TokenLiteralNode().Column,
		Level:   SeverityStyle,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1009",
		Title: "Use `((...))` for C-style arithmetic",
		Description: "The `((...))` construct in Zsh allows for C-style arithmetic. " +
			"It is generally more efficient and readable than using `expr` or other " +
			"external commands for arithmetic.",
		Severity: SeverityStyle,
		Check:    checkZC1009,
	})
}

func checkZC1009(node ast.Node) []Violation {
	violations := []Violation{}

	if cmd, ok := node.(*ast.SimpleCommand); ok {
		if ident, ok := cmd.Name.(*ast.Identifier); ok {
			if ident.Value == "expr" {
				violations = append(violations, Violation{
					KataID:  "ZC1009",
					Message: "Use `((...))` for C-style arithmetic instead of `expr`.",
					Line:    ident.Token.Line,
					Column:  ident.Token.Column,
					Level:   SeverityStyle,
				})
			}
		}
	}

	return violations
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1010",
		Title: "Use [[ ... ]] instead of [ ... ]",
		Description: "Zsh's [[ ... ]] is more powerful and safer than [ ... ]. " +
			"It supports pattern matching, regex, and doesn't require quoting variables to prevent word splitting.",
		Severity: SeverityStyle,
		Check:    checkZC1010,
		Fix:      fixZC1010,
	})
}

// fixZC1010 rewrites a `[ … ]` test command to `[[ … ]]`. The opening
// bracket at the violation's coordinates becomes `[[`; the matching
// closing bracket on the same logical line becomes `]]`. Contents
// stay byte-identical so quoting and expansions are preserved.
//
// Bail when the shape is not a simple `[ … ]` test (e.g. second token
// is not `[`, or the logical line has no closing bracket): a
// malformed test is not safely auto-fixable.
func fixZC1010(node ast.Node, v Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	if cmd.Name == nil || cmd.Name.String() != "[" {
		return nil
	}
	open := LineColToByteOffset(source, v.Line, v.Column)
	if open < 0 || open >= len(source) || source[open] != '[' {
		return nil
	}
	close := findTestCloseBracket(source, open)
	if close < 0 {
		return nil
	}
	return []FixEdit{
		{Line: v.Line, Column: v.Column, Length: 1, Replace: "[["},
		offsetToEdit(source, close, 1, "]]"),
	}
}

// findTestCloseBracket returns the byte offset of the closing `]`
// that terminates a `[ … ]` test opened at open, or -1 when no clean
// close is found before an end-of-statement terminator. Single- and
// double-quoted strings, and `${…}` braces are respected so the scan
// doesn't match `]` inside `"$arr[idx]"`.
func findTestCloseBracket(source []byte, open int) int {
	st := bracketScan{}
	for i := open + 1; i < len(source); i++ {
		if next, ok := st.advance(source, i); ok {
			i = next
			continue
		}
		if st.closedAt(source, i) {
			return i
		}
		if st.terminated(source, i) {
			return -1
		}
	}
	return -1
}

type bracketScan struct {
	inSingle, inDouble bool
	braceDepth         int
}

// advance returns (newIndex, true) when the byte at i drives the scanner
// forward past quoted/escaped material without further classification.
func (s *bracketScan) advance(source []byte, i int) (int, bool) {
	c := source[i]
	switch {
	case s.inSingle:
		if c == '\'' {
			s.inSingle = false
		}
		return i, true
	case s.inDouble:
		if c == '\\' && i+1 < len(source) {
			return i + 1, true
		}
		if c == '"' {
			s.inDouble = false
		}
		return i, true
	}
	return i, false
}

// closedAt reports whether the unquoted byte at i is the matching `]`.
// Side effect: classifies opening / closing braces to track depth.
func (s *bracketScan) closedAt(source []byte, i int) bool {
	switch source[i] {
	case '\'':
		s.inSingle = true
	case '"':
		s.inDouble = true
	case '{':
		s.braceDepth++
	case '}':
		if s.braceDepth > 0 {
			s.braceDepth--
		}
	case ']':
		if s.braceDepth == 0 {
			return true
		}
	}
	return false
}

func (s *bracketScan) terminated(source []byte, i int) bool {
	if s.inSingle || s.inDouble {
		return false
	}
	c := source[i]
	return c == '\n' || c == ';'
}

// offsetToEdit builds a FixEdit whose Line/Column correspond to the
// given byte offset inside source. Used when a Fix already has a
// byte offset but the FixEdit type expects 1-based coordinates.
func offsetToEdit(source []byte, offset, length int, replace string) FixEdit {
	line := 1
	col := 1
	for i := 0; i < offset && i < len(source); i++ {
		if source[i] == '\n' {
			line++
			col = 1
			continue
		}
		col++
	}
	return FixEdit{Line: line, Column: col, Length: length, Replace: replace}
}

func checkZC1010(node ast.Node) []Violation {
	violations := []Violation{}

	if cmd, ok := node.(*ast.SimpleCommand); ok {
		// Check if command name is "["
		if cmd.Name.String() == "[" {
			violations = append(violations, Violation{
				KataID:  "ZC1010",
				Message: "Use `[[ ... ]]` instead of `[ ... ]` or `test`. `[[` is safer and more powerful.",
				Line:    cmd.Token.Line,
				Column:  cmd.Token.Column,
				Level:   SeverityStyle,
			})
		}
	}

	return violations
}

var plumbingCommands = map[string]string{
	"rev-parse":    "git-rev-parse",
	"update-ref":   "git-update-ref",
	"symbolic-ref": "git-symbolic-ref",
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1011",
		Title: "Use `git` porcelain commands instead of plumbing commands",
		Description: "Plumbing commands in `git` are designed for scripting and can be unstable. " +
			"Porcelain commands are designed for interactive use and are more stable.",
		Severity: SeverityInfo,
		Check:    checkZC1011,
	})
}

func checkZC1011(node ast.Node) []Violation {
	violations := []Violation{}

	if cmd, ok := node.(*ast.SimpleCommand); ok {
		if ident, ok := cmd.Name.(*ast.Identifier); ok {
			if ident.Value == "git" {
				for _, arg := range cmd.Arguments {
					val := getArgValue(arg)
					if _, ok := plumbingCommands[val]; ok {
						violations = append(violations, Violation{
							KataID:  "ZC1011",
							Message: "Avoid using `git` plumbing commands in scripts. They are not guaranteed to be stable.",
							Line:    ident.Token.Line,
							Column:  ident.Token.Column,
							Level:   SeverityInfo,
						})
					}
				}
			}
		}
	}

	return violations
}

func getArgValue(expr ast.Expression) string {
	switch e := expr.(type) {
	case *ast.Identifier:
		return e.Value
	case *ast.StringLiteral:
		return e.Value
	case *ast.InfixExpression:
		return getArgValue(e.Left) + e.Operator + getArgValue(e.Right)
	case *ast.ConcatenatedExpression:
		var out string
		for _, p := range e.Parts {
			out += getArgValue(p)
		}
		return out
	case *ast.PrefixExpression:
		// Special handling for prefix expressions in command args to avoid parens
		return e.Operator + getArgValue(e.Right)
	}
	return ""
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1012",
		Title: "Use `read -r` to prevent backslash escaping",
		Description: "By default, `read` interprets backslashes as escape characters. " +
			"Use `read -r` to treat backslashes literally, which is usually what you want.",
		Severity: SeverityStyle,
		Check:    checkZC1012,
		Fix:      fixZC1012,
	})
}

// fixZC1012 inserts ` -r` directly after the `read` command name.
// Existing flags are left untouched (`read -p "x" VAR` becomes
// `read -r -p "x" VAR`) so the fix is order-preserving and idempotent
// on a second pass (the re-parse will see `-r` and detection won't
// fire).
func fixZC1012(node ast.Node, v Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	if cmd.Name.String() != "read" {
		return nil
	}
	nameOffset := LineColToByteOffset(source, v.Line, v.Column)
	if nameOffset < 0 {
		return nil
	}
	nameLen := IdentLenAt(source, nameOffset)
	if nameLen != len("read") {
		return nil
	}
	insertAt := nameOffset + nameLen
	insertLine, insertCol := byteOffsetToLineColZC1012(source, insertAt)
	if insertLine < 0 {
		return nil
	}
	return []FixEdit{{
		Line:    insertLine,
		Column:  insertCol,
		Length:  0,
		Replace: " -r",
	}}
}

// byteOffsetToLineColZC1012 converts a byte offset to a 1-based
// (line, column). Kept kata-local to avoid exposing a shared helper
// that the rest of the package does not yet need.
func byteOffsetToLineColZC1012(source []byte, offset int) (int, int) {
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

func checkZC1012(node ast.Node) []Violation {
	violations := []Violation{}

	if cmd, ok := node.(*ast.SimpleCommand); ok {
		if cmd.Name.String() == "read" {
			hasR := false
			for _, arg := range cmd.Arguments {
				s := arg.String()

				// Handle PrefixExpression String() format: "(-r)" -> "-r"
				s = strings.Trim(s, "()")

				if len(s) > 0 && s[0] == '-' {
					if strings.Contains(s, "r") {
						hasR = true
						break
					}
				}
			}

			if !hasR {
				violations = append(violations, Violation{
					KataID:  "ZC1012",
					Message: "Use `read -r` to read input without interpreting backslashes.",
					Line:    cmd.Token.Line,
					Column:  cmd.Token.Column,
					Level:   SeverityStyle,
				})
			}
		}
	}

	return violations
}

func init() {
	RegisterKata(ast.LetStatementNode, Kata{
		ID:    "ZC1013",
		Title: "Use `((...))` for arithmetic operations instead of `let`",
		Description: "The `let` command is a shell builtin, but the `((...))` syntax is more portable " +
			"and generally preferred for arithmetic operations in Zsh.",
		Severity: SeverityInfo,
		Check:    checkZC1013,
		Fix:      fixZC1013,
	})
}

// fixZC1013 rewrites `let NAME=EXPR` -> `(( NAME = EXPR ))`. The
// replacement spans from the `let` keyword to the end of the logical
// line (first `;`, `\n`, or EOF). The original expression text is
// preserved byte-identical so operator precedence and side effects
// match the source. Multi-assignment forms (`let a=1 b=2`) are not
// attempted — the Fix bails when the AST does not have a single
// Name/Value pair.
func fixZC1013(node ast.Node, v Violation, source []byte) []FixEdit {
	stmt, ok := node.(*ast.LetStatement)
	if !ok {
		return nil
	}
	if stmt.Name == nil || stmt.Value == nil {
		return nil
	}
	// Defer to ZC1032 when the value matches the C-style
	// increment/decrement shape so the rewrite produces the
	// idiomatic `(( i++ ))` / `(( i-- ))` form rather than the
	// generic `(( i = i+1 ))` form. Both fixes span the same source
	// range, so emitting both would lose ZC1032's narrower output to
	// the conflict resolver's deterministic tie-break.
	if _, increment := zc1032Op(stmt); increment {
		return nil
	}
	start := LineColToByteOffset(source, v.Line, v.Column)
	if start < 0 {
		return nil
	}
	// Scan for the end of the statement: first semicolon, newline,
	// or EOF.
	end := start
	for end < len(source) {
		b := source[end]
		if b == '\n' || b == ';' {
			break
		}
		end++
	}
	// `let NAME=EXPR` — after the keyword the source is NAME=EXPR.
	// Split on the first `=` to honour the original spelling (the
	// AST stores Name separately but Value.String() can lose
	// whitespace or parentheses that matter).
	prefix := "let "
	if start+len(prefix) > end {
		return nil
	}
	body := string(source[start+len(prefix) : end])
	eq := -1
	for i := 0; i < len(body); i++ {
		if body[i] == '=' {
			eq = i
			break
		}
	}
	if eq < 0 {
		return nil
	}
	name := body[:eq]
	rhs := body[eq+1:]
	return []FixEdit{{
		Line:    v.Line,
		Column:  v.Column,
		Length:  end - start,
		Replace: "(( " + name + " = " + rhs + " ))",
	}}
}

func checkZC1013(node ast.Node) []Violation {
	stmt, ok := node.(*ast.LetStatement)
	if !ok {
		return nil
	}

	return []Violation{{
		KataID:  "ZC1013",
		Message: "Use `((...))` for arithmetic operations instead of `let`.",
		Line:    stmt.TokenLiteralNode().Line,
		Column:  stmt.TokenLiteralNode().Column,
		Level:   SeverityInfo,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1014",
		Title: "Use `git switch` or `git restore` instead of `git checkout`",
		Description: "The `git checkout` command can be ambiguous. `git switch` is used for switching " +
			"branches and `git restore` is used for restoring files. Using these more specific commands " +
			"can make your scripts clearer and less error-prone.",
		Severity: SeverityInfo,
		Check:    checkZC1014,
	})
}

func checkZC1014(node ast.Node) []Violation {
	violations := []Violation{}

	if cmd, ok := node.(*ast.SimpleCommand); ok {
		if name, ok := cmd.Name.(*ast.Identifier); ok && name.Value == "git" {
			if len(cmd.Arguments) > 0 {
				if arg, ok := cmd.Arguments[0].(*ast.Identifier); ok && arg.Value == "checkout" {
					violations = append(violations, Violation{
						KataID:  "ZC1014",
						Message: "Use `git switch` or `git restore` instead of the ambiguous `git checkout`.",
						Line:    name.Token.Line,
						Column:  name.Token.Column,
						Level:   SeverityInfo,
					})
				}
			}
		}
	}

	return violations
}

func init() {
	RegisterKata(ast.CommandSubstitutionNode, Kata{
		ID:    "ZC1015",
		Title: "Use `$(...)` for command substitution instead of backticks",
		Description: "The `$(...)` syntax is the modern, recommended way to perform command substitution. " +
			"It is more readable and can be nested easily, unlike backticks.",
		Severity: SeverityStyle,
		Check:    checkZC1015,
		Fix:      fixZC1002,
	})
}

func checkZC1015(node ast.Node) []Violation {
	violations := []Violation{}

	if cs, ok := node.(*ast.CommandSubstitution); ok {
		violations = append(violations, Violation{
			KataID:  "ZC1015",
			Message: "Use `$(...)` for command substitution instead of backticks.",
			Line:    cs.Token.Line,
			Column:  cs.Token.Column,
			Level:   SeverityStyle,
		})
	}

	return violations
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1016",
		Title: "Use `read -s` when reading sensitive information",
		Description: "When asking for passwords or secrets, use `read -s` to prevent " +
			"the input from being echoed to the terminal.",
		Severity: SeverityStyle,
		Check:    checkZC1016,
		Fix:      fixZC1016,
	})
}

// fixZC1016 inserts ` -s` after the `read` command name. The detector
// gates on the absence of `-s` in any flag bundle, so the insertion
// is idempotent on a re-run.
func fixZC1016(node ast.Node, v Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	if cmd.Name == nil || cmd.Name.String() != "read" {
		return nil
	}
	nameOff := LineColToByteOffset(source, v.Line, v.Column)
	if nameOff < 0 || IdentLenAt(source, nameOff) != len("read") {
		return nil
	}
	insertAt := nameOff + len("read")
	insLine, insCol := offsetLineColZC1016(source, insertAt)
	if insLine < 0 {
		return nil
	}
	return []FixEdit{{
		Line:    insLine,
		Column:  insCol,
		Length:  0,
		Replace: " -s",
	}}
}

func offsetLineColZC1016(source []byte, offset int) (int, int) {
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

func checkZC1016(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	if cmd.Name.String() != "read" {
		return nil
	}

	hasS := false
	sensitiveVars := []string{"password", "passwd", "pwd", "secret", "token", "key", "api_key"}

	// Check flags
	for _, arg := range cmd.Arguments {
		argStr := arg.String()
		// Remove quotes if present
		argStr = strings.Trim(argStr, "\"'")
		if strings.HasPrefix(argStr, "-") {
			if strings.Contains(argStr, "s") {
				hasS = true
			}
		}
	}

	if hasS {
		return nil
	}

	violations := []Violation{}

	for _, arg := range cmd.Arguments {
		// Skip flags
		argStr := arg.String()
		argStrClean := strings.Trim(argStr, "\"'")
		if strings.HasPrefix(argStrClean, "-") {
			continue
		}

		// Handle Zsh read syntax: variable?prompt
		parts := strings.Split(argStr, "?")
		varName := strings.TrimSpace(parts[0])
		varName = strings.Trim(varName, "'\"")

		varLower := strings.ToLower(varName)
		isSensitive := false
		for _, s := range sensitiveVars {
			if strings.Contains(varLower, s) {
				isSensitive = true
				break
			}
		}

		if isSensitive {
			violations = append(violations, Violation{
				KataID:  "ZC1016",
				Message: "Use `read -s` to hide input when reading sensitive variable '" + varName + "'.",
				Line:    cmd.TokenLiteralNode().Line,
				Column:  cmd.TokenLiteralNode().Column,
				Level:   SeverityStyle,
			})
		}
	}

	return violations
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1017",
		Title: "Use `print -r` to print strings literally",
		Description: "The `print` command interprets backslash escape sequences by default. " +
			"To print a string literally, use the `-r` option.",
		Severity: SeverityStyle,
		Check:    checkZC1017,
		Fix:      fixZC1017,
	})
}

// fixZC1017 inserts ` -r` directly after the `print` command name.
// Existing flags are left in place, mirroring ZC1012's `read -r`
// insertion: `print "x"` becomes `print -r "x"`, `print -n "x"`
// becomes `print -r -n "x"`. Idempotent on a second pass — once
// `-r` appears among the flags the detector no longer fires.
func fixZC1017(node ast.Node, v Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	name, ok := cmd.Name.(*ast.Identifier)
	if !ok || name.Value != "print" {
		return nil
	}
	nameOffset := LineColToByteOffset(source, v.Line, v.Column)
	if nameOffset < 0 {
		return nil
	}
	nameLen := IdentLenAt(source, nameOffset)
	if nameLen != len("print") {
		return nil
	}
	insertAt := nameOffset + nameLen
	insLine, insCol := byteOffsetToLineColZC1017(source, insertAt)
	if insLine < 0 {
		return nil
	}
	return []FixEdit{{
		Line:    insLine,
		Column:  insCol,
		Length:  0,
		Replace: " -r",
	}}
}

func byteOffsetToLineColZC1017(source []byte, offset int) (int, int) {
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

func checkZC1017(node ast.Node) []Violation {
	violations := []Violation{}

	if cmd, ok := node.(*ast.SimpleCommand); ok {
		if name, ok := cmd.Name.(*ast.Identifier); ok && name.Value == "print" {
			hasRFlag := false
			for _, arg := range cmd.Arguments {
				argStr := arg.String()
				argStr = strings.Trim(argStr, "\"'")
				if strings.HasPrefix(argStr, "-") && strings.Contains(argStr, "r") {
					hasRFlag = true
					break
				}
			}
			if !hasRFlag {
				violations = append(violations, Violation{
					KataID:  "ZC1017",
					Message: "Use `print -r` to print strings literally.",
					Line:    name.Token.Line,
					Column:  name.Token.Column,
					Level:   SeverityStyle,
				})
			}
		}
	}

	return violations
}

// Issue #343: ZC1018 fires on the same input as the canonical
// ZC1009 with overlapping advice. Stubbed to a no-op so legacy
// `disabled_katas` lists keep parsing; the detection now lives in
// ZC1009.

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:          "ZC1018",
		Title:       "Superseded by ZC1009 — retired duplicate",
		Severity:    SeverityStyle,
		Description: "Retained as a no-op stub so legacy `.zshellcheckrc` files that disable this ID keep parsing. See https://github.com/afadesigns/zshellcheck/issues/343 for context; the canonical detection lives in ZC1009.",
		Check:       checkZC1018,
	})
}

func checkZC1018(ast.Node) []Violation {
	return nil
}

// Issue #342: ZC1019 fires on the same input as the canonical
// ZC1005 with overlapping advice. Stubbed to a no-op so legacy
// `disabled_katas` lists keep parsing; the detection now lives in
// ZC1005.

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:          "ZC1019",
		Title:       "Superseded by ZC1005 — retired duplicate",
		Severity:    SeverityStyle,
		Description: "Retained as a no-op stub so legacy `.zshellcheckrc` files that disable this ID keep parsing. See https://github.com/afadesigns/zshellcheck/issues/342 for context; the canonical detection lives in ZC1005.",
		Check:       checkZC1019,
	})
}

func checkZC1019(ast.Node) []Violation {
	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1020",
		Title: "Use `[[ ... ]]` for tests instead of `test`",
		Description: "The `test` command is an external command and may not be available on all systems. " +
			"The `[[...]]` construct is a Zsh keyword, offering safer and more powerful conditional " +
			"expressions than the traditional `test` command.",
		Severity: SeverityStyle,
		Check:    checkZC1020,
		Fix:      fixZC1293,
	})
}

func checkZC1020(node ast.Node) []Violation {
	violations := []Violation{}

	if cmd, ok := node.(*ast.SimpleCommand); ok {
		if ident, ok := cmd.Name.(*ast.Identifier); ok {
			if ident.Value == "test" {
				violations = append(violations, Violation{
					KataID:  "ZC1020",
					Message: "Use `[[ ... ]]` for tests instead of `test`.",
					Line:    ident.Token.Line,
					Column:  ident.Token.Column,
					Level:   SeverityStyle,
				})
			}
		}
	}

	return violations
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1021",
		Title: "Use symbolic permissions with `chmod` instead of octal",
		Description: "Symbolic permissions (e.g., `u+x`) are more readable and less error-prone than " +
			"octal permissions (e.g., `755`).",
		Severity: SeverityStyle,
		Check:    checkZC1021,
	})
}

func checkZC1021(node ast.Node) []Violation {
	violations := []Violation{}

	if cmd, ok := node.(*ast.SimpleCommand); ok {
		if name, ok := cmd.Name.(*ast.Identifier); ok && name.Value == "chmod" {
			for _, arg := range cmd.Arguments {
				if _, ok := arg.(*ast.IntegerLiteral); ok {
					violations = append(violations, Violation{
						KataID:  "ZC1021",
						Message: "Use symbolic permissions with `chmod` instead of octal.",
						Line:    name.Token.Line,
						Column:  name.Token.Column,
						Level:   SeverityStyle,
					})
					break
				}
			}
		}
	}

	return violations
}

func init() {
	RegisterKata(ast.LetStatementNode, Kata{
		ID:    "ZC1022",
		Title: "Use `$((...))` for arithmetic expansion",
		Description: "The `$((...))` syntax is the modern, recommended way to perform arithmetic expansion. " +
			"It is more readable and can be nested easily, unlike `let`.",
		Severity: SeverityStyle,
		Check:    checkZC1022,
		// Reuse ZC1013's `let NAME=EXPR` → `(( NAME = EXPR ))` rewrite.
		// For a standalone arithmetic statement the `(( ))` command
		// form is the right shape; the `$((...))` text in the message
		// reads as the broader "use Zsh arithmetic" recommendation.
		Fix: fixZC1013,
	})
}

func checkZC1022(node ast.Node) []Violation {
	violations := []Violation{}

	if let, ok := node.(*ast.LetStatement); ok {
		violations = append(violations, Violation{
			KataID:  "ZC1022",
			Message: "Use `$((...))` for arithmetic expansion instead of `let`.",
			Line:    let.Token.Line,
			Column:  let.Token.Column,
			Level:   SeverityStyle,
		})
	}

	return violations
}

// Per issue #345: ZC1023–ZC1029, ZC1033, ZC1035 were identical copies of
// the canonical ZC1022 `let` detection, each emitting the same violation
// on the same input. Ten fires per `let` line inflated user-visible noise.
// Project rule ("once committed, fix — don't remove") keeps the IDs alive
// as no-op stubs so existing `disabled_katas` lists that reference them
// keep parsing, but the duplicate detection is retired in favour of
// ZC1022.

func init() {
	RegisterKata(ast.LetStatementNode, Kata{
		ID:          "ZC1023",
		Title:       "Superseded by ZC1022 — retired duplicate `let` detector",
		Severity:    SeverityStyle,
		Description: "Retained as a no-op stub so legacy `.zshellcheckrc` files that disable this ID keep parsing. The canonical `let` → `$((...))` guidance lives in ZC1022; see https://github.com/afadesigns/zshellcheck/issues/345.",
		Check:       checkZC1023,
	})
}

func checkZC1023(ast.Node) []Violation {
	return nil
}

// Per issue #345: ZC1023–ZC1029, ZC1033, ZC1035 were identical copies of
// the canonical ZC1022 `let` detection, each emitting the same violation
// on the same input. Ten fires per `let` line inflated user-visible noise.
// Project rule ("once committed, fix — don't remove") keeps the IDs alive
// as no-op stubs so existing `disabled_katas` lists that reference them
// keep parsing, but the duplicate detection is retired in favour of
// ZC1022.

func init() {
	RegisterKata(ast.LetStatementNode, Kata{
		ID:          "ZC1024",
		Title:       "Superseded by ZC1022 — retired duplicate `let` detector",
		Severity:    SeverityStyle,
		Description: "Retained as a no-op stub so legacy `.zshellcheckrc` files that disable this ID keep parsing. The canonical `let` → `$((...))` guidance lives in ZC1022; see https://github.com/afadesigns/zshellcheck/issues/345.",
		Check:       checkZC1024,
	})
}

func checkZC1024(ast.Node) []Violation {
	return nil
}

// Per issue #345: ZC1023–ZC1029, ZC1033, ZC1035 were identical copies of
// the canonical ZC1022 `let` detection, each emitting the same violation
// on the same input. Ten fires per `let` line inflated user-visible noise.
// Project rule ("once committed, fix — don't remove") keeps the IDs alive
// as no-op stubs so existing `disabled_katas` lists that reference them
// keep parsing, but the duplicate detection is retired in favour of
// ZC1022.

func init() {
	RegisterKata(ast.LetStatementNode, Kata{
		ID:          "ZC1025",
		Title:       "Superseded by ZC1022 — retired duplicate `let` detector",
		Severity:    SeverityStyle,
		Description: "Retained as a no-op stub so legacy `.zshellcheckrc` files that disable this ID keep parsing. The canonical `let` → `$((...))` guidance lives in ZC1022; see https://github.com/afadesigns/zshellcheck/issues/345.",
		Check:       checkZC1025,
	})
}

func checkZC1025(ast.Node) []Violation {
	return nil
}

// Per issue #345: ZC1023–ZC1029, ZC1033, ZC1035 were identical copies of
// the canonical ZC1022 `let` detection, each emitting the same violation
// on the same input. Ten fires per `let` line inflated user-visible noise.
// Project rule ("once committed, fix — don't remove") keeps the IDs alive
// as no-op stubs so existing `disabled_katas` lists that reference them
// keep parsing, but the duplicate detection is retired in favour of
// ZC1022.

func init() {
	RegisterKata(ast.LetStatementNode, Kata{
		ID:          "ZC1026",
		Title:       "Superseded by ZC1022 — retired duplicate `let` detector",
		Severity:    SeverityStyle,
		Description: "Retained as a no-op stub so legacy `.zshellcheckrc` files that disable this ID keep parsing. The canonical `let` → `$((...))` guidance lives in ZC1022; see https://github.com/afadesigns/zshellcheck/issues/345.",
		Check:       checkZC1026,
	})
}

func checkZC1026(ast.Node) []Violation {
	return nil
}

// Per issue #345: ZC1023–ZC1029, ZC1033, ZC1035 were identical copies of
// the canonical ZC1022 `let` detection, each emitting the same violation
// on the same input. Ten fires per `let` line inflated user-visible noise.
// Project rule ("once committed, fix — don't remove") keeps the IDs alive
// as no-op stubs so existing `disabled_katas` lists that reference them
// keep parsing, but the duplicate detection is retired in favour of
// ZC1022.

func init() {
	RegisterKata(ast.LetStatementNode, Kata{
		ID:          "ZC1027",
		Title:       "Superseded by ZC1022 — retired duplicate `let` detector",
		Severity:    SeverityStyle,
		Description: "Retained as a no-op stub so legacy `.zshellcheckrc` files that disable this ID keep parsing. The canonical `let` → `$((...))` guidance lives in ZC1022; see https://github.com/afadesigns/zshellcheck/issues/345.",
		Check:       checkZC1027,
	})
}

func checkZC1027(ast.Node) []Violation {
	return nil
}

// Per issue #345: ZC1023–ZC1029, ZC1033, ZC1035 were identical copies of
// the canonical ZC1022 `let` detection, each emitting the same violation
// on the same input. Ten fires per `let` line inflated user-visible noise.
// Project rule ("once committed, fix — don't remove") keeps the IDs alive
// as no-op stubs so existing `disabled_katas` lists that reference them
// keep parsing, but the duplicate detection is retired in favour of
// ZC1022.

func init() {
	RegisterKata(ast.LetStatementNode, Kata{
		ID:          "ZC1028",
		Title:       "Superseded by ZC1022 — retired duplicate `let` detector",
		Severity:    SeverityStyle,
		Description: "Retained as a no-op stub so legacy `.zshellcheckrc` files that disable this ID keep parsing. The canonical `let` → `$((...))` guidance lives in ZC1022; see https://github.com/afadesigns/zshellcheck/issues/345.",
		Check:       checkZC1028,
	})
}

func checkZC1028(ast.Node) []Violation {
	return nil
}

// Per issue #345: ZC1023–ZC1029, ZC1033, ZC1035 were identical copies of
// the canonical ZC1022 `let` detection, each emitting the same violation
// on the same input. Ten fires per `let` line inflated user-visible noise.
// Project rule ("once committed, fix — don't remove") keeps the IDs alive
// as no-op stubs so existing `disabled_katas` lists that reference them
// keep parsing, but the duplicate detection is retired in favour of
// ZC1022.

func init() {
	RegisterKata(ast.LetStatementNode, Kata{
		ID:          "ZC1029",
		Title:       "Superseded by ZC1022 — retired duplicate `let` detector",
		Severity:    SeverityStyle,
		Description: "Retained as a no-op stub so legacy `.zshellcheckrc` files that disable this ID keep parsing. The canonical `let` → `$((...))` guidance lives in ZC1022; see https://github.com/afadesigns/zshellcheck/issues/345.",
		Check:       checkZC1029,
	})
}

func checkZC1029(ast.Node) []Violation {
	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1030",
		Title: "Use `printf` instead of `echo`",
		Description: "The `echo` command's behavior can be inconsistent across different shells and " +
			"environments, especially with flags and escape sequences. `printf` provides more reliable " +
			"and portable string formatting.",
		Severity: SeverityStyle,
		Check:    checkZC1030,
	})
}

func checkZC1030(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	if cmd.Name.TokenLiteral() != "echo" {
		return nil
	}

	// Defer to ZC1037 if any argument is a variable.
	for _, arg := range cmd.Arguments {
		if ident, ok := arg.(*ast.Identifier); ok {
			if ident.Token.Type == "VARIABLE" {
				return nil
			}
		}
	}

	return []Violation{
		{
			KataID:  "ZC1030",
			Message: "Use `printf` for more reliable and portable string formatting instead of `echo`.",
			Line:    cmd.Token.Line,
			Column:  cmd.Token.Column,
			Level:   SeverityStyle,
		},
	}
}

func init() {
	RegisterKata(ast.ShebangNode, Kata{
		ID:    "ZC1031",
		Title: "Use `#!/usr/bin/env zsh` for portability",
		Description: "Using `#!/usr/bin/env zsh` is more portable than `#!/bin/zsh` because it searches " +
			"for the `zsh` executable in the user's `PATH`.",
		Severity: SeverityInfo,
		Check:    checkZC1031,
		Fix:      fixZC1031,
	})
}

// fixZC1031 rewrites `#!/bin/zsh` to `#!/usr/bin/env zsh` in the
// shebang line. Span-aware: replaces the whole `#!/bin/zsh` run as a
// single edit at column 1, line 1.
func fixZC1031(node ast.Node, v Violation, source []byte) []FixEdit {
	shebang, ok := node.(*ast.Shebang)
	if !ok {
		return nil
	}
	if shebang.Path != "#!/bin/zsh" {
		return nil
	}
	return []FixEdit{{
		Line:    v.Line,
		Column:  v.Column,
		Length:  len("#!/bin/zsh"),
		Replace: "#!/usr/bin/env zsh",
	}}
}

func checkZC1031(node ast.Node) []Violation {
	violations := []Violation{}

	if shebang, ok := node.(*ast.Shebang); ok {
		if shebang.Path == "#!/bin/zsh" {
			violations = append(violations, Violation{
				KataID:  "ZC1031",
				Message: "Use `#!/usr/bin/env zsh` for portability instead of `#!/bin/zsh`.",
				Line:    1,
				Column:  1,
				Level:   SeverityInfo,
			})
		}
	}

	return violations
}

func init() {
	RegisterKata(ast.LetStatementNode, Kata{
		ID:    "ZC1032",
		Title: "Use `((...))` for C-style incrementing",
		Description: "Instead of `let i=i+1` or `let i=i-1`, you can use the more concise and idiomatic " +
			"C-style increment `(( i++ ))` / decrement `(( i-- ))` in Zsh.",
		Severity: SeverityStyle,
		Check:    checkZC1032,
		Fix:      fixZC1032,
	})
}

// zc1032Op classifies the let-statement value as either an increment
// (`name + 1`) or decrement (`name - 1`). Returns the C-style suffix
// (`++` / `--`) and true on match; empty string and false otherwise.
//
// The Zsh lexer treats `-` as an identifier byte, so `let i=i-1`
// parses as a single Identifier value `i-1` rather than an
// InfixExpression. The detector handles both shapes:
//   - InfixExpression: `i + 1` (operator `+`, integer literal 1).
//   - Identifier with `NAME-1` literal: synthetic decrement form.
func zc1032Op(stmt *ast.LetStatement) (string, bool) {
	if stmt == nil || stmt.Name == nil {
		return "", false
	}
	if infix, ok := stmt.Value.(*ast.InfixExpression); ok {
		leftIdent, ok := infix.Left.(*ast.Identifier)
		if !ok {
			return "", false
		}
		rightInt, ok := infix.Right.(*ast.IntegerLiteral)
		if !ok {
			return "", false
		}
		if stmt.Name.Value != leftIdent.Value || rightInt.Value != 1 {
			return "", false
		}
		switch infix.Operator {
		case "+":
			return "++", true
		case "-":
			return "--", true
		}
		return "", false
	}
	if ident, ok := stmt.Value.(*ast.Identifier); ok {
		want := stmt.Name.Value + "-1"
		if ident.Value == want {
			return "--", true
		}
	}
	return "", false
}

// fixZC1032 rewrites `let NAME=NAME+1` / `let NAME=NAME-1` into the
// C-style `(( NAME++ ))` / `(( NAME-- ))` arithmetic form. The
// replacement spans from the `let` keyword to the end of the logical
// line (first `;`, `\n`, or EOF). Bails when the AST does not match
// the increment/decrement shape so the fix stays unambiguous and
// idempotent on a re-run (the rewritten form is no longer a
// LetStatement).
func fixZC1032(node ast.Node, v Violation, source []byte) []FixEdit {
	stmt, ok := node.(*ast.LetStatement)
	if !ok {
		return nil
	}
	suffix, ok := zc1032Op(stmt)
	if !ok {
		return nil
	}
	start := LineColToByteOffset(source, v.Line, v.Column)
	if start < 0 {
		return nil
	}
	end := start
	for end < len(source) {
		b := source[end]
		if b == '\n' || b == ';' {
			break
		}
		end++
	}
	return []FixEdit{{
		Line:    v.Line,
		Column:  v.Column,
		Length:  end - start,
		Replace: "(( " + stmt.Name.Value + suffix + " ))",
	}}
}

func checkZC1032(node ast.Node) []Violation {
	letStmt, ok := node.(*ast.LetStatement)
	if !ok {
		return nil
	}
	suffix, ok := zc1032Op(letStmt)
	if !ok {
		return nil
	}
	msg := "Use `(( " + letStmt.Name.Value + suffix + " ))` for C-style "
	if suffix == "++" {
		msg += "incrementing instead of `let " + letStmt.Name.Value + "=" + letStmt.Name.Value + "+1`."
	} else {
		msg += "decrementing instead of `let " + letStmt.Name.Value + "=" + letStmt.Name.Value + "-1`."
	}
	return []Violation{{
		KataID:  "ZC1032",
		Message: msg,
		Line:    letStmt.Token.Line,
		Column:  letStmt.Token.Column,
		Level:   SeverityStyle,
	}}
}

// Per issue #345: ZC1023–ZC1029, ZC1033, ZC1035 were identical copies of
// the canonical ZC1022 `let` detection, each emitting the same violation
// on the same input. Ten fires per `let` line inflated user-visible noise.
// Project rule ("once committed, fix — don't remove") keeps the IDs alive
// as no-op stubs so existing `disabled_katas` lists that reference them
// keep parsing, but the duplicate detection is retired in favour of
// ZC1022.

func init() {
	RegisterKata(ast.LetStatementNode, Kata{
		ID:          "ZC1033",
		Title:       "Superseded by ZC1022 — retired duplicate `let` detector",
		Severity:    SeverityStyle,
		Description: "Retained as a no-op stub so legacy `.zshellcheckrc` files that disable this ID keep parsing. The canonical `let` → `$((...))` guidance lives in ZC1022; see https://github.com/afadesigns/zshellcheck/issues/345.",
		Check:       checkZC1033,
	})
}

func checkZC1033(ast.Node) []Violation {
	return nil
}

func init() {
	RegisterKata(ast.ExpressionStatementNode, Kata{
		ID:    "ZC1034",
		Title: "Use `command -v` instead of `which`",
		Description: "`which` is an external command and may not be available or consistent across all " +
			"systems. `command -v` is a POSIX standard and a shell builtin, making it more portable " +
			"and reliable for checking if a command exists.",
		Severity: SeverityStyle,
		Check:    checkZC1034,
		Fix:      fixZC1034,
	})
}

// fixZC1034 rewrites `which` to `command -v` at the command name
// position inside an ExpressionStatement. Single replacement — arguments
// stay untouched.
func fixZC1034(node ast.Node, v Violation, source []byte) []FixEdit {
	es, ok := node.(*ast.ExpressionStatement)
	if !ok {
		return nil
	}
	cmd, ok := es.Expression.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	if cmd.Name == nil || cmd.Name.TokenLiteral() != "which" {
		return nil
	}
	return []FixEdit{{
		Line:    v.Line,
		Column:  v.Column,
		Length:  len("which"),
		Replace: "command -v",
	}}
}

func checkZC1034(node ast.Node) []Violation {
	violations := []Violation{}

	if es, ok := node.(*ast.ExpressionStatement); ok {
		if cmd, ok := es.Expression.(*ast.SimpleCommand); ok {
			if cmd.Name.TokenLiteral() == "which" {
				violations = append(violations, Violation{
					KataID:  "ZC1034",
					Message: "Use `command -v` instead of `which` for portability.",
					Line:    cmd.Token.Line,
					Column:  cmd.Token.Column,
					Level:   SeverityStyle,
				})
			}
		}
	}

	return violations
}

// Per issue #345: ZC1023–ZC1029, ZC1033, ZC1035 were identical copies of
// the canonical ZC1022 `let` detection, each emitting the same violation
// on the same input. Ten fires per `let` line inflated user-visible noise.
// Project rule ("once committed, fix — don't remove") keeps the IDs alive
// as no-op stubs so existing `disabled_katas` lists that reference them
// keep parsing, but the duplicate detection is retired in favour of
// ZC1022.

func init() {
	RegisterKata(ast.LetStatementNode, Kata{
		ID:          "ZC1035",
		Title:       "Superseded by ZC1022 — retired duplicate `let` detector",
		Severity:    SeverityStyle,
		Description: "Retained as a no-op stub so legacy `.zshellcheckrc` files that disable this ID keep parsing. The canonical `let` → `$((...))` guidance lives in ZC1022; see https://github.com/afadesigns/zshellcheck/issues/345.",
		Check:       checkZC1035,
	})
}

func checkZC1035(ast.Node) []Violation {
	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1036",
		Title: "Prefer `[[ ... ]]` over `test` command",
		Description: "The `[[ ... ]]` construct is a more powerful and safer alternative to the `test` " +
			"command (or `[ ... ]`) for conditional expressions in modern shells. It handles word " +
			"splitting and globbing more intuitively and supports advanced features like regex matching.",
		Severity: SeverityStyle,
		Check:    checkZC1036,
		Fix:      fixZC1293,
	})
}

func checkZC1036(node ast.Node) []Violation {
	violations := []Violation{}

	if cmd, ok := node.(*ast.SimpleCommand); ok {
		if cmd.Name.TokenLiteral() == "test" {
			violations = append(violations, Violation{
				KataID:  "ZC1036",
				Message: "Prefer `[[ ... ]]` over `test` command for conditional expressions.",
				Line:    cmd.Token.Line,
				Column:  cmd.Token.Column,
				Level:   SeverityStyle,
			})
		}
	}

	return violations
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1037",
		Title: "Use 'print -r --' for variable expansion",
		Description: "Using 'echo' to print strings containing variables can lead to unexpected behavior " +
			"if the variable contains special characters or flags. A safer, more reliable alternative " +
			"is 'print -r --'.",
		Severity: SeverityStyle,
		Check:    checkZC1037,
		// Reuse ZC1092's `echo` → `print -r --` rewrite. The detector
		// here is a stricter subset (only fires when echo prints a
		// variable expansion) but the rewrite is identical; the
		// conflict resolver dedupes overlapping edits.
		Fix: fixZC1092,
	})
}

func checkZC1037(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	if cmd.Name.TokenLiteral() != "echo" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		if ident, ok := arg.(*ast.Identifier); ok && ident.Token.Type == token.VARIABLE {
			return []Violation{
				{
					KataID:  "ZC1037",
					Message: "Use 'print -r --' instead of 'echo' to reliably print variable expansions.",
					Line:    cmd.Token.Line,
					Column:  cmd.Token.Column,
					Level:   SeverityStyle,
				},
			}
		}
		if str, ok := arg.(*ast.StringLiteral); ok && strings.Contains(str.Value, "$") {
			return []Violation{
				{
					KataID:  "ZC1037",
					Message: "Use 'print -r --' instead of 'echo' to reliably print variable expansions.",
					Line:    cmd.Token.Line,
					Column:  cmd.Token.Column,
					Level:   SeverityStyle,
				},
			}
		}
	}

	return nil
}

func init() {
	RegisterKata(ast.InfixExpressionNode, Kata{
		ID:    "ZC1038",
		Title: "Avoid useless use of cat",
		Description: "Using `cat file | command` is unnecessary and inefficient. " +
			"Most commands can read from a file directly, e.g., `command file`. " +
			"If not, you can use input redirection: `command < file`.",
		Severity: SeverityStyle,
		Check:    checkZC1038,
	})
}

func checkZC1038(node ast.Node) []Violation {
	violations := []Violation{}

	infix, ok := node.(*ast.InfixExpression)
	if !ok {
		return violations
	}

	if infix.Operator != "|" {
		return violations
	}

	cmd, ok := infix.Left.(*ast.SimpleCommand)
	if !ok {
		return violations
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return violations
	}

	if ident.Value != "cat" {
		return violations
	}

	// cat must have exactly one argument to be considered a "useless use" in this context.
	// cat without args reads from stdin (valid pipe).
	// cat with multiple args concatenates (valid use).
	if len(cmd.Arguments) == 1 {
		violations = append(violations, Violation{
			KataID: "ZC1038",
			Message: "Avoid useless use of cat. " +
				"Prefer `command file` or `command < file` over `cat file | command`.",
			Line:   ident.Token.Line,
			Column: ident.Token.Column,
			Level:  SeverityStyle,
		})
	}

	return violations
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1039",
		Title: "Avoid `rm` with root path",
		Description: "Running `rm` on the root directory `/` is dangerous. " +
			"Ensure you are not deleting the entire filesystem.",
		Severity: SeverityWarning,
		Check:    checkZC1039,
	})
}

func checkZC1039(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	// Check if command is rm
	if cmdName, ok := cmd.Name.(*ast.Identifier); !ok || cmdName.Value != "rm" {
		return nil
	}

	violations := []Violation{}

	for _, arg := range cmd.Arguments {
		// Bare `/` argument arrives as an Identifier with Value "/"
		// after the SLASH prefix registration; quoted forms arrive
		// as StringLiteral. Cover both shapes.
		var val string
		var line, col int
		switch n := arg.(type) {
		case *ast.StringLiteral:
			val = strings.Trim(n.Value, "\"'")
			line, col = n.Token.Line, n.Token.Column
		case *ast.Identifier:
			val = n.Value
			line, col = n.Token.Line, n.Token.Column
		}
		if val == "/" {
			violations = append(violations, Violation{
				KataID:  "ZC1039",
				Message: "Avoid `rm` on the root directory `/`. This is highly dangerous.",
				Line:    line,
				Column:  col,
				Level:   SeverityWarning,
			})
		}
	}

	return violations
}

func init() {
	RegisterKata(ast.ForLoopStatementNode, Kata{
		ID:    "ZC1040",
		Title: "Use (N) nullglob qualifier for globs in loops",
		Description: "In Zsh, a glob that matches nothing (e.g., `*.txt`) will cause an error by default. " +
			"Use the `(N)` glob qualifier to make it null (empty) if no matches found, preventing the error.",
		Severity: SeverityStyle,
		Check:    checkZC1040,
		Fix:      fixZC1040,
	})
}

// fixZC1040 appends `(N)` after a glob pattern in a `for` loop item
// list, turning `for f in *.txt` into `for f in *.txt(N)` so an
// empty match produces an empty iterator instead of an error.
// Span scanning ends at the first unescaped whitespace / delimiter.
func fixZC1040(_ ast.Node, v Violation, source []byte) []FixEdit {
	start := LineColToByteOffset(source, v.Line, v.Column)
	if start < 0 || start >= len(source) {
		return nil
	}
	argLen := unquotedArgLen(source, start)
	if argLen == 0 {
		return nil
	}
	end := start + argLen
	endLine, endCol := offsetLineColZC1040(source, end)
	if endLine < 0 {
		return nil
	}
	return []FixEdit{{
		Line:    endLine,
		Column:  endCol,
		Length:  0,
		Replace: "(N)",
	}}
}

func offsetLineColZC1040(source []byte, offset int) (int, int) {
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

func checkZC1040(node ast.Node) []Violation {
	loop, ok := node.(*ast.ForLoopStatement)
	if !ok {
		return nil
	}

	// Only check "for i in items..." style loops, not arithmetic loops
	if loop.Items == nil {
		return nil
	}

	violations := []Violation{}

	for _, item := range loop.Items {
		// We are looking for string literals that look like globs (contain *, ?, etc)
		// but do NOT contain (N) or (N-...) qualifiers.

		val := getStringValue(item)

		// If it is quoted, it is NOT a glob expansion.
		if len(val) > 0 && (val[0] == '"' || val[0] == '\'') {
			continue
		}

		if isGlob(val) && !hasNullGlobQualifier(val) {
			violations = append(violations, Violation{
				KataID: "ZC1040",
				Message: "Glob pattern '" + val + "' may error if no files match. " +
					"Append '(N)' to enable nullglob behavior: '" + val + "(N)'",
				Line:   item.TokenLiteralNode().Line,
				Column: item.TokenLiteralNode().Column,
				Level:  SeverityStyle,
			})
		}
	}

	return violations
}

func getStringValue(node ast.Node) string {
	switch n := node.(type) {
	case *ast.StringLiteral:
		return n.Value
	case *ast.ConcatenatedExpression:
		var sb strings.Builder
		for _, p := range n.Parts {
			sb.WriteString(getStringValue(p))
		}
		return sb.String()
	case *ast.Identifier:
		return n.Value
	case *ast.GroupedExpression:
		return "(" + getStringValue(n.Expression) + ")"
	case *ast.ArrayLiteral:
		var sb strings.Builder
		sb.WriteString("(")
		for i, el := range n.Elements {
			if i > 0 {
				sb.WriteString(" ")
			}
			sb.WriteString(getStringValue(el))
		}
		sb.WriteString(")")
		return sb.String()
	default:
		// Fallback for operators treated as literals (like *)
		return n.TokenLiteral()
	}
}

func isGlob(s string) bool {
	// Simple check for common glob characters
	return strings.ContainsAny(s, "*?[]")
}

func hasNullGlobQualifier(s string) bool {
	// Check for (N) at the end. Zsh qualifiers are at the end.
	// This is a naive check.
	return strings.Contains(s, "(N)") || strings.Contains(s, "(N") // (N) or (N...)
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1041",
		Title: "Do not use variables in printf format string",
		Description: "Using variables in `printf` format strings allows for format string attacks and unexpected behavior " +
			"if the variable contains `%`. Use `printf '%s' \"$var\"` instead.",
		Severity: SeverityStyle,
		Check:    checkZC1041,
	})
}

func checkZC1041(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	// Check if command is printf
	if cmdName, ok := cmd.Name.(*ast.Identifier); !ok || cmdName.Value != "printf" {
		return nil
	}

	if len(cmd.Arguments) == 0 {
		return nil
	}

	firstArg := cmd.Arguments[0]

	// The first argument should be a static StringLiteral.
	// If it is an Identifier ($var), ConcatenatedExpression ("$var"), or CommandSubstitution, warn.
	// Note: A StringLiteral might still contain interpolation if the lexer didn't split it,
	// but generally in this AST, StringLiteral is safe/static or single-quoted.
	// We warn if it's NOT a StringLiteral.

	_, isStringLiteral := firstArg.(*ast.StringLiteral)

	if !isStringLiteral {
		violations := []Violation{{
			KataID:  "ZC1041",
			Message: "Do not use variables in printf format string. Use `printf '..%s..' \"$var\"` instead.",
			Line:    firstArg.TokenLiteralNode().Line,
			Column:  firstArg.TokenLiteralNode().Column,
			Level:   SeverityStyle,
		}}
		return violations
	}

	// Even if it is a StringLiteral, it might be "$var" (interpolation).
	// We should inspect the value.
	if str, ok := firstArg.(*ast.StringLiteral); ok {
		val := str.Value
		// If it contains $ and is not single-quoted, it's likely a variable.
		// Heuristic: if it starts with " and contains $, it's risky.
		if len(val) > 0 && val[0] == '"' {
			// Check for $ not escaped? The lexer hands us the raw string usually.
			// Simple check: if it has unescaped $, flag it.
			// This is a basic heuristic.
			for i := 0; i < len(val); i++ {
				if val[i] == '$' && (i == 0 || val[i-1] != '\\') {
					violations := []Violation{{
						KataID:  "ZC1041",
						Message: "Do not use variables in printf format string. Use `printf '..%s..' \"$var\"` instead.",
						Line:    firstArg.TokenLiteralNode().Line,
						Column:  firstArg.TokenLiteralNode().Column,
						Level:   SeverityStyle,
					}}
					return violations
				}
			}
		}
	}

	return nil
}

func init() {
	RegisterKata(ast.ForLoopStatementNode, Kata{
		ID:    "ZC1042",
		Title: "Use \"$@\" to iterate over arguments",
		Description: "`$*` joins all arguments into a single string, which is rarely what you want in a loop. " +
			"Use `\"$@\"` to iterate over each argument individually.",
		Severity: SeverityStyle,
		Check:    checkZC1042,
	})
}

func checkZC1042(node ast.Node) []Violation {
	loop, ok := node.(*ast.ForLoopStatement)
	if !ok {
		return nil
	}

	// Only check "for i in items..." style loops
	if loop.Items == nil {
		return nil
	}

	violations := []Violation{}

	for _, item := range loop.Items {
		// Check for "$*" (quoted) or $* (unquoted)

		// Helper to get raw value structure would be useful, but let's inspect manually.
		// 1. Unquoted $* -> Identifier with Value="$*"
		// 2. Quoted "$*" -> StringLiteral (if handled by lexer as one token) or ConcatenatedExpression?

		// In our parser/lexer, variables inside quotes often result in StringLiteral if simple,
		// or if interpolated, we need to check the parts.
		// However, "$*" is special.

		found := false

		if ident, ok := item.(*ast.Identifier); ok {
			if ident.Value == "$*" {
				found = true
			}
		} else if str, ok := item.(*ast.StringLiteral); ok {
			// Check if it *contains* $* inside quotes.
			// Note: Our Lexer.readString now preserves quotes.
			// If Value is `"$"` that's bad.
			if strings.Contains(str.Value, "$*") {
				found = true
			}
		} else if concat, ok := item.(*ast.ConcatenatedExpression); ok {
			// Check parts for identifier $*
			for _, part := range concat.Parts {
				if ident, ok := part.(*ast.Identifier); ok && ident.Value == "$*" {
					found = true
					break
				}
				// Or string part containing it
				if str, ok := part.(*ast.StringLiteral); ok && strings.Contains(str.Value, "$*") {
					found = true
					break
				}
			}
		}

		if found {
			violations = append(violations, Violation{
				KataID: "ZC1042",
				Message: "Use \"$@\" instead of \"$*\" (or $*) to iterate over arguments. " +
					"\"$*\" merges arguments into a single string.",
				Line:   item.TokenLiteralNode().Line,
				Column: item.TokenLiteralNode().Column,
				Level:  SeverityStyle,
			})
		}
	}

	return violations
}

func init() {
	RegisterKata(ast.FunctionDefinitionNode, Kata{
		ID:    "ZC1043",
		Title: "Use `local` for variables in functions",
		Description: "Variables defined in functions are global by default in Zsh. " +
			"Use `local` to scope them to the function.",
		Severity: SeverityStyle,
		Check:    checkZC1043,
		Fix:      fixZC1043,
	})
}

// fixZC1043 prepends `local ` to the unscoped assignment the detector
// flagged. The violation's Line/Column points at the assignment LHS;
// inserting `local ` there yields `local NAME=value`. On re-run the
// detector recognises `local …` as a declaration and skips the line,
// so the rewrite is idempotent.
func fixZC1043(_ ast.Node, v Violation, source []byte) []FixEdit {
	off := LineColToByteOffset(source, v.Line, v.Column)
	if off < 0 || off >= len(source) {
		return nil
	}
	// Defensive idempotency guard: refuse to insert if a declaration
	// keyword already sits at the violation column.
	for _, prefix := range []string{"local ", "typeset ", "declare ", "integer ", "float ", "readonly "} {
		end := off + len(prefix)
		if end <= len(source) && string(source[off:end]) == prefix {
			return nil
		}
	}
	return []FixEdit{{
		Line:    v.Line,
		Column:  v.Column,
		Length:  0,
		Replace: "local ",
	}}
}

var zc1043LocalDecls = map[string]struct{}{
	"local": {}, "typeset": {}, "declare": {},
	"integer": {}, "float": {}, "readonly": {},
}

func checkZC1043(node ast.Node) []Violation {
	funcDef, ok := node.(*ast.FunctionDefinition)
	if !ok {
		return nil
	}
	violations := []Violation{}
	locals := make(map[string]bool)
	ast.Walk(funcDef.Body, func(n ast.Node) bool {
		if _, ok := n.(*ast.FunctionDefinition); ok && n != funcDef {
			return false
		}
		zc1043HarvestLocals(n, locals)
		if v, ok := zc1043UnscopedAssign(n, locals); ok {
			violations = append(violations, v)
		}
		return true
	})
	return violations
}

func zc1043HarvestLocals(n ast.Node, locals map[string]bool) {
	cmd, ok := n.(*ast.SimpleCommand)
	if !ok {
		return
	}
	if _, hit := zc1043LocalDecls[cmd.Name.String()]; !hit {
		return
	}
	for _, arg := range cmd.Arguments {
		argStr := arg.String()
		if argStr == "" || argStr[0] == '-' {
			continue
		}
		name := argStr
		for i, c := range argStr {
			if c == '=' {
				name = argStr[:i]
				break
			}
		}
		locals[name] = true
	}
}

func zc1043UnscopedAssign(n ast.Node, locals map[string]bool) (Violation, bool) {
	exprStmt, ok := n.(*ast.ExpressionStatement)
	if !ok {
		return Violation{}, false
	}
	assign, ok := exprStmt.Expression.(*ast.InfixExpression)
	if !ok || assign.Operator != "=" {
		return Violation{}, false
	}
	ident, ok := assign.Left.(*ast.Identifier)
	if !ok || locals[ident.Value] {
		return Violation{}, false
	}
	rhs := ""
	if assign.Right != nil {
		rhs = assign.Right.String()
	}
	return Violation{
		KataID: "ZC1043",
		Message: "Variable '" + ident.Value + "' is assigned without 'local'. It will be global. " +
			"Use `local " + ident.Value + "=" + rhs + "`.",
		Line:   ident.Token.Line,
		Column: ident.Token.Column,
		Level:  SeverityStyle,
	}, true
}

func init() {
	RegisterKata(ast.ProgramNode, Kata{
		ID:    "ZC1044",
		Title: "Check for unchecked `cd` commands",
		Description: "`cd` failures should be handled to avoid executing commands in the wrong directory. " +
			"Use `cd ... || return` (or `exit`).",
		Severity: SeverityWarning,
		Check:    checkZC1044,
	})
}

func checkZC1044(node ast.Node) []Violation {
	// We only run on ProgramNode, but we do a full context-aware traversal.
	// Since main.go calls Check on every node, we might get called for Program.
	// But we might also want to ensure we don't double check if we traverse children.
	// Actually, if we register ONLY for ProgramNode, we are called once per file.
	// BUT, standard ast.Walk visits all nodes and calls Check.
	// If Check returns violations, they are added.
	// So if we implement a full walker here, it works.

	violations := []Violation{}

	walkZC1044(node, false, &violations)

	return violations
}

func walkZC1044(node ast.Node, isChecked bool, violations *[]Violation) {
	if node == nil {
		return
	}

	// Handle nil pointers inside interface
	if v := reflect.ValueOf(node); v.Kind() == reflect.Ptr && v.IsNil() {
		return
	}

	switch n := node.(type) {
	case *ast.Program:
		for _, stmt := range n.Statements {
			walkZC1044(stmt, false, violations)
		}
	case *ast.BlockStatement:
		for i, stmt := range n.Statements {
			check := false
			// Only the last statement in a block inherits the "checked" status of the block itself
			// (e.g. the last command in an `if` condition determines the result).
			// Previous statements are unchecked unless they have their own handling.
			if isChecked && i == len(n.Statements)-1 {
				check = true
			}
			walkZC1044(stmt, check, violations)
		}
	case *ast.GroupedExpression:
		walkZC1044(n.Expression, isChecked, violations)
	case *ast.IfStatement:
		walkZC1044(n.Condition, true, violations) // Condition is checked
		walkZC1044(n.Consequence, false, violations)
		walkZC1044(n.Alternative, false, violations)
	case *ast.WhileLoopStatement:
		walkZC1044(n.Condition, true, violations) // Condition is checked
		walkZC1044(n.Body, false, violations)
	case *ast.ForLoopStatement:
		walkZC1044(n.Init, false, violations)
		walkZC1044(n.Condition, true, violations) // Loop condition checked (arithmetic)
		walkZC1044(n.Post, false, violations)
		for _, item := range n.Items {
			walkZC1044(item, false, violations)
		}
		walkZC1044(n.Body, false, violations)
	case *ast.ExpressionStatement:
		walkZC1044(n.Expression, isChecked, violations)
	case *ast.InfixExpression:
		switch n.Operator {
		case "||":
			walkZC1044(n.Left, true, violations) // Left checked by Right
			walkZC1044(n.Right, isChecked, violations)
		case "&&":
			walkZC1044(n.Left, isChecked, violations) // Left inherits check
			walkZC1044(n.Right, isChecked, violations)
		default:
			walkZC1044(n.Left, false, violations)
			walkZC1044(n.Right, false, violations)
		}
	case *ast.PrefixExpression:
		if n.Operator == "!" {
			walkZC1044(n.Right, true, violations) // Negation checks it
		} else {
			walkZC1044(n.Right, false, violations)
		}
	case *ast.FunctionDefinition:
		walkZC1044(n.Body, false, violations)
	case *ast.SimpleCommand:
		checkCommandZC1044(n, isChecked, violations)
		// Also walk args? Arguments might contain subshells etc.
		for _, arg := range n.Arguments {
			walkZC1044(arg, false, violations)
		}
	case *ast.CommandSubstitution:
		walkZC1044(n.Command, false, violations) // Inner command starts unchecked

	// Recursion for other nodes
	default:
		// Use generic walk or manual?
		// Generic walk doesn't pass state.
		// We must implement all nodes that contain statements/expressions.
		// ...
	}
}

func checkCommandZC1044(cmd *ast.SimpleCommand, isChecked bool, violations *[]Violation) {
	if isChecked {
		return
	}
	if name, ok := cmd.Name.(*ast.Identifier); ok && name.Value == "cd" {
		*violations = append(*violations, Violation{
			KataID:  "ZC1044",
			Message: "Use `cd ... || return` (or `exit`) in case cd fails.",
			Line:    name.Token.Line,
			Column:  name.Token.Column,
			Level:   SeverityWarning,
		})
	}
}

func init() {
	kata := Kata{
		ID:    "ZC1045",
		Title: "Declare and assign separately to avoid masking return values",
		Description: "Declaring a variable with `local var=$(cmd)` masks the return value of `cmd`. " +
			"The `local` command returns 0 (success) even if `cmd` fails. " +
			"Declare the variable first (`local var`), then assign it (`var=$(cmd)`).",
		Severity: SeverityInfo,
		Check:    checkZC1045,
	}
	RegisterKata(ast.SimpleCommandNode, kata)
	RegisterKata(ast.DeclarationStatementNode, kata)
}

func checkZC1045(node ast.Node) []Violation {
	violations := []Violation{}

	// Check SimpleCommand (local, readonly)
	if cmd, ok := node.(*ast.SimpleCommand); ok {
		name := cmd.Name.String()
		if name == "local" || name == "readonly" {
			for _, arg := range cmd.Arguments {
				if hasCommandSubstitutionAssignment(arg) {
					violations = append(violations, Violation{
						KataID: "ZC1045",
						Message: "Declare and assign separately to avoid masking return values. " +
							"`" + name + " var=$(cmd)` masks the exit code of `cmd`.",
						Line:   arg.TokenLiteralNode().Line,
						Column: arg.TokenLiteralNode().Column,
						Level:  SeverityInfo,
					})
				}
			}
		}
	}

	// Check DeclarationStatement (typeset, declare)
	if decl, ok := node.(*ast.DeclarationStatement); ok {
		// Command is "typeset" or "declare"
		for _, assign := range decl.Assignments {
			if assign.Value != nil && isCommandSubstitution(assign.Value) {
				violations = append(violations, Violation{
					KataID: "ZC1045",
					Message: "Declare and assign separately to avoid masking return values. " +
						"`" + decl.Command + " var=$(cmd)` masks the exit code of `cmd`.",
					Line:   decl.Token.Line,
					Column: decl.Token.Column,
					Level:  SeverityInfo,
				})
			}
		}
	}

	return violations
}

func hasCommandSubstitutionAssignment(arg ast.Expression) bool {
	// Argument structure depends on parsing.
	// Usually ConcatenatedExpression for `var=$(cmd)`: [Identifier(var), StringLiteral(=), DollarParenExpression]
	// Or `var=`cmd``: [Identifier(var), StringLiteral(=), CommandSubstitution]

	concat, ok := arg.(*ast.ConcatenatedExpression)
	if !ok {
		return false
	}

	hasEquals := false
	hasCmdSubst := false

	for _, part := range concat.Parts {
		if str, ok := part.(*ast.StringLiteral); ok && str.Value == "=" {
			hasEquals = true
			continue
		}

		if hasEquals {
			// Check if RHS has command substitution
			if isCommandSubstitution(part) {
				hasCmdSubst = true
			}
		}
	}

	return hasEquals && hasCmdSubst
}

func isCommandSubstitution(node ast.Node) bool {
	switch n := node.(type) {
	case *ast.CommandSubstitution, *ast.DollarParenExpression:
		return true
	case *ast.ConcatenatedExpression:
		return zc1045ConcatHasSub(n)
	case *ast.StringLiteral:
		return zc1045StringHasSub(n.Value)
	}
	return false
}

func zc1045ConcatHasSub(n *ast.ConcatenatedExpression) bool {
	for _, p := range n.Parts {
		if isCommandSubstitution(p) {
			return true
		}
	}
	return false
}

// zc1045StringHasSub scans a double-quoted string literal for embedded
// `$(...)` or backtick command substitutions, ignoring backslash-
// escaped bytes.
func zc1045StringHasSub(val string) bool {
	if len(val) < 2 || val[0] != '"' || val[len(val)-1] != '"' {
		return false
	}
	for i := 0; i < len(val); i++ {
		switch val[i] {
		case '\\':
			i++
		case '`':
			return true
		case '$':
			if i+1 < len(val) && val[i+1] == '(' {
				return true
			}
		}
	}
	return false
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1046",
		Title: "Avoid `eval`",
		Description: "`eval` is dangerous as it executes arbitrary code. " +
			"Use arrays, parameter expansion, or other constructs instead.",
		Severity: SeverityWarning,
		Check:    checkZC1046,
	})
}

func checkZC1046(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	name := cmd.Name.String()

	// Check for direct 'eval'
	if name == "eval" {
		return []Violation{{
			KataID:  "ZC1046",
			Message: "Avoid `eval`. It allows execution of arbitrary code and is hard to debug.",
			Line:    cmd.Token.Line,
			Column:  cmd.Token.Column,
			Level:   SeverityWarning,
		}}
	}

	// Check for 'builtin eval' or 'command eval'
	if (name == "builtin" || name == "command") && len(cmd.Arguments) > 0 {
		arg := cmd.Arguments[0]
		if arg.String() == "eval" {
			return []Violation{{
				KataID:  "ZC1046",
				Message: "Avoid `eval`. It allows execution of arbitrary code and is hard to debug.",
				Line:    arg.TokenLiteralNode().Line,
				Column:  arg.TokenLiteralNode().Column,
				Level:   SeverityWarning,
			}}
		}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1047",
		Title: "Avoid `sudo` in scripts",
		Description: "Using `sudo` in scripts is generally discouraged. It makes the script interactive and less portable. " +
			"Run the script as root or use `sudo` to invoke the script.",
		Severity: SeverityWarning,
		Check:    checkZC1047,
	})
}

func checkZC1047(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	// Check if command is sudo
	if name, ok := cmd.Name.(*ast.Identifier); ok && name.Value == "sudo" {
		return []Violation{{
			KataID:  "ZC1047",
			Message: "Avoid `sudo` in scripts. Run the entire script as root if privileges are required.",
			Line:    name.Token.Line,
			Column:  name.Token.Column,
			Level:   SeverityWarning,
		}}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1048",
		Title: "Avoid `source` with relative paths",
		Description: "Sourcing a file with a relative path (e.g. `source ./lib.zsh`) depends on the current " +
			"working directory. Use `${0:a:h}/lib.zsh` to source relative to the script location.",
		Severity: SeverityStyle,
		Check:    checkZC1048,
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
	val := getStringValue(arg)

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
			Level:   SeverityStyle,
		}}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1049",
		Title: "Prefer functions over aliases",
		Description: "Aliases are expanded at parse time and can be confusing in scripts. " +
			"Use functions for more predictable behavior.",
		Severity: SeverityStyle,
		Check:    checkZC1049,
	})
}

func checkZC1049(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	// Check if command is alias
	if name, ok := cmd.Name.(*ast.Identifier); ok && name.Value == "alias" {
		return []Violation{{
			KataID: "ZC1049",
			Message: "Prefer functions over aliases. " +
				"Aliases are expanded at parse time and can behave unexpectedly in scripts.",
			Line:   name.Token.Line,
			Column: name.Token.Column,
			Level:  SeverityStyle,
		}}
	}

	return nil
}

func init() {
	RegisterKata(ast.ForLoopStatementNode, Kata{
		ID:    "ZC1050",
		Title: "Avoid iterating over `ls` output",
		Description: "Iterating over `ls` output is fragile because filenames can contain spaces and newlines. " +
			"Use globs (e.g. `for f in *.txt`) instead.",
		Severity: SeverityStyle,
		Check:    checkZC1050,
	})
}

func checkZC1050(node ast.Node) []Violation {
	loop, ok := node.(*ast.ForLoopStatement)
	if !ok {
		return nil
	}

	// Check loop items
	if loop.Items == nil {
		return nil
	}

	violations := []Violation{}

	for _, item := range loop.Items {
		// Check for $(ls ...) or `ls ...`
		cmd := getCommandFromSubstitution(item)
		if cmd != nil {
			if simpleCmd, ok := cmd.(*ast.SimpleCommand); ok {
				if name, ok := simpleCmd.Name.(*ast.Identifier); ok && name.Value == "ls" {
					violations = append(violations, Violation{
						KataID: "ZC1050",
						Message: "Avoid iterating over `ls` output. " +
							"Use globs (e.g. `*.txt`) to handle filenames with spaces correctly.",
						Line:   item.TokenLiteralNode().Line,
						Column: item.TokenLiteralNode().Column,
						Level:  SeverityStyle,
					})
				}
			}
		}
	}

	return violations
}

func getCommandFromSubstitution(node ast.Node) ast.Node {
	switch n := node.(type) {
	case *ast.CommandSubstitution:
		return n.Command
	case *ast.DollarParenExpression:
		return n.Command
	case *ast.ConcatenatedExpression:
		// Check if any part is a substitution of ls
		for _, part := range n.Parts {
			if cmd := getCommandFromSubstitution(part); cmd != nil {
				return cmd
			}
		}
	}
	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1051",
		Title: "Quote variables in `rm` to avoid globbing",
		Description: "`rm $VAR` is dangerous if `$VAR` contains spaces or glob characters. " +
			"Quote the variable (`rm \"$VAR\"`) to ensure safe deletion.",
		Severity: SeverityWarning,
		Check:    checkZC1051,
		Fix:      fixZC1051,
	})
}

// fixZC1051 wraps an unquoted `$VAR` argument in double-quotes.
// Two edits: one `"` before the arg and one after. Arg span is
// measured from source — we scan forward from the arg's token
// position until the first unescaped whitespace / delimiter,
// honouring `{…}` / `[…]` / `(…)` nesting so expansions like
// `${var[1]}`, `$(cmd)`, `${arr[@]}` stay whole.
func fixZC1051(_ ast.Node, v Violation, source []byte) []FixEdit {
	start := LineColToByteOffset(source, v.Line, v.Column)
	if start < 0 || start >= len(source) {
		return nil
	}
	argLen := unquotedArgLen(source, start)
	if argLen == 0 {
		return nil
	}
	endLine, endCol := offsetLineColZC1051(source, start+argLen)
	if endLine < 0 {
		return nil
	}
	return []FixEdit{
		{Line: v.Line, Column: v.Column, Length: 0, Replace: `"`},
		{Line: endLine, Column: endCol, Length: 0, Replace: `"`},
	}
}

// unquotedArgLen returns the byte length of a shell-word starting
// at offset. Honours brace / paren / bracket nesting so
// `${arr[$i]}` and `$(cmd (sub))` stay whole, and stops on the
// first top-level delimiter.
func unquotedArgLen(source []byte, offset int) int {
	if offset < 0 || offset >= len(source) {
		return 0
	}
	st := unquotedArgScan{}
	for n := 0; offset+n < len(source); n++ {
		c := source[offset+n]
		if st.atTopLevel() && unquotedArgIsTerminator(c) {
			return n
		}
		st.advance(c)
	}
	return len(source) - offset
}

type unquotedArgScan struct {
	braceDepth, parenDepth, bracketDepth int
}

func (s *unquotedArgScan) atTopLevel() bool {
	return s.braceDepth == 0 && s.parenDepth == 0 && s.bracketDepth == 0
}

func (s *unquotedArgScan) advance(c byte) {
	switch c {
	case '{':
		s.braceDepth++
	case '}':
		if s.braceDepth > 0 {
			s.braceDepth--
		}
	case '(':
		s.parenDepth++
	case ')':
		if s.parenDepth > 0 {
			s.parenDepth--
		}
	case '[':
		s.bracketDepth++
	case ']':
		if s.bracketDepth > 0 {
			s.bracketDepth--
		}
	}
}

func unquotedArgIsTerminator(c byte) bool {
	switch c {
	case ' ', '\t', '\n', ';', '|', '&', '>', '<', ')':
		return true
	}
	return false
}

func offsetLineColZC1051(source []byte, offset int) (int, int) {
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

func checkZC1051(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	// Check if command is rm
	if name, ok := cmd.Name.(*ast.Identifier); !ok || name.Value != "rm" {
		return nil
	}

	violations := []Violation{}

	for _, arg := range cmd.Arguments {
		isUnquoted := false

		switch n := arg.(type) {
		case *ast.Identifier:
			// $VAR
			if len(n.Value) > 0 && n.Value[0] == '$' {
				isUnquoted = true
			}
		case *ast.PrefixExpression:
			// $var (if parsed as prefix)
			if n.Operator == "$" {
				isUnquoted = true
			}
		case *ast.ArrayAccess:
			// ${var[...]} unquoted
			// Zsh DOES NOT split unquoted variable expansions by default!
			// BUT it DOES glob them.
			// `rm $var`. If var="a b", it deletes "a b" (one file).
			// If var="*", it expands to all files.
			// So checking for globbing safety is key.
			// `rm \"$var\"` prevents globbing.
			isUnquoted = true
		case *ast.DollarParenExpression:
			// $(...)
			isUnquoted = true
		}

		if isUnquoted {
			violations = append(violations, Violation{
				KataID:  "ZC1051",
				Message: "Unquoted variable in `rm`. Quote it to prevent globbing (e.g. `rm \"$VAR\"`).",
				Line:    arg.TokenLiteralNode().Line,
				Column:  arg.TokenLiteralNode().Column,
				Level:   SeverityWarning,
			})
		}
	}

	return violations
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1052",
		Title: "Avoid `sed -i` for portability",
		Description: "`sed -i` usage varies between GNU/Linux and macOS/BSD. " +
			"macOS requires an extension argument (e.g. `sed -i ''`), while GNU does not. " +
			"Use a temporary file and `mv`, or `perl -i`, for portability.",
		Severity: SeverityStyle,
		Check:    checkZC1052,
	})
}

func checkZC1052(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	if name, ok := cmd.Name.(*ast.Identifier); !ok || name.Value != "sed" {
		return nil
	}

	violations := []Violation{}

	for _, arg := range cmd.Arguments {
		argStr := arg.String()
		argStr = strings.Trim(argStr, "\"'")
		if strings.HasPrefix(argStr, "-") {
			if argStr == "-i" {
				violations = append(violations, Violation{
					KataID:  "ZC1052",
					Message: "`sed -i` is non-portable (GNU vs BSD differences). Use `perl -i` or a temporary file.",
					Line:    arg.TokenLiteralNode().Line,
					Column:  arg.TokenLiteralNode().Column,
					Level:   SeverityStyle,
				})
			}
		}
	}

	return violations
}

func init() {
	RegisterKata(ast.IfStatementNode, Kata{
		ID:    "ZC1053",
		Title: "Silence `grep` output in conditions",
		Description: "Using `grep` in a condition prints matches to stdout. " +
			"Use `grep -q` (or `> /dev/null`) to silence output if you only care about the exit code.",
		Severity: SeverityStyle,
		Check:    checkZC1053,
		Fix:      fixZC1053,
	})
	RegisterKata(ast.WhileLoopStatementNode, Kata{
		ID:    "ZC1053",
		Title: "Silence `grep` output in conditions",
		Description: "Using `grep` in a condition prints matches to stdout. " +
			"Use `grep -q` (or `> /dev/null`) to silence output if you only care about the exit code.",
		Severity: SeverityStyle,
		Check:    checkZC1053,
		Fix:      fixZC1053,
	})
}

// fixZC1053 inserts ` -q` directly after the grep / egrep / fgrep /
// zgrep command name reported at the violation column. Idempotent —
// once `-q` is present the detector's hasQuiet check short-circuits
// and the kata no longer fires. Defensive byte-match guard refuses
// to insert unless the source at the offset is one of the recognised
// grep variants followed by whitespace.
func fixZC1053(_ ast.Node, v Violation, source []byte) []FixEdit {
	off := LineColToByteOffset(source, v.Line, v.Column)
	if off < 0 {
		return nil
	}
	var name string
	for _, n := range []string{"grep", "egrep", "fgrep", "zgrep"} {
		end := off + len(n)
		if end > len(source) {
			continue
		}
		if string(source[off:end]) != n {
			continue
		}
		// Boundary: next byte must be whitespace, newline, or end of file.
		if end < len(source) {
			c := source[end]
			if c != ' ' && c != '\t' && c != '\n' {
				continue
			}
		}
		name = n
		break
	}
	if name == "" {
		return nil
	}
	insertAt := off + len(name)
	line, col := offsetLineColZC1053(source, insertAt)
	if line < 0 {
		return nil
	}
	return []FixEdit{{
		Line:    line,
		Column:  col,
		Length:  0,
		Replace: " -q",
	}}
}

func offsetLineColZC1053(source []byte, offset int) (int, int) {
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

func checkZC1053(node ast.Node) []Violation {
	violations := []Violation{}

	var condition ast.Node

	switch n := node.(type) {
	case *ast.IfStatement:
		condition = n.Condition
	case *ast.WhileLoopStatement:
		condition = n.Condition
	default:
		return nil
	}

	if condition == nil {
		return nil
	}

	walkZC1053(condition, false, &violations)

	return violations
}

func walkZC1053(node ast.Node, isSilenced bool, violations *[]Violation) {
	if node == nil {
		return
	}
	switch n := node.(type) {
	case *ast.BlockStatement:
		for _, stmt := range n.Statements {
			walkZC1053(stmt, isSilenced, violations)
		}
	case *ast.ExpressionStatement:
		walkZC1053(n.Expression, isSilenced, violations)
	case *ast.InfixExpression:
		zc1053WalkInfix(n, isSilenced, violations)
	case *ast.PrefixExpression:
		if n.Operator == "!" {
			walkZC1053(n.Right, isSilenced, violations)
		}
	case *ast.Redirection:
		walkZC1053(n.Left, isSilenced || zc1053SilencesStdout(n), violations)
	case *ast.SimpleCommand:
		checkCommandZC1053(n, isSilenced, violations)
	case *ast.GroupedExpression:
		walkZC1053(n.Expression, isSilenced, violations)
	}
}

func zc1053WalkInfix(n *ast.InfixExpression, isSilenced bool, violations *[]Violation) {
	if n.Operator == "|" {
		// Left side of pipe is silenced (stdout goes to pipe).
		walkZC1053(n.Left, true, violations)
		walkZC1053(n.Right, isSilenced, violations)
		return
	}
	walkZC1053(n.Left, isSilenced, violations)
	walkZC1053(n.Right, isSilenced, violations)
}

func zc1053SilencesStdout(n *ast.Redirection) bool {
	switch n.Operator {
	case ">", ">>", "&>":
		return isDevNull(n.Right)
	}
	return false
}

func checkCommandZC1053(cmd *ast.SimpleCommand, isSilenced bool, violations *[]Violation) {
	if isSilenced {
		return
	}

	if name, ok := cmd.Name.(*ast.Identifier); ok {
		if name.Value == "grep" || name.Value == "egrep" || name.Value == "fgrep" || name.Value == "zgrep" {
			// Check args for -q, --quiet, --silent
			hasQuiet := false
			for _, arg := range cmd.Arguments {
				argStr := arg.String()
				argStr = strings.Trim(argStr, "\"'")
				if strings.HasPrefix(argStr, "-") {
					if argStr == "-q" || argStr == "--quiet" || argStr == "--silent" {
						hasQuiet = true
						break
					}
					// Check for combined flags e.g. -rq
					if !strings.HasPrefix(argStr, "--") && strings.Contains(argStr, "q") {
						hasQuiet = true
						break
					}
				}
			}

			if !hasQuiet {
				*violations = append(*violations, Violation{
					KataID:  "ZC1053",
					Message: "Silence `grep` output in conditions. Use `grep -q` or redirect to `/dev/null`.",
					Line:    name.Token.Line,
					Column:  name.Token.Column,
					Level:   SeverityStyle,
				})
			}
		}
	}
}

func isDevNull(node ast.Node) bool {
	val := getStringValueZC1053(node)
	// Remove quotes
	if len(val) >= 2 && (val[0] == '"' || val[0] == '\'') {
		val = val[1 : len(val)-1]
	}
	return val == "/dev/null"
}

func getStringValueZC1053(node ast.Node) string {
	switch n := node.(type) {
	case *ast.StringLiteral:
		return n.Value
	case *ast.ConcatenatedExpression:
		var sb strings.Builder
		for _, p := range n.Parts {
			sb.WriteString(getStringValueZC1053(p))
		}
		return sb.String()
	}
	return ""
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:          "ZC1054",
		Title:       "Use POSIX classes in regex/glob",
		Description: "Ranges like `[a-z]` are locale-dependent. Use `[[:lower:]]` or `[a-z]` with `LC_ALL=C` to be explicit.",
		Severity:    SeverityStyle,
		Check:       checkZC1054,
	})
}

var rangeRegex = regexp.MustCompile(`\[[a-zA-Z0-9]-[a-zA-Z0-9]\]`)

func checkZC1054(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	violations := []Violation{}

	for _, arg := range cmd.Arguments {
		val := getStringValueZC1054(arg)
		if rangeRegex.MatchString(val) {
			// Avoid flagging if it looks like a POSIX class like [[:lower:]]
			// But regex `\[[a-z]-[a-z]\]` matches `[a-z]` but not `[[:lower:]]`
			// Wait, `[[:lower:]]` contains `[:` which is not `[a-z]-[a-z]`.
			// So it should be safe.

			violations = append(violations, Violation{
				KataID:  "ZC1054",
				Message: "Ranges like `[a-z]` are locale-dependent. Use POSIX classes like `[[:lower:]]` or `[[:digit:]]`.",
				Line:    arg.TokenLiteralNode().Line,
				Column:  arg.TokenLiteralNode().Column,
				Level:   SeverityStyle,
			})
		}
	}

	return violations
}

func getStringValueZC1054(node ast.Node) string {
	switch n := node.(type) {
	case *ast.StringLiteral:
		return n.Value
	case *ast.Identifier:
		return n.Value
	case *ast.PrefixExpression:
		// Reconstruct -z, !foo, etc.
		// Right is Expression.
		return n.Operator + getStringValueZC1054(n.Right)
	case *ast.ConcatenatedExpression:
		var sb strings.Builder
		for _, p := range n.Parts {
			sb.WriteString(getStringValueZC1054(p))
		}
		return sb.String()
	}
	return node.String() // Fallback to String() for other types
}

func init() {
	RegisterKata(ast.InfixExpressionNode, Kata{
		ID:    "ZC1055",
		Title: "Use `[[ -n/-z ]]` for empty string checks",
		Description: "Comparing with empty string is less idiomatic than using `[[ -z $var ]]` (is empty) " +
			"or `[[ -n $var ]]` (is not empty).",
		Severity: SeverityStyle,
		Check:    checkZC1055,
		Fix:      fixZC1055,
	})
}

func checkZC1055(node ast.Node) []Violation {
	expr, ok := node.(*ast.InfixExpression)
	if !ok {
		return nil
	}

	// Check for == "" or != ""
	if expr.Operator != "==" && expr.Operator != "!=" {
		return nil
	}

	// Check if either side is empty string literal
	isEmptyString := func(n ast.Node) bool {
		if str, ok := n.(*ast.StringLiteral); ok {
			// Check for "" or ''
			val := str.Value
			return val == `""` || val == `''`
		}
		return false
	}

	if isEmptyString(expr.Left) || isEmptyString(expr.Right) {
		opSuggestion := "-z"
		if expr.Operator == "!=" {
			opSuggestion = "-n"
		}

		return []Violation{{
			KataID:  "ZC1055",
			Message: "Use `[[ " + opSuggestion + " ... ]]` instead of comparing with empty string.",
			Line:    expr.TokenLiteralNode().Line,
			Column:  expr.TokenLiteralNode().Column,
			Level:   SeverityStyle,
		}}
	}

	return nil
}

// fixZC1055 rewrites `$var == ""` / `$var != ""` into `-z $var` /
// `-n $var` respectively. The span covers the full infix expression
// so `[[ $var == "" ]]` ends up as `[[ -z $var ]]`. Handles both
// left-side and right-side empty-string positions.
func fixZC1055(node ast.Node, v Violation, source []byte) []FixEdit {
	expr, ok := node.(*ast.InfixExpression)
	if !ok {
		return nil
	}
	if expr.Operator != "==" && expr.Operator != "!=" {
		return nil
	}
	isEmpty := func(n ast.Node) (bool, int) {
		str, ok := n.(*ast.StringLiteral)
		if !ok {
			return false, 0
		}
		if str.Value == `""` || str.Value == `''` {
			return true, len(str.Value)
		}
		return false, 0
	}

	var varNode ast.Node
	var emptyNode ast.Node
	var emptyLen int
	if ok, n := isEmpty(expr.Left); ok {
		varNode = expr.Right
		emptyNode = expr.Left
		emptyLen = n
	} else if ok, n := isEmpty(expr.Right); ok {
		varNode = expr.Left
		emptyNode = expr.Right
		emptyLen = n
	} else {
		return nil
	}

	varExpr, vok := varNode.(ast.Expression)
	emptyExpr, eok := emptyNode.(ast.Expression)
	if !vok || !eok {
		return nil
	}
	var varTok, emptyTok token.Token = varExpr.TokenLiteralNode(), emptyExpr.TokenLiteralNode()
	if varTok.Line == 0 || emptyTok.Line == 0 {
		return nil
	}
	varOffset := LineColToByteOffset(source, varTok.Line, varTok.Column)
	emptyOffset := LineColToByteOffset(source, emptyTok.Line, emptyTok.Column)
	if varOffset < 0 || emptyOffset < 0 {
		return nil
	}

	// Determine span start / end based on which side is the empty
	// literal. The span covers varNode + operator + emptyNode in
	// source order.
	start := varOffset
	if emptyOffset < start {
		start = emptyOffset
	}
	end := emptyOffset + emptyLen
	if varEnd := varOffset + identOrVarLen(source, varOffset); varEnd > end {
		end = varEnd
	}

	op := "-z"
	if expr.Operator == "!=" {
		op = "-n"
	}

	varText := string(source[varOffset : varOffset+identOrVarLen(source, varOffset)])
	line, col := v.Line, v.Column
	if startLine, startCol := byteOffsetToLineColZC1055(source, start); startLine > 0 {
		line, col = startLine, startCol
	}
	return []FixEdit{{
		Line:    line,
		Column:  col,
		Length:  end - start,
		Replace: op + " " + varText,
	}}
}

// identOrVarLen returns the byte length of an identifier or variable
// token that starts at offset. Variables may begin with `$`, `${`,
// or a plain identifier run. We scan until whitespace / delimiter so
// composite words like `$var.ext` stay together.
func identOrVarLen(source []byte, offset int) int {
	if offset < 0 || offset >= len(source) {
		return 0
	}
	n := 0
	depth := 0
	for offset+n < len(source) {
		c := source[offset+n]
		if c == ' ' || c == '\t' || c == '\n' {
			break
		}
		if depth == 0 {
			switch c {
			case ';', '|', '&', ')', ']', '}':
				return n
			}
		}
		if c == '{' || c == '(' {
			depth++
		} else if c == '}' || c == ')' {
			if depth > 0 {
				depth--
			}
		}
		n++
	}
	return n
}

func byteOffsetToLineColZC1055(source []byte, offset int) (int, int) {
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

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1056",
		Title: "Avoid `$((...))` as a statement",
		Description: "Using `$((...))` as a statement tries to execute the result as a command. " +
			"Use `((...))` for arithmetic evaluation/assignment.",
		Severity: SeverityStyle,
		Check:    checkZC1056,
	})
}

func checkZC1056(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	// Check if the command Name is a DollarParenExpression (arithmetic)
	var dpe *ast.DollarParenExpression

	if d, ok := cmd.Name.(*ast.DollarParenExpression); ok {
		dpe = d
	} else if concat, ok := cmd.Name.(*ast.ConcatenatedExpression); ok {
		if len(concat.Parts) == 1 {
			if d, ok := concat.Parts[0].(*ast.DollarParenExpression); ok {
				dpe = d
			}
		}
	}

	if dpe == nil {
		return nil
	}

	// Check if it is an arithmetic expression, not a command substitution.
	// Our parser distinguishes:
	// $(( ... )) -> Command is usually Infix/Prefix/Identifier/Integer/Grouped
	// $( ... )   -> Command is usually SimpleCommand (via parseCommandList)

	isArithmetic := true

	switch dpe.Command.(type) {
	case *ast.SimpleCommand:
		// $(cmd)
		isArithmetic = false
	case *ast.ConcatenatedExpression:
		// $(cmd arg)
		isArithmetic = false
	}

	if isArithmetic {
		return []Violation{{
			KataID:  "ZC1056",
			Message: "Avoid `$((...))` as a statement. It executes the result. Use `((...))` for arithmetic.",
			Line:    dpe.TokenLiteralNode().Line,
			Column:  dpe.TokenLiteralNode().Column,
			Level:   SeverityStyle,
		}}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1057",
		Title: "Avoid `ls` in assignments",
		Description: "Assigning the output of `ls` to a variable is fragile. " +
			"Use globs or arrays (e.g. `files=(*)`) to handle filenames correctly.",
		Severity: SeverityStyle,
		Check:    checkZC1057,
	})
	RegisterKata(ast.InfixExpressionNode, Kata{
		ID:    "ZC1057",
		Title: "Avoid `ls` in assignments",
		Description: "Assigning the output of `ls` to a variable is fragile. " +
			"Use globs or arrays (e.g. `files=(*)`) to handle filenames correctly.",
		Severity: SeverityStyle,
		Check:    checkZC1057,
	})
}

func checkZC1057(node ast.Node) []Violation {
	violations := []Violation{}

	checkAssignment := func(expr ast.Expression) {
		// Check if expr is an assignment to `ls` output
		// Usually ConcatenatedExpression: [Ident(var), String(=), DollarParen(ls)]
		if concat, ok := expr.(*ast.ConcatenatedExpression); ok {
			hasEquals := false
			for _, part := range concat.Parts {
				if str, ok := part.(*ast.StringLiteral); ok && str.Value == "=" {
					hasEquals = true
					continue
				}
				if hasEquals {
					// Check RHS for ls substitution
					if isLsSubstitution(part) {
						violations = append(violations, Violation{
							KataID:  "ZC1057",
							Message: "Avoid assigning `ls` output to a variable. Use globs (e.g. `files=(*)`) instead.",
							Line:    part.TokenLiteralNode().Line,
							Column:  part.TokenLiteralNode().Column,
							Level:   SeverityStyle,
						})
					}
				}
			}
		}
	}

	switch n := node.(type) {
	case *ast.SimpleCommand:
		checkAssignment(n.Name)
		for _, arg := range n.Arguments {
			checkAssignment(arg)
		}
	case *ast.InfixExpression:
		if n.Operator == "=" {
			if isLsSubstitution(n.Right) {
				violations = append(violations, Violation{
					KataID:  "ZC1057",
					Message: "Avoid assigning `ls` output to a variable. Use globs (e.g. `files=(*)`) instead.",
					Line:    n.TokenLiteralNode().Line,
					Column:  n.TokenLiteralNode().Column,
					Level:   SeverityStyle,
				})
			}
		}
	}

	return violations
}

func isLsSubstitution(node ast.Node) bool {
	// Reuse logic from ZC1050?
	// ZC1050 `getCommandFromSubstitution` returns the command.
	// We check if command is `ls`.

	var cmd ast.Node

	switch n := node.(type) {
	case *ast.CommandSubstitution:
		cmd = n.Command
	case *ast.DollarParenExpression:
		cmd = n.Command
	default:
		return false
	}

	// Check if cmd is `ls`
	// cmd can be SimpleCommand or Infix (pipeline).
	// If simple command `ls ...`
	if simple, ok := cmd.(*ast.SimpleCommand); ok {
		if name, ok := simple.Name.(*ast.Identifier); ok && name.Value == "ls" {
			return true
		}
	}
	return false
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1058",
		Title: "Avoid `sudo` with redirection",
		Description: "Redirecting output of `sudo` (e.g. `sudo cmd > /file`) fails if the current user " +
			"doesn't have permission. Use `| sudo tee /file` instead.",
		Severity: SeverityStyle,
		Check:    checkZC1058,
	})
}

func checkZC1058(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	if name, ok := cmd.Name.(*ast.Identifier); !ok || name.Value != "sudo" {
		return nil
	}

	violations := []Violation{}

	for _, arg := range cmd.Arguments {
		argStr := arg.String()
		argStr = strings.Trim(argStr, "\"'")
		if argStr == ">" || argStr == ">>" {
			violations = append(violations, Violation{
				KataID:  "ZC1058",
				Message: "Redirecting `sudo` output happens as the current user. Use `| sudo tee file` to write with privileges.",
				Line:    arg.TokenLiteralNode().Line,
				Column:  arg.TokenLiteralNode().Column,
				Level:   SeverityStyle,
			})
		}
	}

	return violations
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1059",
		Title: "Use `${var:?}` for `rm` arguments",
		Description: "Deleting a directory based on a variable is dangerous if the variable is empty or unset. " +
			"Use `${var:?}` to fail if empty, or check explicitly.",
		Severity: SeverityWarning,
		Check:    checkZC1059,
	})
}

func checkZC1059(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	if name, ok := cmd.Name.(*ast.Identifier); !ok || name.Value != "rm" {
		return nil
	}

	violations := []Violation{}

	for _, arg := range cmd.Arguments {
		isUnsafeVar := false

		switch n := arg.(type) {
		case *ast.PrefixExpression:
			if n.Operator == "$" {
				isUnsafeVar = true // $VAR
			}
		case *ast.ArrayAccess:
			// ${VAR}. Check if it has modifiers?
			// Parser for ArrayAccess currently parses ${VAR} as ArrayAccess with Index=nil.
			// It does NOT parse modifiers like :?.
			// If the source has ${VAR:?}, the parser might fail or parse differently?
			// Current parser implementation for ArrayAccess:
			// Expects IDENT. Then optional [index]. Then }.
			// It does NOT handle : modifiers.
			// So ${VAR:?} would likely fail parsing or be parsed incorrectly.
			// If parser fails, we can't check it.
			// Assuming parser parses simple ${VAR}, we flag it.
			isUnsafeVar = true
		case *ast.StringLiteral:
			// "$VAR".
			// If value is exactly "$VAR" or "${VAR}".
			// If value contains other things, it's safer (e.g. "$VAR/foo").
			// But "$VAR/" is dangerous too if VAR is empty.
			// For now, focus on exact variable.
			if isSimpleVariableString(n.Value) {
				isUnsafeVar = true
			}
		}

		if isUnsafeVar {
			violations = append(violations, Violation{
				KataID:  "ZC1059",
				Message: "Use `${var:?}` or ensure the variable is set before using it in `rm`.",
				Line:    arg.TokenLiteralNode().Line,
				Column:  arg.TokenLiteralNode().Column,
				Level:   SeverityWarning,
			})
		}
	}

	return violations
}

func isSimpleVariableString(s string) bool {
	// Check if string is "$VAR" or "${VAR}" (quoted)
	// Quotes are included in StringLiteral value.
	// "$VAR" -> len >= 4. e.g. "$V"
	if len(s) < 4 {
		return false
	}
	if s[0] != '"' || s[len(s)-1] != '"' {
		return false
	}
	inner := s[1 : len(s)-1]
	if len(inner) < 2 || inner[0] != '$' {
		return false
	}
	// Check if rest is valid identifier char (naive)
	// OR ${...}
	if inner[1] == '{' {
		// Must end with }
		if inner[len(inner)-1] != '}' {
			return false
		}
		return true // Assume ${...}
	}
	// $VAR
	return true
}

func init() {
	RegisterKata(ast.InfixExpressionNode, Kata{
		ID:    "ZC1060",
		Title: "Avoid `ps | grep` without exclusion",
		Description: "`ps | grep pattern` often matches the grep process itself. " +
			"Use `grep [p]attern`, `pgrep`, or exclude grep with `grep -v grep`.",
		Severity: SeverityStyle,
		Check:    checkZC1060,
	})
}

func checkZC1060(node ast.Node) []Violation {
	pipe, ok := node.(*ast.InfixExpression)
	if !ok || pipe.Operator != "|" {
		return nil
	}

	// Check if left command is `ps`
	if !isCommandName(pipe.Left, "ps") {
		return nil
	}

	// Check if right command is `grep`
	if !isCommandName(pipe.Right, "grep") {
		return nil
	}

	// Check if grep arguments exclude the grep process
	// Strategies:
	// 1. `grep -v grep` (chained pipe?)
	//    If pipe.Right is `grep`, we only see `grep ...`.
	//    If user does `ps | grep foo | grep -v grep`, the parsing structure is `(ps | grep foo) | grep -v grep`.
	//    So we are looking at `ps | grep foo`. The parent pipe handles the exclusion?
	//    We can't see the parent here easily.
	//    BUT, `ps | grep foo` is inherently risky unless `foo` uses `[]`.
	// 2. Pattern uses `[]`. e.g. `grep [f]oo`.

	// We inspect `grep` arguments.
	cmd, ok := pipe.Right.(*ast.SimpleCommand)
	if !ok {
		return nil // complex command
	}

	hasExclusion := false

	for _, arg := range cmd.Arguments {
		// Check if arg contains `[` and `]`.
		val := getStringValueZC1060(arg)
		// Naive check for `[...]` pattern
		// If it starts with - (flag), ignore unless it is -v grep?
		// But we only check THIS grep.

		if len(val) > 0 && val[0] != '-' {
			// Assume this is the pattern
			// Check for brackets
			for i := 0; i < len(val); i++ {
				if val[i] == '[' {
					// Look for closing ]
					for j := i + 1; j < len(val); j++ {
						if val[j] == ']' {
							hasExclusion = true
							break
						}
					}
				}
			}
		}
	}

	if !hasExclusion {
		return []Violation{{
			KataID:  "ZC1060",
			Message: "`ps | grep pattern` matches the grep process itself. Use `grep [p]attern` to exclude the grep process.",
			Line:    pipe.TokenLiteralNode().Line,
			Column:  pipe.TokenLiteralNode().Column,
			Level:   SeverityStyle,
		}}
	}

	return nil
}

func isCommandName(node ast.Node, name string) bool {
	if cmd, ok := node.(*ast.SimpleCommand); ok {
		if ident, ok := cmd.Name.(*ast.Identifier); ok {
			return ident.Value == name
		}
	}
	return false
}

func getStringValueZC1060(node ast.Node) string {
	switch n := node.(type) {
	case *ast.StringLiteral:
		return n.Value
	case *ast.ConcatenatedExpression:
		// Simplify
		return ""
	case *ast.Identifier:
		return n.Value
	}
	return ""
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:          "ZC1061",
		Title:       "Prefer `{start..end}` over `seq`",
		Description: "Using `seq` creates an external process. Zsh supports integer range expansion natively: `{1..10}`.",
		Severity:    SeverityStyle,
		Check:       checkZC1061,
		Fix:         fixZC1061,
	})
}

// fixZC1061 rewrites `seq M` / `seq M N` / `seq M S N` with integer
// literal arguments into Zsh's brace range expansion `{M..N}` or
// `{M..N..S}`. Forms with `-s sep` separator, variable arguments,
// or floats are left alone because the rewrite semantics differ.
func fixZC1061(node ast.Node, v Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "seq" {
		return nil
	}
	if len(cmd.Arguments) == 0 || len(cmd.Arguments) > 3 {
		return nil
	}
	nums := make([]string, 0, len(cmd.Arguments))
	for _, arg := range cmd.Arguments {
		s := arg.String()
		if !isAllDigits(s) {
			return nil
		}
		nums = append(nums, s)
	}
	var rng string
	switch len(nums) {
	case 1:
		rng = "{1.." + nums[0] + "}"
	case 2:
		rng = "{" + nums[0] + ".." + nums[1] + "}"
	case 3:
		rng = "{" + nums[0] + ".." + nums[2] + ".." + nums[1] + "}"
	}
	// Replace from the `seq` command name through the last argument
	// so the whole invocation becomes a single brace expansion.
	start := LineColToByteOffset(source, v.Line, v.Column)
	if start < 0 {
		return nil
	}
	lastArg := cmd.Arguments[len(cmd.Arguments)-1]
	lastTok := lastArg.TokenLiteralNode()
	lastOff := LineColToByteOffset(source, lastTok.Line, lastTok.Column)
	if lastOff < 0 {
		return nil
	}
	end := lastOff + len(lastTok.Literal)
	if end <= start {
		return nil
	}
	return []FixEdit{{
		Line:    v.Line,
		Column:  v.Column,
		Length:  end - start,
		Replace: rng,
	}}
}

func isAllDigits(s string) bool {
	if len(s) == 0 {
		return false
	}
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	return true
}

func checkZC1061(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	// Check if command is seq
	if name, ok := cmd.Name.(*ast.Identifier); ok && name.Value == "seq" {
		return []Violation{{
			KataID:  "ZC1061",
			Message: "Prefer `{start..end}` range expansion over `seq`. It is built-in and faster.",
			Line:    name.Token.Line,
			Column:  name.Token.Column,
			Level:   SeverityStyle,
		}}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:          "ZC1062",
		Title:       "Prefer `grep -E` over `egrep`",
		Description: "`egrep` is deprecated. Use `grep -E` instead.",
		Severity:    SeverityInfo,
		Check:       checkZC1062,
		Fix:         fixZC1062,
	})
}

// fixZC1062 rewrites `egrep` to `grep -E` at the command name
// position. Single replacement — arguments stay untouched.
func fixZC1062(node ast.Node, v Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "egrep" {
		return nil
	}
	return []FixEdit{{
		Line:    v.Line,
		Column:  v.Column,
		Length:  len("egrep"),
		Replace: "grep -E",
	}}
}

func checkZC1062(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	if name, ok := cmd.Name.(*ast.Identifier); ok && name.Value == "egrep" {
		return []Violation{{
			KataID:  "ZC1062",
			Message: "`egrep` is deprecated. Use `grep -E` instead.",
			Line:    name.Token.Line,
			Column:  name.Token.Column,
			Level:   SeverityInfo,
		}}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:          "ZC1063",
		Title:       "Prefer `grep -F` over `fgrep`",
		Description: "`fgrep` is deprecated. Use `grep -F` instead.",
		Severity:    SeverityInfo,
		Check:       checkZC1063,
		Fix:         fixZC1063,
	})
}

// fixZC1063 rewrites `fgrep` to `grep -F` at the command name
// position. Single replacement — arguments stay untouched.
func fixZC1063(node ast.Node, v Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "fgrep" {
		return nil
	}
	return []FixEdit{{
		Line:    v.Line,
		Column:  v.Column,
		Length:  len("fgrep"),
		Replace: "grep -F",
	}}
}

func checkZC1063(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	if name, ok := cmd.Name.(*ast.Identifier); ok && name.Value == "fgrep" {
		return []Violation{{
			KataID:  "ZC1063",
			Message: "`fgrep` is deprecated. Use `grep -F` instead.",
			Line:    name.Token.Line,
			Column:  name.Token.Column,
			Level:   SeverityInfo,
		}}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1064",
		Title: "Prefer `command -v` over `type`",
		Description: "`type` output format varies and is not POSIX standard for checking existence. " +
			"`command -v` is quieter and standard.",
		Severity: SeverityInfo,
		Check:    checkZC1064,
		Fix:      fixZC1064,
	})
}

// fixZC1064 rewrites `type` to `command -v` at the command name
// position. Single replacement — arguments stay untouched.
func fixZC1064(node ast.Node, v Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "type" {
		return nil
	}
	return []FixEdit{{
		Line:    v.Line,
		Column:  v.Column,
		Length:  len("type"),
		Replace: "command -v",
	}}
}

func checkZC1064(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	if name, ok := cmd.Name.(*ast.Identifier); ok && name.Value == "type" {
		return []Violation{{
			KataID:  "ZC1064",
			Message: "Prefer `command -v` over `type`. `type` output is not stable/standard for checking command existence.",
			Line:    name.Token.Line,
			Column:  name.Token.Column,
			Level:   SeverityInfo,
		}}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1065",
		Title: "Ensure spaces around `[` and `[[`",
		Description: "`[condition]` is parsed as a command named `[condition]`, which likely doesn't exist. " +
			"Add spaces: `[ condition ]`.",
		Severity: SeverityError,
		Check:    checkZC1065,
	})
	// Register for DoubleBracketExpression to check `[[foo]]`
	RegisterKata(ast.DoubleBracketExpressionNode, Kata{
		ID:          "ZC1065",
		Title:       "Ensure spaces around `[` and `[[`",
		Description: "`[[condition]]` is parsed incorrectly. Add spaces: `[[ condition ]]`.",
		Severity:    SeverityError,
		Check:       checkZC1065,
	})
}

func checkZC1065(node ast.Node) []Violation {
	violations := []Violation{}

	switch n := node.(type) {
	case *ast.SimpleCommand:
		if n.Name.String() == "[" {
			// Check first arg for preceding space
			if len(n.Arguments) > 0 {
				firstArg := n.Arguments[0]
				if !firstArg.TokenLiteralNode().HasPrecedingSpace {
					violations = append(violations, Violation{
						KataID:  "ZC1065",
						Message: "Missing space after `[`. Use `[ condition ]`.",
						Line:    n.Token.Line,
						Column:  n.Token.Column,
						Level:   SeverityError,
					})
				}
			}
		}
	case *ast.DoubleBracketExpression:
		// Check first expression
		if len(n.Elements) > 0 {
			firstExp := n.Elements[0]
			if !firstExp.TokenLiteralNode().HasPrecedingSpace {
				violations = append(violations, Violation{
					KataID:  "ZC1065",
					Message: "Missing space after `[[`. Use `[[ condition ]]`.",
					Line:    n.Token.Line,
					Column:  n.Token.Column,
					Level:   SeverityError,
				})
			}
		}
	}

	return violations
}

func init() {
	RegisterKata(ast.ForLoopStatementNode, Kata{
		ID:    "ZC1066",
		Title: "Avoid iterating over `cat` output",
		Description: "Iterating over `cat` output is fragile because lines can contain spaces. " +
			"Use `while IFS= read -r line; do ... done < file` or `($(<file))` array expansion.",
		Severity: SeverityStyle,
		Check:    checkZC1066,
	})
}

func checkZC1066(node ast.Node) []Violation {
	loop, ok := node.(*ast.ForLoopStatement)
	if !ok {
		return nil
	}

	if loop.Items == nil {
		return nil
	}

	violations := []Violation{}

	for _, item := range loop.Items {
		// Check for $(cat ...) or `cat ...`
		// Reuse ZC1050 logic but for `cat`
		cmd := getCommandFromSubstitutionZC1066(item)
		if cmd != nil {
			if simpleCmd, ok := cmd.(*ast.SimpleCommand); ok {
				if name, ok := simpleCmd.Name.(*ast.Identifier); ok && name.Value == "cat" {
					violations = append(violations, Violation{
						KataID:  "ZC1066",
						Message: "Avoid iterating over `cat` output. Use `while read` loop or `($(<file))` for line-based iteration.",
						Line:    item.TokenLiteralNode().Line,
						Column:  item.TokenLiteralNode().Column,
						Level:   SeverityStyle,
					})
				}
			}
		}
	}

	return violations
}

func getCommandFromSubstitutionZC1066(node ast.Node) ast.Node {
	switch n := node.(type) {
	case *ast.CommandSubstitution:
		return n.Command
	case *ast.DollarParenExpression:
		return n.Command
	case *ast.ConcatenatedExpression:
		// Check if any part is a substitution of cat
		for _, part := range n.Parts {
			if cmd := getCommandFromSubstitutionZC1066(part); cmd != nil {
				return cmd
			}
		}
	}
	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1067",
		Title: "Separate `export` and assignment to avoid masking return codes",
		Description: "Running `export var=$(cmd)` masks the return code of `cmd`. " +
			"The exit status will be that of `export` (usually 0). " +
			"Declare the variable first or export it after assignment.",
		Severity: SeverityStyle,
		Check:    checkZC1067,
	})
}

func checkZC1067(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	// Check if command is "export"
	name := cmd.Name.String()
	if name != "export" {
		return nil
	}

	violations := []Violation{}

	for _, arg := range cmd.Arguments {
		if containsSubstitutionAfterEquals(arg) {
			violations = append(violations, Violation{
				KataID: "ZC1067",
				Message: "Exporting and assigning a command substitution in one step masks the return value. " +
					"Use `var=$(cmd); export var`.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityStyle,
			})
		}
	}

	return violations
}

func containsSubstitutionAfterEquals(expr ast.Expression) bool {
	if stringIndex(expr.String(), "=") < 0 {
		return false
	}

	// Now check if it contains a command substitution
	return containsSubst(expr)
}

func stringIndex(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

func isSubstitution(n ast.Node) bool {
	switch n.(type) {
	case *ast.CommandSubstitution, *ast.DollarParenExpression:
		return true
	}
	return false
}

func containsSubst(n ast.Node) bool {
	found := false
	ast.Walk(n, func(node ast.Node) bool {
		if isSubstitution(node) {
			found = true
			return false
		}
		return true
	})
	return found
}

func init() {
	kata := Kata{
		ID:    "ZC1068",
		Title: "Use `add-zsh-hook` instead of defining hook functions directly",
		Description: "Defining special functions like `precmd`, `preexec`, `chpwd`, etc. directly overwrites any " +
			"previously defined hooks. Use `autoload -Uz add-zsh-hook; add-zsh-hook <hook> <function>` " +
			"to append to the hook list safely.",
		Severity: SeverityInfo,
		Check:    checkZC1068,
	}
	RegisterKata(ast.FunctionDefinitionNode, kata)
	RegisterKata(ast.FunctionLiteralNode, kata)
}

func checkZC1068(node ast.Node) []Violation {
	var name string
	var tokenLine, tokenCol int

	switch n := node.(type) {
	case *ast.FunctionDefinition:
		if n.Name == nil {
			return nil
		}
		name = n.Name.Value
		tokenLine = n.Token.Line
		tokenCol = n.Token.Column
	case *ast.FunctionLiteral:
		if n.Name == nil {
			return nil
		}
		name = n.Name.Value
		tokenLine = n.Token.Line
		tokenCol = n.Token.Column
	default:
		return nil
	}

	// List of special hook functions in Zsh
	specialHooks := map[string]bool{
		"precmd":             true,
		"preexec":            true,
		"chpwd":              true,
		"periodic":           true,
		"zshaddhistory":      true,
		"zshexit":            true,
		"zsh_directory_name": true,
	}

	if specialHooks[name] {
		return []Violation{
			{
				KataID: "ZC1068",
				Message: "Defining `" + name + "` directly overwrites existing hooks. " +
					"Use `autoload -Uz add-zsh-hook; add-zsh-hook " + name + " my_func` instead.",
				Line:   tokenLine,
				Column: tokenCol,
				Level:  SeverityInfo,
			},
		}
	}

	return nil
}

func init() {
	RegisterKata(ast.ProgramNode, Kata{
		ID:    "ZC1069",
		Title: "Avoid `local` outside of functions",
		Description: "The `local` builtin can only be used inside functions. " +
			"Using it in the global scope causes an error.",
		Severity: SeverityInfo,
		Check:    checkZC1069,
		Fix:      fixZC1069,
	})
}

// fixZC1069 rewrites `local` to `typeset` when used at file scope.
// `typeset` works in both function and global contexts, so the
// rewrite is safe wherever the detector fires. Single-edit name
// swap at the violation column. Idempotent — a re-run sees
// `typeset`, not `local`. Defensive byte-match guard.
func fixZC1069(_ ast.Node, v Violation, source []byte) []FixEdit {
	off := LineColToByteOffset(source, v.Line, v.Column)
	if off < 0 || off+len("local") > len(source) {
		return nil
	}
	if string(source[off:off+len("local")]) != "local" {
		return nil
	}
	return []FixEdit{{
		Line:    v.Line,
		Column:  v.Column,
		Length:  len("local"),
		Replace: "typeset",
	}}
}

func checkZC1069(node ast.Node) []Violation {
	program, ok := node.(*ast.Program)
	if !ok {
		return nil
	}
	w := zc1069Walker{}
	w.walk(program, false)
	return w.violations
}

type zc1069Walker struct {
	violations []Violation
}

func (w *zc1069Walker) walk(n ast.Node, inFunction bool) {
	if n == nil {
		return
	}
	w.recordIfBareLocal(n, inFunction)
	w.descendChildren(n, inFunction)
}

func (w *zc1069Walker) recordIfBareLocal(n ast.Node, inFunction bool) {
	cmd, ok := n.(*ast.SimpleCommand)
	if !ok {
		return
	}
	name, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return
	}
	if name.Value != "local" || inFunction {
		return
	}
	w.violations = append(w.violations, Violation{
		KataID: "ZC1069",
		Message: "`local` can only be used inside functions. " +
			"Use `typeset`, `declare`, or just assignment for global variables.",
		Line:   name.Token.Line,
		Column: name.Token.Column,
		Level:  SeverityInfo,
	})
}

func (w *zc1069Walker) descendChildren(n ast.Node, inFunction bool) {
	switch t := n.(type) {
	case *ast.Program:
		w.walkStatements(t.Statements, inFunction)
	case *ast.BlockStatement:
		w.walkStatements(t.Statements, inFunction)
	case *ast.IfStatement:
		w.walk(t.Condition, inFunction)
		w.walk(t.Consequence, inFunction)
		w.walk(t.Alternative, inFunction)
	case *ast.ForLoopStatement:
		w.walkForLoop(t, inFunction)
	case *ast.WhileLoopStatement:
		w.walk(t.Condition, inFunction)
		w.walk(t.Body, inFunction)
	case *ast.FunctionDefinition:
		w.walk(t.Name, inFunction)
		w.walk(t.Body, true)
	case *ast.FunctionLiteral:
		w.walkFunctionLiteral(t)
	case *ast.SimpleCommand:
		w.walk(t.Name, inFunction)
		w.walkExpressions(t.Arguments, inFunction)
	default:
		w.descendOtherChildren(n, inFunction)
	}
}

func (w *zc1069Walker) descendOtherChildren(n ast.Node, inFunction bool) {
	switch t := n.(type) {
	case *ast.ExpressionStatement:
		w.walk(t.Expression, inFunction)
	case *ast.InfixExpression:
		w.walk(t.Left, inFunction)
		w.walk(t.Right, inFunction)
	case *ast.PrefixExpression:
		w.walk(t.Right, inFunction)
	case *ast.PostfixExpression:
		w.walk(t.Left, inFunction)
	case *ast.GroupedExpression:
		w.walk(t.Expression, inFunction)
	case *ast.CaseStatement:
		w.walkCaseStatement(t, inFunction)
	case *ast.ConcatenatedExpression:
		w.walkExpressions(t.Parts, inFunction)
	case *ast.CommandSubstitution:
		w.walk(t.Command, inFunction)
	case *ast.DollarParenExpression:
		w.walk(t.Command, inFunction)
	case *ast.Subshell:
		w.walk(t.Command, inFunction)
	}
}

func (w *zc1069Walker) walkStatements(stmts []ast.Statement, inFunction bool) {
	for _, s := range stmts {
		w.walk(s, inFunction)
	}
}

func (w *zc1069Walker) walkExpressions(exprs []ast.Expression, inFunction bool) {
	for _, e := range exprs {
		w.walk(e, inFunction)
	}
}

func (w *zc1069Walker) walkForLoop(t *ast.ForLoopStatement, inFunction bool) {
	w.walk(t.Init, inFunction)
	w.walk(t.Condition, inFunction)
	w.walk(t.Post, inFunction)
	w.walkExpressions(t.Items, inFunction)
	w.walk(t.Body, inFunction)
}

func (w *zc1069Walker) walkFunctionLiteral(t *ast.FunctionLiteral) {
	for _, p := range t.Params {
		w.walk(p, false)
	}
	w.walk(t.Body, true)
}

func (w *zc1069Walker) walkCaseStatement(t *ast.CaseStatement, inFunction bool) {
	w.walk(t.Value, inFunction)
	for _, clause := range t.Clauses {
		w.walkExpressions(clause.Patterns, inFunction)
		w.walk(clause.Body, inFunction)
	}
}

func init() {
	RegisterKata(ast.FunctionDefinitionNode, Kata{
		ID:    "ZC1070",
		Title: "Use `builtin` or `command` to avoid infinite recursion in wrapper functions",
		Description: "When defining a wrapper function with the same name as a builtin or command (e.g., `cd`), " +
			"calling the command directly inside the function causes infinite recursion. " +
			"Use `builtin cd` or `command cd`.",
		Severity: SeverityWarning,
		Check:    checkZC1070,
	})
}

func checkZC1070(node ast.Node) []Violation {
	funcDef, ok := node.(*ast.FunctionDefinition)
	if !ok {
		return nil
	}
	if funcDef.Name == nil || funcDef.Body == nil {
		return nil
	}

	name := funcDef.Name.String()

	// Only check for common builtins/commands to avoid flagging valid recursive algorithms
	targets := map[string]bool{
		"cd": true, "echo": true, "printf": true, "read": true, "source": true, ".": true,
		"eval": true, "exec": true, "exit": true, "export": true, "kill": true,
		"local": true, "pwd": true, "return": true, "set": true, "shift": true,
		"test": true, "trap": true, "typeset": true, "umask": true, "unset": true, "wait": true,
		"ls": true, "grep": true, "mkdir": true, "rm": true, "mv": true, "cp": true, "git": true,
		"dirs": true, "popd": true, "pushd": true,
	}

	if !targets[name] {
		return nil
	}

	violations := []Violation{}

	// Walk body to find self-calls
	ast.Walk(funcDef.Body, func(n ast.Node) bool {
		// Don't recurse into nested functions (they mask the name)
		if _, ok := n.(*ast.FunctionDefinition); ok && n != funcDef {
			return false
		}
		if _, ok := n.(*ast.FunctionLiteral); ok {
			return false
		}

		if cmd, ok := n.(*ast.SimpleCommand); ok {
			if cmd.Name == nil {
				return true
			}
			cmdName := cmd.Name.String()
			if cmdName == name {
				// Found self-call.
				// Check if it is "builtin name" or "command name" is handled?
				// SimpleCommand "builtin" with arg "name".
				// But here `cmd.Name` IS `name`.
				// So `builtin cd` -> Name="builtin", Args=["cd"]
				// `cd` -> Name="cd".

				// If Name == function name, it IS a recursive call.
				// Unless it is `builtin` or `command`?
				// If I write `builtin cd`, the parser sees Name="builtin".
				// So if Name matches `name`, it is NOT `builtin` or `command`.

				// Exception: `command` might not be a keyword in parser?
				// `command -v cd` -> Name="command".

				// So if `cmdName == name`, it is a direct call.

				violations = append(violations, Violation{
					KataID: "ZC1070",
					Message: "Recursive call to `" + name + "` inside `" + name + "`. " +
						"Use `builtin " + name + "` or `command " + name + "` to invoke the underlying command.",
					Line:   cmd.Token.Line,
					Column: cmd.Token.Column,
					Level:  SeverityWarning,
				})
			}
		}
		return true
	})

	return violations
}

func init() {
	RegisterKata(ast.InfixExpressionNode, Kata{
		ID:          "ZC1071",
		Title:       "Use `+=` for appending to arrays",
		Description: "Appending to an array using `arr=($arr ...)` is verbose and slower. Use `arr+=(...)` instead.",
		Severity:    SeverityWarning,
		Check:       checkZC1071,
	})
}

func checkZC1071(node ast.Node) []Violation {
	infix, ok := node.(*ast.InfixExpression)
	if !ok || infix.Operator != "=" {
		return nil
	}
	ident, ok := infix.Left.(*ast.Identifier)
	if !ok {
		return nil
	}
	arrayLit, ok := infix.Right.(*ast.ArrayLiteral)
	if !ok {
		return nil
	}
	found := false
	checkNode := func(n ast.Node) bool {
		if found {
			return false
		}
		if zc1071SelfReferences(n, ident.Value) {
			found = true
			return false
		}
		return true
	}
	for _, elem := range arrayLit.Elements {
		if found {
			break
		}
		ast.Walk(elem, checkNode)
	}

	if found {
		leftToken := infix.Left.TokenLiteralNode()
		return []Violation{{
			KataID: "ZC1071",
			Message: "Appending to an array using `arr=($arr ...)` is verbose and slower. " +
				"Use `arr+=(...)` instead.",
			Line:   leftToken.Line,
			Column: leftToken.Column,
			Level:  SeverityWarning,
		}}
	}

	return nil
}

func zc1071SelfReferences(n ast.Node, varName string) bool {
	switch v := n.(type) {
	case *ast.ArrayAccess:
		id, ok := v.Left.(*ast.Identifier)
		return ok && id.Value == varName
	case *ast.Identifier:
		return v.Value == "$"+varName || v.Value == "${"+varName+"}"
	case *ast.PrefixExpression:
		if v.Operator != "$" {
			return false
		}
		id, ok := v.Right.(*ast.Identifier)
		return ok && id.Value == varName
	}
	return false
}

func init() {
	RegisterKata(ast.InfixExpressionNode, Kata{
		ID:    "ZC1072",
		Title: "Use `awk` instead of `grep | awk`",
		Description: "`grep pattern | awk '{...}'` is inefficient. " +
			"Use `awk '/pattern/ {...}'` to combine matching and processing in a single process.",
		Severity: SeverityStyle,
		Check:    checkZC1072,
	})
}

func checkZC1072(node ast.Node) []Violation {
	pipe, ok := node.(*ast.InfixExpression)
	if !ok || pipe.Operator != "|" {
		return nil
	}

	// Check left command is grep
	grepCmd, ok := pipe.Left.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	if !isCommandName(grepCmd, "grep") {
		return nil
	}

	// Check right command is awk/gawk/mawk
	awkCmd, ok := pipe.Right.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	if !isCommandName(awkCmd, "awk") && !isCommandName(awkCmd, "gawk") && !isCommandName(awkCmd, "mawk") {
		return nil
	}

	// Check grep flags. If flags are complex (like -r, -v, -l), we might skip warning.
	// But `grep | awk` is almost always replaceable.
	// Only if grep does something awk can't easily do (like -r recursive search) should we allow it?
	// Awk doesn't do recursive directory search by default.
	// So if grep has `-r` or `-R`, it's valid.

	if hasRecursiveFlag(grepCmd) {
		return nil
	}

	return []Violation{{
		KataID:  "ZC1072",
		Message: "Use `awk '/pattern/ {...}'` instead of `grep pattern | awk '{...}'` to avoid a pipeline.",
		Line:    pipe.TokenLiteralNode().Line,
		Column:  pipe.TokenLiteralNode().Column,
		Level:   SeverityStyle,
	}}
}

func hasRecursiveFlag(cmd *ast.SimpleCommand) bool {
	for _, arg := range cmd.Arguments {
		val := arg.String()
		// Remove quotes
		val = strings.Trim(val, "\"'")
		if strings.HasPrefix(val, "-") {
			if val == "-r" || val == "-R" || val == "--recursive" {
				return true
			}
			// Combined flags e.g. -rn (must ensure it's not a long flag like --recursive checked above)
			if !strings.HasPrefix(val, "--") && (strings.Contains(val, "r") || strings.Contains(val, "R")) {
				return true
			}
		}
	}
	return false
}

func init() {
	RegisterKata(ast.ArithmeticCommandNode, Kata{
		ID:          "ZC1073",
		Title:       "Unnecessary use of `$` in arithmetic expressions",
		Description: "Variables in `((...))` do not need `$` prefix. Use `(( var > 0 ))` instead of `(( $var > 0 ))`.",
		Severity:    SeverityStyle,
		Check:       checkZC1073,
		Fix:         fixZC1073,
	})
}

// fixZC1073 deletes the leading `$` from a variable used inside
// `(( … ))`. The violation coordinates already point at the `$`
// byte, so a single zero-replacement edit removes it. A second pass
// won't re-trigger because the identifier no longer carries `$`.
func fixZC1073(_ ast.Node, v Violation, source []byte) []FixEdit {
	off := LineColToByteOffset(source, v.Line, v.Column)
	if off < 0 || off >= len(source) || source[off] != '$' {
		return nil
	}
	return []FixEdit{{
		Line:    v.Line,
		Column:  v.Column,
		Length:  1,
		Replace: "",
	}}
}

func checkZC1073(node ast.Node) []Violation {
	cmd, ok := node.(*ast.ArithmeticCommand)
	if !ok {
		return nil
	}

	if cmd.Expression == nil {
		return nil
	}

	var violations []Violation

	ast.Walk(cmd.Expression, func(n ast.Node) bool {
		// Check for PrefixExpression with '$'
		if prefix, ok := n.(*ast.PrefixExpression); ok && prefix.Operator == "$" {
			if ident, ok := prefix.Right.(*ast.Identifier); ok {
				if isUserVariable(ident.Value) {
					violations = append(violations, Violation{
						KataID:  "ZC1073",
						Message: "Unnecessary use of `$` in arithmetic expressions. Use `(( var ))` instead of `(( $var ))`.",
						Line:    prefix.Token.Line,
						Column:  prefix.Token.Column,
						Level:   SeverityStyle,
					})
				}
			}
			return true
		}

		// Check for Identifier starting with '$' (if lexer emits VARIABLE)
		if ident, ok := n.(*ast.Identifier); ok {
			if len(ident.Value) > 1 && ident.Value[0] == '$' {
				varName := ident.Value[1:]
				if isUserVariable(varName) {
					violations = append(violations, Violation{
						KataID:  "ZC1073",
						Message: "Unnecessary use of `$` in arithmetic expressions. Use `(( var ))` instead of `(( $var ))`.",
						Line:    ident.Token.Line,
						Column:  ident.Token.Column,
						Level:   SeverityStyle,
					})
				}
			}
		}

		return true
	})

	return violations
}

func isUserVariable(name string) bool {
	if len(name) == 0 {
		return false
	}

	first := name[0]
	if !isAlpha(first) && first != '_' {
		return false
	}

	for i := 1; i < len(name); i++ {
		if !isAlphaNumeric(name[i]) && name[i] != '_' {
			return false
		}
	}

	return true
}

func isAlpha(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z')
}

func isAlphaNumeric(b byte) bool {
	return isAlpha(b) || (b >= '0' && b <= '9')
}

func init() {
	kata := Kata{
		ID:    "ZC1074",
		Title: "Prefer modifiers :h/:t over dirname/basename",
		Description: "Zsh provides modifiers like `:h` (head/dirname) and `:t` (tail/basename) " +
			"that are faster and more idiomatic than spawning external commands.",
		Severity: SeverityStyle,
		Check:    checkZC1074,
	}
	RegisterKata(ast.CommandSubstitutionNode, kata)
	RegisterKata(ast.DollarParenExpressionNode, kata)
}

func checkZC1074(node ast.Node) []Violation {
	var command ast.Node

	switch n := node.(type) {
	case *ast.CommandSubstitution:
		command = n.Command
	case *ast.DollarParenExpression:
		command = n.Command
	default:
		return nil
	}

	// Check if command is "dirname" or "basename"
	if cmd, ok := command.(*ast.SimpleCommand); ok {
		cmdName := cmd.Name.String()
		if cmdName == "dirname" {
			return []Violation{{
				KataID:  "ZC1074",
				Message: "Use '${var:h}' instead of '$(dirname $var)'. Modifiers are faster and built-in.",
				Line:    node.TokenLiteralNode().Line,
				Column:  node.TokenLiteralNode().Column,
				Level:   SeverityStyle,
			}}
		}
		if cmdName == "basename" {
			return []Violation{{
				KataID:  "ZC1074",
				Message: "Use '${var:t}' instead of '$(basename $var)'. Modifiers are faster and built-in.",
				Line:    node.TokenLiteralNode().Line,
				Column:  node.TokenLiteralNode().Column,
				Level:   SeverityStyle,
			}}
		}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1075",
		Title: "Quote variable expansions to prevent globbing",
		Description: "Unquoted variable expansions in Zsh are subject to globbing (filename generation). " +
			"If the variable contains characters like `*` or `?`, it might match files unexpectedly. " +
			"Use quotes `\"$var\"` to prevent this.",
		Severity: SeverityWarning,
		Check:    checkZC1075,
	})
}

func checkZC1075(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	violations := []Violation{}

	for _, arg := range cmd.Arguments {
		// Check if argument is a simple identifier starting with $
		// or a braced expression ${...} that is NOT inside a string literal.
		// The parser might wrap these in different ways.

		// If it's a bare IdentifierNode (variable expansion), it's unquoted.
		if ident, ok := arg.(*ast.Identifier); ok {
			// Identifiers that start with $ are variable expansions
			if len(ident.Value) > 0 && ident.Value[0] == '$' {
				violations = append(violations, Violation{
					KataID:  "ZC1075",
					Message: "Unquoted variable expansion '" + ident.Value + "' is subject to globbing. Quote it: \"" + ident.Value + "\".",
					Line:    ident.Token.Line,
					Column:  ident.Token.Column,
					Level:   SeverityWarning,
				})
			}
		} else if _, ok := arg.(*ast.ArrayAccess); ok {
			// Array access ${arr[idx]} is also subject to globbing if unquoted
			violations = append(violations, Violation{
				KataID:  "ZC1075",
				Message: "Unquoted array access is subject to globbing. Quote it.",
				Line:    arg.TokenLiteralNode().Line,
				Column:  arg.TokenLiteralNode().Column,
				Level:   SeverityWarning,
			})
		} else if _, ok := arg.(*ast.InvalidArrayAccess); ok {
			_ = ok
			// $arr[idx] - ZC1001 flags this, but it's also unquoted.
			// Let ZC1001 handle the syntax error, but ZC1075 could also flag globbing.
			// We'll skip to reduce noise.
		}

		// Note: StringLiteral arguments are quoted, so we don't check them.
		// But ConcatenatedExpression might contain unquoted parts.
		// e.g. $var/foo
		if concat, ok := arg.(*ast.ConcatenatedExpression); ok {
			for _, part := range concat.Parts {
				if ident, ok := part.(*ast.Identifier); ok {
					if len(ident.Value) > 0 && ident.Value[0] == '$' {
						violations = append(violations, Violation{
							KataID:  "ZC1075",
							Message: "Unquoted variable expansion '" + ident.Value + "' in concatenated string is subject to globbing.",
							Line:    ident.Token.Line,
							Column:  ident.Token.Column,
							Level:   SeverityWarning,
						})
					}
				}
			}
		}
	}

	return violations
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1076",
		Title: "Use `autoload -Uz` for lazy loading",
		Description: "When using `autoload`, prefer `-Uz` to ensure standard Zsh behavior (no alias expansion, zsh style). " +
			"`-U` prevents alias expansion, and `-z` ensures Zsh style autoloading.",
		Severity: SeverityStyle,
		Check:    checkZC1076,
		Fix:      fixZC1076,
	})
}

// fixZC1076 inserts ` -Uz` after the `autoload` command name. Only
// fires when neither `U` nor `z` are already present; the detector
// already gates on that. Idempotent on re-run once both flags exist.
func fixZC1076(node ast.Node, v Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	if cmd.Name.String() != "autoload" {
		return nil
	}
	nameOffset := LineColToByteOffset(source, v.Line, v.Column)
	if nameOffset < 0 {
		return nil
	}
	nameLen := IdentLenAt(source, nameOffset)
	if nameLen != len("autoload") {
		return nil
	}
	insertAt := nameOffset + nameLen
	insLine, insCol := offsetLineColZC1076(source, insertAt)
	if insLine < 0 {
		return nil
	}
	return []FixEdit{{
		Line:    insLine,
		Column:  insCol,
		Length:  0,
		Replace: " -Uz",
	}}
}

func offsetLineColZC1076(source []byte, offset int) (int, int) {
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

func checkZC1076(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	if cmd.Name.String() != "autoload" {
		return nil
	}

	hasU := false
	hasZ := false

	for _, arg := range cmd.Arguments {
		argStr := arg.String()
		argStr = strings.Trim(argStr, "\"'")
		if strings.HasPrefix(argStr, "-") {
			if strings.Contains(argStr, "U") {
				hasU = true
			}
			if strings.Contains(argStr, "z") {
				hasZ = true
			}
		}
	}

	if !hasU || !hasZ {
		return []Violation{{
			KataID:  "ZC1076",
			Message: "Use `autoload -Uz` to ensure consistent and safe function loading.",
			Line:    cmd.TokenLiteralNode().Line,
			Column:  cmd.TokenLiteralNode().Column,
			Level:   SeverityStyle,
		}}
	}

	return nil
}

func init() {
	RegisterKata(ast.CommandSubstitutionNode, Kata{
		ID:    "ZC1077",
		Title: "Prefer `${var:u/l}` over `tr` for case conversion",
		Description: "Using `tr` in a pipeline for simple case conversion is slower than using " +
			"Zsh's built-in parameter expansion flags `:u` (upper) and `:l` (lower).",
		Severity: SeverityStyle,
		Check:    checkZC1077,
	})
	RegisterKata(ast.DollarParenExpressionNode, Kata{
		ID:    "ZC1077",
		Title: "Prefer `${var:u/l}` over `tr` for case conversion",
		Description: "Using `tr` in a pipeline for simple case conversion is slower than using " +
			"Zsh's built-in parameter expansion flags `:u` (upper) and `:l` (lower).",
		Severity: SeverityStyle,
		Check:    checkZC1077,
	})
}

func checkZC1077(node ast.Node) []Violation {
	rightCmd, ok := zc1077TrPipeline(node)
	if !ok || len(rightCmd.Arguments) < 2 {
		return nil
	}
	arg1 := rightCmd.Arguments[0].String()
	arg2 := rightCmd.Arguments[1].String()

	if zc1077IsUpperPair(arg1, arg2) {
		return zc1077Hit(node, "u", "uppercase")
	}
	if zc1077IsLowerPair(arg1, arg2) {
		return zc1077Hit(node, "l", "lowercase")
	}
	return nil
}

// zc1077TrPipeline returns the right-hand `tr` command of a `cmd | tr`
// pipeline embedded in either a backtick or `$()` substitution.
func zc1077TrPipeline(node ast.Node) (*ast.SimpleCommand, bool) {
	var command ast.Node
	switch n := node.(type) {
	case *ast.CommandSubstitution:
		command = n.Command
	case *ast.DollarParenExpression:
		command = n.Command
	default:
		return nil, false
	}
	infix, ok := command.(*ast.InfixExpression)
	if !ok || infix.Operator != "|" {
		return nil, false
	}
	rightCmd, ok := infix.Right.(*ast.SimpleCommand)
	if !ok || rightCmd.Name.String() != "tr" {
		return nil, false
	}
	return rightCmd, true
}

func zc1077IsUpperPair(a, b string) bool {
	return (checkTrPattern(a, "a-z") && checkTrPattern(b, "A-Z")) ||
		(checkTrPattern(a, "[:lower:]") && checkTrPattern(b, "[:upper:]"))
}

func zc1077IsLowerPair(a, b string) bool {
	return (checkTrPattern(a, "A-Z") && checkTrPattern(b, "a-z")) ||
		(checkTrPattern(a, "[:upper:]") && checkTrPattern(b, "[:lower:]"))
}

func zc1077Hit(node ast.Node, flag, label string) []Violation {
	return []Violation{{
		KataID:  "ZC1077",
		Message: "Use `${var:" + flag + "}` instead of `tr` for " + label + " conversion. It is faster and built-in.",
		Line:    node.TokenLiteralNode().Line,
		Column:  node.TokenLiteralNode().Column,
		Level:   SeverityStyle,
	}}
}

func checkTrPattern(arg, pattern string) bool {
	// Remove quotes
	stripped := arg
	if len(arg) >= 2 && ((arg[0] == '"' && arg[len(arg)-1] == '"') || (arg[0] == '\'' && arg[len(arg)-1] == '\'')) {
		stripped = arg[1 : len(arg)-1]
	}

	// Simple containment check - robust enough for standard patterns
	// We check if the core pattern exists
	// e.g. 'a-z' matches "a-z", 'a-z', [a-z]

	// For strictness, let's just check substring
	return stripped == pattern || stripped == "["+pattern+"]"
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1078",
		Title: "Quote `$@` and `$*` when passing arguments",
		Description: "Using unquoted `$@` or `$*` splits arguments by IFS (usually space). " +
			"Use `\"$@\"` to preserve the original argument grouping, or `\"$*\"` to join them into a single string.",
		Severity: SeverityWarning,
		Check:    checkZC1078,
		Fix:      fixZC1078,
	})
}

// fixZC1078 wraps an unquoted `$@` / `$*` argument in double-quotes.
// Both tokens are exactly two bytes; the two-edit insertion always
// surrounds the same 2-byte run.
func fixZC1078(_ ast.Node, v Violation, source []byte) []FixEdit {
	start := LineColToByteOffset(source, v.Line, v.Column)
	if start < 0 || start+2 > len(source) {
		return nil
	}
	if source[start] != '$' || (source[start+1] != '@' && source[start+1] != '*') {
		return nil
	}
	endLine, endCol := offsetLineColZC1078(source, start+2)
	if endLine < 0 {
		return nil
	}
	return []FixEdit{
		{Line: v.Line, Column: v.Column, Length: 0, Replace: `"`},
		{Line: endLine, Column: endCol, Length: 0, Replace: `"`},
	}
}

func offsetLineColZC1078(source []byte, offset int) (int, int) {
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

func checkZC1078(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	violations := []Violation{}

	for _, arg := range cmd.Arguments {
		// Check string representation to catch various parsed forms of $@ and $*
		// unquoted $@ might be parsed as Identifier "$@" -> String() == "$@"
		// unquoted $* might be parsed as GroupedExpression -> String() == "($*)"
		// or other variations depending on parser state (e.g. PrefixExpression)

		s := arg.String()

		// Removing parens from GroupedExpression string representation for checking
		// (Note: String() adds parens for GroupedExpression)
		if len(s) >= 2 && s[0] == '(' && s[len(s)-1] == ')' {
			s = s[1 : len(s)-1]
		}

		if s == "$@" || s == "$*" {
			violations = append(violations, Violation{
				KataID:  "ZC1078",
				Message: "Unquoted " + s + " splits arguments. Use \"" + s + "\" to preserve structure.",
				Line:    arg.TokenLiteralNode().Line,
				Column:  arg.TokenLiteralNode().Column,
				Level:   SeverityWarning,
			})
		}
	}

	return violations
}

func init() {
	RegisterKata(ast.DoubleBracketExpressionNode, Kata{
		ID:    "ZC1079",
		Title: "Quote RHS of `==` in `[[ ... ]]` to prevent pattern matching",
		Description: "In `[[ ... ]]`, unquoted variable expansions on the right-hand side of `==` or `!=` " +
			"are treated as patterns (globbing). If you intend to compare strings literally, quote the variable.",
		Severity: SeverityWarning,
		Check:    checkZC1079,
		Fix:      fixZC1079,
	})
}

// fixZC1079 wraps an unquoted RHS variable reference inside `[[ … ]]`
// with double-quotes. Two edits: one `"` before the RHS token, one
// after. RHS span is measured from source so `${arr[$i]}` and
// `${var:-default}` stay whole. When the sibling LHS is an empty
// string literal, ZC1055's `-z` / `-n` rewrite takes priority and
// this fix no-ops to avoid overlapping edits.
func fixZC1079(node ast.Node, v Violation, source []byte) []FixEdit {
	if dbe, ok := node.(*ast.DoubleBracketExpression); ok {
		for _, el := range dbe.Elements {
			infix, ok := el.(*ast.InfixExpression)
			if !ok {
				continue
			}
			if infix.Operator != "==" && infix.Operator != "=" && infix.Operator != "!=" {
				continue
			}
			if isEmptyStringLiteral(infix.Left) || isEmptyStringLiteral(infix.Right) {
				// ZC1055 owns this rewrite; skip.
				return nil
			}
		}
	}
	start := LineColToByteOffset(source, v.Line, v.Column)
	if start < 0 || start >= len(source) {
		return nil
	}
	argLen := unquotedArgLen(source, start)
	if argLen == 0 {
		return nil
	}
	endOff := start + argLen
	endLine, endCol := offsetLineColZC1079(source, endOff)
	if endLine < 0 {
		return nil
	}
	return []FixEdit{
		{Line: v.Line, Column: v.Column, Length: 0, Replace: `"`},
		{Line: endLine, Column: endCol, Length: 0, Replace: `"`},
	}
}

func isEmptyStringLiteral(n ast.Node) bool {
	str, ok := n.(*ast.StringLiteral)
	if !ok {
		return false
	}
	return str.Value == `""` || str.Value == `''`
}

func offsetLineColZC1079(source []byte, offset int) (int, int) {
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

var zc1079EqualityOps = map[string]struct{}{"==": {}, "=": {}, "!=": {}}

func checkZC1079(node ast.Node) []Violation {
	dbe, ok := node.(*ast.DoubleBracketExpression)
	if !ok {
		return nil
	}
	violations := []Violation{}
	for _, expr := range dbe.Elements {
		infix, ok := expr.(*ast.InfixExpression)
		if !ok {
			continue
		}
		if _, hit := zc1079EqualityOps[infix.Operator]; !hit {
			continue
		}
		if tok := zc1079UnquotedVar(infix.Right); tok != nil {
			violations = append(violations, Violation{
				KataID:  "ZC1079",
				Message: "Unquoted RHS matches as pattern. Quote to force string comparison: `\"$var\"`.",
				Line:    tok.TokenLiteralNode().Line,
				Column:  tok.TokenLiteralNode().Column,
				Level:   SeverityWarning,
			})
		}
	}
	return violations
}

// zc1079UnquotedVar returns the token-bearing node when expr resolves
// to an unquoted variable / array reference; nil otherwise.
func zc1079UnquotedVar(expr ast.Expression) ast.Node {
	switch r := expr.(type) {
	case *ast.Identifier:
		if len(r.Value) > 0 && r.Value[0] == '$' {
			return r
		}
	case *ast.ArrayAccess, *ast.InvalidArrayAccess:
		return r.(ast.Node)
	case *ast.ConcatenatedExpression:
		for _, part := range r.Parts {
			if ident, ok := part.(*ast.Identifier); ok && len(ident.Value) > 0 && ident.Value[0] == '$' {
				return ident
			}
		}
	}
	return nil
}

func init() {
	RegisterKata(ast.ForLoopStatementNode, Kata{
		ID:    "ZC1080",
		Title: "Use `(N)` nullglob qualifier for globs in loops",
		Description: "In Zsh, if a glob matches no files, it throws an error by default. " +
			"When iterating over a glob in a `for` loop, use the `(N)` glob qualifier to allow it to match nothing (nullglob).",
		Severity: SeverityStyle,
		Check:    checkZC1080,
	})
}

func checkZC1080(node ast.Node) []Violation {
	loop, ok := node.(*ast.ForLoopStatement)
	if !ok {
		return nil
	}

	// Only check for-each loops: for i in items...
	if loop.Items == nil {
		return nil // C-style loop
	}

	violations := []Violation{}

	for _, item := range loop.Items {
		if hasGlobChars(item) {
			s := item.String()
			if !strings.Contains(s, "(N)") && !strings.Contains(s, "N") {
				violations = append(violations, Violation{
					KataID:  "ZC1080",
					Message: "Glob '" + s + "' will error if no matches found. Append `(N)` to make it nullglob.",
					Line:    item.TokenLiteralNode().Line,
					Column:  item.TokenLiteralNode().Column,
					Level:   SeverityStyle,
				})
			}
		}
	}

	return violations
}

func hasGlobChars(node ast.Node) bool {
	// Check AST nodes to see if they contain exposed glob characters (*, ?, [)
	switch n := node.(type) {
	case *ast.StringLiteral:
		// Check for unquoted glob chars in the literal
		// If parsed as StringLiteral, it might be single quoted (no glob) or simple word (glob).
		// Lexer strips quotes from Literal? No, value usually keeps them.
		val := n.Value
		if len(val) >= 2 && (val[0] == '\'' || val[0] == '"') {
			return false // Quoted strings don't glob
		}
		return strings.ContainsAny(val, "*?[]")
	case *ast.Identifier:
		// Identifiers don't glob unless they contain * (which usually makes them NOT identifiers but string/prefix)
		// But Parser might be lenient.
		return strings.ContainsAny(n.Value, "*?[]")
	case *ast.PrefixExpression:
		// *, ? prefix operators
		if n.Operator == "*" || n.Operator == "?" {
			return true
		}
		// Recursive check
		return hasGlobChars(n.Right)
	case *ast.ConcatenatedExpression:
		for _, part := range n.Parts {
			if hasGlobChars(part) {
				return true
			}
		}
		return false
	case *ast.ArrayAccess:
		return false // Array access ${...} is not a file glob
	case *ast.SimpleCommand:
		// [ char range ] is parsed as SimpleCommand sometimes?
		// No, usually Concatenated or StringLiteral if [ is treated as literal.
		// If [ is SimpleCommand name (e.g. `[` test command), it's not a glob.
	}
	return false
}

func init() {
	RegisterKata(ast.InfixExpressionNode, Kata{
		ID:    "ZC1081",
		Title: "Use `${#var}` to get string length instead of `wc -c`",
		Description: "Using `echo $var | wc -c` involves a subshell and external command overhead. " +
			"Zsh has a built-in operator `${#var}` to get the length of a string instantly.",
		Severity: SeverityStyle,
		Check:    checkZC1081,
	})
}

func checkZC1081(node ast.Node) []Violation {
	infix, ok := node.(*ast.InfixExpression)
	if !ok {
		return nil
	}

	if infix.Operator != "|" {
		return nil
	}

	// Check Right side: wc -c or wc -m
	rightCmd, ok := infix.Right.(*ast.SimpleCommand)
	if !ok || rightCmd.Name.String() != "wc" {
		return nil
	}

	isCharCount := false
	for _, arg := range rightCmd.Arguments {
		s := arg.String()
		if strings.Contains(s, "-c") || strings.Contains(s, "-m") {
			isCharCount = true
			break
		}
	}

	if !isCharCount {
		return nil
	}

	// Check Left side: echo ... or printf ...
	leftCmd, ok := infix.Left.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	cmdName := leftCmd.Name.String()
	if cmdName == "echo" || cmdName == "print" || cmdName == "printf" {
		return []Violation{{
			KataID:  "ZC1081",
			Message: "Use `${#var}` to get string length. Pipeline to `wc` is inefficient.",
			Line:    infix.TokenLiteralNode().Line,
			Column:  infix.TokenLiteralNode().Column,
			Level:   SeverityStyle,
		}}
	}

	return nil
}

func init() {
	RegisterKata(ast.InfixExpressionNode, Kata{
		ID:    "ZC1082",
		Title: "Prefer `${var//old/new}` over `sed` for simple replacements",
		Description: "Using `sed` for simple string replacement is slower than Zsh's built-in " +
			"parameter expansion. Use `${var/old/new}` (replace first) or `${var//old/new}` (replace all).",
		Severity: SeverityStyle,
		Check:    checkZC1082,
	})
}

func checkZC1082(node ast.Node) []Violation {
	infix, ok := node.(*ast.InfixExpression)
	if !ok || infix.Operator != "|" {
		return nil
	}

	// Check Right side: sed
	rightCmd, ok := infix.Right.(*ast.SimpleCommand)
	if !ok || rightCmd.Name.String() != "sed" {
		return nil
	}

	// Check Left side: echo/printf/print
	leftCmd, ok := infix.Left.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	cmdName := leftCmd.Name.String()
	if cmdName != "echo" && cmdName != "print" && cmdName != "printf" {
		return nil
	}

	// Analyze sed arguments
	for _, arg := range rightCmd.Arguments {
		argStr := arg.String()
		// Remove quotes
		argStr = strings.Trim(argStr, "\"'")

		// Look for s/old/new/ or s/old/new/g
		// Basic check: starts with s/
		if strings.HasPrefix(argStr, "s/") || strings.HasPrefix(argStr, "s|") || strings.HasPrefix(argStr, "s@") {
			// It's a substitution
			return []Violation{{
				KataID:  "ZC1082",
				Message: "Use `${var//old/new}` for string replacement. Pipeline to `sed` is inefficient.",
				Line:    infix.TokenLiteralNode().Line,
				Column:  infix.TokenLiteralNode().Column,
				Level:   SeverityStyle,
			}}
		}
	}

	return nil
}

func init() {
	RegisterKata(ast.ConcatenatedExpressionNode, Kata{
		ID:    "ZC1083",
		Title: "Brace expansion limits cannot be variables",
		Description: "Brace expansion `{x..y}` happens before variable expansion. " +
			"`{1..$n}` will not work. Use `seq` or `for ((...))`.",
		Severity: SeverityError,
		Check:    checkZC1083,
	})
	RegisterKata(ast.StringLiteralNode, Kata{
		ID:    "ZC1083",
		Title: "Brace expansion limits cannot be variables",
		Description: "Brace expansion `{x..y}` happens before variable expansion. " +
			"`{1..$n}` will not work. Use `seq` or `for ((...))`.",
		Severity: SeverityError,
		Check:    checkZC1083,
	})
}

func checkZC1083(node ast.Node) []Violation {
	if strNode, ok := node.(*ast.StringLiteral); ok {
		return zc1083CheckStringLiteral(strNode)
	}
	concat, ok := node.(*ast.ConcatenatedExpression)
	if !ok {
		return nil
	}
	return zc1083CheckConcat(concat)
}

func zc1083CheckStringLiteral(s *ast.StringLiteral) []Violation {
	v := s.Value
	if !strings.Contains(v, "{") || !strings.Contains(v, "..") || !strings.Contains(v, "$") {
		return nil
	}
	return zc1083Hit(s)
}

func zc1083CheckConcat(concat *ast.ConcatenatedExpression) []Violation {
	scan := zc1083ScanParts(concat.Parts)
	if scan.startIdx == -1 {
		return nil
	}
	if !zc1083HasIndexAfter(scan.dotDotIndices, scan.startIdx) {
		return nil
	}
	if !zc1083HasIndexAfter(scan.varIndices, scan.startIdx) {
		return nil
	}
	return zc1083Hit(concat)
}

type zc1083Scan struct {
	startIdx       int
	dotDotIndices  []int
	varIndices     []int
}

func zc1083ScanParts(parts []ast.Expression) zc1083Scan {
	scan := zc1083Scan{startIdx: -1}
	lastPartWasDot := false
	for i, part := range parts {
		if strNode, ok := part.(*ast.StringLiteral); ok {
			zc1083ScanString(&scan, &lastPartWasDot, strNode.Value, i)
			continue
		}
		lastPartWasDot = false
		if _, isInt := part.(*ast.IntegerLiteral); isInt {
			continue
		}
		if idNode, isIdent := part.(*ast.Identifier); isIdent && strings.Contains(idNode.Value, "..") {
			scan.dotDotIndices = append(scan.dotDotIndices, i)
		}
		scan.varIndices = append(scan.varIndices, i)
	}
	return scan
}

func zc1083ScanString(scan *zc1083Scan, lastPartWasDot *bool, val string, i int) {
	if strings.Contains(val, "{") && scan.startIdx == -1 {
		scan.startIdx = i
	}
	switch {
	case strings.Contains(val, ".."):
		scan.dotDotIndices = append(scan.dotDotIndices, i)
		*lastPartWasDot = false
	case val == ".":
		if *lastPartWasDot {
			scan.dotDotIndices = append(scan.dotDotIndices, i-1)
			*lastPartWasDot = false
		} else {
			*lastPartWasDot = true
		}
	default:
		*lastPartWasDot = false
	}
}

func zc1083HasIndexAfter(indices []int, after int) bool {
	for _, idx := range indices {
		if idx > after {
			return true
		}
	}
	return false
}

func zc1083Hit(node ast.Node) []Violation {
	tok := node.TokenLiteralNode()
	return []Violation{{
		KataID:  "ZC1083",
		Message: "Brace expansion limits cannot be variables. `{...$var...}` is treated as a literal string. Use `seq` or `for ((...))`.",
		Line:    tok.Line,
		Column:  tok.Column,
		Level:   SeverityError,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1084",
		Title: "Quote globs in `find` commands",
		Description: "Unquoted globs in `find` commands are expanded by the shell before `find` runs. " +
			"If files match, `find` receives the list of files instead of the pattern. " +
			"Quote arguments to `-name`, `-path`, etc.",
		Severity: SeverityWarning,
		Check:    checkZC1084,
		Fix:      fixZC1084,
	})
}

// fixZC1084 wraps an unquoted `find` glob argument in single-quotes
// so the shell passes the pattern through verbatim. The violation
// column already points at the pattern arg start. Span scanning
// respects `[…]` / `{…}` so character classes and alternations
// stay whole.
func fixZC1084(_ ast.Node, v Violation, source []byte) []FixEdit {
	start := LineColToByteOffset(source, v.Line, v.Column)
	if start < 0 || start >= len(source) {
		return nil
	}
	argLen := unquotedArgLen(source, start)
	if argLen == 0 {
		return nil
	}
	endLine, endCol := offsetLineColZC1084(source, start+argLen)
	if endLine < 0 {
		return nil
	}
	return []FixEdit{
		{Line: v.Line, Column: v.Column, Length: 0, Replace: `'`},
		{Line: endLine, Column: endCol, Length: 0, Replace: `'`},
	}
}

func offsetLineColZC1084(source []byte, offset int) (int, int) {
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
								Level:   SeverityWarning,
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
							Level:   SeverityWarning,
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
				Level:   SeverityWarning,
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
	switch n := node.(type) {
	case *ast.SimpleCommand:
		return n.Name.String() == "["
	case *ast.IndexExpression:
		return true
	case *ast.StringLiteral:
		return isUnquotedGlobString(n)
	case *ast.ConcatenatedExpression:
		return isUnquotedGlobConcat(n)
	case *ast.PrefixExpression:
		return n.Operator == "*" || n.Operator == "?"
	}
	return false
}

func isUnquotedGlobString(s *ast.StringLiteral) bool {
	tok := s.TokenLiteralNode()
	if tok.Type == token.STRING {
		return false
	}
	return isGlobToken(tok)
}

func isUnquotedGlobConcat(concat *ast.ConcatenatedExpression) bool {
	escaped := false
	for _, part := range concat.Parts {
		hit, nextEscaped := isUnquotedGlobConcatPart(part, escaped)
		if hit {
			return true
		}
		escaped = nextEscaped
	}
	return false
}

// isUnquotedGlobConcatPart returns (matched, escapedNext). matched signals
// the concat overall is an unquoted glob; escapedNext propagates the
// backslash-escape state to the next part.
func isUnquotedGlobConcatPart(part ast.Expression, escaped bool) (bool, bool) {
	switch n := part.(type) {
	case *ast.StringLiteral:
		tok := n.TokenLiteralNode()
		if escaped {
			return false, false
		}
		if tok.Literal == "\\" {
			return false, true
		}
		return isGlobToken(tok), false
	case *ast.SimpleCommand:
		return n.Name.String() == "[", false
	case *ast.PrefixExpression:
		if escaped {
			return false, false
		}
		return n.Operator == "*" || n.Operator == "?", false
	}
	return false, false
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

func init() {
	RegisterKata(ast.ForLoopStatementNode, Kata{
		ID:    "ZC1085",
		Title: "Quote variable expansions in `for` loops",
		Description: "Unquoted variable expansions in `for` loops are split by IFS (usually spaces). " +
			"This often leads to iterating over words instead of lines or array elements. Quote the expansion to preserve structure.",
		Severity: SeverityWarning,
		Check:    checkZC1085,
		Fix:      fixZC1085,
	})
}

// fixZC1085 wraps an unquoted expansion in a `for` loop item list
// with double-quotes. Two-edit insertion at the span start and end.
// Span uses the shared unquotedArgLen scanner so `${arr[@]}`,
// `$(cmd args)`, `${var:-default}` all stay whole.
func fixZC1085(_ ast.Node, v Violation, source []byte) []FixEdit {
	start := LineColToByteOffset(source, v.Line, v.Column)
	if start < 0 || start >= len(source) {
		return nil
	}
	argLen := unquotedArgLen(source, start)
	if argLen == 0 {
		return nil
	}
	endLine, endCol := offsetLineColZC1085(source, start+argLen)
	if endLine < 0 {
		return nil
	}
	return []FixEdit{
		{Line: v.Line, Column: v.Column, Length: 0, Replace: `"`},
		{Line: endLine, Column: endCol, Length: 0, Replace: `"`},
	}
}

func offsetLineColZC1085(source []byte, offset int) (int, int) {
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

func checkZC1085(node ast.Node) []Violation {
	loop, ok := node.(*ast.ForLoopStatement)
	if !ok {
		return nil
	}

	// If Items is nil or empty, it's either C-style or implicit `in "$@"`, ignore
	if len(loop.Items) == 0 {
		return nil
	}

	var violations []Violation

	for _, item := range loop.Items {
		if isUnquotedExpansion(item) {
			violations = append(violations, Violation{
				KataID:  "ZC1085",
				Message: "Unquoted variable expansion in for loop. This will split on IFS (usually space). Quote it to iterate over lines or array elements.",
				Line:    item.TokenLiteralNode().Line,
				Column:  item.TokenLiteralNode().Column,
				Level:   SeverityWarning,
			})
		}
	}

	return violations
}

func isUnquotedExpansion(expr ast.Expression) bool {
	// Check for Identifier (e.g. $var)
	if id, ok := expr.(*ast.Identifier); ok {
		return id.TokenLiteralNode().Type == "VARIABLE"
	}

	// Check for ArrayAccess (e.g. ${arr[@]})
	if _, ok := expr.(*ast.ArrayAccess); ok {
		return true
	}

	// Check for DollarParenExpression (e.g. $(cmd))
	if _, ok := expr.(*ast.DollarParenExpression); ok {
		return true
	}

	// Check for CommandSubstitution (e.g. `cmd`)
	if _, ok := expr.(*ast.CommandSubstitution); ok {
		return true
	}

	// Check for ConcatenatedExpression
	if concat, ok := expr.(*ast.ConcatenatedExpression); ok {
		inQuotes := false
		for _, part := range concat.Parts {
			if str, ok := part.(*ast.StringLiteral); ok {
				if str.Value == "\"" {
					inQuotes = !inQuotes
					continue
				}
				// Single quotes technically shouldn't appear here if parsed as StringLiteral?
			}

			if !inQuotes {
				if isUnquotedExpansion(part) {
					return true
				}
			}
		}
	}

	return false
}

func init() {
	RegisterKata(ast.FunctionDefinitionNode, Kata{
		ID:    "ZC1086",
		Title: "Prefer `func() { ... }` over `function func { ... }`",
		Description: "The `function` keyword is optional in Zsh and non-standard in POSIX sh. " +
			"Using `func() { ... }` is more portable and consistent.",
		Severity: SeverityStyle,
		Check:    checkZC1086,
		Fix:      fixZC1086,
	})
	RegisterKata(ast.FunctionLiteralNode, Kata{
		ID:    "ZC1086",
		Title: "Prefer `func() { ... }` over `function func { ... }`",
		Description: "The `function` keyword is optional in Zsh and non-standard in POSIX sh. " +
			"Using `func() { ... }` is more portable and consistent.",
		Severity: SeverityStyle,
		Check:    checkZC1086,
		Fix:      fixZC1086,
	})
}

// fixZC1086 rewrites `function name [()] { body }` to the portable
// `name() { body }` form. Deletes the `function ` prefix and, when
// the source doesn't already carry `()` after the name, inserts it.
func fixZC1086(node ast.Node, v Violation, source []byte) []FixEdit {
	var name string
	switch n := node.(type) {
	case *ast.FunctionLiteral:
		if n.TokenLiteral() != "function" || n.Name == nil {
			return nil
		}
		name = n.Name.Value
	case *ast.FunctionDefinition:
		if n.TokenLiteral() != "function" || n.Name == nil {
			return nil
		}
		name = n.Name.Value
	default:
		return nil
	}
	if name == "" {
		return nil
	}
	kwOffset := LineColToByteOffset(source, v.Line, v.Column)
	if kwOffset < 0 || kwOffset+len("function ") > len(source) {
		return nil
	}
	if string(source[kwOffset:kwOffset+len("function")]) != "function" {
		return nil
	}
	// Delete `function` + the whitespace run that follows, up to the
	// name start. Find the name's byte position by scanning forward
	// past the whitespace.
	i := kwOffset + len("function")
	for i < len(source) && (source[i] == ' ' || source[i] == '\t') {
		i++
	}
	if i+len(name) > len(source) || string(source[i:i+len(name)]) != name {
		return nil
	}
	edits := []FixEdit{{
		Line:    v.Line,
		Column:  v.Column,
		Length:  i - kwOffset,
		Replace: "",
	}}
	// After the name, check whether `()` already follows. Skip any
	// whitespace between name and peek byte; Zsh allows
	// `function foo ()` with a space.
	after := i + len(name)
	j := after
	for j < len(source) && (source[j] == ' ' || source[j] == '\t') {
		j++
	}
	if j >= len(source) || source[j] != '(' {
		afterLine, afterCol := offsetLineColZC1086(source, after)
		if afterLine < 0 {
			return nil
		}
		edits = append(edits, FixEdit{
			Line:    afterLine,
			Column:  afterCol,
			Length:  0,
			Replace: "()",
		})
	}
	return edits
}

func offsetLineColZC1086(source []byte, offset int) (int, int) {
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

func checkZC1086(node ast.Node) []Violation {
	// Case 1: function my_func { ... } -> Parsed as FunctionLiteralNode
	if funcLit, ok := node.(*ast.FunctionLiteral); ok {
		if funcLit.TokenLiteral() == "function" {
			return []Violation{
				{
					KataID:  "ZC1086",
					Message: "Prefer `func() { ... }` over `function func { ... }` for portability and consistency.",
					Line:    funcLit.TokenLiteralNode().Line,
					Column:  funcLit.TokenLiteralNode().Column,
					Level:   SeverityStyle,
				},
			}
		}
	}

	// Case 2: my_func() { ... } -> Parsed as FunctionDefinitionNode
	if funcDef, ok := node.(*ast.FunctionDefinition); ok {
		if funcDef.TokenLiteral() == "function" {
			return []Violation{
				{
					KataID:  "ZC1086",
					Message: "Prefer `func() { ... }` over `function func { ... }` for portability and consistency.",
					Line:    funcDef.TokenLiteralNode().Line,
					Column:  funcDef.TokenLiteralNode().Column,
					Level:   SeverityStyle,
				},
			}
		}
	}

	return nil
}

func init() {
	// Register for SimpleCommand (to check args)
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1087",
		Title: "Output redirection overwrites input file",
		Description: "Redirecting output to a file that is also being read as input causes the file to be truncated before it is read. " +
			"Use a temporary file or `sponge`.",
		Severity: SeverityError,
		Check:    checkZC1087,
	})
	// Register for Pipeline (|) to detect clobbering across pipe
	RegisterKata(ast.InfixExpressionNode, Kata{
		ID:    "ZC1087",
		Title: "Output redirection overwrites input file",
		Description: "Redirecting output to a file that is also being read as input causes the file to be truncated before it is read. " +
			"Use a temporary file or `sponge`.",
		Severity: SeverityError,
		Check:    checkZC1087,
	})
}

func checkZC1087(node ast.Node) []Violation {
	// Case 1: SimpleCommand (checking args for > file)
	if cmd, ok := node.(*ast.SimpleCommand); ok {
		inputs := collectInputs(cmd)
		outputs := collectOutputs(cmd)

		for _, output := range outputs {
			for _, input := range inputs {
				if input == output {
					return []Violation{
						{
							KataID:  "ZC1087",
							Message: "Output redirection overwrites input file `" + output + "`. The file is truncated before reading.",
							Line:    cmd.TokenLiteralNode().Line,
							Column:  cmd.TokenLiteralNode().Column,
							Level:   SeverityError,
						},
					}
				}
			}
		}
		return nil
	}

	// Case 2: Pipeline (cmd1 | cmd2)
	if infix, ok := node.(*ast.InfixExpression); ok {
		if infix.Operator != "|" {
			return nil
		}

		// Left side inputs
		inputs := collectInputs(infix.Left)
		// Right side outputs
		outputs := collectOutputs(infix.Right)

		for _, output := range outputs {
			for _, input := range inputs {
				if input == output {
					return []Violation{
						{
							KataID:  "ZC1087",
							Message: "Output redirection overwrites input file `" + output + "`. The file is truncated before reading.",
							Line:    infix.TokenLiteralNode().Line,
							Column:  infix.TokenLiteralNode().Column,
							Level:   SeverityError,
						},
					}
				}
			}
		}
	}

	return nil
}

func collectInputs(node ast.Node) []string {
	var inputs []string

	ast.Walk(node, func(n ast.Node) bool {
		if n == nil {
			return true
		}

		if cmd, ok := n.(*ast.SimpleCommand); ok {
			for i := 0; i < len(cmd.Arguments); i++ {
				arg := cmd.Arguments[i].String()
				switch arg {
				case "<":
					if i+1 < len(cmd.Arguments) {
						inputs = append(inputs, cmd.Arguments[i+1].String())
						i++
					}
				case ">", ">>", ">|", "&>":
					// Skip output redirection
					i++
				default:
					// Assume args are inputs unless they are flags
					if len(arg) > 0 && arg[0] != '-' {
						inputs = append(inputs, arg)
					}
				}
			}
		}
		return true
	})

	return inputs
}

func collectOutputs(node ast.Node) []string {
	var outputs []string
	ast.Walk(node, func(n ast.Node) bool {
		if n == nil {
			return true
		}
		if cmd, ok := n.(*ast.SimpleCommand); ok {
			for i := 0; i < len(cmd.Arguments); i++ {
				arg := cmd.Arguments[i].String()
				// Only output redirection that truncates: > or >|
				// Ignore append >>, &>>, etc. unless we want to catch clobbering there too?
				// Kata description says "truncated". >> does not truncate.
				if arg == ">" || arg == ">|" {
					if i+1 < len(cmd.Arguments) {
						outputs = append(outputs, cmd.Arguments[i+1].String())
						i++
					}
				} else if arg == ">>" || arg == "&>" || arg == "&>>" {
					// Skip operator and file
					i++
				}
			}
		}
		return true
	})
	return outputs
}

func init() {
	RegisterKata(ast.ProgramNode, Kata{
		ID:    "ZC1088",
		Title: "Subshell isolates state changes",
		Description: "Commands inside `( ... )` run in a subshell. " +
			"State changes like `cd`, `export`, or variable assignments are lost when the subshell exits. " +
			"Use `{ ... }` for grouping if you want to preserve state changes.",
		Severity: SeverityWarning,
		Check:    checkZC1088,
	})
}

func checkZC1088(node ast.Node) []Violation {
	// We perform a context-aware traversal from the root
	v := &zc1088Visitor{violations: []Violation{}}
	v.traverse(node, false)
	return v.violations
}

type zc1088Visitor struct {
	violations []Violation
}

func (v *zc1088Visitor) traverse(node ast.Node, expectsStatus bool) {
	if node == nil {
		return
	}

	// Handle typed nil interfaces
	if t, ok := node.(*ast.BlockStatement); ok && t == nil {
		return
	}
	if t, ok := node.(*ast.IfStatement); ok && t == nil {
		return
	}
	if t, ok := node.(*ast.WhileLoopStatement); ok && t == nil {
		return
	}
	if t, ok := node.(*ast.ExpressionStatement); ok && t == nil {
		return
	}
	if t, ok := node.(*ast.InfixExpression); ok && t == nil {
		return
	}
	if t, ok := node.(*ast.PrefixExpression); ok && t == nil {
		return
	}
	if t, ok := node.(*ast.GroupedExpression); ok && t == nil {
		return
	}
	if t, ok := node.(*ast.Program); ok && t == nil {
		return
	}
	if t, ok := node.(*ast.Subshell); ok && t == nil {
		return
	}

	switch n := node.(type) {
	case *ast.Program:
		for _, stmt := range n.Statements {
			v.traverse(stmt, false)
		}
	case *ast.Subshell:
		// ( ... ) in statement context
		if !expectsStatus {
			v.checkSubshell(n)
		}
		// Traverse content (BlockStatement) with expectsStatus=false for children
		v.traverse(n.Command, false)

	case *ast.BlockStatement:
		// Normal block { ... } or implicit (e.g. If Condition)
		// If this block expects status, only the LAST statement expects status.
		for i, stmt := range n.Statements {
			isLast := i == len(n.Statements)-1
			v.traverse(stmt, expectsStatus && isLast)
		}
	case *ast.IfStatement:
		v.traverse(n.Condition, true)
		v.traverse(n.Consequence, false)
		v.traverse(n.Alternative, false)
	case *ast.WhileLoopStatement:
		v.traverse(n.Condition, true)
		v.traverse(n.Body, false)
	case *ast.ExpressionStatement:
		v.traverse(n.Expression, expectsStatus)
	case *ast.InfixExpression:
		if n.Operator == "&&" || n.Operator == "||" {
			v.traverse(n.Left, true)
			v.traverse(n.Right, true)
		} else {
			v.traverse(n.Left, false)
			v.traverse(n.Right, false)
		}
	case *ast.PrefixExpression:
		if n.Operator == "!" {
			v.traverse(n.Right, true)
		} else {
			v.traverse(n.Right, false)
		}
	case *ast.GroupedExpression:
		// ( ... ) expression context
		if !expectsStatus {
			v.checkGroupedExpression(n)
		}
		v.traverse(n.Expression, false)
	default:
		// Do nothing for other nodes
	}
}

func (v *zc1088Visitor) checkSubshell(sub *ast.Subshell) {
	if v.isStateChangeOnly(sub.Command) {
		v.violations = append(v.violations, Violation{
			KataID:  "ZC1088",
			Message: "Subshell `( ... )` isolates state changes. The changes (e.g. `cd`, variable assignment) will be lost. Use `{ ... }` to preserve them, or add commands that use the changed state.",
			Line:    sub.TokenLiteralNode().Line,
			Column:  sub.TokenLiteralNode().Column,
			Level:   SeverityWarning,
		})
	}
}

func (v *zc1088Visitor) checkGroupedExpression(group *ast.GroupedExpression) {
	// Similar logic for GroupedExpression if it wraps state changes
	// But GroupedExpression wraps Expression.
	// SimpleCommand is Expression.
	if v.isStateChangeOnly(group.Expression) {
		v.violations = append(v.violations, Violation{
			KataID:  "ZC1088",
			Message: "Subshell `( ... )` isolates state changes. The changes (e.g. `cd`, variable assignment) will be lost. Use `{ ... }` to preserve them, or add commands that use the changed state.",
			Line:    group.TokenLiteralNode().Line,
			Column:  group.TokenLiteralNode().Column,
			Level:   SeverityWarning,
		})
	}
}

func (v *zc1088Visitor) isStateChangeOnly(node ast.Node) bool {
	hasStateChange := false
	hasSideEffect := false

	ast.Walk(node, func(n ast.Node) bool {
		if n == nil || n == node {
			return true
		}

		if hasSideEffect {
			return false
		} // Optimization

		if cmd, ok := n.(*ast.SimpleCommand); ok {
			name := cmd.Name.String()
			if isStateChanger(name) {
				hasStateChange = true
			} else {
				hasSideEffect = true
			}
		} else if infix, ok := n.(*ast.InfixExpression); ok {
			if infix.Operator == "=" {
				hasStateChange = true
			}
			// e.g. a && b. If a is state change, b might be side effect.
			// We just traverse.
		} else if _, ok := n.(*ast.IfStatement); ok {
			hasSideEffect = true
		} else if _, ok := n.(*ast.ForLoopStatement); ok {
			hasSideEffect = true
		} else if _, ok := n.(*ast.WhileLoopStatement); ok {
			hasSideEffect = true
		}
		return true
	})

	return hasStateChange && !hasSideEffect
}

func isStateChanger(name string) bool {
	switch name {
	case "cd", "export", "unset", "alias", "unalias", "declare", "typeset", "local", "shift":
		return true
	}
	return false
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1089",
		Title: "Redirection order matters (`2>&1 > file`)",
		Description: "Redirecting stderr to stdout (`2>&1`) before redirecting stdout to a file (`> file`) " +
			"means stderr goes to the *original* stdout (usually tty), not the file. " +
			"Use `> file 2>&1` or `&> file` to redirect both.",
		Severity: SeverityError,
		Check:    checkZC1089,
	})
}

func checkZC1089(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	idx2to1 := -1
	idxRedirect := -1
	var redirectArg ast.Expression

	for i, arg := range cmd.Arguments {
		s := arg.String()
		if s == "2>&1" {
			if idx2to1 == -1 {
				idx2to1 = i
			}
		} else if s == ">" || s == ">>" {
			// Found redirection operator
			if idxRedirect == -1 {
				idxRedirect = i
				redirectArg = arg
			}
		}
	}

	if idx2to1 != -1 && idxRedirect != -1 && idx2to1 < idxRedirect {
		return []Violation{
			{
				KataID:  "ZC1089",
				Message: "Redirection order matters. `2>&1 > file` does not redirect stderr to file. Use `> file 2>&1` instead.",
				Line:    redirectArg.TokenLiteralNode().Line,
				Column:  redirectArg.TokenLiteralNode().Column,
				Level:   SeverityError,
			},
		}
	}

	return nil
}

func init() {
	RegisterKata(ast.DoubleBracketExpressionNode, Kata{
		ID:    "ZC1090",
		Title: "Quoted regex pattern in `=~`",
		Description: "Quoting the pattern on the right side of `=~` forces literal string matching in Zsh/Bash. " +
			"Regex metacharacters inside quotes will be matched literally. " +
			"Remove quotes to enable regex matching, or use `==` for literal string comparison.",
		Severity: SeverityWarning,
		Check:    checkZC1090,
	})
}

func checkZC1090(node ast.Node) []Violation {
	expr, ok := node.(*ast.DoubleBracketExpression)
	if !ok {
		return nil
	}

	var violations []Violation

	for _, e := range expr.Elements {
		infix, ok := e.(*ast.InfixExpression)
		if !ok {
			continue
		}

		if infix.Operator != "=~" {
			continue
		}

		// Check Right operand
		checkOperand(infix.Right, infix, &violations)
	}

	return violations
}

func checkOperand(node ast.Expression, infix *ast.InfixExpression, violations *[]Violation) {
	switch n := node.(type) {
	case *ast.StringLiteral:
		if containsRegexMeta(n.Value) {
			*violations = append(*violations, Violation{
				KataID:  "ZC1090",
				Message: "Quoted regex pattern matches literally. Remove quotes to enable regex matching.",
				Line:    n.TokenLiteralNode().Line,
				Column:  n.TokenLiteralNode().Column,
				Level:   SeverityWarning,
			})
		}
	case *ast.ConcatenatedExpression:
		for _, part := range n.Parts {
			if sl, ok := part.(*ast.StringLiteral); ok {
				if containsRegexMeta(sl.Value) {
					*violations = append(*violations, Violation{
						KataID:  "ZC1090",
						Message: "Quoted regex pattern matches literally. Remove quotes from the regex part.",
						Line:    sl.TokenLiteralNode().Line,
						Column:  sl.TokenLiteralNode().Column,
						Level:   SeverityWarning,
					})
					return // One violation per expression is enough
				}
			}
		}
	}
}

func containsRegexMeta(s string) bool {
	// Check for regex metacharacters that are likely intended as regex but broken by quotes.
	// ^ $ * + ? [ ( |
	// We exclude . because it's common in text.
	// We exclude $ because it's used for variables (and my parser keeps it in StringLiteral?).
	// Wait, "$var" literal is "$var" or "value"?
	// Parser stores raw literal including quotes usually?
	// Lexer `readString` returns content WITH quotes.
	// So `s` includes quotes!
	// `containsRegexMeta` should check INSIDE quotes.

	if len(s) < 2 {
		return false
	}
	// Strip quotes
	content := s[1 : len(s)-1]

	for _, char := range content {
		switch char {
		case '^', '*', '+', '?', '[', '(', '|':
			return true
			// case '$':
			// 	// $ might be variable. Don't flag.
			// case '.':
			//  // . is common. Don't flag "file.txt".
		}
	}
	// Check for $ at end? `foo$` -> regex end anchor.
	// if strings.HasSuffix(content, "$") { ... } - Removed empty block
	return false
}

func init() {
	RegisterKata(ast.DoubleBracketExpressionNode, Kata{
		ID:    "ZC1091",
		Title: "Use `((...))` for arithmetic comparisons in `[[...]]`",
		Description: "The `[[ ... ]]` construct is primarily for string comparisons and file tests. " +
			"For arithmetic comparisons (`-eq`, `-lt`, etc.), use the dedicated arithmetic context `(( ... ))`. " +
			"It is cleaner and strictly numeric.",
		Severity: SeverityStyle,
		Check:    checkZC1091,
		Fix:      fixZC1091,
	})
}

// fixZC1091 rewrites a bracket conditional that uses dashed
// comparison operators into arithmetic form. Example:
// `[[ x -lt 10 ]]` → `(( x < 10 ))`. Only fires when exactly one
// recognised operator appears inside the brackets to keep the
// rewrite unambiguous.
func fixZC1091(node ast.Node, _ Violation, source []byte) []FixEdit {
	dbe, ok := node.(*ast.DoubleBracketExpression)
	if !ok {
		return nil
	}
	openOff, openLine, openCol, ok := zc1091OpenBracket(source, dbe)
	if !ok {
		return nil
	}
	closeOff := findDoubleBracketClose(source, openOff+2)
	if closeOff < 0 {
		return nil
	}
	infix, ok := zc1091SingleArithOp(dbe)
	if !ok {
		return nil
	}
	closeLine, closeCol := offsetLineColZC1091(source, closeOff)
	if closeLine < 0 {
		return nil
	}
	return []FixEdit{
		{Line: openLine, Column: openCol, Length: 2, Replace: "(("},
		{Line: infix.Token.Line, Column: infix.Token.Column, Length: len(infix.Operator), Replace: arithCmpReplacements[infix.Operator]},
		{Line: closeLine, Column: closeCol, Length: 2, Replace: "))"},
	}
}

func zc1091OpenBracket(source []byte, dbe *ast.DoubleBracketExpression) (off, line, col int, ok bool) {
	line = dbe.Token.Line
	col = dbe.Token.Column
	off = LineColToByteOffset(source, line, col)
	if off < 0 {
		return 0, 0, 0, false
	}
	if off > 0 && source[off] == '[' && source[off-1] == '[' {
		off--
		col--
	}
	if off+2 > len(source) || source[off] != '[' || source[off+1] != '[' {
		return 0, 0, 0, false
	}
	return off, line, col, true
}

func zc1091SingleArithOp(dbe *ast.DoubleBracketExpression) (*ast.InfixExpression, bool) {
	var found *ast.InfixExpression
	for _, el := range dbe.Elements {
		infix, ok := el.(*ast.InfixExpression)
		if !ok {
			continue
		}
		if _, hit := arithCmpReplacements[infix.Operator]; !hit {
			continue
		}
		if found != nil {
			return nil, false
		}
		found = infix
	}
	return found, found != nil
}

// findDoubleBracketClose scans source for the matching `]]` that
// closes the `[[` just before `start`. Honours `[…]` nesting so
// character classes like `[:alnum:]` don't trip the scan.
func findDoubleBracketClose(source []byte, start int) int {
	depth := 0
	for i := start; i < len(source)-1; i++ {
		switch source[i] {
		case '\\':
			i++
		case '[':
			depth++
		case ']':
			if depth > 0 {
				depth--
				continue
			}
			if source[i+1] == ']' {
				return i
			}
		}
	}
	return -1
}

func offsetLineColZC1091(source []byte, offset int) (int, int) {
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

func checkZC1091(node ast.Node) []Violation {
	dbe, ok := node.(*ast.DoubleBracketExpression)
	if !ok {
		return nil
	}

	var violations []Violation

	visitor := func(n ast.Node) bool {
		if infix, ok := n.(*ast.InfixExpression); ok {
			switch infix.Operator {
			case "-eq", "-ne", "-lt", "-le", "-gt", "-ge":
				violations = append(violations, Violation{
					KataID:  "ZC1091",
					Message: "Use `(( ... ))` for arithmetic comparisons. e.g. `(( a < b ))` instead of `[[ a -lt b ]]`.",
					Line:    infix.TokenLiteralNode().Line,
					Column:  infix.TokenLiteralNode().Column,
					Level:   SeverityStyle,
				})
			}
		}
		return true
	}

	for _, expr := range dbe.Elements {
		ast.Walk(expr, visitor)
	}

	return violations
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1092",
		Title: "Prefer `print` or `printf` over `echo` in Zsh",
		Description: "In Zsh, `echo` behavior can vary significantly based on options like `BSD_ECHO`. " +
			"`print` is a builtin with consistent behavior and more features. " +
			"For formatted output, `printf` is preferred.",
		Severity: SeverityWarning,
		Check:    checkZC1092,
		Fix:      fixZC1092,
	})
}

// fixZC1092 rewrites plain `echo ARGS...` -> `print -r -- ARGS...`.
// Only the no-flag form is auto-fixed. When the first argument starts
// with `-` the command is using BSD-style flags (-n / -e / -E) whose
// translation to print differs per flag and is deferred to human
// review. The replacement covers only the command name — arguments
// stay byte-identical so quoting and expansions are preserved.
func fixZC1092(node ast.Node, v Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	if cmd.Name == nil || cmd.Name.String() != "echo" {
		return nil
	}
	// Skip the flagged forms; print's flag semantics differ.
	if len(cmd.Arguments) > 0 {
		if first := cmd.Arguments[0].String(); len(first) > 0 && first[0] == '-' {
			return nil
		}
	}
	return []FixEdit{{
		Line:    v.Line,
		Column:  v.Column,
		Length:  len("echo"),
		Replace: "print -r --",
	}}
}

func checkZC1092(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	if cmd.Name == nil {
		return nil
	}

	if cmd.Name.String() == "echo" {
		// Check if it's just a simple echo or if flags are involved
		// If flags are used (like -n, -e), print is definitely better.
		// Even without flags, print is idiomatic Zsh.

		// We can be slightly lenient and only warn if flags are present OR if it contains backslashes?
		// The prompt suggests "Prefer 'print' over 'echo'". Let's be strict for now as it's "Platinum Standard".

		msg := "Prefer `print` over `echo`. `echo` behavior varies. `print` is the Zsh builtin. Especially with flags, `print -n` or `print -r` is more reliable."

		return []Violation{{
			KataID:  "ZC1092",
			Message: msg,
			Line:    cmd.Token.Line,
			Column:  cmd.Token.Column,
			Level:   SeverityWarning,
		}}
	}

	return nil
}

// Issue #341: ZC1093 fires on the same input as the canonical
// ZC1038 with overlapping advice. Stubbed to a no-op so legacy
// `disabled_katas` lists keep parsing; the detection now lives in
// ZC1038.

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:          "ZC1093",
		Title:       "Superseded by ZC1038 — retired duplicate",
		Severity:    SeverityStyle,
		Description: "Retained as a no-op stub so legacy `.zshellcheckrc` files that disable this ID keep parsing. See https://github.com/afadesigns/zshellcheck/issues/341 for context; the canonical detection lives in ZC1038.",
		Check:       checkZC1093,
	})
}

func checkZC1093(ast.Node) []Violation {
	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1094",
		Title: "Use parameter expansion instead of `sed` for simple substitutions",
		Description: "For simple string substitutions on variables, use Zsh parameter expansion " +
			"`${var//pattern/replacement}` instead of piping through `sed`. It avoids spawning an external process.",
		Severity: SeverityStyle,
		Check:    checkZC1094,
	})
}

func checkZC1094(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok || CommandIdentifier(cmd) != "sed" || len(cmd.Arguments) != 1 {
		return nil
	}
	val := cmd.Arguments[0].String()
	if val != "" && val[0] == '-' {
		return nil
	}
	if !zc1094IsSimpleSubst(val) {
		return nil
	}
	return []Violation{{
		KataID: "ZC1094",
		Message: "Use `${var//pattern/replacement}` instead of piping through `sed` for simple substitutions. " +
			"Parameter expansion avoids spawning an external process.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityStyle,
	}}
}

// zc1094IsSimpleSubst reports whether v looks like a `s/pat/rep/`
// expression, optionally surrounded by single or double quotes.
func zc1094IsSimpleSubst(v string) bool {
	if zc1094IsRawSubst(v) {
		return true
	}
	if len(v) < 6 || (v[0] != '\'' && v[0] != '"') {
		return false
	}
	return zc1094IsRawSubst(v[1 : len(v)-1])
}

func zc1094IsRawSubst(v string) bool {
	if len(v) < 4 || v[0] != 's' {
		return false
	}
	switch v[1] {
	case '/', '|', '#':
		return true
	}
	return false
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1095",
		Title: "Use `repeat N` for simple repetition",
		Description: "Zsh provides `repeat N do ... done` for running a block a fixed number of times. " +
			"It is cleaner than `for i in {1..N}` or C-style for loops when the iterator variable is unused.",
		Severity: SeverityStyle,
		Check:    checkZC1095,
		// Reuse the seq → {start..end} rewrite from ZC1061. The detector
		// here fires on a single-numeric-arg `seq N`, which fixZC1061
		// rewrites to `{1..N}` — exactly the brace expansion this kata
		// suggests for `for i in {1..N}`.
		Fix: fixZC1061,
	})
}

func checkZC1095(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "seq" {
		return nil
	}

	// Flag bare `seq N` calls (often used in `for i in $(seq N)`)
	// Only flag if seq has exactly one numeric argument
	if len(cmd.Arguments) != 1 {
		return nil
	}

	arg := cmd.Arguments[0].String()
	for _, ch := range arg {
		if ch < '0' || ch > '9' {
			return nil
		}
	}

	return []Violation{{
		KataID: "ZC1095",
		Message: "Use `repeat N do ... done` or `for i in {1..N}` instead of `seq N`. " +
			"Zsh has built-in constructs for repetition that avoid spawning an external process.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityStyle,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1096",
		Title: "Warn on `bc` for simple arithmetic",
		Description: "Zsh has built-in support for floating point arithmetic using `(( ... ))` or `$(( ... ))`. " +
			"Using `bc` is often unnecessary and slower.",
		Severity: SeverityStyle,
		Check:    checkZC1096,
	})
}

func checkZC1096(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	if cmd.Name != nil && cmd.Name.String() == "bc" {
		return []Violation{{
			KataID:  "ZC1096",
			Message: "Zsh supports floating point arithmetic natively. You often don't need `bc`.",
			Line:    cmd.Token.Line,
			Column:  cmd.Token.Column,
			Level:   SeverityStyle,
		}}
	}

	return nil
}

func init() {
	RegisterKata(ast.FunctionDefinitionNode, Kata{
		ID:    "ZC1097",
		Title: "Declare loop variables as `local` in functions",
		Description: "Loop variables in `for` loops are global by default in Zsh functions. " +
			"Use `local` to scope them to the function before the loop.",
		Severity: SeverityStyle,
		Check:    checkZC1097,
	})
}

func checkZC1097(node ast.Node) []Violation {
	funcDef, ok := node.(*ast.FunctionDefinition)
	if !ok {
		return nil
	}

	violations := []Violation{}
	locals := make(map[string]bool)

	ast.Walk(funcDef.Body, func(n ast.Node) bool {
		if _, ok := n.(*ast.FunctionDefinition); ok && n != funcDef {
			return false
		}
		zc1043HarvestLocals(n, locals)
		zc1097HarvestDeclLocals(n, locals)
		if v, ok := zc1097UnscopedLoopVar(n, locals); ok {
			violations = append(violations, v)
		}
		return true
	})

	return violations
}

func zc1097HarvestDeclLocals(n ast.Node, locals map[string]bool) {
	decl, ok := n.(*ast.DeclarationStatement)
	if !ok {
		return
	}
	for _, assign := range decl.Assignments {
		if assign.Name != nil {
			locals[assign.Name.String()] = true
		}
	}
}

func zc1097UnscopedLoopVar(n ast.Node, locals map[string]bool) (Violation, bool) {
	loop, ok := n.(*ast.ForLoopStatement)
	if !ok || loop.Name == nil || locals[loop.Name.Value] {
		return Violation{}, false
	}
	return Violation{
		KataID: "ZC1097",
		Message: "Loop variable '" + loop.Name.Value + "' is used without 'local'. It will be global. " +
			"Use `local " + loop.Name.Value + "` before the loop.",
		Line:   loop.Name.Token.Line,
		Column: loop.Name.Token.Column,
		Level:  SeverityStyle,
	}, true
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1098",
		Title: "Use `(q)` flag for quoting variables in eval",
		Description: "When constructing a command string for `eval`, use the `(q)` flag (or `(qq)`, `(q-)`) to safely quote variables " +
			"and prevent command injection.",
		Severity: SeverityStyle,
		Check:    checkZC1098,
	})
}

func checkZC1098(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	if cmd.Name != nil && cmd.Name.String() == "eval" {
		for _, arg := range cmd.Arguments {
			// Check if the argument string contains '$' and NOT '(q)'
			argStr := arg.String()
			// Very rough heuristic. Real parsing inside the string would be better but complex.
			// If arg contains '$' and not `(q)`, warn.
			// Also skip if it contains `(qq)` or `(q-)`.

			// We need to handle the case where user wrote `${(q)var}`.
			// arg.String() would be `${(q)var}`.

			// If we find `$` but no `(q`, warn.

			// Check for variable usage
			hasVar := false
			for i := 0; i < len(argStr); i++ {
				if argStr[i] == '$' {
					hasVar = true
					break
				}
			}

			if hasVar {
				// Check for quoting flags
				if !containsFlag(argStr) {
					return []Violation{{
						KataID:  "ZC1098",
						Message: "Use the `(q)` flag (or `(qq)`, `(q-)`) when using variables in `eval` to prevent injection.",
						Line:    cmd.Token.Line,
						Column:  cmd.Token.Column,
						Level:   SeverityStyle,
					}}
				}
			}
		}
	}

	return nil
}

func containsFlag(s string) bool {
	// Simple check for (q), (qq), (q-)
	// This is not perfect (could be inside a string literal), but good enough for a linter warning.
	// We look for `(q` pattern after `$`.
	// e.g. `${(q)var}` or `$var[(q)...]`? No, flags are at start of expansion.
	// `${(q)...}` or `$(...)` (command subst is also dangerous in eval without q).

	// Let's just check if the string contains "(q".
	for i := 0; i < len(s)-1; i++ {
		if s[i] == '(' && s[i+1] == 'q' {
			return true
		}
	}
	return false
}

func init() {
	RegisterKata(ast.InfixExpressionNode, Kata{
		ID:    "ZC1099",
		Title: "Use `(f)` flag to split lines instead of `while read`",
		Description: "Zsh provides the `(f)` parameter expansion flag to split a string into lines. " +
			"Iterating over `${(f)variable}` is often cleaner and faster than piping to `while read`.",
		Severity: SeverityStyle,
		Check:    checkZC1099,
	})
}

func checkZC1099(node ast.Node) []Violation {
	infix, ok := node.(*ast.InfixExpression)
	if !ok {
		return nil
	}

	if infix.Operator == "|" {
		if whileLoop, ok := infix.Right.(*ast.WhileLoopStatement); ok {
			foundRead := false
			for _, stmt := range whileLoop.Condition.(*ast.BlockStatement).Statements {
				if exprStmt, ok := stmt.(*ast.ExpressionStatement); ok {
					if simpleCmd, ok := exprStmt.Expression.(*ast.SimpleCommand); ok && simpleCmd.Name != nil && simpleCmd.Name.String() == "read" {
						foundRead = true
						break
					}
				}
			}

			if foundRead {
				return []Violation{{
					KataID:  "ZC1099",
					Message: "Consider using `for line in ${(f)variable}` instead of `... | while read line`. It's faster and cleaner in Zsh.",
					Line:    infix.Token.Line,
					Column:  infix.Token.Column,
					Level:   SeverityStyle,
				}}
			}
		}
	}
	return nil
}
