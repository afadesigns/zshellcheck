// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
	"github.com/afadesigns/zshellcheck/pkg/token"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1100",
		Title: "Use parameter expansion instead of `dirname`/`basename`",
		Description: "Zsh parameter expansion `${var%/*}` (dirname) and `${var##*/}` (basename) " +
			"avoid spawning external processes for simple path manipulation.",
		Severity: SeverityStyle,
		Check:    checkZC1100,
	})
}

func checkZC1100(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	name := ident.Value
	if name != "dirname" && name != "basename" {
		return nil
	}

	// Only flag simple single-argument calls
	// basename with -s or -a flags is more complex
	for _, arg := range cmd.Arguments {
		val := arg.String()
		if len(val) > 0 && val[0] == '-' {
			return nil
		}
	}

	if len(cmd.Arguments) != 1 {
		return nil
	}

	var msg string
	if name == "dirname" {
		msg = "Use `${var%/*}` instead of `dirname` to extract the directory path. " +
			"Parameter expansion avoids spawning an external process."
	} else {
		msg = "Use `${var##*/}` instead of `basename` to extract the filename. " +
			"Parameter expansion avoids spawning an external process."
	}

	return []Violation{{
		KataID:  "ZC1100",
		Message: msg,
		Line:    cmd.Token.Line,
		Column:  cmd.Token.Column,
		Level:   SeverityStyle,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1101",
		Title: "Use `$(( ))` instead of `bc` for simple arithmetic",
		Description: "Zsh supports arithmetic expansion with `$(( ))` and floating point via `zmodload zsh/mathfunc`. " +
			"Avoid piping to `bc` for simple calculations.",
		Severity: SeverityStyle,
		Check:    checkZC1101,
	})
}

func checkZC1101(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "bc" {
		return nil
	}

	// bc with file arguments is a valid external use
	for _, arg := range cmd.Arguments {
		val := arg.String()
		if len(val) > 0 && val[0] != '-' {
			return nil
		}
	}

	return []Violation{{
		KataID: "ZC1101",
		Message: "Use `$(( ))` for arithmetic instead of `bc`. " +
			"Zsh arithmetic expansion avoids spawning an external process.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityStyle,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1102",
		Title: "Redirecting output of `sudo` doesn't work as expected",
		Description: "Redirections are performed by the current shell before `sudo` is started. " +
			"So `sudo echo > /root/file` will try to open `/root/file` as the current user, failing. " +
			"Use `echo ... | sudo tee file` or `sudo sh -c 'echo ... > file'`.",
		Severity: SeverityStyle,
		Check:    checkZC1102,
	})
}

func checkZC1102(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	// Check if the command name is 'sudo'
	if cmd.Name != nil && cmd.Name.String() == "sudo" {
		// Scan arguments for output redirection operators
		for _, arg := range cmd.Arguments {
			argStr := arg.String()
			if argStr == ">" || argStr == ">>" {
				return []Violation{{
					KataID:  "ZC1102",
					Message: "Redirection happens before `sudo`. This will likely fail permission checks. Use `| sudo tee`.",
					Line:    cmd.TokenLiteralNode().Line,
					Column:  cmd.TokenLiteralNode().Column,
					Level:   SeverityStyle,
				}}
			}
		}
	}

	return nil
}

func init() {
	RegisterKata(ast.InfixExpressionNode, Kata{
		ID:    "ZC1103",
		Title: "Suggest `path` array instead of `$PATH` string manipulation (direct assignment)",
		Description: "Zsh automatically maps the `$PATH` environment variable to the `$path` array. " +
			"Modifying `$path` is cleaner and less error-prone than manipulating the colon-separated `$PATH` string.",
		Severity: SeverityStyle,
		Check:    checkZC1103,
	})
}

func checkZC1103(node ast.Node) []Violation {
	infixExp, ok := node.(*ast.InfixExpression)
	if !ok {
		return nil
	}

	if infixExp.Operator == "=" {
		if ident, ok := infixExp.Left.(*ast.Identifier); ok && ident.Value == "PATH" {
			// Check if the right-hand side is an old-style PATH manipulation
			if strings.Contains(infixExp.Right.String(), "$PATH") {
				return []Violation{{
					KataID:  "ZC1103",
					Message: "Use the `path` array instead of manually manipulating the `$PATH` string.",
					Line:    infixExp.Token.Line,
					Column:  infixExp.Token.Column,
					Level:   SeverityStyle,
				}}
			}
		}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1104",
		Title: "Suggest `path` array instead of `export PATH` string manipulation",
		Description: "Zsh automatically maps the `$PATH` environment variable to the `$path` array. " +
			"Modifying `$path` is cleaner and less error-prone than manipulating the colon-separated `$PATH` string.",
		Severity: SeverityStyle,
		Check:    checkZC1104,
	})
}

func checkZC1104(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	// Check for `export PATH=...`
	if cmd.Name != nil && cmd.Name.String() == "export" {
		for _, arg := range cmd.Arguments {
			argStr := arg.String()
			if strings.HasPrefix(argStr, "PATH=") {
				return []Violation{{
					KataID:  "ZC1104",
					Message: "Use the `path` array instead of manually manipulating the `$PATH` string.",
					Line:    cmd.Token.Line,
					Column:  cmd.Token.Column,
					Level:   SeverityStyle,
				}}
			}
		}
	}

	return nil
}

func init() {
	RegisterKata(ast.ArithmeticCommandNode, Kata{
		ID:    "ZC1105",
		Title: "Avoid nested arithmetic expansions for clarity",
		Description: "While Zsh supports nested arithmetic expansions like `(( $((...)) ))`, " +
			"they can make code harder to read and reason about. Prefer flatter expressions " +
			"or temporary variables for intermediate results to improve clarity.",
		Severity: SeverityStyle,
		Check:    checkZC1105,
	})
}

func checkZC1105(node ast.Node) []Violation {
	arithCmd, ok := node.(*ast.ArithmeticCommand)
	if !ok {
		return nil
	}

	// Check if the expression contains a nested arithmetic expansion
	// A simplified check: if the string representation contains another $(( or ((
	exprString := arithCmd.Expression.String()
	if strings.Contains(exprString, "$((") || strings.Contains(exprString, "((") {
		return []Violation{{
			KataID:  "ZC1105",
			Message: "Avoid nested arithmetic expansions. Use intermediate variables for clarity.",
			Line:    arithCmd.Token.Line,
			Column:  arithCmd.Token.Column,
			Level:   SeverityStyle,
		}}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1106",
		Title: "Avoid `set -x` in production scripts for sensitive data exposure",
		Description: "Using `set -x` (xtrace) in production environments can expose sensitive " +
			"information, such as API keys or passwords, in logs. While useful for debugging, " +
			"it should be avoided in production. Consider using targeted debugging or secure logging.",
		Severity: SeverityStyle,
		Check:    checkZC1106,
	})
}

func checkZC1106(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	if cmd.Name != nil && cmd.Name.String() == "set" {
		for _, arg := range cmd.Arguments {
			argStr := arg.String()
			argStr = strings.Trim(argStr, "\"'")
			if strings.HasPrefix(argStr, "-") {
				// Check for -x flag explicitly or combined flags like -eux
				if strings.Contains(argStr, "x") {
					return []Violation{{
						KataID:  "ZC1106",
						Message: "Avoid `set -x` in production scripts to prevent sensitive data exposure.",
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

func init() {
	RegisterKata(ast.DoubleBracketExpressionNode, Kata{
		ID:          "ZC1107",
		Title:       "Use (( ... )) for arithmetic conditions",
		Description: "Use `(( ... ))` for arithmetic comparisons instead of `[[ ... -gt ... ]]`. The double parenthesis syntax supports standard math operators (`>`, `<`, `==`, `!=`) and is optimized.",
		Severity:    SeverityStyle,
		Check:       checkZC1107DoubleBracket,
		Fix:         fixZC1091,
	})

	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:          "ZC1107",
		Title:       "Use (( ... )) for arithmetic conditions",
		Description: "Use `(( ... ))` for arithmetic comparisons instead of `[ ... -eq ... ]`. The double parenthesis syntax supports standard math operators (`>`, `<`, `==`, `!=`) and is optimized.",
		Severity:    SeverityStyle,
		Check:       checkZC1107SimpleCommand,
	})
}

func checkZC1107DoubleBracket(node ast.Node) []Violation {
	dbe := node.(*ast.DoubleBracketExpression)
	var violations []Violation

	// Helper to check infix expressions recursively
	check := func(n ast.Node) bool {
		if infix, ok := n.(*ast.InfixExpression); ok {
			switch infix.Operator {
			case "-eq", "-ne", "-lt", "-le", "-gt", "-ge":
				violations = append(violations, Violation{
					KataID:  "ZC1107",
					Message: "Prefer `(( ... ))` for arithmetic comparisons (e.g., `(( a > b ))`) over `[[ ... ]]` with flags like `" + infix.Operator + "`.",
					Line:    infix.TokenLiteralNode().Line,
					Column:  infix.TokenLiteralNode().Column,
					Level:   SeverityStyle,
				})
			}
		}
		return true
	}

	// Walk the elements of the double bracket expression
	for _, el := range dbe.Elements {
		ast.Walk(el, check)
	}

	return violations
}

func checkZC1107SimpleCommand(node ast.Node) []Violation {
	cmd := node.(*ast.SimpleCommand)

	// Check if command is "[" or "test"
	cmdName := cmd.Name.TokenLiteral()
	if cmdName != "[" && cmdName != "test" {
		return nil
	}

	var violations []Violation
	for _, arg := range cmd.Arguments {
		argText := arg.TokenLiteral()
		switch argText {
		case "-eq", "-ne", "-lt", "-le", "-gt", "-ge":
			violations = append(violations, Violation{
				KataID:  "ZC1107",
				Message: "Prefer `(( ... ))` for arithmetic comparisons (e.g., `(( a > b ))`) over `[ ... ]` with flags like `" + argText + "`.",
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
		ID:    "ZC1108",
		Title: "Use Zsh case conversion instead of `tr`",
		Description: "Zsh provides `${(U)var}` for uppercase and `${(L)var}` for lowercase. " +
			"Avoid piping through `tr '[:lower:]' '[:upper:]'` for simple case conversion.",
		Severity: SeverityStyle,
		Check:    checkZC1108,
	})
}

func checkZC1108(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "tr" {
		return nil
	}

	if len(cmd.Arguments) != 2 {
		return nil
	}

	arg1 := strings.Trim(cmd.Arguments[0].String(), "'\"")
	arg2 := strings.Trim(cmd.Arguments[1].String(), "'\"")

	isLowerToUpper := (arg1 == "[:lower:]" && arg2 == "[:upper:]") ||
		(arg1 == "a-z" && arg2 == "A-Z")
	isUpperToLower := (arg1 == "[:upper:]" && arg2 == "[:lower:]") ||
		(arg1 == "A-Z" && arg2 == "a-z")

	if !isLowerToUpper && !isUpperToLower {
		return nil
	}

	var suggestion string
	if isLowerToUpper {
		suggestion = "`${(U)var}`"
	} else {
		suggestion = "`${(L)var}`"
	}

	return []Violation{{
		KataID: "ZC1108",
		Message: "Use " + suggestion + " for case conversion instead of `tr`. " +
			"Zsh parameter expansion flags avoid spawning an external process.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityStyle,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1109",
		Title: "Use parameter expansion instead of `cut` for field extraction",
		Description: "For simple field extraction from variables, use Zsh parameter expansion " +
			"like `${var%%:*}` or `${(s.:.)var}` instead of piping through `cut`.",
		Severity: SeverityStyle,
		Check:    checkZC1109,
	})
}

func checkZC1109(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "cut" {
		return nil
	}

	// Only flag simple cut with -d and -f flags and no file argument
	hasDelimiter := false
	hasField := false
	hasFileArg := false

	for _, arg := range cmd.Arguments {
		val := strings.Trim(arg.String(), "'\"")
		switch {
		case strings.HasPrefix(val, "-d"), strings.HasPrefix(val, "--delimiter"):
			hasDelimiter = true
		case strings.HasPrefix(val, "-f"), strings.HasPrefix(val, "--fields"):
			hasField = true
		case len(val) > 0 && val[0] != '-':
			hasFileArg = true
		}
	}

	if hasFileArg || !hasDelimiter || !hasField {
		return nil
	}

	return []Violation{{
		KataID: "ZC1109",
		Message: "Use Zsh parameter expansion for field extraction instead of `cut`. " +
			"`${var%%delim*}` or `${(s.delim.)var}` avoid spawning an external process.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityStyle,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1110",
		Title: "Use Zsh subscripts instead of `head -1` or `tail -1`",
		Description: "Zsh array subscripts `${lines[1]}` and `${lines[-1]}` can extract the first or last " +
			"element without spawning `head` or `tail` as external processes.",
		Severity: SeverityStyle,
		Check:    checkZC1110,
	})
}

func checkZC1110(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	name := CommandIdentifier(cmd)
	if name != "head" && name != "tail" {
		return nil
	}
	hasFile, isSingleLine := zc1110ScanArgs(cmd)
	if hasFile || !isSingleLine {
		return nil
	}
	suggestion := "`${lines[1]}`"
	if name == "tail" {
		suggestion = "`${lines[-1]}`"
	}
	return []Violation{{
		KataID: "ZC1110",
		Message: "Use " + suggestion + " instead of `" + name + " -1`. " +
			"Zsh array subscripts avoid spawning an external process.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityStyle,
	}}
}

func zc1110ScanArgs(cmd *ast.SimpleCommand) (hasFile, isSingleLine bool) {
	skip := false
	for i, arg := range cmd.Arguments {
		if skip {
			skip = false
			continue
		}
		val := arg.String()
		switch {
		case val == "-1", val == "-n1":
			isSingleLine = true
		case val == "-n" && i+1 < len(cmd.Arguments) && cmd.Arguments[i+1].String() == "1":
			isSingleLine = true
			skip = true
		case len(val) > 0 && val[0] != '-':
			hasFile = true
		}
	}
	return
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1111",
		Title: "Avoid `xargs` for simple command invocation",
		Description: "Zsh can iterate arrays directly with `for` loops or use `${(f)...}` to split " +
			"command output by newlines. Avoid `xargs` when processing lines one at a time.",
		Severity: SeverityStyle,
		Check:    checkZC1111,
	})
}

func checkZC1111(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "xargs" {
		return nil
	}

	// Only flag simple xargs without complex flags
	// -0, -P (parallel), -I (replace string), -L are complex uses — skip them
	for _, arg := range cmd.Arguments {
		val := arg.String()
		if len(val) > 1 && val[0] == '-' {
			switch {
			case val == "-0", val == "--null":
				return nil
			case val == "-P", val == "--max-procs":
				return nil
			case val == "-I", val == "--replace":
				return nil
			case val == "-L", val == "--max-lines":
				return nil
			case val == "-p", val == "--interactive":
				return nil
			}
		}
	}

	return []Violation{{
		KataID: "ZC1111",
		Message: "Consider using Zsh array iteration instead of `xargs`. " +
			"`for item in ${(f)$(cmd)}` splits output by newlines without spawning xargs.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityStyle,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1112",
		Title: "Avoid `grep -c` — use Zsh pattern matching for counting",
		Description: "For counting matches in a variable, use Zsh `${#${(f)...}}` or array filtering " +
			"with `${(M)array:#pattern}` instead of piping through `grep -c`.",
		Severity: SeverityStyle,
		Check:    checkZC1112,
	})
}

func checkZC1112(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "grep" {
		return nil
	}

	// Only flag grep -c without file arguments (pipeline use)
	hasCountFlag := false
	hasFileAfterPattern := false
	patternSeen := false

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if len(val) > 0 && val[0] == '-' {
			if val == "-c" || val == "--count" {
				hasCountFlag = true
			}
		} else {
			if patternSeen {
				hasFileAfterPattern = true
				break
			}
			patternSeen = true
		}
	}

	if !hasCountFlag || hasFileAfterPattern {
		return nil
	}

	return []Violation{{
		KataID: "ZC1112",
		Message: "Use Zsh array filtering `${(M)array:#pattern}` or `${#${(f)...}}` for counting " +
			"instead of `grep -c`. Avoids spawning an external process.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityStyle,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1113",
		Title: "Use `${var:A}` instead of `realpath` or `readlink -f`",
		Description: "Zsh provides the `:A` modifier to resolve a path to its absolute form, " +
			"following symlinks. Avoid spawning `realpath` or `readlink -f` as external processes.",
		Severity: SeverityStyle,
		Check:    checkZC1113,
	})
}

func checkZC1113(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	name := ident.Value

	if name == "realpath" {
		for _, arg := range cmd.Arguments {
			val := arg.String()
			if len(val) > 1 && val[0] == '-' && val != "-s" {
				return nil
			}
		}
		return []Violation{{
			KataID: "ZC1113",
			Message: "Use `${var:A}` instead of `realpath` to resolve absolute paths. " +
				"Zsh path modifiers avoid spawning an external process.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityStyle,
		}}
	}

	if name == "readlink" {
		hasResolveFlag := false
		for _, arg := range cmd.Arguments {
			val := arg.String()
			if val == "-f" || val == "-e" || val == "-m" {
				hasResolveFlag = true
			}
		}
		if !hasResolveFlag {
			return nil
		}
		return []Violation{{
			KataID: "ZC1113",
			Message: "Use `${var:A}` instead of `readlink -f` to resolve absolute paths. " +
				"Zsh path modifiers avoid spawning an external process.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityStyle,
		}}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1114",
		Title: "Consider Zsh `=(...)` for temporary files",
		Description: "Zsh `=(cmd)` creates a temporary file with the command output that is automatically " +
			"cleaned up. Consider this instead of manual `mktemp` and cleanup patterns.",
		Severity: SeverityStyle,
		Check:    checkZC1114,
	})
}

func checkZC1114(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "mktemp" {
		return nil
	}

	// Skip mktemp -d (directory creation — no Zsh equivalent)
	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-d" {
			return nil
		}
	}

	return []Violation{{
		KataID: "ZC1114",
		Message: "Consider using Zsh `=(cmd)` for temporary files instead of `mktemp`. " +
			"Zsh auto-cleans temporary files created with `=(...)` process substitution.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityStyle,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1115",
		Title: "Use Zsh string manipulation instead of `rev`",
		Description: "Zsh can reverse strings using parameter expansion. " +
			"Avoid spawning `rev` as an external process for simple string reversal.",
		Severity: SeverityStyle,
		Check:    checkZC1115,
	})
}

func checkZC1115(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "rev" {
		return nil
	}

	// Only flag rev without file arguments (pipeline use)
	for _, arg := range cmd.Arguments {
		val := arg.String()
		if len(val) > 0 && val[0] != '-' {
			return nil
		}
	}

	return []Violation{{
		KataID: "ZC1115",
		Message: "Use Zsh string manipulation instead of `rev`. " +
			"Parameter expansion can reverse strings without spawning an external process.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityStyle,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1116",
		Title: "Use Zsh multios instead of `tee`",
		Description: "Zsh `setopt multios` allows redirecting output to multiple files with " +
			"`cmd > file1 > file2`. Avoid spawning `tee` for simple output duplication.",
		Severity: SeverityStyle,
		Check:    checkZC1116,
	})
}

func checkZC1116(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "tee" {
		return nil
	}

	// Only flag simple tee without -a (append) or -i (ignore interrupt)
	// tee -a is append mode which multios handles differently
	for _, arg := range cmd.Arguments {
		val := arg.String()
		if len(val) > 0 && val[0] == '-' {
			return nil
		}
	}

	// Must have at least one file argument
	if len(cmd.Arguments) == 0 {
		return nil
	}

	return []Violation{{
		KataID: "ZC1116",
		Message: "Use Zsh multios (`setopt multios`) instead of `tee`. " +
			"With multios, `cmd > file1 > file2` writes to both files without spawning tee.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityStyle,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1117",
		Title: "Use `&!` or `disown` instead of `nohup`",
		Description: "Zsh provides `&!` (shorthand for `& disown`) to run a command in the background " +
			"immune to hangups. Avoid spawning `nohup` as an external process.",
		Severity: SeverityStyle,
		Check:    checkZC1117,
	})
}

func checkZC1117(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "nohup" {
		return nil
	}

	return []Violation{{
		KataID: "ZC1117",
		Message: "Use `cmd &!` or `cmd & disown` instead of `nohup cmd &`. " +
			"Zsh `&!` is a built-in shorthand that avoids spawning nohup.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityStyle,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1118",
		Title: "Use `print -rn` instead of `echo -n`",
		Description: "The behavior of `echo -n` varies across shells and platforms. " +
			"In Zsh, `print -rn` is the reliable way to output text without a trailing newline.",
		Severity: SeverityStyle,
		Check:    checkZC1118,
		Fix:      fixZC1118,
	})
}

// fixZC1118 collapses `echo -n` (with any whitespace between) into
// `print -rn`. Spans the `echo` name, intervening whitespace, and
// the `-n` flag in a single edit; remaining arguments stay in place.
func fixZC1118(node ast.Node, v Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "echo" {
		return nil
	}
	if len(cmd.Arguments) == 0 {
		return nil
	}
	nameOff := LineColToByteOffset(source, v.Line, v.Column)
	if nameOff < 0 || nameOff+len("echo") > len(source) {
		return nil
	}
	if string(source[nameOff:nameOff+len("echo")]) != "echo" {
		return nil
	}
	// Walk forward past whitespace, then expect `-n`.
	i := nameOff + len("echo")
	for i < len(source) && (source[i] == ' ' || source[i] == '\t') {
		i++
	}
	if i+2 > len(source) || source[i] != '-' || source[i+1] != 'n' {
		return nil
	}
	end := i + 2
	return []FixEdit{{
		Line:    v.Line,
		Column:  v.Column,
		Length:  end - nameOff,
		Replace: "print -rn",
	}}
}

func checkZC1118(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "echo" {
		return nil
	}

	if len(cmd.Arguments) == 0 {
		return nil
	}

	firstArg := cmd.Arguments[0].String()
	if firstArg == "-n" {
		return []Violation{{
			KataID: "ZC1118",
			Message: "Use `print -rn` instead of `echo -n`. " +
				"`echo -n` behavior varies across shells; `print -rn` is the reliable Zsh idiom.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityStyle,
		}}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1119",
		Title: "Use `$EPOCHSECONDS` instead of `date +%s`",
		Description: "Zsh provides `$EPOCHSECONDS` and `$EPOCHREALTIME` via `zsh/datetime` module. " +
			"Avoid spawning `date` for simple Unix timestamp retrieval.",
		Severity: SeverityStyle,
		Check:    checkZC1119,
	})
}

func checkZC1119(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "date" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		val := strings.Trim(arg.String(), "'\"")
		if val == "+%s" || val == "+%s%N" {
			return []Violation{{
				KataID: "ZC1119",
				Message: "Use `$EPOCHSECONDS` or `$EPOCHREALTIME` (via `zmodload zsh/datetime`) " +
					"instead of `date +%s`. Avoids spawning an external process.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityStyle,
			}}
		}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1120",
		Title: "Use `$PWD` instead of `pwd`",
		Description: "Zsh maintains `$PWD` as a built-in variable tracking the current directory. " +
			"Avoid spawning `pwd` as an external process.",
		Severity: SeverityStyle,
		Check:    checkZC1120,
	})
}

func checkZC1120(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "pwd" {
		return nil
	}

	// pwd -P (physical) resolves symlinks — $PWD may not
	// Only flag bare pwd or pwd -L (logical, same as $PWD)
	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-P" {
			return nil
		}
	}

	return []Violation{{
		KataID: "ZC1120",
		Message: "Use `$PWD` instead of `pwd`. " +
			"Zsh maintains `$PWD` as a built-in variable, avoiding an external process.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityStyle,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1121",
		Title: "Use `$HOST` instead of `hostname`",
		Description: "Zsh provides `$HOST` as a built-in variable containing the hostname. " +
			"Avoid spawning `hostname` as an external process.",
		Severity: SeverityStyle,
		Check:    checkZC1121,
	})
}

func checkZC1121(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "hostname" {
		return nil
	}

	// hostname with flags like -f, -I, -d does more than $HOST
	for _, arg := range cmd.Arguments {
		val := arg.String()
		if len(val) > 0 && val[0] == '-' {
			return nil
		}
	}

	return []Violation{{
		KataID: "ZC1121",
		Message: "Use `$HOST` instead of `hostname`. " +
			"Zsh maintains `$HOST` as a built-in variable, avoiding an external process.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityStyle,
	}}
}

func init() {
	RegisterKata(ast.IdentifierNode, Kata{
		ID:    "ZC1122",
		Title: "Use `$USER` instead of `whoami`",
		Description: "Zsh provides `$USER` as a built-in variable containing the current username. " +
			"Avoid spawning `whoami` as an external process.",
		Severity: SeverityStyle,
		Check:    checkZC1122,
	})
}

func checkZC1122(node ast.Node) []Violation {
	ident, ok := node.(*ast.Identifier)
	if !ok || ident == nil || ident.Value != "whoami" {
		return nil
	}

	return []Violation{{
		KataID: "ZC1122",
		Message: "Use `$USER` instead of `whoami`. " +
			"Zsh maintains `$USER` as a built-in variable, avoiding an external process.",
		Line:   ident.Token.Line,
		Column: ident.Token.Column,
		Level:  SeverityStyle,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1123",
		Title: "Use `$OSTYPE` instead of `uname`",
		Description: "Zsh provides `$OSTYPE` (e.g., `linux-gnu`, `darwin`) as a built-in variable. " +
			"Avoid spawning `uname` for simple OS detection.",
		Severity: SeverityStyle,
		Check:    checkZC1123,
	})
}

func checkZC1123(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "uname" {
		return nil
	}

	// Only flag simple uname, uname -s, uname -o (OS type detection)
	// Skip uname -r, -m, -a, -n, -p which provide different info
	if len(cmd.Arguments) == 0 {
		return []Violation{{
			KataID: "ZC1123",
			Message: "Use `$OSTYPE` instead of `uname` for OS detection. " +
				"Zsh maintains `$OSTYPE` as a built-in variable.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityStyle,
		}}
	}

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-s" || val == "-o" {
			return []Violation{{
				KataID: "ZC1123",
				Message: "Use `$OSTYPE` instead of `uname -s` for OS detection. " +
					"Zsh maintains `$OSTYPE` as a built-in variable.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityStyle,
			}}
		}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1124",
		Title: "Use `: > file` instead of `cat /dev/null > file` to truncate",
		Description: "Truncating a file with `cat /dev/null > file` spawns an unnecessary process. " +
			"Use `: > file` or simply `> file` in Zsh.",
		Severity: SeverityStyle,
		Check:    checkZC1124,
		Fix:      fixZC1124,
	})
}

// fixZC1124 replaces the `cat /dev/null` prefix with the `:` builtin.
// The redirection and anything following stays in place:
// `cat /dev/null > file` becomes `: > file`. Only fires when
// `/dev/null` is the first argument — the detector already requires
// that shape.
func fixZC1124(node ast.Node, v Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	if len(cmd.Arguments) == 0 {
		return nil
	}
	var devNull ast.Expression
	for _, arg := range cmd.Arguments {
		if arg.String() == "/dev/null" {
			devNull = arg
			break
		}
	}
	if devNull == nil {
		return nil
	}
	nameOff := LineColToByteOffset(source, v.Line, v.Column)
	if nameOff < 0 || nameOff+3 > len(source) {
		return nil
	}
	if string(source[nameOff:nameOff+3]) != "cat" {
		return nil
	}
	argTok := devNull.TokenLiteralNode()
	argOff := LineColToByteOffset(source, argTok.Line, argTok.Column)
	if argOff < 0 {
		return nil
	}
	end := argOff + len("/dev/null")
	if end > len(source) || string(source[argOff:end]) != "/dev/null" {
		return nil
	}
	return []FixEdit{{
		Line:    v.Line,
		Column:  v.Column,
		Length:  end - nameOff,
		Replace: ":",
	}}
}

func checkZC1124(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "cat" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "/dev/null" {
			return []Violation{{
				KataID: "ZC1124",
				Message: "Use `: > file` instead of `cat /dev/null > file` to truncate. " +
					"The `:` builtin avoids spawning cat.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityStyle,
			}}
		}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1125",
		Title: "Avoid `echo | grep` for string matching",
		Description: "Using `echo $var | grep pattern` spawns two unnecessary processes. " +
			"Use Zsh `[[ $var =~ pattern ]]` or `[[ $var == *pattern* ]]` for string matching.",
		Severity: SeverityStyle,
		Check:    checkZC1125,
	})
}

func checkZC1125(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "grep" {
		return nil
	}

	// Only flag grep with -q (quiet) and no file argument
	// grep -q is typically used for string matching in conditionals
	hasQuiet := false
	hasFile := false
	patternSeen := false

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if len(val) > 0 && val[0] == '-' {
			if val == "-q" {
				hasQuiet = true
			}
		} else {
			if patternSeen {
				hasFile = true
				break
			}
			patternSeen = true
		}
	}

	if !hasQuiet || hasFile {
		return nil
	}

	return []Violation{{
		KataID: "ZC1125",
		Message: "Use `[[ $var =~ pattern ]]` or `[[ $var == *pattern* ]]` instead of piping " +
			"through `grep -q`. Zsh pattern matching avoids spawning external processes.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityStyle,
	}}
}

func init() {
	RegisterKata(ast.InfixExpressionNode, Kata{
		ID:    "ZC1126",
		Title: "Use `sort -u` instead of `sort | uniq`",
		Description: "`sort | uniq` spawns two processes when `sort -u` does the same in one. " +
			"Use `sort -u` to deduplicate sorted output efficiently.",
		Severity: SeverityStyle,
		Check:    checkZC1126,
		Fix:      fixZC1126,
	})
}

// fixZC1126 collapses `sort ... | uniq` into `sort -u ...`. Uses a
// single span-replacement from just after the `sort` command name
// through the end of `uniq`, rewriting the region to ` -u` +
// whatever sort args sit between the name and the pipe. Only fires
// when `uniq` has no flags (ZC1126's detector already guards that).
func fixZC1126(node ast.Node, v Violation, source []byte) []FixEdit {
	pipe, ok := node.(*ast.InfixExpression)
	if !ok || pipe.Operator != "|" {
		return nil
	}
	sortCmd, ok := pipe.Left.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	uniqCmd, ok := pipe.Right.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	sortTok := sortCmd.TokenLiteralNode()
	sortNameOff := LineColToByteOffset(source, sortTok.Line, sortTok.Column)
	if sortNameOff < 0 {
		return nil
	}
	sortNameLen := IdentLenAt(source, sortNameOff)
	if sortNameLen == 0 {
		return nil
	}
	spanStart := sortNameOff + sortNameLen

	// Find the pipe byte and walk back past trailing whitespace.
	pipeOff := LineColToByteOffset(source, pipe.Token.Line, pipe.Token.Column)
	if pipeOff < 0 || source[pipeOff] != '|' {
		return nil
	}
	argsEnd := pipeOff
	for argsEnd > spanStart && (source[argsEnd-1] == ' ' || source[argsEnd-1] == '\t') {
		argsEnd--
	}
	middle := string(source[spanStart:argsEnd])

	// End of uniq: the identifier itself; detector forbids flags.
	uniqTok := uniqCmd.TokenLiteralNode()
	uniqOff := LineColToByteOffset(source, uniqTok.Line, uniqTok.Column)
	uniqLen := IdentLenAt(source, uniqOff)
	if uniqOff < 0 || uniqLen == 0 {
		return nil
	}
	spanEnd := uniqOff + uniqLen

	replace := " -u" + middle
	startLine, startCol := offsetLineColZC1126(source, spanStart)
	if startLine < 0 {
		return nil
	}
	return []FixEdit{{
		Line:    startLine,
		Column:  startCol,
		Length:  spanEnd - spanStart,
		Replace: replace,
	}}
}

func offsetLineColZC1126(source []byte, offset int) (int, int) {
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

func checkZC1126(node ast.Node) []Violation {
	pipe, ok := node.(*ast.InfixExpression)
	if !ok || pipe.Operator != "|" {
		return nil
	}

	sortCmd, ok := pipe.Left.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	if !isCommandName(sortCmd, "sort") {
		return nil
	}

	uniqCmd, ok := pipe.Right.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	if !isCommandName(uniqCmd, "uniq") {
		return nil
	}

	// If uniq has flags like -c (count), -d (duplicates), skip
	for _, arg := range uniqCmd.Arguments {
		val := arg.String()
		if len(val) > 0 && val[0] == '-' {
			return nil
		}
	}

	return []Violation{{
		KataID: "ZC1126",
		Message: "Use `sort -u` instead of `sort | uniq`. " +
			"Combining into one command avoids an unnecessary pipeline.",
		Line:   pipe.TokenLiteralNode().Line,
		Column: pipe.TokenLiteralNode().Column,
		Level:  SeverityStyle,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1127",
		Title: "Avoid `ls` for counting files",
		Description: "Using `ls | wc -l` to count files spawns unnecessary processes. " +
			"Use Zsh glob qualifiers: `files=(*(N)); echo ${#files}` for file counting.",
		Severity: SeverityStyle,
		Check:    checkZC1127,
	})
}

func checkZC1127(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "ls" {
		return nil
	}

	// Flag ls -1 (single column listing, typically used for counting)
	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-1" {
			return []Violation{{
				KataID: "ZC1127",
				Message: "Use Zsh glob qualifiers `files=(*(N)); echo ${#files}` instead of `ls -1 | wc -l`. " +
					"Avoids spawning external processes for file counting.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityStyle,
			}}
		}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1128",
		Title: "Use `> file` instead of `touch file` for creation",
		Description: "If the goal is to create an empty file, `> file` does it without " +
			"spawning `touch`. Use `touch` only when you need to update timestamps.",
		Severity: SeverityStyle,
		Check:    checkZC1128,
		Fix:      fixZC1128,
	})
}

// fixZC1128 rewrites `touch file` into `> file`. Detector already
// guards against flagged forms (timestamp updates) and multi-arg
// invocations, so the fix covers only the single-file case.
func fixZC1128(node ast.Node, v Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	if len(cmd.Arguments) != 1 {
		return nil
	}
	nameOff := LineColToByteOffset(source, v.Line, v.Column)
	if nameOff < 0 || nameOff+len("touch") > len(source) {
		return nil
	}
	if string(source[nameOff:nameOff+len("touch")]) != "touch" {
		return nil
	}
	return []FixEdit{{
		Line:    v.Line,
		Column:  v.Column,
		Length:  len("touch"),
		Replace: ">",
	}}
}

func checkZC1128(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "touch" {
		return nil
	}

	// Skip touch with flags (timestamps: -t, -d, -r, -a, -m)
	for _, arg := range cmd.Arguments {
		val := arg.String()
		if len(val) > 0 && val[0] == '-' {
			return nil
		}
	}

	// Only flag touch with a single file argument
	if len(cmd.Arguments) != 1 {
		return nil
	}

	return []Violation{{
		KataID: "ZC1128",
		Message: "Use `> file` instead of `touch file` to create an empty file. " +
			"This avoids spawning an external process.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityStyle,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1129",
		Title: "Use Zsh `stat` module instead of `wc -c` for file size",
		Description: "Zsh's `zstat` (via `zmodload zsh/stat`) provides file size without " +
			"spawning `wc`. Use `zstat +size file` for efficient file size queries.",
		Severity: SeverityStyle,
		Check:    checkZC1129,
	})
}

func checkZC1129(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "wc" {
		return nil
	}

	hasCharFlag := false
	hasFile := false

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-c" {
			hasCharFlag = true
		} else if len(val) > 0 && val[0] != '-' {
			hasFile = true
		}
	}

	if !hasCharFlag || !hasFile {
		return nil
	}

	return []Violation{{
		KataID: "ZC1129",
		Message: "Use `zstat +size file` (via `zmodload zsh/stat`) instead of `wc -c file`. " +
			"Avoids reading the entire file for a simple size query.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityStyle,
	}}
}

func init() {
	RegisterKata(ast.InfixExpressionNode, Kata{
		ID:    "ZC1131",
		Title: "Avoid `cat file | while read` — use redirection",
		Description: "`cat file | while read line` spawns an unnecessary cat process " +
			"and runs the loop in a subshell. Use `while read line; do ...; done < file` instead.",
		Severity: SeverityStyle,
		Check:    checkZC1131,
	})
}

func checkZC1131(node ast.Node) []Violation {
	pipe, ok := node.(*ast.InfixExpression)
	if !ok || pipe.Operator != "|" {
		return nil
	}

	catCmd, ok := pipe.Left.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := catCmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "cat" {
		return nil
	}

	// cat must have exactly one file argument and no flags
	if len(catCmd.Arguments) != 1 {
		return nil
	}
	for _, arg := range catCmd.Arguments {
		val := arg.String()
		if len(val) > 0 && val[0] == '-' {
			return nil
		}
	}

	// Right side should involve 'while' or 'read'
	if right, ok := pipe.Right.(*ast.SimpleCommand); ok {
		rightIdent, ok := right.Name.(*ast.Identifier)
		if ok && rightIdent.Value == "read" {
			return []Violation{{
				KataID: "ZC1131",
				Message: "Use `while read line; do ...; done < file` instead of `cat file | while read line`. " +
					"Avoids unnecessary cat and subshell from the pipe.",
				Line:   pipe.TokenLiteralNode().Line,
				Column: pipe.TokenLiteralNode().Column,
				Level:  SeverityStyle,
			}}
		}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1132",
		Title: "Use Zsh pattern extraction instead of `grep -o`",
		Description: "For extracting matching parts from variables, use Zsh `${(M)var:#pattern}` " +
			"or `${match[1]}` with `=~` instead of piping through `grep -o`.",
		Severity: SeverityStyle,
		Check:    checkZC1132,
	})
}

func checkZC1132(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "grep" {
		return nil
	}

	hasOnlyMatching := false
	hasFile := false
	patternSeen := false

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if len(val) > 0 && val[0] == '-' {
			if val == "-o" {
				hasOnlyMatching = true
			}
		} else {
			if patternSeen {
				hasFile = true
				break
			}
			patternSeen = true
		}
	}

	if !hasOnlyMatching || hasFile {
		return nil
	}

	return []Violation{{
		KataID: "ZC1132",
		Message: "Use Zsh pattern extraction `${(M)var:#pattern}` or `[[ $var =~ regex ]] && echo $match[1]` " +
			"instead of piping through `grep -o`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityStyle,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1133",
		Title: "Avoid `kill -9` — use `kill` first, then escalate",
		Description: "`kill -9` (SIGKILL) cannot be caught or ignored. Always try `kill` (SIGTERM) first " +
			"to allow the process to clean up, then use `kill -9` only as a last resort.",
		Severity: SeverityStyle,
		Check:    checkZC1133,
	})
}

func checkZC1133(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "kill" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-9" || val == "-KILL" || val == "-SIGKILL" {
			return []Violation{{
				KataID: "ZC1133",
				Message: "Avoid `kill -9` as a first resort. Use `kill` (SIGTERM) first " +
					"to allow graceful shutdown, then escalate to `kill -9` if needed.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityStyle,
			}}
		}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1134",
		Title: "Avoid `sleep` in tight loops",
		Description: "Using `sleep` inside a loop for polling creates busy-wait patterns. " +
			"Consider `inotifywait`, `zle`, or event-driven approaches instead.",
		Severity: SeverityStyle,
		Check:    checkZC1134,
	})
}

func checkZC1134(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "sleep" {
		return nil
	}

	// Flag sleep with very short intervals (0, 0.1, 0.5, 1)
	if len(cmd.Arguments) != 1 {
		return nil
	}

	val := cmd.Arguments[0].String()
	if val == "0" || val == "0.1" || val == "0.01" || val == "0.5" {
		return []Violation{{
			KataID: "ZC1134",
			Message: "Avoid `sleep " + val + "` in loops. Short sleep intervals suggest busy-waiting. " +
				"Consider event-driven alternatives like `inotifywait` or `zle -F`.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityStyle,
		}}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1135",
		Title: "Avoid `env VAR=val cmd` — use inline assignment",
		Description: "Zsh supports inline environment variable assignment with `VAR=val cmd`. " +
			"Avoid spawning `env` for simple variable-prefixed command execution.",
		Severity: SeverityStyle,
		Check:    checkZC1135,
		Fix:      fixZC1135,
	})
}

// fixZC1135 strips the `env ` prefix from `env VAR=val cmd`. Detector
// already forbids `env` flags, so the remaining args form a valid
// inline-assignment command.
func fixZC1135(node ast.Node, v Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "env" {
		return nil
	}
	nameOff := LineColToByteOffset(source, v.Line, v.Column)
	if nameOff < 0 || nameOff+len("env") > len(source) {
		return nil
	}
	if string(source[nameOff:nameOff+len("env")]) != "env" {
		return nil
	}
	// Span covers `env` plus the whitespace that follows it.
	end := nameOff + len("env")
	for end < len(source) && (source[end] == ' ' || source[end] == '\t') {
		end++
	}
	return []FixEdit{{
		Line:    v.Line,
		Column:  v.Column,
		Length:  end - nameOff,
		Replace: "",
	}}
}

func checkZC1135(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "env" {
		return nil
	}

	// Only flag env with VAR=val patterns followed by a command
	// Skip env -i (clean environment), env -u (unset), env -S
	for _, arg := range cmd.Arguments {
		val := arg.String()
		if len(val) > 0 && val[0] == '-' {
			return nil
		}
	}

	// Check if any argument contains = (env var assignment)
	for _, arg := range cmd.Arguments {
		val := arg.String()
		if strings.Contains(val, "=") {
			return []Violation{{
				KataID: "ZC1135",
				Message: "Use inline `VAR=val cmd` instead of `env VAR=val cmd`. " +
					"Zsh supports inline env assignment without spawning env.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityStyle,
			}}
		}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1136",
		Title: "Avoid `rm -rf` without safeguard",
		Description: "`rm -rf` with a variable path is dangerous if the variable is empty. " +
			"Always validate the path or use `${var:?}` to fail on empty values.",
		Severity: SeverityWarning,
		Check:    checkZC1136,
	})
}

func checkZC1136(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "rm" {
		return nil
	}

	hasRecursiveForce := false
	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-rf" || val == "-fr" {
			hasRecursiveForce = true
			break
		}
	}

	if !hasRecursiveForce {
		return nil
	}

	// Check if any argument is a bare variable (unprotected)
	for _, arg := range cmd.Arguments {
		val := arg.String()
		if len(val) > 0 && val[0] == '$' {
			return []Violation{{
				KataID: "ZC1136",
				Message: "Avoid `rm -rf $var` without safeguards. Use `rm -rf ${var:?}` " +
					"to abort if the variable is empty, preventing accidental deletion.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1137",
		Title: "Avoid hardcoded `/tmp` paths",
		Description: "Hardcoded `/tmp` paths are predictable and may cause race conditions " +
			"or symlink attacks. Use `mktemp` or Zsh `=(...)` for safe temp files.",
		Severity: SeverityStyle,
		Check:    checkZC1137,
	})
}

func checkZC1137(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	// Skip mktemp itself (it creates temp files properly)
	cmdIdent, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if cmdIdent.Value == "mktemp" || cmdIdent.Value == "cd" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		val := arg.String()
		// Match /tmp/something with a predictable name
		if len(val) > 5 && val[:5] == "/tmp/" && val != "/tmp" {
			// Skip if it uses a variable (dynamic path)
			hasVar := false
			for _, ch := range val {
				if ch == '$' {
					hasVar = true
					break
				}
			}
			if !hasVar {
				return []Violation{{
					KataID: "ZC1137",
					Message: "Avoid hardcoded `/tmp/` paths. Use `mktemp` or Zsh `=(cmd)` " +
						"for temp files to prevent race conditions and symlink attacks.",
					Line:   cmd.Token.Line,
					Column: cmd.Token.Column,
					Level:  SeverityStyle,
				}}
			}
		}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1139",
		Title: "Avoid `source` with URL — use local files",
		Description: "Sourcing scripts from URLs (curl | source) is a security risk. " +
			"Download, verify, then source local files.",
		Severity: SeverityWarning,
		Check:    checkZC1139,
	})
}

func checkZC1139(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	if ident.Value != "source" && ident.Value != "." {
		return nil
	}

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if len(val) > 8 && (val[:8] == "https://" || val[:7] == "http://") {
			return []Violation{{
				KataID: "ZC1139",
				Message: "Avoid sourcing scripts from URLs. Download, verify integrity, " +
					"then source from a local path to prevent supply-chain attacks.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1140",
		Title: "Use `command -v` instead of `hash` for command existence",
		Description: "`hash cmd` is a POSIX way to check command existence but provides " +
			"poor error messages. Use `command -v cmd` for cleaner checks in Zsh.",
		Severity: SeverityStyle,
		Check:    checkZC1140,
		Fix:      fixZC1140,
	})
}

// fixZC1140 rewrites `hash cmd` to `command -v cmd`. Single-edit
// command-name replacement — arguments stay intact. Detector gates
// on flagged forms like `hash -r`, so those stay as-is.
func fixZC1140(node ast.Node, v Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "hash" {
		return nil
	}
	return []FixEdit{{
		Line:    v.Line,
		Column:  v.Column,
		Length:  len("hash"),
		Replace: "command -v",
	}}
}

func checkZC1140(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "hash" {
		return nil
	}

	// Only flag bare hash (command existence check)
	// hash -r (rehash) is a different valid use
	for _, arg := range cmd.Arguments {
		val := strings.Trim(arg.String(), "'\"")
		if len(val) > 0 && val[0] == '-' {
			return nil
		}
	}

	if len(cmd.Arguments) != 1 {
		return nil
	}

	return []Violation{{
		KataID: "ZC1140",
		Message: "Use `command -v cmd` instead of `hash cmd` for command existence checks. " +
			"`command -v` provides clearer semantics in Zsh.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityStyle,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1141",
		Title: "Avoid `curl | sh` pattern",
		Description: "Piping curl output to sh/bash/zsh is a security risk. Download first, " +
			"verify integrity (checksum or signature), then execute.",
		Severity: SeverityWarning,
		Check:    checkZC1141,
	})
}

func checkZC1141(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "curl" {
		return nil
	}

	// Check for -s or -sSL flags which suggest piping intent
	hasSilent := false
	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-s" || val == "-sS" || val == "-sSL" || val == "-sL" {
			hasSilent = true
		}
	}

	if !hasSilent {
		return nil
	}

	return []Violation{{
		KataID: "ZC1141",
		Message: "Avoid `curl -s URL | sh`. Download the script first, verify its integrity, " +
			"then execute. Piping directly from the internet is a supply-chain risk.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.InfixExpressionNode, Kata{
		ID:       "ZC1142",
		Title:    "Avoid chained `grep | grep` — combine patterns",
		Severity: SeverityStyle,
		Description: "Chaining `grep pattern1 | grep pattern2` spawns multiple processes. " +
			"Use `grep -E 'p1.*p2|p2.*p1'` or `awk` for multi-pattern matching.",
		Check: checkZC1142,
	})
}

func checkZC1142(node ast.Node) []Violation {
	pipe, ok := node.(*ast.InfixExpression)
	if !ok || pipe.Operator != "|" {
		return nil
	}

	leftCmd, ok := pipe.Left.(*ast.SimpleCommand)
	if !ok || !isCommandName(leftCmd, "grep") {
		return nil
	}

	rightCmd, ok := pipe.Right.(*ast.SimpleCommand)
	if !ok || !isCommandName(rightCmd, "grep") {
		return nil
	}

	return []Violation{{
		KataID: "ZC1142",
		Message: "Avoid chaining `grep | grep`. Combine into a single `grep -E` with alternation " +
			"or use `awk` for multi-pattern matching to reduce pipeline processes.",
		Line:   pipe.TokenLiteralNode().Line,
		Column: pipe.TokenLiteralNode().Column,
		Level:  SeverityStyle,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1143",
		Title:    "Avoid `set -e` — use explicit error handling",
		Severity: SeverityInfo,
		Description: "`set -e` (errexit) has surprising behavior in Zsh with conditionals, " +
			"pipes, and subshells. Use explicit `|| return` or `|| exit` for reliable error handling.",
		Check: checkZC1143,
	})
}

func checkZC1143(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "set" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-e" || val == "-o" {
			// Check for set -o errexit pattern
			if val == "-e" {
				return []Violation{{
					KataID: "ZC1143",
					Message: "Avoid `set -e`. It has surprising behavior with conditionals and subshells in Zsh. " +
						"Use explicit error handling with `cmd || return 1` instead.",
					Line:   cmd.Token.Line,
					Column: cmd.Token.Column,
					Level:  SeverityInfo,
				}}
			}
		}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1144",
		Title:    "Avoid `trap` with signal numbers — use names",
		Severity: SeverityInfo,
		Description: "Signal numbers vary across platforms. Use signal names like " +
			"`SIGTERM`, `SIGINT`, `EXIT` instead of numeric values for portability.",
		Check: checkZC1144,
		Fix:   fixZC1144,
	})
}

// zc1144SignalNames maps POSIX signal numbers to their canonical
// names. Numbers above 31 (realtime signals) aren't included because
// their names vary across platforms and are rarely used with `trap`.
var zc1144SignalNames = map[string]string{
	"1":  "HUP",
	"2":  "INT",
	"3":  "QUIT",
	"4":  "ILL",
	"5":  "TRAP",
	"6":  "ABRT",
	"7":  "BUS",
	"8":  "FPE",
	"9":  "KILL",
	"10": "USR1",
	"11": "SEGV",
	"12": "USR2",
	"13": "PIPE",
	"14": "ALRM",
	"15": "TERM",
	"17": "CHLD",
	"18": "CONT",
	"19": "STOP",
	"20": "TSTP",
	"21": "TTIN",
	"22": "TTOU",
	"23": "URG",
	"24": "XCPU",
	"25": "XFSZ",
	"26": "VTALRM",
	"27": "PROF",
	"28": "WINCH",
	"29": "IO",
	"30": "PWR",
	"31": "SYS",
}

// fixZC1144 replaces numeric signal arguments in a `trap` call with
// their canonical names. Each numeric arg becomes a separate edit at
// that arg's position. Unknown numbers stay untouched.
func fixZC1144(node ast.Node, _ Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "trap" {
		return nil
	}
	var edits []FixEdit
	for i := 1; i < len(cmd.Arguments); i++ {
		arg := cmd.Arguments[i]
		val := arg.String()
		name, ok := zc1144SignalNames[val]
		if !ok {
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
			Replace: name,
		})
	}
	return edits
}

func checkZC1144(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "trap" {
		return nil
	}

	if len(cmd.Arguments) < 2 {
		return nil
	}

	// Check last arguments for numeric signal values
	for i := 1; i < len(cmd.Arguments); i++ {
		val := cmd.Arguments[i].String()
		// Numeric signals: 1-31
		isNumeric := len(val) > 0
		for _, ch := range val {
			if ch < '0' || ch > '9' {
				isNumeric = false
				break
			}
		}
		if isNumeric && val != "0" {
			return []Violation{{
				KataID: "ZC1144",
				Message: "Use signal names (`SIGTERM`, `SIGINT`, `EXIT`) instead of numbers in `trap`. " +
					"Signal numbers vary across platforms.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityInfo,
			}}
		}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1145",
		Title:    "Avoid `tr -d` for character deletion — use parameter expansion",
		Severity: SeverityStyle,
		Description: "For simple character deletion from variables, use Zsh `${var//char/}` " +
			"instead of piping through `tr -d`.",
		Check: checkZC1145,
	})
}

func checkZC1145(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "tr" {
		return nil
	}

	if len(cmd.Arguments) < 2 {
		return nil
	}

	firstArg := cmd.Arguments[0].String()
	if firstArg == "-d" {
		// Check second arg is a simple single char
		secondArg := cmd.Arguments[1].String()
		if len(secondArg) <= 3 { // Simple char like 'x' or " "
			return []Violation{{
				KataID: "ZC1145",
				Message: "Use `${var//char/}` instead of piping through `tr -d`. " +
					"Parameter expansion is faster for simple character deletion.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityStyle,
			}}
		}
	}

	return nil
}

func init() {
	RegisterKata(ast.InfixExpressionNode, Kata{
		ID:       "ZC1146",
		Title:    "Avoid `cat file | awk` — pass file to awk directly",
		Severity: SeverityStyle,
		Description: "`cat file | awk` spawns an unnecessary cat process. " +
			"Pass the file directly as `awk '...' file`.",
		Check: checkZC1146,
		Fix:   fixZC1146,
	})
}

// fixZC1146 collapses `cat FILE | tool [args]` into `tool [args] FILE`.
// One span replacement runs from the start of `cat` through the end of
// the right-hand command; the replacement is the right-hand source
// verbatim with ` FILE` appended. Only fires when the cat command has
// exactly one filename argument (the detector already guards that).
func fixZC1146(node ast.Node, _ Violation, source []byte) []FixEdit {
	_, catCmd, rightCmd, _, ok := zc1146Pipe(node)
	if !ok {
		return nil
	}
	catStart, ok := zc1146Offset(source, catCmd.TokenLiteralNode())
	if !ok {
		return nil
	}
	fileLit, _, ok := zc1146ArgSlice(source, catCmd.Arguments[0])
	if !ok {
		return nil
	}
	rightStart, ok := zc1146Offset(source, rightCmd.TokenLiteralNode())
	if !ok {
		return nil
	}
	rightEnd, ok := zc1146RightEnd(source, rightCmd, rightStart)
	if !ok {
		return nil
	}
	startLine, startCol := offsetLineColZC1146(source, catStart)
	if startLine < 0 {
		return nil
	}
	return []FixEdit{{
		Line:    startLine,
		Column:  startCol,
		Length:  rightEnd - catStart,
		Replace: string(source[rightStart:rightEnd]) + " " + fileLit,
	}}
}

func zc1146Offset(source []byte, tok token.Token) (int, bool) {
	off := LineColToByteOffset(source, tok.Line, tok.Column)
	return off, off >= 0
}

// zc1146ArgSlice returns the literal text of arg as it appears in
// source plus the offset, or ok=false when the AST coordinates do not
// line up with the source bytes.
func zc1146ArgSlice(source []byte, arg ast.Expression) (lit string, off int, ok bool) {
	tok := arg.TokenLiteralNode()
	off, ok = zc1146Offset(source, tok)
	if !ok {
		return "", 0, false
	}
	lit = arg.String()
	if off+len(lit) > len(source) || string(source[off:off+len(lit)]) != lit {
		return "", 0, false
	}
	return lit, off, true
}

func zc1146RightEnd(source []byte, rightCmd *ast.SimpleCommand, rightStart int) (int, bool) {
	rightIdent, ok := rightCmd.Name.(*ast.Identifier)
	if !ok {
		return 0, false
	}
	end := rightStart + len(rightIdent.Value)
	if n := len(rightCmd.Arguments); n > 0 {
		lastArg := rightCmd.Arguments[n-1]
		laOff, ok := zc1146Offset(source, lastArg.TokenLiteralNode())
		if !ok {
			return 0, false
		}
		end = laOff + len(lastArg.String())
	}
	if end > len(source) || end < rightStart {
		return 0, false
	}
	return end, true
}

func offsetLineColZC1146(source []byte, offset int) (int, int) {
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

var zc1146FileTakers = map[string]struct{}{
	"awk":  {},
	"sed":  {},
	"sort": {},
	"head": {},
	"tail": {},
}

func checkZC1146(node ast.Node) []Violation {
	pipe, _, _, name, ok := zc1146Pipe(node)
	if !ok {
		return nil
	}
	if _, hit := zc1146FileTakers[name]; !hit {
		return nil
	}
	return []Violation{{
		KataID: "ZC1146",
		Message: "Pass the file directly to `" + name + "` instead of `cat file | " + name + "`. " +
			"Most text-processing tools accept file arguments.",
		Line:   pipe.TokenLiteralNode().Line,
		Column: pipe.TokenLiteralNode().Column,
		Level:  SeverityStyle,
	}}
}

// zc1146Pipe destructures `cat FILE | NAME [args]` into its parts and
// reports whether the cat side is well-formed (single non-flag arg).
func zc1146Pipe(node ast.Node) (pipe *ast.InfixExpression, catCmd, rightCmd *ast.SimpleCommand, name string, ok bool) {
	pipe, isPipe := node.(*ast.InfixExpression)
	if !isPipe || pipe.Operator != "|" {
		return nil, nil, nil, "", false
	}
	catCmd, isCat := pipe.Left.(*ast.SimpleCommand)
	if !isCat || !isCommandName(catCmd, "cat") || len(catCmd.Arguments) != 1 {
		return nil, nil, nil, "", false
	}
	if first := catCmd.Arguments[0].String(); first != "" && first[0] == '-' {
		return nil, nil, nil, "", false
	}
	rightCmd, isRight := pipe.Right.(*ast.SimpleCommand)
	if !isRight {
		return nil, nil, nil, "", false
	}
	rightIdent, isIdent := rightCmd.Name.(*ast.Identifier)
	if !isIdent {
		return nil, nil, nil, "", false
	}
	return pipe, catCmd, rightCmd, rightIdent.Value, true
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1147",
		Title:    "Avoid `mkdir` without `-p` for nested paths",
		Severity: SeverityInfo,
		Description: "Using `mkdir` without `-p` fails if parent directories don't exist. " +
			"Use `mkdir -p` to create the full path safely.",
		Check: checkZC1147,
		Fix:   fixZC1147,
	})
}

// fixZC1147 inserts ` -p` after the `mkdir` command name so nested
// paths survive missing intermediates. Detector already gates on
// absence of `-p` and presence of a nested path.
func fixZC1147(node ast.Node, v Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "mkdir" {
		return nil
	}
	nameOff := LineColToByteOffset(source, v.Line, v.Column)
	if nameOff < 0 {
		return nil
	}
	nameLen := IdentLenAt(source, nameOff)
	if nameLen != len("mkdir") {
		return nil
	}
	insertAt := nameOff + nameLen
	insLine, insCol := offsetLineColZC1147(source, insertAt)
	if insLine < 0 {
		return nil
	}
	return []FixEdit{{
		Line:    insLine,
		Column:  insCol,
		Length:  0,
		Replace: " -p",
	}}
}

func offsetLineColZC1147(source []byte, offset int) (int, int) {
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

func checkZC1147(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "mkdir" {
		return nil
	}

	hasParentFlag := false
	hasNestedPath := false

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-p" {
			hasParentFlag = true
		}
		// Check for paths with multiple slashes (nested)
		if len(val) > 0 && val[0] != '-' {
			slashCount := 0
			for _, ch := range val {
				if ch == '/' {
					slashCount++
				}
			}
			if slashCount >= 2 {
				hasNestedPath = true
			}
		}
	}

	if hasParentFlag || !hasNestedPath {
		return nil
	}

	return []Violation{{
		KataID: "ZC1147",
		Message: "Use `mkdir -p` when creating nested directories. " +
			"Without `-p`, `mkdir` fails if parent directories don't exist.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityInfo,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1148",
		Title:    "Use `compdef` instead of `compctl` for completions",
		Severity: SeverityInfo,
		Description: "`compctl` is the old Zsh completion system. " +
			"Use `compdef` with the new completion system (`compsys`) for modern Zsh.",
		Check: checkZC1148,
	})
}

func checkZC1148(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "compctl" {
		return nil
	}

	return []Violation{{
		KataID: "ZC1148",
		Message: "Use `compdef` instead of `compctl`. The `compctl` system is deprecated; " +
			"use `compinit` and `compdef` for modern Zsh completions.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityInfo,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1149",
		Title:    "Avoid `echo` for error messages — use `>&2`",
		Severity: SeverityInfo,
		Description: "Error messages should go to stderr, not stdout. " +
			"Use `print -u2` or `echo ... >&2` to ensure errors are properly separated.",
		Check: checkZC1149,
	})
}

func checkZC1149(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	if ident.Value != "echo" && ident.Value != "printf" && ident.Value != "print" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		val := arg.String()
		// Check for error-like messages
		if len(val) > 5 {
			clean := val
			if len(clean) > 2 && (clean[0] == '\'' || clean[0] == '"') {
				clean = clean[1 : len(clean)-1]
			}
			if len(clean) >= 5 && (clean[:5] == "Error" || clean[:5] == "error" || clean[:5] == "ERROR") {
				return []Violation{{
					KataID: "ZC1149",
					Message: "Error messages should go to stderr. Use `print -u2` or append `>&2` " +
						"to separate error output from normal stdout.",
					Line:   cmd.Token.Line,
					Column: cmd.Token.Column,
					Level:  SeverityInfo,
				}}
			}
		}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1151",
		Title:    "Avoid `cat -A` — use `print -v` or od for non-printable characters",
		Severity: SeverityStyle,
		Description: "`cat -A` shows non-printable characters but varies across platforms. " +
			"Use Zsh `print -v` or `od -c` for reliable non-printable character inspection.",
		Check: checkZC1151,
	})
}

func checkZC1151(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "cat" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-A" || val == "-v" || val == "-e" {
			return []Violation{{
				KataID: "ZC1151",
				Message: "Avoid `cat " + val + "` for inspecting non-printable characters. " +
					"Use `od -c` or `hexdump -C` for reliable cross-platform output.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityStyle,
			}}
		}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1152",
		Title:    "Use Zsh PCRE module instead of `grep -P`",
		Severity: SeverityStyle,
		Description: "`grep -P` (Perl regex) is not available on all platforms (e.g., macOS). " +
			"Use `zmodload zsh/pcre` and `pcre_compile`/`pcre_match` for portable PCRE matching.",
		Check: checkZC1152,
	})
}

func checkZC1152(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "grep" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-P" {
			return []Violation{{
				KataID: "ZC1152",
				Message: "Avoid `grep -P` — it's unavailable on macOS. Use `zmodload zsh/pcre` " +
					"with `pcre_compile`/`pcre_match` or `grep -E` for portable regex matching.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityStyle,
			}}
		}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1153",
		Title:    "Use `cmp -s` instead of `diff` for equality check",
		Severity: SeverityStyle,
		Description: "When only checking if two files are identical (not viewing differences), " +
			"`cmp -s` is faster than `diff` as it stops at the first difference.",
		Check: checkZC1153,
		Fix:   fixZC1153,
	})
}

// fixZC1153 rewrites `diff -q FILE1 FILE2` into `cmp -s FILE1 FILE2`.
// Two non-overlapping edits: the command name (`diff` → `cmp`) and the
// quiet flag (`-q` → `-s`). Other arguments stay byte-identical.
// Idempotent because the detector gates on `diff -q` literal presence.
func fixZC1153(node ast.Node, v Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "diff" {
		return nil
	}
	var dashQ ast.Expression
	for _, arg := range cmd.Arguments {
		if arg.String() == "-q" {
			dashQ = arg
			break
		}
	}
	if dashQ == nil {
		return nil
	}
	nameOff := LineColToByteOffset(source, v.Line, v.Column)
	if nameOff < 0 || nameOff+len("diff") > len(source) {
		return nil
	}
	if string(source[nameOff:nameOff+len("diff")]) != "diff" {
		return nil
	}
	dashTok := dashQ.TokenLiteralNode()
	dashOff := LineColToByteOffset(source, dashTok.Line, dashTok.Column)
	if dashOff < 0 || dashOff+2 > len(source) {
		return nil
	}
	if string(source[dashOff:dashOff+2]) != "-q" {
		return nil
	}
	return []FixEdit{
		{Line: v.Line, Column: v.Column, Length: len("diff"), Replace: "cmp"},
		{Line: dashTok.Line, Column: dashTok.Column, Length: 2, Replace: "-s"},
	}
}

func checkZC1153(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "diff" {
		return nil
	}

	// Only flag diff -q (quiet) which is used for equality checks
	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-q" {
			return []Violation{{
				KataID: "ZC1153",
				Message: "Use `cmp -s file1 file2` instead of `diff -q`. " +
					"`cmp -s` is faster for equality checks as it stops at the first difference.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityStyle,
			}}
		}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1154",
		Title:    "Use `find -exec {} +` instead of `find -exec {} \\;`",
		Severity: SeverityStyle,
		Description: "`find -exec cmd {} \\;` runs cmd once per file. " +
			"`find -exec cmd {} +` batches files into fewer invocations, improving performance.",
		Check: checkZC1154,
	})
}

func checkZC1154(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "find" {
		return nil
	}

	for i, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-exec" {
			// Check if the exec block ends with \; (not +)
			for j := i + 1; j < len(cmd.Arguments); j++ {
				endVal := cmd.Arguments[j].String()
				if endVal == ";" {
					return []Violation{{
						KataID: "ZC1154",
						Message: "Use `find -exec cmd {} +` instead of `find -exec cmd {} \\;`. " +
							"The `+` form batches files for fewer process invocations.",
						Line:   cmd.Token.Line,
						Column: cmd.Token.Column,
						Level:  SeverityStyle,
					}}
				}
				if endVal == "+" {
					break
				}
			}
		}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1155",
		Title:    "Use `whence -a` instead of `which -a`",
		Severity: SeverityInfo,
		Description: "`which -a` may be an external command on some systems. " +
			"Zsh builtin `whence -a` reliably lists all command locations.",
		Check: checkZC1155,
		Fix:   fixZC1155,
	})
}

// fixZC1155 rewrites the `which` command name to `whence`, leaving the
// `-a` flag and any other arguments in place. Detector already
// guarantees the shape (which + -a anywhere in argv).
func fixZC1155(node ast.Node, v Violation, _ []byte) []FixEdit {
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

func checkZC1155(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "which" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-a" {
			return []Violation{{
				KataID: "ZC1155",
				Message: "Use `whence -a` instead of `which -a`. " +
					"Zsh `whence` is a reliable builtin for listing all command locations.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityInfo,
			}}
		}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1156",
		Title:    "Avoid `ln` without `-s` for symlinks",
		Severity: SeverityInfo,
		Description: "Hard links (`ln` without `-s`) share inodes and can cause confusion. " +
			"Prefer symbolic links (`ln -s`) unless you specifically need hard links.",
		Check: checkZC1156,
	})
}

func checkZC1156(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "ln" {
		return nil
	}

	hasSymlink := false
	hasForce := false
	fileCount := 0

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-s" || val == "-sf" || val == "-snf" {
			hasSymlink = true
		}
		if val == "-f" {
			hasForce = true
		}
		if len(val) > 0 && val[0] != '-' {
			fileCount++
		}
	}

	_ = hasForce

	if hasSymlink || fileCount < 2 {
		return nil
	}

	return []Violation{{
		KataID: "ZC1156",
		Message: "Use `ln -s` for symbolic links instead of hard links. " +
			"Hard links share inodes and don't work across filesystems.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityInfo,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1157",
		Title:    "Avoid `strings` command — use Zsh `${(ps:\\0:)var}`",
		Severity: SeverityStyle,
		Description: "The `strings` command extracts printable strings from binaries. " +
			"For simple filtering, Zsh parameter expansion with `(ps:\\0:)` can split on null bytes.",
		Check: checkZC1157,
	})
}

func checkZC1157(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "strings" {
		return nil
	}

	// Only flag simple strings without special flags
	for _, arg := range cmd.Arguments {
		val := arg.String()
		if len(val) > 0 && val[0] == '-' {
			return nil
		}
	}

	if len(cmd.Arguments) != 1 {
		return nil
	}

	return []Violation{{
		KataID: "ZC1157",
		Message: "Consider Zsh parameter expansion for string extraction from variables. " +
			"`strings` is typically needed only for binary file analysis.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityStyle,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1158",
		Title:    "Avoid `chown -R` without `--no-dereference`",
		Severity: SeverityWarning,
		Description: "`chown -R` follows symlinks by default, potentially changing ownership " +
			"outside the intended directory. Use `--no-dereference` or `-h` to avoid this.",
		Check: checkZC1158,
	})
}

func checkZC1158(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "chown" {
		return nil
	}

	hasRecursive := false
	hasSafe := false

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-R" {
			hasRecursive = true
		}
		if val == "-h" || val == "--no-dereference" {
			hasSafe = true
		}
	}

	if hasRecursive && !hasSafe {
		return []Violation{{
			KataID: "ZC1158",
			Message: "Use `chown -Rh` or `chown -R --no-dereference` to prevent following " +
				"symlinks during recursive ownership changes.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityWarning,
		}}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1159",
		Title:    "Avoid `tar` without explicit compression flag",
		Severity: SeverityInfo,
		Description: "Use explicit compression flags (`-z` for gzip, `-j` for bzip2, `-J` for xz) " +
			"instead of relying on `tar` auto-detection for clarity and portability.",
		Check: checkZC1159,
	})
}

var (
	zc1159CreateFlags = map[string]struct{}{"-c": {}, "-cf": {}, "cf": {}}
	zc1159CompressionFlags = map[string]struct{}{
		"-z": {}, "-j": {}, "-J": {},
		"--gzip": {}, "--bzip2": {}, "--xz": {},
	}
)

func checkZC1159(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok || CommandIdentifier(cmd) != "tar" {
		return nil
	}
	hasCreate, hasCompression := zc1159ScanFlags(cmd.Arguments)
	if !hasCreate || hasCompression {
		return nil
	}
	return []Violation{{
		KataID: "ZC1159",
		Message: "Specify an explicit compression flag (`-z`, `-j`, `-J`) when creating tar archives. " +
			"Relying on auto-detection reduces clarity and portability.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityInfo,
	}}
}

func zc1159ScanFlags(args []ast.Expression) (hasCreate, hasCompression bool) {
	for _, arg := range args {
		val := arg.String()
		if _, hit := zc1159CreateFlags[val]; hit {
			hasCreate = true
		}
		if _, hit := zc1159CompressionFlags[val]; hit {
			hasCompression = true
		}
		c, z := zc1159BundleFlags(val)
		hasCreate = hasCreate || c
		hasCompression = hasCompression || z
	}
	return
}

func zc1159BundleFlags(val string) (hasCreate, hasCompression bool) {
	if len(val) <= 1 || val[0] == '-' {
		return
	}
	for _, ch := range val {
		switch ch {
		case 'c':
			hasCreate = true
		case 'z', 'j', 'J':
			hasCompression = true
		}
	}
	return
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1160",
		Title:    "Prefer `curl` over `wget` for portability",
		Severity: SeverityStyle,
		Description: "`wget` is not installed by default on macOS. " +
			"`curl` is available on virtually all Unix systems and is more portable.",
		Check: checkZC1160,
	})
}

func checkZC1160(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "wget" {
		return nil
	}

	return []Violation{{
		KataID: "ZC1160",
		Message: "Prefer `curl` over `wget` for portability. " +
			"`curl` is pre-installed on macOS and most Linux distributions.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityStyle,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1161",
		Title:    "Avoid `openssl` for simple hashing — use Zsh modules",
		Severity: SeverityStyle,
		Description: "For simple SHA/MD5 hashing, Zsh provides `zmodload zsh/sha256` and " +
			"`zmodload zsh/md5`. Avoid spawning `openssl` or `sha256sum` for basic hash operations.",
		Check: checkZC1161,
	})
}

func checkZC1161(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	name := ident.Value
	if name != "sha256sum" && name != "sha1sum" && name != "md5sum" && name != "md5" {
		return nil
	}

	// Only flag when used without file arguments (pipeline usage)
	for _, arg := range cmd.Arguments {
		val := arg.String()
		if len(val) > 0 && val[0] != '-' {
			return nil
		}
	}

	return []Violation{{
		KataID: "ZC1161",
		Message: "Consider `zmodload zsh/sha256` or `zmodload zsh/md5` for hash operations. " +
			"Zsh modules avoid spawning external hashing processes.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityStyle,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1162",
		Title:    "Use `cp -a` instead of `cp -r` to preserve attributes",
		Severity: SeverityInfo,
		Description: "`cp -r` copies recursively but may not preserve permissions, timestamps, " +
			"or symlinks. Use `cp -a` (archive mode) to preserve all attributes.",
		Check: checkZC1162,
		Fix:   fixZC1162,
	})
}

// fixZC1162 rewrites `cp -r` / `cp -R` to `cp -a`. Single-edit
// replacement of the recursive flag; surrounding args stay put.
func fixZC1162(node ast.Node, _ Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "cp" {
		return nil
	}
	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val != "-r" && val != "-R" {
			continue
		}
		tok := arg.TokenLiteralNode()
		off := LineColToByteOffset(source, tok.Line, tok.Column)
		if off < 0 || off+2 > len(source) {
			return nil
		}
		if string(source[off:off+2]) != val {
			return nil
		}
		return []FixEdit{{
			Line:    tok.Line,
			Column:  tok.Column,
			Length:  2,
			Replace: "-a",
		}}
	}
	return nil
}

func checkZC1162(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "cp" {
		return nil
	}

	hasRecursive := false
	hasArchive := false

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-r" || val == "-R" {
			hasRecursive = true
		}
		if val == "-a" || val == "-rp" || val == "-Rp" {
			hasArchive = true
		}
	}

	if hasRecursive && !hasArchive {
		return []Violation{{
			KataID: "ZC1162",
			Message: "Use `cp -a` instead of `cp -r` to preserve permissions, timestamps, and symlinks. " +
				"Archive mode ensures a faithful copy.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityInfo,
		}}
	}

	return nil
}

func init() {
	RegisterKata(ast.InfixExpressionNode, Kata{
		ID:       "ZC1163",
		Title:    "Use `grep -m 1` instead of `grep | head -1`",
		Severity: SeverityStyle,
		Description: "`grep pattern | head -1` spawns two processes when `grep -m 1` does the same. " +
			"The `-m` flag stops after the first match, avoiding the pipeline.",
		Check: checkZC1163,
		Fix:   fixZC1163,
	})
}

// fixZC1163 collapses `grep PAT | head -1` into `grep -m 1 PAT`. Span
// runs from just after the `grep` command name through the end of the
// `head -1` invocation; the replacement preserves every original grep
// argument verbatim and drops the pipe + head suffix in one edit. Only
// fires for the `-1` / `-n1` shapes the detector already guards.
var zc1163FirstFlags = map[string]struct{}{"-1": {}, "-n1": {}}

func fixZC1163(node ast.Node, _ Violation, source []byte) []FixEdit {
	grepCmd, headCmd, pipe, ok := zc1163Pipeline(node)
	if !ok {
		return nil
	}
	spanStart, ok := zc1163GrepArgsStart(source, grepCmd)
	if !ok {
		return nil
	}
	middle, ok := zc1163GrepArgsSlice(source, pipe, spanStart)
	if !ok {
		return nil
	}
	spanEnd, ok := zc1163HeadEnd(source, headCmd)
	if !ok || spanEnd <= spanStart {
		return nil
	}
	startLine, startCol := offsetLineColZC1163(source, spanStart)
	if startLine < 0 {
		return nil
	}
	return []FixEdit{{
		Line:    startLine,
		Column:  startCol,
		Length:  spanEnd - spanStart,
		Replace: " -m 1" + middle,
	}}
}

func zc1163Pipeline(node ast.Node) (*ast.SimpleCommand, *ast.SimpleCommand, *ast.InfixExpression, bool) {
	pipe, ok := node.(*ast.InfixExpression)
	if !ok || pipe.Operator != "|" {
		return nil, nil, nil, false
	}
	grepCmd, ok := pipe.Left.(*ast.SimpleCommand)
	if !ok || !isCommandName(grepCmd, "grep") {
		return nil, nil, nil, false
	}
	headCmd, ok := pipe.Right.(*ast.SimpleCommand)
	if !ok || !isCommandName(headCmd, "head") || len(headCmd.Arguments) == 0 {
		return nil, nil, nil, false
	}
	if !HasArgFlag(headCmd, zc1163FirstFlags) {
		return nil, nil, nil, false
	}
	return grepCmd, headCmd, pipe, true
}

func zc1163GrepArgsStart(source []byte, grepCmd *ast.SimpleCommand) (int, bool) {
	tok := grepCmd.TokenLiteralNode()
	off := LineColToByteOffset(source, tok.Line, tok.Column)
	if off < 0 {
		return 0, false
	}
	n := IdentLenAt(source, off)
	if n == 0 {
		return 0, false
	}
	return off + n, true
}

func zc1163GrepArgsSlice(source []byte, pipe *ast.InfixExpression, spanStart int) (string, bool) {
	pipeOff := LineColToByteOffset(source, pipe.Token.Line, pipe.Token.Column)
	if pipeOff < 0 || pipeOff >= len(source) || source[pipeOff] != '|' {
		return "", false
	}
	end := pipeOff
	for end > spanStart && (source[end-1] == ' ' || source[end-1] == '\t') {
		end--
	}
	return string(source[spanStart:end]), true
}

func zc1163HeadEnd(source []byte, headCmd *ast.SimpleCommand) (int, bool) {
	last := headCmd.Arguments[len(headCmd.Arguments)-1]
	tok := last.TokenLiteralNode()
	off := LineColToByteOffset(source, tok.Line, tok.Column)
	if off < 0 {
		return 0, false
	}
	return off + len(last.String()), true
}

func offsetLineColZC1163(source []byte, offset int) (int, int) {
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

func checkZC1163(node ast.Node) []Violation {
	pipe, ok := node.(*ast.InfixExpression)
	if !ok || pipe.Operator != "|" {
		return nil
	}

	grepCmd, ok := pipe.Left.(*ast.SimpleCommand)
	if !ok || !isCommandName(grepCmd, "grep") {
		return nil
	}

	headCmd, ok := pipe.Right.(*ast.SimpleCommand)
	if !ok || !isCommandName(headCmd, "head") {
		return nil
	}

	// Check head has -1 or -n 1
	for _, arg := range headCmd.Arguments {
		val := arg.String()
		if val == "-1" || val == "-n1" {
			return []Violation{{
				KataID: "ZC1163",
				Message: "Use `grep -m 1` instead of `grep | head -1`. " +
					"The `-m` flag stops after the first match without a pipeline.",
				Line:   pipe.TokenLiteralNode().Line,
				Column: pipe.TokenLiteralNode().Column,
				Level:  SeverityStyle,
			}}
		}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1164",
		Title:    "Avoid `sed -n 'Np'` — use Zsh array subscript",
		Severity: SeverityStyle,
		Description: "Extracting a specific line with `sed -n 'Np'` spawns a process. " +
			"Use Zsh array subscript `${lines[N]}` after splitting with `${(f)...}`.",
		Check: checkZC1164,
	})
}

func checkZC1164(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "sed" {
		return nil
	}

	if len(cmd.Arguments) < 2 {
		return nil
	}

	first := cmd.Arguments[0].String()
	if first != "-n" {
		return nil
	}

	// Check if second arg matches pattern like '3p', '10p', etc.
	second := strings.Trim(cmd.Arguments[1].String(), "'\"")
	if len(second) >= 2 && second[len(second)-1] == 'p' {
		allDigits := true
		for _, ch := range second[:len(second)-1] {
			if ch < '0' || ch > '9' {
				allDigits = false
				break
			}
		}
		if allDigits && len(cmd.Arguments) == 2 {
			return []Violation{{
				KataID: "ZC1164",
				Message: "Use Zsh array subscript `${lines[N]}` instead of `sed -n 'Np'`. " +
					"Split input with `${(f)...}` then index directly.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityStyle,
			}}
		}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1165",
		Title:    "Use Zsh parameter expansion for simple `awk` field extraction",
		Severity: SeverityStyle,
		Description: "Simple `awk '{print $1}'` or `awk '{print $NF}'` can often be replaced with " +
			"Zsh parameter expansion `${var%% *}` (first field) or `${var##* }` (last field).",
		Check: checkZC1165,
	})
}

func checkZC1165(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "awk" {
		return nil
	}

	// Only flag awk with a single print statement and no file argument
	if len(cmd.Arguments) != 1 {
		return nil
	}

	arg := strings.Trim(cmd.Arguments[0].String(), "'\"")
	if arg == "{print $1}" || arg == "{print $NF}" {
		return []Violation{{
			KataID: "ZC1165",
			Message: "Use Zsh parameter expansion (`${var%% *}` or `${var##* }`) instead of " +
				"`awk '{print $1}'` for simple field extraction without spawning awk.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityStyle,
		}}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1166",
		Title:    "Avoid `grep -i` for case-insensitive match — use `(#i)` glob flag",
		Severity: SeverityStyle,
		Description: "Zsh provides the `(#i)` glob flag for case-insensitive matching. " +
			"For variable matching, use `[[ $var == (#i)pattern ]]` instead of piping through grep -i.",
		Check: checkZC1166,
	})
}

func checkZC1166(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "grep" {
		return nil
	}

	hasCaseInsensitive := false
	hasFile := false
	patternSeen := false

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if len(val) > 0 && val[0] == '-' {
			if val == "-i" {
				hasCaseInsensitive = true
			}
		} else {
			if patternSeen {
				hasFile = true
				break
			}
			patternSeen = true
		}
	}

	if !hasCaseInsensitive || hasFile {
		return nil
	}

	return []Violation{{
		KataID: "ZC1166",
		Message: "Use Zsh `(#i)` glob flag for case-insensitive matching instead of piping through `grep -i`. " +
			"Example: `[[ $var == (#i)pattern ]]`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityStyle,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1167",
		Title:    "Avoid `timeout` command — use Zsh `TMOUT` or `zsh/sched`",
		Severity: SeverityStyle,
		Description: "`timeout` is not available on all systems (macOS lacks it by default). " +
			"Use Zsh `TMOUT` variable or `zmodload zsh/sched` for timeout functionality.",
		Check: checkZC1167,
	})
}

func checkZC1167(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "timeout" {
		return nil
	}

	return []Violation{{
		KataID: "ZC1167",
		Message: "Avoid `timeout` — it's unavailable on macOS. Use Zsh `TMOUT` variable " +
			"or `zmodload zsh/sched` for portable timeout functionality.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityStyle,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1168",
		Title:    "Use `${(f)...}` instead of `readarray`/`mapfile`",
		Severity: SeverityStyle,
		Description: "`readarray` and `mapfile` are Bash builtins not available in Zsh. " +
			"Use Zsh `${(f)...}` parameter expansion flag to split output into an array by newlines.",
		Check: checkZC1168,
	})
}

func checkZC1168(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	if ident.Value != "readarray" && ident.Value != "mapfile" {
		return nil
	}

	return []Violation{{
		KataID: "ZC1168",
		Message: "Use Zsh `${(f)$(cmd)}` instead of `" + ident.Value + "`. " +
			"`readarray`/`mapfile` are Bash builtins not available in Zsh.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityStyle,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1169",
		Title:    "Avoid `install` for simple copy+chmod — use `cp` then `chmod`",
		Severity: SeverityStyle,
		Description: "`install` command is less common and may confuse readers. " +
			"For clarity, use separate `cp` and `chmod` commands or `install` only in Makefiles.",
		Check: checkZC1169,
	})
}

func checkZC1169(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "install" {
		return nil
	}

	// Only flag install with -m (mode) flag in scripts
	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-m" {
			return []Violation{{
				KataID: "ZC1169",
				Message: "Consider using `cp` + `chmod` instead of `install -m`. " +
					"Separate commands are clearer in shell scripts.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityStyle,
			}}
		}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1170",
		Title:    "Avoid `pushd`/`popd` without `-q` flag",
		Severity: SeverityStyle,
		Description: "`pushd` and `popd` print the directory stack by default, cluttering output. " +
			"Use `-q` flag to suppress output in scripts.",
		Check: checkZC1170,
		Fix:   fixZC1170,
	})
}

// fixZC1170 inserts ` -q` after `pushd` or `popd` so the directory
// stack output is suppressed in scripts.
func fixZC1170(node ast.Node, v Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || (ident.Value != "pushd" && ident.Value != "popd") {
		return nil
	}
	nameOff := LineColToByteOffset(source, v.Line, v.Column)
	if nameOff < 0 {
		return nil
	}
	nameLen := IdentLenAt(source, nameOff)
	if nameLen != len(ident.Value) {
		return nil
	}
	insertAt := nameOff + nameLen
	insLine, insCol := offsetLineColZC1170(source, insertAt)
	if insLine < 0 {
		return nil
	}
	return []FixEdit{{
		Line:    insLine,
		Column:  insCol,
		Length:  0,
		Replace: " -q",
	}}
}

func offsetLineColZC1170(source []byte, offset int) (int, int) {
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

func checkZC1170(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	if ident.Value != "pushd" && ident.Value != "popd" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-q" {
			return nil
		}
	}

	return []Violation{{
		KataID: "ZC1170",
		Message: "Use `" + ident.Value + " -q` to suppress directory stack output in scripts. " +
			"Without `-q`, the stack is printed on every call.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityStyle,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1171",
		Title:    "Use `print` instead of `echo -e` for escape sequences",
		Severity: SeverityStyle,
		Description: "`echo -e` behavior varies across shells and platforms. " +
			"In Zsh, `print` natively interprets escape sequences and is more reliable.",
		Check: checkZC1171,
		Fix:   fixZC1171,
	})
}

// fixZC1171 collapses `echo -e` into `print`. Span covers the
// command name, intervening whitespace, and the `-e` flag; remaining
// arguments stay in place.
func fixZC1171(node ast.Node, v Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "echo" {
		return nil
	}
	if len(cmd.Arguments) == 0 {
		return nil
	}
	nameOff := LineColToByteOffset(source, v.Line, v.Column)
	if nameOff < 0 || nameOff+len("echo") > len(source) {
		return nil
	}
	if string(source[nameOff:nameOff+len("echo")]) != "echo" {
		return nil
	}
	i := nameOff + len("echo")
	for i < len(source) && (source[i] == ' ' || source[i] == '\t') {
		i++
	}
	if i+2 > len(source) || source[i] != '-' || source[i+1] != 'e' {
		return nil
	}
	end := i + 2
	return []FixEdit{{
		Line:    v.Line,
		Column:  v.Column,
		Length:  end - nameOff,
		Replace: "print",
	}}
}

func checkZC1171(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "echo" {
		return nil
	}

	if len(cmd.Arguments) == 0 {
		return nil
	}

	first := cmd.Arguments[0].String()
	if first == "-e" {
		return []Violation{{
			KataID: "ZC1171",
			Message: "Use `print` instead of `echo -e`. Zsh `print` natively interprets " +
				"escape sequences and is more portable than `echo -e`.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityStyle,
		}}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1172",
		Title:    "Use `read -A` instead of Bash `read -a` for arrays",
		Severity: SeverityInfo,
		Description: "Bash uses `read -a` to read into an array, but Zsh uses `read -A`. " +
			"Using `-a` in Zsh reads into a scalar, not an array.",
		Check: checkZC1172,
		Fix:   fixZC1172,
	})
}

// fixZC1172 swaps the lowercase `-a` flag for the uppercase `-A` Zsh
// equivalent. Single-byte replacement at the argument's column.
// Idempotent: a re-run sees `-A`, not `-a`, so the detector won't
// fire. Defensive byte-match guard refuses to insert unless the
// source at the offset is literally `-a`.
func fixZC1172(node ast.Node, _ Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "read" {
		return nil
	}
	for _, arg := range cmd.Arguments {
		if arg.String() != "-a" {
			continue
		}
		tok := arg.TokenLiteralNode()
		off := LineColToByteOffset(source, tok.Line, tok.Column)
		if off < 0 || off+2 > len(source) {
			return nil
		}
		if string(source[off:off+2]) != "-a" {
			return nil
		}
		return []FixEdit{{
			Line:    tok.Line,
			Column:  tok.Column,
			Length:  2,
			Replace: "-A",
		}}
	}
	return nil
}

func checkZC1172(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "read" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-a" {
			return []Violation{{
				KataID: "ZC1172",
				Message: "Use `read -A` instead of `read -a` in Zsh. " +
					"The `-a` flag is Bash syntax; Zsh uses `-A` to read into arrays.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityInfo,
			}}
		}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1173",
		Title:    "Avoid `column` command — use Zsh `print -C` for columnar output",
		Severity: SeverityStyle,
		Description: "Zsh `print -C N` formats output into N columns natively. " +
			"Avoid spawning `column` as an external process for simple tabulation.",
		Check: checkZC1173,
	})
}

func checkZC1173(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "column" {
		return nil
	}

	// Only flag simple column usage (column -t is the most common)
	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-t" {
			return []Violation{{
				KataID: "ZC1173",
				Message: "Use Zsh `print -C N` for columnar output instead of `column -t`. " +
					"The `print` builtin formats columns without spawning an external process.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityStyle,
			}}
		}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1174",
		Title:    "Use Zsh `${(j:delim:)}` instead of `paste -sd`",
		Severity: SeverityStyle,
		Description: "Zsh `${(j:delim:)array}` joins array elements with a delimiter. " +
			"Avoid spawning `paste` for simple field joining from variables.",
		Check: checkZC1174,
	})
}

func checkZC1174(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "paste" {
		return nil
	}

	hasSD := false
	hasFile := false

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-sd" || val == "-s" {
			hasSD = true
		}
		if len(val) > 0 && val[0] != '-' {
			hasFile = true
		}
	}

	if hasSD && !hasFile {
		return []Violation{{
			KataID: "ZC1174",
			Message: "Use Zsh `${(j:delim:)array}` to join array elements instead of `paste -sd`. " +
				"Parameter expansion avoids spawning an external process.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityStyle,
		}}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1175",
		Title:    "Avoid `tput` for simple ANSI colors — use Zsh `%F{color}`",
		Severity: SeverityStyle,
		Description: "Zsh prompt expansion `%F{red}` and `%f` handle colors natively. " +
			"Avoid spawning `tput` for simple color output in prompts and scripts.",
		Check: checkZC1175,
	})
}

func checkZC1175(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "tput" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "setaf" || val == "setab" || val == "sgr0" || val == "bold" {
			return []Violation{{
				KataID: "ZC1175",
				Message: "Use Zsh `%F{color}` / `%f` or `$fg[color]` / `$reset_color` instead of `tput`. " +
					"Zsh handles ANSI colors natively without spawning external processes.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityStyle,
			}}
		}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1176",
		Title:    "Use `zparseopts` instead of `getopt`/`getopts`",
		Severity: SeverityStyle,
		Description: "Zsh provides `zparseopts` for powerful option parsing with long options, " +
			"arrays, and defaults. Avoid `getopt`/`getopts` which are less capable in Zsh.",
		Check: checkZC1176,
	})
}

func checkZC1176(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	if ident.Value != "getopt" && ident.Value != "getopts" {
		return nil
	}

	return []Violation{{
		KataID: "ZC1176",
		Message: "Use Zsh `zparseopts` instead of `" + ident.Value + "`. " +
			"`zparseopts` supports long options, arrays, and is the native Zsh approach.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityStyle,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1177",
		Title:    "Avoid `id -u` — use Zsh `$UID` or `$EUID`",
		Severity: SeverityStyle,
		Description: "Zsh provides `$UID` and `$EUID` as built-in variables for user/effective " +
			"user ID. Avoid spawning `id` for simple UID checks.",
		Check: checkZC1177,
	})
}

func checkZC1177(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "id" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-u" || val == "-un" {
			return []Violation{{
				KataID: "ZC1177",
				Message: "Use `$UID` or `$EUID` instead of `id -u`. " +
					"Zsh provides these as built-in variables.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityStyle,
			}}
		}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1178",
		Title:    "Avoid `stty` for terminal size — use Zsh `$COLUMNS`/`$LINES`",
		Severity: SeverityStyle,
		Description: "Zsh maintains `$COLUMNS` and `$LINES` as built-in variables tracking " +
			"terminal dimensions. Avoid spawning `stty` or `tput` for size queries.",
		Check: checkZC1178,
	})
}

func checkZC1178(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "stty" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "size" {
			return []Violation{{
				KataID: "ZC1178",
				Message: "Use `$COLUMNS` and `$LINES` instead of `stty size`. " +
					"Zsh tracks terminal dimensions as built-in variables.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityStyle,
			}}
		}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1179",
		Title:    "Use Zsh `strftime` instead of `date` for formatting",
		Severity: SeverityStyle,
		Description: "Zsh provides `strftime` via `zmodload zsh/datetime` for date formatting. " +
			"Avoid spawning `date` for simple timestamp formatting.",
		Check: checkZC1179,
	})
}

func checkZC1179(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "date" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		val := strings.Trim(arg.String(), "'\"")
		if len(val) > 1 && val[0] == '+' && val != "+%s" && val != "+%s%N" {
			return []Violation{{
				KataID: "ZC1179",
				Message: "Use `strftime` (via `zmodload zsh/datetime`) instead of `date +" + val[1:] + "`. " +
					"Zsh date formatting avoids spawning an external process.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityStyle,
			}}
		}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1180",
		Title:    "Avoid `pgrep` for own background jobs — use Zsh job control",
		Severity: SeverityInfo,
		Description: "For managing your own background jobs, use Zsh job control (`jobs`, `kill %N`, " +
			"`fg`, `bg`) instead of `pgrep`/`pkill` which search system-wide.",
		Check: checkZC1180,
	})
}

func checkZC1180(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	if ident.Value != "pgrep" && ident.Value != "pkill" {
		return nil
	}

	// Only flag simple pgrep/pkill without complex flags
	for _, arg := range cmd.Arguments {
		val := arg.String()
		if len(val) > 1 && val[0] == '-' && val != "-f" {
			return nil
		}
	}

	return []Violation{{
		KataID: "ZC1180",
		Message: "For own background jobs, use Zsh job control (`jobs`, `kill %N`) instead of `" +
			ident.Value + "`. Job control is more precise for script-spawned processes.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityInfo,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1181",
		Title:    "Avoid `xdg-open`/`open` — use `$BROWSER` for portability",
		Severity: SeverityInfo,
		Description: "`xdg-open` is Linux-only, `open` is macOS-only. " +
			"Use `$BROWSER` or check `$OSTYPE` for cross-platform URL/file opening.",
		Check: checkZC1181,
	})
}

func checkZC1181(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	if ident.Value != "xdg-open" && ident.Value != "open" {
		return nil
	}

	if ident.Value == "open" {
		// open with flags like -a, -e is macOS-specific file opening
		for _, arg := range cmd.Arguments {
			val := arg.String()
			if len(val) > 0 && val[0] == '-' {
				return nil // Likely intentional macOS usage
			}
		}
	}

	return []Violation{{
		KataID: "ZC1181",
		Message: "Use `$BROWSER` or check `$OSTYPE` instead of `" + ident.Value + "` for portable " +
			"URL/file opening across Linux and macOS.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityInfo,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1182",
		Title:    "Avoid `nc`/`netcat` for HTTP — use `curl` or `zsh/net/tcp`",
		Severity: SeverityWarning,
		Description: "`nc`/`netcat` for HTTP requests is fragile and lacks TLS support. " +
			"Use `curl` or Zsh `zsh/net/tcp` module for reliable network operations.",
		Check: checkZC1182,
	})
}

func checkZC1182(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	if ident.Value != "nc" && ident.Value != "netcat" && ident.Value != "ncat" {
		return nil
	}

	return []Violation{{
		KataID: "ZC1182",
		Message: "Avoid `" + ident.Value + "` for network operations in scripts. Use `curl` for HTTP " +
			"or `zmodload zsh/net/tcp` for raw TCP connections with TLS support.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1183",
		Title:    "Use Zsh glob qualifiers instead of `ls -t` for file ordering",
		Severity: SeverityStyle,
		Description: "Zsh glob qualifiers like `*(om[1])` (newest) or `*(Om[1])` (oldest) " +
			"order files without spawning `ls`. Avoid `ls -t | head` patterns.",
		Check: checkZC1183,
	})
}

func checkZC1183(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "ls" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-t" || val == "-tr" || val == "-lt" || val == "-ltr" {
			return []Violation{{
				KataID: "ZC1183",
				Message: "Use Zsh glob qualifiers `*(om[1])` for newest file or `*(Om[1])` for oldest " +
					"instead of `ls -t`. Glob qualifiers avoid spawning external processes.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityStyle,
			}}
		}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1184",
		Title:    "Avoid `diff -u` for patch generation — use `git diff` when in a repo",
		Severity: SeverityStyle,
		Description: "When working within a git repository, `git diff` provides better context, " +
			"color output, and integration. Use `diff -u` only for non-repo file comparisons.",
		Check: checkZC1184,
	})
}

func checkZC1184(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "diff" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-u" || val == "--unified" {
			return []Violation{{
				KataID: "ZC1184",
				Message: "Consider `git diff` instead of `diff -u` when working in a repository. " +
					"`git diff` provides better context and integration.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityStyle,
			}}
		}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1185",
		Title:    "Use Zsh `${#${(z)var}}` instead of `wc -w` for word count",
		Severity: SeverityStyle,
		Description: "Zsh `${(z)var}` splits a string into words and `${#...}` counts them. " +
			"Avoid piping through `wc -w` for simple word counting from variables.",
		Check: checkZC1185,
	})
}

func checkZC1185(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "wc" {
		return nil
	}

	hasWordFlag := false
	hasFile := false

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-w" {
			hasWordFlag = true
		} else if len(val) > 0 && val[0] != '-' {
			hasFile = true
		}
	}

	if hasWordFlag && !hasFile {
		return []Violation{{
			KataID: "ZC1185",
			Message: "Use Zsh `${#${(z)var}}` for word counting instead of piping through `wc -w`. " +
				"Parameter expansion avoids spawning an external process.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityStyle,
		}}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1186",
		Title:    "Use `unset -v` or `unset -f` for explicit unsetting",
		Severity: SeverityInfo,
		Description: "Bare `unset name` is ambiguous — it unsets variables first, then functions. " +
			"Use `unset -v` for variables or `unset -f` for functions to be explicit.",
		Check: checkZC1186,
	})
}

func checkZC1186(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "unset" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-v" || val == "-f" {
			return nil
		}
	}

	if len(cmd.Arguments) == 0 {
		return nil
	}

	return []Violation{{
		KataID: "ZC1186",
		Message: "Use `unset -v name` for variables or `unset -f name` for functions. " +
			"Bare `unset` is ambiguous about what is being removed.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityInfo,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1187",
		Title:    "Avoid `notify-send` without fallback — check availability first",
		Severity: SeverityInfo,
		Description: "`notify-send` is Linux-only (libnotify). For portable notifications, " +
			"check `$OSTYPE` and fall back to `osascript` on macOS or `print` as default.",
		Check: checkZC1187,
	})
}

func checkZC1187(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "notify-send" {
		return nil
	}

	return []Violation{{
		KataID: "ZC1187",
		Message: "Wrap `notify-send` with an `$OSTYPE` check or `command -v` guard. " +
			"It is Linux-only and will fail silently on macOS.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityInfo,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1188",
		Title:    "Use Zsh `path+=()` instead of `export PATH=$PATH:dir`",
		Severity: SeverityStyle,
		Description: "Zsh ties the `path` array to `$PATH`. Use `path+=(dir)` to append " +
			"directories cleanly instead of string manipulation with `export PATH=`.",
		Check: checkZC1188,
	})
}

func checkZC1188(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "export" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if len(val) > 5 && val[:5] == "PATH=" {
			return []Violation{{
				KataID: "ZC1188",
				Message: "Use `path+=(dir)` instead of `export PATH=$PATH:dir`. " +
					"Zsh ties the `path` array to `$PATH` for cleaner manipulation.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityStyle,
			}}
		}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1189",
		Title:    "Avoid `source /dev/stdin` — use direct evaluation",
		Severity: SeverityWarning,
		Description: "`source /dev/stdin` is fragile and platform-dependent. " +
			"Use `eval \"$(cmd)\"` or direct command execution instead.",
		Check: checkZC1189,
	})
}

func checkZC1189(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	if ident.Value != "source" && ident.Value != "." {
		return nil
	}

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "/dev/stdin" || val == "/proc/self/fd/0" {
			return []Violation{{
				KataID: "ZC1189",
				Message: "Avoid `source /dev/stdin`. Use `eval \"$(cmd)\"` for direct evaluation. " +
					"`/dev/stdin` sourcing is fragile across platforms.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}

	return nil
}

func init() {
	RegisterKata(ast.InfixExpressionNode, Kata{
		ID:       "ZC1190",
		Title:    "Combine chained `grep -v` into single invocation",
		Severity: SeverityStyle,
		Description: "`grep -v p1 | grep -v p2` spawns two processes. " +
			"Use `grep -v -e p1 -e p2` to combine exclusions in one invocation.",
		Check: checkZC1190,
		Fix:   fixZC1190,
	})
}

// fixZC1190 collapses `grep -v p1 | grep -v p2` into a single
// `grep -v -e p1 -e p2`. Only fires when each grep has exactly one
// non-flag pattern argument and at most the lone `-v` flag — keeps the
// rewrite safe in the presence of trailing FILE / additional flags.
func fixZC1190(node ast.Node, _ Violation, source []byte) []FixEdit {
	pipe, ok := node.(*ast.InfixExpression)
	if !ok || pipe.Operator != "|" {
		return nil
	}
	left, ok := pipe.Left.(*ast.SimpleCommand)
	if !ok || !isCommandName(left, "grep") {
		return nil
	}
	right, ok := pipe.Right.(*ast.SimpleCommand)
	if !ok || !isCommandName(right, "grep") {
		return nil
	}
	leftPat, leftOk := zc1190SinglePattern(left)
	rightPat, rightOk := zc1190SinglePattern(right)
	if !leftOk || !rightOk {
		return nil
	}

	leftTok := left.TokenLiteralNode()
	leftStart := LineColToByteOffset(source, leftTok.Line, leftTok.Column)
	if leftStart < 0 {
		return nil
	}
	if len(right.Arguments) == 0 {
		return nil
	}
	lastArg := right.Arguments[len(right.Arguments)-1]
	laTok := lastArg.TokenLiteralNode()
	laOff := LineColToByteOffset(source, laTok.Line, laTok.Column)
	if laOff < 0 {
		return nil
	}
	laLit := lastArg.String()
	end := laOff + len(laLit)
	if end > len(source) {
		return nil
	}
	startLine, startCol := offsetLineColZC1190(source, leftStart)
	if startLine < 0 {
		return nil
	}
	return []FixEdit{{
		Line:    startLine,
		Column:  startCol,
		Length:  end - leftStart,
		Replace: "grep -v -e " + leftPat + " -e " + rightPat,
	}}
}

// zc1190SinglePattern returns the lone non-flag argument of a
// `grep -v PAT` invocation. Returns ok=false when args contain
// extras, multiple flags, or no pattern at all.
func zc1190SinglePattern(cmd *ast.SimpleCommand) (string, bool) {
	pattern := ""
	hasV := false
	for _, arg := range cmd.Arguments {
		v := arg.String()
		switch {
		case v == "-v":
			hasV = true
		case len(v) > 0 && v[0] == '-':
			return "", false
		default:
			if pattern != "" {
				return "", false
			}
			pattern = v
		}
	}
	if !hasV || pattern == "" {
		return "", false
	}
	return pattern, true
}

func offsetLineColZC1190(source []byte, offset int) (int, int) {
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

func checkZC1190(node ast.Node) []Violation {
	pipe, ok := node.(*ast.InfixExpression)
	if !ok || pipe.Operator != "|" {
		return nil
	}

	leftCmd, ok := pipe.Left.(*ast.SimpleCommand)
	if !ok || !isCommandName(leftCmd, "grep") {
		return nil
	}

	rightCmd, ok := pipe.Right.(*ast.SimpleCommand)
	if !ok || !isCommandName(rightCmd, "grep") {
		return nil
	}

	leftHasV := false
	rightHasV := false

	for _, arg := range leftCmd.Arguments {
		if arg.String() == "-v" {
			leftHasV = true
		}
	}
	for _, arg := range rightCmd.Arguments {
		if arg.String() == "-v" {
			rightHasV = true
		}
	}

	if leftHasV && rightHasV {
		return []Violation{{
			KataID: "ZC1190",
			Message: "Combine `grep -v p1 | grep -v p2` into `grep -v -e p1 -e p2`. " +
				"A single invocation avoids an unnecessary pipeline.",
			Line:   pipe.TokenLiteralNode().Line,
			Column: pipe.TokenLiteralNode().Column,
			Level:  SeverityStyle,
		}}
	}

	return nil
}

func init() {
	RegisterKata(ast.IdentifierNode, Kata{
		ID:       "ZC1191",
		Title:    "Avoid `clear` command — use ANSI escape sequences",
		Severity: SeverityStyle,
		Description: "`clear` spawns an external process for screen clearing. " +
			"Use `print -n '\\e[2J\\e[H'` for faster terminal clearing.",
		Check: checkZC1191,
		Fix:   fixZC1191,
	})
}

// fixZC1191 rewrites a bare `clear` identifier into the equivalent
// ANSI-escape `print` invocation, avoiding the external process. The
// `$'...'` quoting is required so the lexer interprets the escape
// codes; plain single quotes pass them through literally. The `-rn`
// flag-bundle matches the canonical `print -rn` form ZShellCheck
// recommends elsewhere (see ZC1017, ZC1118), so the rewrite is
// idempotent on re-run.
func fixZC1191(node ast.Node, v Violation, _ []byte) []FixEdit {
	ident, ok := node.(*ast.Identifier)
	if !ok || ident == nil || ident.Value != "clear" {
		return nil
	}
	return []FixEdit{{
		Line:    v.Line,
		Column:  v.Column,
		Length:  len("clear"),
		Replace: "print -rn $'\\e[2J\\e[H'",
	}}
}

func checkZC1191(node ast.Node) []Violation {
	ident, ok := node.(*ast.Identifier)
	if !ok || ident == nil || ident.Value != "clear" {
		return nil
	}

	return []Violation{{
		KataID: "ZC1191",
		Message: "Use `print -n '\\e[2J\\e[H'` instead of `clear`. " +
			"ANSI escape sequences avoid spawning an external process.",
		Line:   ident.Token.Line,
		Column: ident.Token.Column,
		Level:  SeverityStyle,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1192",
		Title:    "Avoid `sleep 0` — it is a no-op external process",
		Severity: SeverityInfo,
		Description: "`sleep 0` spawns an external process that does nothing. " +
			"Remove it or use `:` if an explicit no-op is needed.",
		Check: checkZC1192,
		Fix:   fixZC1192,
	})
}

// fixZC1192 rewrites the no-op `sleep 0` invocation into `:`, the
// builtin no-op. Span covers the command name through the `0` arg.
func fixZC1192(node ast.Node, v Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	if len(cmd.Arguments) != 1 {
		return nil
	}
	zeroArg := cmd.Arguments[0]
	if zeroArg.String() != "0" {
		return nil
	}
	nameOff := LineColToByteOffset(source, v.Line, v.Column)
	if nameOff < 0 || nameOff+len("sleep") > len(source) {
		return nil
	}
	if string(source[nameOff:nameOff+len("sleep")]) != "sleep" {
		return nil
	}
	argTok := zeroArg.TokenLiteralNode()
	argOff := LineColToByteOffset(source, argTok.Line, argTok.Column)
	if argOff < 0 || argOff+1 > len(source) || source[argOff] != '0' {
		return nil
	}
	end := argOff + 1
	return []FixEdit{{
		Line:    v.Line,
		Column:  v.Column,
		Length:  end - nameOff,
		Replace: ":",
	}}
}

func checkZC1192(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "sleep" {
		return nil
	}

	if len(cmd.Arguments) == 1 && cmd.Arguments[0].String() == "0" {
		return []Violation{{
			KataID: "ZC1192",
			Message: "Remove `sleep 0` — it spawns a process that does nothing. " +
				"Use `:` if an explicit no-op is needed.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityInfo,
		}}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1193",
		Title:    "Avoid `rm -i` in non-interactive scripts",
		Severity: SeverityWarning,
		Description: "`rm -i` prompts for confirmation which hangs in non-interactive scripts. " +
			"Remove the `-i` flag or use `rm -f` for scripts that run unattended.",
		Check: checkZC1193,
	})
}

func checkZC1193(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "rm" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-i" || val == "-ri" || val == "-ir" {
			return []Violation{{
				KataID: "ZC1193",
				Message: "Avoid `rm -i` in scripts — it prompts interactively and will hang " +
					"in non-interactive execution. Remove `-i` or use explicit checks instead.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1194",
		Title:    "Avoid `sed` with multiple `-e` — use a single script",
		Severity: SeverityStyle,
		Description: "Multiple `sed -e 's/a/b/' -e 's/c/d/'` can be combined into " +
			"`sed 's/a/b/; s/c/d/'` for cleaner syntax and fewer shell word splits.",
		Check: checkZC1194,
	})
}

func checkZC1194(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "sed" {
		return nil
	}

	eCount := 0
	for _, arg := range cmd.Arguments {
		if arg.String() == "-e" {
			eCount++
		}
	}

	if eCount >= 2 {
		return []Violation{{
			KataID: "ZC1194",
			Message: "Combine multiple `sed -e` expressions into a single script: " +
				"`sed 's/a/b/; s/c/d/'` is cleaner than multiple `-e` flags.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityStyle,
		}}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1195",
		Title:    "Avoid overly permissive `umask` values",
		Severity: SeverityWarning,
		Description: "`umask 000` or `umask 0000` creates world-writable files by default. " +
			"Use `umask 022` or more restrictive values for security.",
		Check: checkZC1195,
	})
}

func checkZC1195(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "umask" {
		return nil
	}

	if len(cmd.Arguments) != 1 {
		return nil
	}

	val := cmd.Arguments[0].String()
	if val == "000" || val == "0000" || val == "0" {
		return []Violation{{
			KataID: "ZC1195",
			Message: "Avoid `umask 000` — it creates world-writable files. " +
				"Use `umask 022` or `umask 077` for secure default permissions.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityWarning,
		}}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1196",
		Title:    "Avoid `cat` for reading single file into variable",
		Severity: SeverityStyle,
		Description: "Use Zsh `$(<file)` instead of `$(cat file)` to read file contents. " +
			"`$(<file)` is a Zsh builtin that avoids spawning cat.",
		Check: checkZC1196,
	})
}

func checkZC1196(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "less" {
		return nil
	}

	// less without flags in a script is likely a mistake
	hasFlags := false
	for _, arg := range cmd.Arguments {
		val := arg.String()
		if len(val) > 0 && val[0] == '-' {
			hasFlags = true
		}
	}

	if !hasFlags && len(cmd.Arguments) > 0 {
		return []Violation{{
			KataID: "ZC1196",
			Message: "Avoid `less` in scripts — it requires interactive terminal input. " +
				"Use `cat` or redirect output to a pager only when `$TERM` is available.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityStyle,
		}}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1197",
		Title:    "Avoid `more` in scripts — use `cat` or pager check",
		Severity: SeverityStyle,
		Description: "`more` requires an interactive terminal and will hang in scripts. " +
			"Use `cat` for output or check `$TERM` before invoking a pager.",
		Check: checkZC1197,
	})
}

func checkZC1197(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "more" {
		return nil
	}

	hasFlags := false
	for _, arg := range cmd.Arguments {
		val := arg.String()
		if len(val) > 0 && val[0] == '-' {
			hasFlags = true
		}
	}

	if !hasFlags && len(cmd.Arguments) > 0 {
		return []Violation{{
			KataID: "ZC1197",
			Message: "Avoid `more` in scripts — it requires an interactive terminal. " +
				"Use `cat` for output or check `[[ -t 1 ]]` before paging.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityStyle,
		}}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1198",
		Title:    "Avoid interactive editors in scripts",
		Severity: SeverityWarning,
		Description: "`nano`, `vi`, and `vim` require interactive terminals and will hang " +
			"in non-interactive scripts. Use `sed -i` or `ed` for scripted editing.",
		Check: checkZC1198,
	})
}

func checkZC1198(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	name := ident.Value
	if name != "nano" && name != "vi" && name != "vim" && name != "emacs" {
		return nil
	}

	return []Violation{{
		KataID: "ZC1198",
		Message: "Avoid `" + name + "` in scripts — interactive editors hang without a terminal. " +
			"Use `sed -i` or `ed` for scripted file editing.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1199",
		Title:    "Avoid `telnet` in scripts — use `curl` or `zsh/net/tcp`",
		Severity: SeverityWarning,
		Description: "`telnet` is interactive and sends data in plain text. " +
			"Use `curl` for HTTP or `zmodload zsh/net/tcp` for port checks in scripts.",
		Check: checkZC1199,
	})
}

func checkZC1199(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "telnet" {
		return nil
	}

	return []Violation{{
		KataID: "ZC1199",
		Message: "Avoid `telnet` in scripts — it is interactive and insecure. " +
			"Use `curl` for HTTP checks or `zmodload zsh/net/tcp` for port testing.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}
