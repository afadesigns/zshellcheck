// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"regexp"
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.IdentifierNode, Kata{
		ID:       "ZC1300",
		Title:    "Avoid `$BASH_VERSINFO` — use `$ZSH_VERSION` in Zsh",
		Severity: SeverityWarning,
		Description: "`$BASH_VERSINFO` is a Bash-specific array containing version components. " +
			"In Zsh, use `$ZSH_VERSION` (string) or `${(s:.:)ZSH_VERSION}` to split " +
			"it into components for version comparison.",
		Check: checkZC1300,
		Fix:   fixZC1300,
	})
}

// fixZC1300 renames `$BASH_VERSION` / `$BASH_VERSINFO` to the Zsh
// equivalent `$ZSH_VERSION`. The lossy case (BASH_VERSINFO is an
// array, ZSH_VERSION is a string) is the best single-token swap
// available; callers that need components can split the string with
// the `${(s:.:)ZSH_VERSION}` flag.
func fixZC1300(node ast.Node, v Violation, source []byte) []FixEdit {
	ident, ok := node.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident == nil {
		return nil
	}
	switch ident.Value {
	case "$BASH_VERSION", "$BASH_VERSINFO":
		return []FixEdit{{
			Line:    v.Line,
			Column:  v.Column,
			Length:  len(ident.Value),
			Replace: "$ZSH_VERSION",
		}}
	case "BASH_VERSION", "BASH_VERSINFO":
		return []FixEdit{{
			Line:    v.Line,
			Column:  v.Column,
			Length:  len(ident.Value),
			Replace: "ZSH_VERSION",
		}}
	}
	return nil
}

func checkZC1300(node ast.Node) []Violation {
	ident, ok := node.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident == nil {
		return nil
	}

	if ident.Value != "$BASH_VERSINFO" && ident.Value != "BASH_VERSINFO" &&
		ident.Value != "$BASH_VERSION" && ident.Value != "BASH_VERSION" {
		return nil
	}

	return []Violation{{
		KataID:  "ZC1300",
		Message: "Avoid Bash version variables in Zsh — use `$ZSH_VERSION` instead. Bash version variables are undefined in Zsh.",
		Line:    ident.Token.Line,
		Column:  ident.Token.Column,
		Level:   SeverityWarning,
	}}
}

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
	if ident == nil {
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
	if ident == nil {
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

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1302",
		Title:    "Avoid `help` builtin — use `run-help` or `man` in Zsh",
		Severity: SeverityInfo,
		Description: "The `help` command is a Bash builtin that displays builtin help. " +
			"Zsh does not have a `help` builtin. Use `run-help <command>` or " +
			"`man zshbuiltins` for Zsh builtin documentation.",
		Check: checkZC1302,
	})
}

func checkZC1302(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "help" {
		return nil
	}

	return []Violation{{
		KataID:  "ZC1302",
		Message: "Avoid `help` in Zsh — it is a Bash builtin. Use `run-help` or `man zshbuiltins` instead.",
		Line:    cmd.Token.Line,
		Column:  cmd.Token.Column,
		Level:   SeverityInfo,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1303",
		Title:    "Avoid `enable` command — use `zmodload` for Zsh modules",
		Severity: SeverityWarning,
		Description: "The `enable` command is a Bash builtin for enabling/disabling builtins. " +
			"Zsh uses `zmodload` to load and manage modules, and `disable`/`enable` " +
			"have different semantics. Use `zmodload` for module management.",
		Check: checkZC1303,
	})
}

func checkZC1303(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "enable" {
		return nil
	}

	// enable with -f flag loads a builtin from a shared object (Bash-specific)
	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-f" {
			return []Violation{{
				KataID:  "ZC1303",
				Message: "Avoid `enable -f` in Zsh — use `zmodload` to load modules. `enable -f` is Bash-specific.",
				Line:    cmd.Token.Line,
				Column:  cmd.Token.Column,
				Level:   SeverityWarning,
			}}
		}
	}

	return nil
}

func init() {
	RegisterKata(ast.IdentifierNode, Kata{
		ID:       "ZC1304",
		Title:    "Avoid `$BASH_SUBSHELL` — use `$ZSH_SUBSHELL` in Zsh",
		Severity: SeverityWarning,
		Description: "`$BASH_SUBSHELL` tracks subshell nesting depth in Bash. " +
			"Zsh provides `$ZSH_SUBSHELL` as the native equivalent.",
		Check: checkZC1304,
		Fix:   fixZC1304,
	})
}

// fixZC1304 renames the Bash `$BASH_SUBSHELL` identifier to the Zsh
// `$ZSH_SUBSHELL` equivalent. Handles both the dollar-prefixed and
// bare forms.
func fixZC1304(node ast.Node, v Violation, source []byte) []FixEdit {
	ident, ok := node.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident == nil {
		return nil
	}
	switch ident.Value {
	case "$BASH_SUBSHELL":
		return []FixEdit{{
			Line:    v.Line,
			Column:  v.Column,
			Length:  len("$BASH_SUBSHELL"),
			Replace: "$ZSH_SUBSHELL",
		}}
	case "BASH_SUBSHELL":
		return []FixEdit{{
			Line:    v.Line,
			Column:  v.Column,
			Length:  len("BASH_SUBSHELL"),
			Replace: "ZSH_SUBSHELL",
		}}
	}
	return nil
}

func checkZC1304(node ast.Node) []Violation {
	ident, ok := node.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident == nil {
		return nil
	}

	if ident.Value != "$BASH_SUBSHELL" && ident.Value != "BASH_SUBSHELL" {
		return nil
	}

	return []Violation{{
		KataID:  "ZC1304",
		Message: "Avoid `$BASH_SUBSHELL` in Zsh — use `$ZSH_SUBSHELL` instead. `BASH_SUBSHELL` is Bash-specific.",
		Line:    ident.Token.Line,
		Column:  ident.Token.Column,
		Level:   SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.IdentifierNode, Kata{
		ID:       "ZC1305",
		Title:    "Avoid `$COMP_WORDS` — use `$words` in Zsh completion",
		Severity: SeverityWarning,
		Description: "`$COMP_WORDS` is a Bash completion variable containing the words on " +
			"the command line. Zsh completion uses `$words` array for the same purpose.",
		Check: checkZC1305,
		Fix:   fixZC1305,
	})
}

// fixZC1305 renames the Bash `$COMP_WORDS` identifier to the Zsh
// `$words` completion array. Handles both dollar-prefixed and bare
// forms.
func fixZC1305(node ast.Node, v Violation, source []byte) []FixEdit {
	ident, ok := node.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident == nil {
		return nil
	}
	switch ident.Value {
	case "$COMP_WORDS":
		return []FixEdit{{
			Line:    v.Line,
			Column:  v.Column,
			Length:  len("$COMP_WORDS"),
			Replace: "$words",
		}}
	case "COMP_WORDS":
		return []FixEdit{{
			Line:    v.Line,
			Column:  v.Column,
			Length:  len("COMP_WORDS"),
			Replace: "words",
		}}
	}
	return nil
}

func checkZC1305(node ast.Node) []Violation {
	ident, ok := node.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident == nil {
		return nil
	}

	if ident.Value != "$COMP_WORDS" && ident.Value != "COMP_WORDS" {
		return nil
	}

	return []Violation{{
		KataID:  "ZC1305",
		Message: "Avoid `$COMP_WORDS` in Zsh — use `$words` array instead. `COMP_WORDS` is Bash completion-specific.",
		Line:    ident.Token.Line,
		Column:  ident.Token.Column,
		Level:   SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.IdentifierNode, Kata{
		ID:       "ZC1306",
		Title:    "Avoid `$COMP_CWORD` — use `$CURRENT` in Zsh completion",
		Severity: SeverityWarning,
		Description: "`$COMP_CWORD` is a Bash completion variable for the current cursor " +
			"word index. Zsh completion uses `$CURRENT` for the same purpose.",
		Check: checkZC1306,
		Fix:   fixZC1306,
	})
}

// fixZC1306 renames the Bash `$COMP_CWORD` identifier to the Zsh
// `$CURRENT` completion variable.
func fixZC1306(node ast.Node, v Violation, source []byte) []FixEdit {
	ident, ok := node.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident == nil {
		return nil
	}
	switch ident.Value {
	case "$COMP_CWORD":
		return []FixEdit{{
			Line:    v.Line,
			Column:  v.Column,
			Length:  len("$COMP_CWORD"),
			Replace: "$CURRENT",
		}}
	case "COMP_CWORD":
		return []FixEdit{{
			Line:    v.Line,
			Column:  v.Column,
			Length:  len("COMP_CWORD"),
			Replace: "CURRENT",
		}}
	}
	return nil
}

func checkZC1306(node ast.Node) []Violation {
	ident, ok := node.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident == nil {
		return nil
	}

	if ident.Value != "$COMP_CWORD" && ident.Value != "COMP_CWORD" {
		return nil
	}

	return []Violation{{
		KataID:  "ZC1306",
		Message: "Avoid `$COMP_CWORD` in Zsh — use `$CURRENT` instead. `COMP_CWORD` is Bash completion-specific.",
		Line:    ident.Token.Line,
		Column:  ident.Token.Column,
		Level:   SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.IdentifierNode, Kata{
		ID:       "ZC1307",
		Title:    "Avoid `$DIRSTACK` — use `$dirstack` (lowercase) in Zsh",
		Severity: SeverityWarning,
		Description: "`$DIRSTACK` is the Bash form of the directory stack array. " +
			"Zsh uses `$dirstack` (lowercase) for the same purpose.",
		Check: checkZC1307,
		Fix:   fixZC1307,
	})
}

// fixZC1307 renames the Bash `$DIRSTACK` / `DIRSTACK` identifier to
// the Zsh lowercase `$dirstack` / `dirstack` form. Mirrors ZC1301's
// rename pattern.
func fixZC1307(node ast.Node, v Violation, source []byte) []FixEdit {
	ident, ok := node.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident == nil {
		return nil
	}
	switch ident.Value {
	case "$DIRSTACK":
		return []FixEdit{{
			Line:    v.Line,
			Column:  v.Column,
			Length:  len("$DIRSTACK"),
			Replace: "$dirstack",
		}}
	case "DIRSTACK":
		return []FixEdit{{
			Line:    v.Line,
			Column:  v.Column,
			Length:  len("DIRSTACK"),
			Replace: "dirstack",
		}}
	}
	return nil
}

func checkZC1307(node ast.Node) []Violation {
	ident, ok := node.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident == nil {
		return nil
	}

	if ident.Value != "$DIRSTACK" && ident.Value != "DIRSTACK" {
		return nil
	}

	return []Violation{{
		KataID:  "ZC1307",
		Message: "Avoid `$DIRSTACK` in Zsh — use `$dirstack` (lowercase) instead. The uppercase form is Bash-specific.",
		Line:    ident.Token.Line,
		Column:  ident.Token.Column,
		Level:   SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.IdentifierNode, Kata{
		ID:       "ZC1308",
		Title:    "Avoid `$COMP_LINE` — use `$BUFFER` in Zsh completion",
		Severity: SeverityWarning,
		Description: "`$COMP_LINE` is a Bash completion variable containing the full command " +
			"line. Zsh completion uses `$BUFFER` for the current command line content.",
		Check: checkZC1308,
		Fix:   fixZC1308,
	})
}

// fixZC1308 renames the Bash `$COMP_LINE` identifier to the Zsh
// `$BUFFER` completion variable.
func fixZC1308(node ast.Node, v Violation, source []byte) []FixEdit {
	ident, ok := node.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident == nil {
		return nil
	}
	switch ident.Value {
	case "$COMP_LINE":
		return []FixEdit{{
			Line:    v.Line,
			Column:  v.Column,
			Length:  len("$COMP_LINE"),
			Replace: "$BUFFER",
		}}
	case "COMP_LINE":
		return []FixEdit{{
			Line:    v.Line,
			Column:  v.Column,
			Length:  len("COMP_LINE"),
			Replace: "BUFFER",
		}}
	}
	return nil
}

func checkZC1308(node ast.Node) []Violation {
	ident, ok := node.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident == nil {
		return nil
	}

	if ident.Value != "$COMP_LINE" && ident.Value != "COMP_LINE" {
		return nil
	}

	return []Violation{{
		KataID:  "ZC1308",
		Message: "Avoid `$COMP_LINE` in Zsh — use `$BUFFER` instead. `COMP_LINE` is Bash completion-specific.",
		Line:    ident.Token.Line,
		Column:  ident.Token.Column,
		Level:   SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.IdentifierNode, Kata{
		ID:       "ZC1309",
		Title:    "Avoid `$BASH_COMMAND` — not available in Zsh",
		Severity: SeverityWarning,
		Description: "`$BASH_COMMAND` contains the currently executing command in Bash. " +
			"Zsh does not provide a direct equivalent. Use `$ZSH_DEBUG_CMD` in " +
			"debug traps or restructure the logic.",
		Check: checkZC1309,
	})
}

func checkZC1309(node ast.Node) []Violation {
	ident, ok := node.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident == nil {
		return nil
	}

	if ident.Value != "$BASH_COMMAND" && ident.Value != "BASH_COMMAND" {
		return nil
	}

	return []Violation{{
		KataID:  "ZC1309",
		Message: "Avoid `$BASH_COMMAND` in Zsh — it is undefined. Use `$ZSH_DEBUG_CMD` in debug traps if needed.",
		Line:    ident.Token.Line,
		Column:  ident.Token.Column,
		Level:   SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.IdentifierNode, Kata{
		ID:       "ZC1310",
		Title:    "Avoid `$BASH_EXECUTION_STRING` — not available in Zsh",
		Severity: SeverityWarning,
		Description: "`$BASH_EXECUTION_STRING` contains the argument to `bash -c`. " +
			"Zsh does not provide this variable. Access the script argument directly.",
		Check: checkZC1310,
	})
}

func checkZC1310(node ast.Node) []Violation {
	ident, ok := node.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident == nil {
		return nil
	}

	if ident.Value != "$BASH_EXECUTION_STRING" && ident.Value != "BASH_EXECUTION_STRING" {
		return nil
	}

	return []Violation{{
		KataID:  "ZC1310",
		Message: "Avoid `$BASH_EXECUTION_STRING` in Zsh — it is undefined. Access command arguments directly instead.",
		Line:    ident.Token.Line,
		Column:  ident.Token.Column,
		Level:   SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1311",
		Title:    "Avoid `complete` command — use `compdef` in Zsh",
		Severity: SeverityWarning,
		Description: "`complete` is a Bash builtin for registering tab completions. " +
			"Zsh uses `compdef` for completion registration and the `compctl` " +
			"legacy interface. Use `compdef` for the modern Zsh completion system.",
		Check: checkZC1311,
	})
}

func checkZC1311(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "complete" {
		return nil
	}

	return []Violation{{
		KataID:  "ZC1311",
		Message: "Avoid `complete` in Zsh — it is a Bash builtin. Use `compdef` for Zsh completion registration.",
		Line:    cmd.Token.Line,
		Column:  cmd.Token.Column,
		Level:   SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1312",
		Title:    "Avoid `compgen` command — use `compadd` in Zsh",
		Severity: SeverityWarning,
		Description: "`compgen` is a Bash builtin for generating completions. " +
			"Zsh uses `compadd` and the completion system functions for adding " +
			"completion candidates.",
		Check: checkZC1312,
	})
}

func checkZC1312(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "compgen" {
		return nil
	}

	return []Violation{{
		KataID:  "ZC1312",
		Message: "Avoid `compgen` in Zsh — it is a Bash builtin. Use `compadd` or Zsh completion functions instead.",
		Line:    cmd.Token.Line,
		Column:  cmd.Token.Column,
		Level:   SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.IdentifierNode, Kata{
		ID:       "ZC1313",
		Title:    "Avoid `$BASH_ALIASES` — use Zsh `aliases` hash",
		Severity: SeverityWarning,
		Description: "`$BASH_ALIASES` is a Bash associative array of defined aliases. " +
			"Zsh provides the `aliases` associative array for the same purpose.",
		Check: checkZC1313,
		Fix:   fixZC1313,
	})
}

// fixZC1313 renames the Bash `$BASH_ALIASES` identifier to the Zsh
// `$aliases` associative array.
func fixZC1313(node ast.Node, v Violation, source []byte) []FixEdit {
	ident, ok := node.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident == nil {
		return nil
	}
	switch ident.Value {
	case "$BASH_ALIASES":
		return []FixEdit{{
			Line:    v.Line,
			Column:  v.Column,
			Length:  len("$BASH_ALIASES"),
			Replace: "$aliases",
		}}
	case "BASH_ALIASES":
		return []FixEdit{{
			Line:    v.Line,
			Column:  v.Column,
			Length:  len("BASH_ALIASES"),
			Replace: "aliases",
		}}
	}
	return nil
}

func checkZC1313(node ast.Node) []Violation {
	ident, ok := node.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident == nil {
		return nil
	}

	if ident.Value != "$BASH_ALIASES" && ident.Value != "BASH_ALIASES" {
		return nil
	}

	return []Violation{{
		KataID:  "ZC1313",
		Message: "Avoid `$BASH_ALIASES` in Zsh — use the `aliases` associative array instead.",
		Line:    ident.Token.Line,
		Column:  ident.Token.Column,
		Level:   SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.IdentifierNode, Kata{
		ID:       "ZC1314",
		Title:    "Avoid `$BASH_LOADABLES_PATH` — not available in Zsh",
		Severity: SeverityWarning,
		Description: "`$BASH_LOADABLES_PATH` is a Bash variable for loadable builtin search paths. " +
			"Zsh has no equivalent; use `zmodload` with full module names instead.",
		Check: checkZC1314,
	})
}

func checkZC1314(node ast.Node) []Violation {
	ident, ok := node.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident == nil {
		return nil
	}

	if ident.Value != "$BASH_LOADABLES_PATH" && ident.Value != "BASH_LOADABLES_PATH" {
		return nil
	}

	return []Violation{{
		KataID:  "ZC1314",
		Message: "Avoid `$BASH_LOADABLES_PATH` in Zsh — it is undefined. Use `zmodload` with full module names.",
		Line:    ident.Token.Line,
		Column:  ident.Token.Column,
		Level:   SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.IdentifierNode, Kata{
		ID:       "ZC1315",
		Title:    "Avoid `$BASH_COMPAT` — use `emulate` for compatibility in Zsh",
		Severity: SeverityWarning,
		Description: "`$BASH_COMPAT` sets Bash compatibility level. Zsh uses `emulate` " +
			"to control compatibility mode (e.g., `emulate -L sh` for POSIX mode).",
		Check: checkZC1315,
	})
}

func checkZC1315(node ast.Node) []Violation {
	ident, ok := node.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident == nil {
		return nil
	}

	if ident.Value != "$BASH_COMPAT" && ident.Value != "BASH_COMPAT" {
		return nil
	}

	return []Violation{{
		KataID:  "ZC1315",
		Message: "Avoid `$BASH_COMPAT` in Zsh — use `emulate` for shell compatibility mode instead.",
		Line:    ident.Token.Line,
		Column:  ident.Token.Column,
		Level:   SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1316",
		Title:    "Avoid `caller` builtin — use `$funcfiletrace` in Zsh",
		Severity: SeverityWarning,
		Description: "`caller` is a Bash builtin that returns the call stack context. " +
			"Zsh provides `$funcfiletrace`, `$funcstack`, and `$funcsourcetrace` " +
			"for inspecting the call stack.",
		Check: checkZC1316,
	})
}

func checkZC1316(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "caller" {
		return nil
	}

	return []Violation{{
		KataID:  "ZC1316",
		Message: "Avoid `caller` in Zsh — it is a Bash builtin. Use `$funcfiletrace` and `$funcstack` instead.",
		Line:    cmd.Token.Line,
		Column:  cmd.Token.Column,
		Level:   SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.IdentifierNode, Kata{
		ID:       "ZC1317",
		Title:    "Avoid `$BASH_ENV` — use `$ZDOTDIR` and `$ENV` in Zsh",
		Severity: SeverityInfo,
		Description: "`$BASH_ENV` specifies a startup file for non-interactive Bash shells. " +
			"Zsh uses `$ZDOTDIR` to locate `.zshrc` and related files, and `$ENV` " +
			"for POSIX-compatible startup.",
		Check: checkZC1317,
	})
}

func checkZC1317(node ast.Node) []Violation {
	ident, ok := node.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident == nil {
		return nil
	}

	if ident.Value != "$BASH_ENV" && ident.Value != "BASH_ENV" {
		return nil
	}

	return []Violation{{
		KataID:  "ZC1317",
		Message: "Avoid `$BASH_ENV` in Zsh — use `$ZDOTDIR` for Zsh startup file locations instead.",
		Line:    ident.Token.Line,
		Column:  ident.Token.Column,
		Level:   SeverityInfo,
	}}
}

func init() {
	RegisterKata(ast.IdentifierNode, Kata{
		ID:       "ZC1318",
		Title:    "Avoid `$BASH_CMDS` — use `$commands` hash in Zsh",
		Severity: SeverityWarning,
		Description: "`$BASH_CMDS` is a Bash associative array caching command lookups. " +
			"Zsh provides the `$commands` hash for the same purpose, mapping " +
			"command names to their full paths.",
		Check: checkZC1318,
		Fix:   fixZC1318,
	})
}

// fixZC1318 renames the Bash `$BASH_CMDS` identifier to the Zsh
// `$commands` hash.
func fixZC1318(node ast.Node, v Violation, source []byte) []FixEdit {
	ident, ok := node.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident == nil {
		return nil
	}
	switch ident.Value {
	case "$BASH_CMDS":
		return []FixEdit{{
			Line:    v.Line,
			Column:  v.Column,
			Length:  len("$BASH_CMDS"),
			Replace: "$commands",
		}}
	case "BASH_CMDS":
		return []FixEdit{{
			Line:    v.Line,
			Column:  v.Column,
			Length:  len("BASH_CMDS"),
			Replace: "commands",
		}}
	}
	return nil
}

func checkZC1318(node ast.Node) []Violation {
	ident, ok := node.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident == nil {
		return nil
	}

	if ident.Value != "$BASH_CMDS" && ident.Value != "BASH_CMDS" {
		return nil
	}

	return []Violation{{
		KataID:  "ZC1318",
		Message: "Avoid `$BASH_CMDS` in Zsh — use the `$commands` hash for command path lookups instead.",
		Line:    ident.Token.Line,
		Column:  ident.Token.Column,
		Level:   SeverityWarning,
	}}
}

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

func init() {
	RegisterKata(ast.IdentifierNode, Kata{
		ID:       "ZC1320",
		Title:    "Avoid `$BASH_ARGV` — use `$argv` in Zsh",
		Severity: SeverityWarning,
		Description: "`$BASH_ARGV` is a Bash array containing arguments in reverse order. " +
			"Zsh provides `$argv` (or `$@`) for positional parameters.",
		Check: checkZC1320,
		Fix:   fixZC1320,
	})
}

// fixZC1320 rewrites the Bash `$BASH_ARGV` / `BASH_ARGV` identifier to
// the Zsh `$argv` form. Caveat: `$BASH_ARGV` lists args in reverse
// stacking order in Bash; `$argv` is the current-frame positional
// array. Most usages target the current frame and the rewrite is
// correct; deeper stack walks need a hand-port.
func fixZC1320(node ast.Node, v Violation, source []byte) []FixEdit {
	ident, ok := node.(*ast.Identifier)
	if !ok || ident == nil {
		return nil
	}
	switch ident.Value {
	case "$BASH_ARGV":
		return []FixEdit{{
			Line:    v.Line,
			Column:  v.Column,
			Length:  len("$BASH_ARGV"),
			Replace: "$argv",
		}}
	case "BASH_ARGV":
		return []FixEdit{{
			Line:    v.Line,
			Column:  v.Column,
			Length:  len("BASH_ARGV"),
			Replace: "argv",
		}}
	}
	return nil
}

func checkZC1320(node ast.Node) []Violation {
	ident, ok := node.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident == nil {
		return nil
	}

	if ident.Value != "$BASH_ARGV" && ident.Value != "BASH_ARGV" {
		return nil
	}

	return []Violation{{
		KataID:  "ZC1320",
		Message: "Avoid `$BASH_ARGV` in Zsh — use `$argv` or `$@` for positional parameters. `BASH_ARGV` is Bash-specific.",
		Line:    ident.Token.Line,
		Column:  ident.Token.Column,
		Level:   SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.IdentifierNode, Kata{
		ID:       "ZC1321",
		Title:    "Avoid `$BASH_XTRACEFD` — not available in Zsh",
		Severity: SeverityWarning,
		Description: "`$BASH_XTRACEFD` redirects Bash xtrace output to a file descriptor. " +
			"Zsh does not have this variable. Use `exec 2>file` or redirect " +
			"stderr directly for trace output redirection.",
		Check: checkZC1321,
	})
}

func checkZC1321(node ast.Node) []Violation {
	ident, ok := node.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident == nil {
		return nil
	}

	if ident.Value != "$BASH_XTRACEFD" && ident.Value != "BASH_XTRACEFD" {
		return nil
	}

	return []Violation{{
		KataID:  "ZC1321",
		Message: "Avoid `$BASH_XTRACEFD` in Zsh — it is undefined. Redirect stderr directly for xtrace output.",
		Line:    ident.Token.Line,
		Column:  ident.Token.Column,
		Level:   SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.IdentifierNode, Kata{
		ID:       "ZC1322",
		Title:    "Avoid `$COPROC` — Zsh coproc uses different syntax",
		Severity: SeverityWarning,
		Description: "`$COPROC` is a Bash array for coprocess file descriptors. " +
			"Zsh coprocesses use `coproc` keyword with different variable naming " +
			"and `read -p`/`print -p` for I/O.",
		Check: checkZC1322,
	})
}

func checkZC1322(node ast.Node) []Violation {
	ident, ok := node.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident == nil {
		return nil
	}

	if ident.Value != "$COPROC" && ident.Value != "COPROC" {
		return nil
	}

	return []Violation{{
		KataID:  "ZC1322",
		Message: "Avoid `$COPROC` in Zsh — Zsh coprocesses use `read -p`/`print -p` for I/O. `COPROC` is Bash-specific.",
		Line:    ident.Token.Line,
		Column:  ident.Token.Column,
		Level:   SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1323",
		Title:    "Avoid `suspend` builtin — use `kill -STOP $$` in Zsh",
		Severity: SeverityWarning,
		Description: "`suspend` is a Bash builtin that suspends the shell. Zsh does not have " +
			"a `suspend` builtin. Use `kill -STOP $$` or Ctrl-Z for the same effect.",
		Check: checkZC1323,
	})
}

func checkZC1323(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "suspend" {
		return nil
	}

	return []Violation{{
		KataID:  "ZC1323",
		Message: "Avoid `suspend` in Zsh — it is a Bash builtin. Use `kill -STOP $$` if needed.",
		Line:    cmd.Token.Line,
		Column:  cmd.Token.Column,
		Level:   SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.IdentifierNode, Kata{
		ID:       "ZC1324",
		Title:    "Avoid `$PROMPT_COMMAND` — use `precmd` hook in Zsh",
		Severity: SeverityWarning,
		Description: "`$PROMPT_COMMAND` is a Bash variable that executes a command before " +
			"each prompt. Zsh uses the `precmd` hook function for the same purpose.",
		Check: checkZC1324,
	})
}

func checkZC1324(node ast.Node) []Violation {
	ident, ok := node.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident == nil {
		return nil
	}

	if ident.Value != "$PROMPT_COMMAND" && ident.Value != "PROMPT_COMMAND" {
		return nil
	}

	return []Violation{{
		KataID:  "ZC1324",
		Message: "Avoid `$PROMPT_COMMAND` in Zsh — use the `precmd` hook function instead. `PROMPT_COMMAND` is Bash-specific.",
		Line:    ident.Token.Line,
		Column:  ident.Token.Column,
		Level:   SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.IdentifierNode, Kata{
		ID:       "ZC1325",
		Title:    "Avoid `$PS0` — use `preexec` hook in Zsh",
		Severity: SeverityWarning,
		Description: "`$PS0` is a Bash 4.4+ prompt string displayed before command execution. " +
			"Zsh uses the `preexec` hook function for running code before each command.",
		Check: checkZC1325,
	})
}

func checkZC1325(node ast.Node) []Violation {
	ident, ok := node.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident == nil {
		return nil
	}

	if ident.Value != "$PS0" && ident.Value != "PS0" {
		return nil
	}

	return []Violation{{
		KataID:  "ZC1325",
		Message: "Avoid `$PS0` in Zsh — use the `preexec` hook function instead. `PS0` is Bash 4.4+ specific.",
		Line:    ident.Token.Line,
		Column:  ident.Token.Column,
		Level:   SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.IdentifierNode, Kata{
		ID:       "ZC1326",
		Title:    "Avoid `$HISTTIMEFORMAT` — use `fc -li` in Zsh",
		Severity: SeverityInfo,
		Description: "`$HISTTIMEFORMAT` is a Bash variable for formatting history timestamps. " +
			"Zsh stores timestamps automatically when `EXTENDED_HISTORY` is set, " +
			"and displays them with `fc -li` or `history -i`.",
		Check: checkZC1326,
	})
}

func checkZC1326(node ast.Node) []Violation {
	ident, ok := node.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident == nil {
		return nil
	}

	if ident.Value != "$HISTTIMEFORMAT" && ident.Value != "HISTTIMEFORMAT" {
		return nil
	}

	return []Violation{{
		KataID:  "ZC1326",
		Message: "Avoid `$HISTTIMEFORMAT` in Zsh — use `setopt EXTENDED_HISTORY` and `fc -li` instead.",
		Line:    ident.Token.Line,
		Column:  ident.Token.Column,
		Level:   SeverityInfo,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1327",
		Title:    "Avoid `history -c` — Zsh uses different history management",
		Severity: SeverityWarning,
		Description: "`history -c` clears history in Bash. Zsh provides `fc -p` for pushing " +
			"history to a new file and `fc -P` for popping. Use `fc -W` to write and " +
			"`fc -R` to read history files.",
		Check: checkZC1327,
	})
}

func checkZC1327(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "history" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		val := arg.String()
		// `-c` / `-d` are the destructive / anti-forensics flags and are
		// owned by ZC1487; this kata narrows to the Bash-only write/read
		// portability flags (`-w` / `-r` / `-a`).
		if val == "-w" || val == "-r" || val == "-a" {
			return []Violation{{
				KataID:  "ZC1327",
				Message: "Avoid `history " + val + "` in Zsh — Bash history flags differ. Use `fc` commands for Zsh history management.",
				Line:    cmd.Token.Line,
				Column:  cmd.Token.Column,
				Level:   SeverityWarning,
			}}
		}
	}

	return nil
}

func init() {
	RegisterKata(ast.IdentifierNode, Kata{
		ID:       "ZC1328",
		Title:    "Avoid `$HISTCONTROL` — use Zsh `setopt` history options",
		Severity: SeverityInfo,
		Description: "`$HISTCONTROL` is a Bash variable controlling history deduplication. " +
			"Zsh uses `setopt HIST_IGNORE_DUPS`, `HIST_IGNORE_ALL_DUPS`, and " +
			"`HIST_IGNORE_SPACE` for the same functionality.",
		Check: checkZC1328,
	})
}

func checkZC1328(node ast.Node) []Violation {
	ident, ok := node.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident == nil {
		return nil
	}

	if ident.Value != "$HISTCONTROL" && ident.Value != "HISTCONTROL" {
		return nil
	}

	return []Violation{{
		KataID:  "ZC1328",
		Message: "Avoid `$HISTCONTROL` in Zsh — use `setopt HIST_IGNORE_DUPS` and related options instead.",
		Line:    ident.Token.Line,
		Column:  ident.Token.Column,
		Level:   SeverityInfo,
	}}
}

func init() {
	RegisterKata(ast.IdentifierNode, Kata{
		ID:       "ZC1329",
		Title:    "Avoid `$HISTIGNORE` — use `zshaddhistory` hook in Zsh",
		Severity: SeverityInfo,
		Description: "`$HISTIGNORE` is a Bash variable for pattern-based history filtering. " +
			"Zsh uses the `zshaddhistory` hook function and `setopt HIST_IGNORE_SPACE` " +
			"for controlling which commands enter history.",
		Check: checkZC1329,
	})
}

func checkZC1329(node ast.Node) []Violation {
	ident, ok := node.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident == nil {
		return nil
	}

	if ident.Value != "$HISTIGNORE" && ident.Value != "HISTIGNORE" {
		return nil
	}

	return []Violation{{
		KataID:  "ZC1329",
		Message: "Avoid `$HISTIGNORE` in Zsh — use `zshaddhistory` hook for history filtering instead.",
		Line:    ident.Token.Line,
		Column:  ident.Token.Column,
		Level:   SeverityInfo,
	}}
}

func init() {
	RegisterKata(ast.IdentifierNode, Kata{
		ID:       "ZC1330",
		Title:    "Avoid `$INPUTRC` — use `bindkey` in Zsh",
		Severity: SeverityInfo,
		Description: "`$INPUTRC` points to the readline configuration file in Bash. " +
			"Zsh uses `bindkey` and ZLE widgets for key binding configuration, " +
			"not readline.",
		Check: checkZC1330,
	})
}

func checkZC1330(node ast.Node) []Violation {
	ident, ok := node.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident == nil {
		return nil
	}

	if ident.Value != "$INPUTRC" && ident.Value != "INPUTRC" {
		return nil
	}

	return []Violation{{
		KataID:  "ZC1330",
		Message: "Avoid `$INPUTRC` in Zsh — Zsh uses `bindkey` and ZLE, not readline. `INPUTRC` is Bash-specific.",
		Line:    ident.Token.Line,
		Column:  ident.Token.Column,
		Level:   SeverityInfo,
	}}
}

func init() {
	RegisterKata(ast.IdentifierNode, Kata{
		ID:       "ZC1331",
		Title:    "Avoid `$BASH_REMATCH` — use `$match` array in Zsh",
		Severity: SeverityWarning,
		Description: "`$BASH_REMATCH` holds regex capture groups in Bash. Zsh stores " +
			"regex matches in the `$match` array (and `$MATCH` for the full match) " +
			"when using `=~` with `setopt BASH_REMATCH` disabled.",
		Check: checkZC1331,
		Fix:   fixZC1331,
	})
}

// fixZC1331 renames the Bash `$BASH_REMATCH` identifier to the Zsh
// `$match` regex-capture array.
func fixZC1331(node ast.Node, v Violation, source []byte) []FixEdit {
	ident, ok := node.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident == nil {
		return nil
	}
	switch ident.Value {
	case "$BASH_REMATCH":
		return []FixEdit{{
			Line:    v.Line,
			Column:  v.Column,
			Length:  len("$BASH_REMATCH"),
			Replace: "$match",
		}}
	case "BASH_REMATCH":
		return []FixEdit{{
			Line:    v.Line,
			Column:  v.Column,
			Length:  len("BASH_REMATCH"),
			Replace: "match",
		}}
	}
	return nil
}

func checkZC1331(node ast.Node) []Violation {
	ident, ok := node.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident == nil {
		return nil
	}

	if ident.Value != "$BASH_REMATCH" && ident.Value != "BASH_REMATCH" {
		return nil
	}

	return []Violation{{
		KataID:  "ZC1331",
		Message: "Avoid `$BASH_REMATCH` in Zsh — use `$match` array and `$MATCH` for regex captures instead.",
		Line:    ident.Token.Line,
		Column:  ident.Token.Column,
		Level:   SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.IdentifierNode, Kata{
		ID:       "ZC1332",
		Title:    "Avoid `$GLOBIGNORE` — use `setopt EXTENDED_GLOB` in Zsh",
		Severity: SeverityInfo,
		Description: "`$GLOBIGNORE` is a Bash variable for excluding patterns from glob expansion. " +
			"Zsh uses `setopt EXTENDED_GLOB` with the `~` (exclusion) operator or " +
			"`setopt NULL_GLOB` for different glob behavior.",
		Check: checkZC1332,
	})
}

func checkZC1332(node ast.Node) []Violation {
	ident, ok := node.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident == nil {
		return nil
	}

	if ident.Value != "$GLOBIGNORE" && ident.Value != "GLOBIGNORE" {
		return nil
	}

	return []Violation{{
		KataID:  "ZC1332",
		Message: "Avoid `$GLOBIGNORE` in Zsh — use `setopt EXTENDED_GLOB` with `~` operator for glob exclusion.",
		Line:    ident.Token.Line,
		Column:  ident.Token.Column,
		Level:   SeverityInfo,
	}}
}

func init() {
	RegisterKata(ast.IdentifierNode, Kata{
		ID:       "ZC1333",
		Title:    "Avoid `$TIMEFORMAT` — use `$TIMEFMT` in Zsh",
		Severity: SeverityInfo,
		Description: "`$TIMEFORMAT` is the Bash variable for customizing `time` output. " +
			"Zsh uses `$TIMEFMT` for the same purpose, with different format specifiers.",
		Check: checkZC1333,
		Fix:   fixZC1333,
	})
}

// fixZC1333 renames the Bash `$TIMEFORMAT` identifier to the Zsh
// `$TIMEFMT` variable. Format specifiers differ between the two
// shells; the rename preserves the identifier itself but authors
// should still review the format string after conversion.
func fixZC1333(node ast.Node, v Violation, source []byte) []FixEdit {
	ident, ok := node.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident == nil {
		return nil
	}
	switch ident.Value {
	case "$TIMEFORMAT":
		return []FixEdit{{
			Line:    v.Line,
			Column:  v.Column,
			Length:  len("$TIMEFORMAT"),
			Replace: "$TIMEFMT",
		}}
	case "TIMEFORMAT":
		return []FixEdit{{
			Line:    v.Line,
			Column:  v.Column,
			Length:  len("TIMEFORMAT"),
			Replace: "TIMEFMT",
		}}
	}
	return nil
}

func checkZC1333(node ast.Node) []Violation {
	ident, ok := node.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident == nil {
		return nil
	}

	if ident.Value != "$TIMEFORMAT" && ident.Value != "TIMEFORMAT" {
		return nil
	}

	return []Violation{{
		KataID:  "ZC1333",
		Message: "Avoid `$TIMEFORMAT` in Zsh — use `$TIMEFMT` instead. Format specifiers differ between Bash and Zsh.",
		Line:    ident.Token.Line,
		Column:  ident.Token.Column,
		Level:   SeverityInfo,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1334",
		Title:    "Avoid `type -p` — use `whence -p` in Zsh",
		Severity: SeverityWarning,
		Description: "`type -p` is a Bash flag that prints the path of a command. " +
			"Zsh `type` does not support `-p`. Use `whence -p` to get " +
			"the path of an external command in Zsh.",
		Check: checkZC1334,
		Fix:   fixZC1334,
	})
}

// fixZC1334 rewrites `type -p X` / `type -P X` to `whence -p X`. The
// span covers both the `type` command name and the `-p`/`-P` flag in a
// single edit — emitting the wider rewrite ensures it wins over the
// narrower `type` -> `command -v` swap from ZC1064 when both katas fire
// on the same input. Trailing argument(s) stay in place.
func fixZC1334(node ast.Node, v Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "type" {
		return nil
	}
	var flag ast.Expression
	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-p" || val == "-P" {
			flag = arg
			break
		}
	}
	if flag == nil {
		return nil
	}
	nameOff := LineColToByteOffset(source, v.Line, v.Column)
	if nameOff < 0 || nameOff+len("type") > len(source) {
		return nil
	}
	if string(source[nameOff:nameOff+len("type")]) != "type" {
		return nil
	}
	flagTok := flag.TokenLiteralNode()
	flagOff := LineColToByteOffset(source, flagTok.Line, flagTok.Column)
	if flagOff < 0 || flagOff+2 > len(source) {
		return nil
	}
	if string(source[flagOff:flagOff+2]) != "-p" && string(source[flagOff:flagOff+2]) != "-P" {
		return nil
	}
	return []FixEdit{{
		Line:    v.Line,
		Column:  v.Column,
		Length:  flagOff + 2 - nameOff,
		Replace: "whence -p",
	}}
}

func checkZC1334(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "type" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-p" || val == "-P" {
			return []Violation{{
				KataID:  "ZC1334",
				Message: "Avoid `type -p` in Zsh — use `whence -p` to get the command path instead.",
				Line:    cmd.Token.Line,
				Column:  cmd.Token.Column,
				Level:   SeverityWarning,
			}}
		}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1335",
		Title:    "Use Zsh array reversal instead of `tac` for in-memory data",
		Severity: SeverityStyle,
		Description: "`tac` reverses lines from a file or stdin. For in-memory array data, " +
			"Zsh provides `${(Oa)array}` to reverse array element order without " +
			"spawning an external process.",
		Check: checkZC1335,
	})
}

func checkZC1335(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "tac" {
		return nil
	}

	return []Violation{{
		KataID:  "ZC1335",
		Message: "Consider Zsh `${(Oa)array}` for reversing array data instead of piping to `tac`.",
		Line:    cmd.Token.Line,
		Column:  cmd.Token.Column,
		Level:   SeverityStyle,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1336",
		Title:    "Avoid `printenv` — use `typeset -x` or `export` in Zsh",
		Severity: SeverityStyle,
		Description: "`printenv` is an external command for listing environment variables. " +
			"Zsh provides `typeset -x` to list exported variables and `export` " +
			"to display them without spawning a subprocess.",
		Check: checkZC1336,
	})
}

func checkZC1336(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "printenv" {
		return nil
	}

	return []Violation{{
		KataID:  "ZC1336",
		Message: "Avoid `printenv` in Zsh — use `typeset -x` or `export` to list environment variables.",
		Line:    cmd.Token.Line,
		Column:  cmd.Token.Column,
		Level:   SeverityStyle,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1337",
		Title:    "Avoid `fold` command — use Zsh `print -l` with `$COLUMNS`",
		Severity: SeverityStyle,
		Description: "`fold` wraps text to a specified width. Zsh provides `$COLUMNS` for " +
			"terminal width and `print -l` for line-by-line output, reducing " +
			"dependency on external commands.",
		Check: checkZC1337,
	})
}

func checkZC1337(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "fold" {
		return nil
	}

	return []Violation{{
		KataID:  "ZC1337",
		Message: "Consider Zsh `$COLUMNS` and `print` for text wrapping instead of `fold`.",
		Line:    cmd.Token.Line,
		Column:  cmd.Token.Column,
		Level:   SeverityStyle,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1338",
		Title:    "Avoid `seq -s` — use Zsh `${(j:sep:)${(s::)...}}` for joining",
		Severity: SeverityStyle,
		Description: "`seq -s` generates a sequence with a custom separator. Zsh provides " +
			"native brace expansion with `{start..end}` and `${(j:sep:)array}` " +
			"for joining, avoiding an external process.",
		Check: checkZC1338,
	})
}

func checkZC1338(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "seq" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-s" {
			return []Violation{{
				KataID:  "ZC1338",
				Message: "Avoid `seq -s` in Zsh — use `${(j:sep:)array}` with brace expansion for joined sequences.",
				Line:    cmd.Token.Line,
				Column:  cmd.Token.Column,
				Level:   SeverityStyle,
			}}
		}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1339",
		Title:    "Use Zsh `${#${(f)var}}` instead of `wc -l` for line count",
		Severity: SeverityStyle,
		Description: "Zsh `${(f)var}` splits a string into lines and `${#...}` counts them. " +
			"Avoid piping through `wc -l` for simple line counting from variables.",
		Check: checkZC1339,
	})
}

func checkZC1339(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "wc" {
		return nil
	}

	hasLineFlag := false
	hasFile := false

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-l" {
			hasLineFlag = true
		} else if len(val) > 0 && val[0] != '-' {
			hasFile = true
		}
	}

	if hasLineFlag && !hasFile {
		return []Violation{{
			KataID: "ZC1339",
			Message: "Use Zsh `${#${(f)var}}` for line counting instead of piping through `wc -l`. " +
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
		ID:       "ZC1340",
		Title:    "Avoid `shuf` for random array element — use Zsh `$RANDOM`",
		Severity: SeverityStyle,
		Description: "Zsh provides `$RANDOM` and array subscripts to pick random elements " +
			"without spawning `shuf`. For a single random array element, use " +
			"`${array[RANDOM%$#array+1]}`.",
		Check: checkZC1340,
	})
}

func checkZC1340(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "shuf" {
		return nil
	}

	return []Violation{{
		KataID: "ZC1340",
		Message: "Avoid `shuf` for random selection — use Zsh `${array[RANDOM%$#array+1]}` " +
			"with `$RANDOM` for in-shell randomness without spawning an external.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityStyle,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1341",
		Title:    "Use Zsh `*(.x)` glob qualifier instead of `find -executable`",
		Severity: SeverityStyle,
		Description: "Zsh's `*(.x)` glob qualifier matches regular files that are executable. " +
			"Avoid shelling out to `find -executable` when the same selection is one glob away.",
		Check: checkZC1341,
	})
}

func checkZC1341(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "find" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		if arg.String() == "-executable" {
			return []Violation{{
				KataID: "ZC1341",
				Message: "Use Zsh `*(.x)` glob qualifier instead of `find -executable`. " +
					"The `.` restricts to regular files and `x` to the executable bit.",
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
		ID:       "ZC1342",
		Title:    "Use Zsh `*(L0)` glob qualifier instead of `find -empty`",
		Severity: SeverityStyle,
		Description: "Zsh's `*(L0)` glob qualifier matches files with length 0. " +
			"Combine with `.` or `/` to restrict to regular files or directories. " +
			"Avoid shelling out to `find -empty` for the same result.",
		Check: checkZC1342,
	})
}

func checkZC1342(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "find" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		if arg.String() == "-empty" {
			return []Violation{{
				KataID: "ZC1342",
				Message: "Use Zsh `*(L0)` glob qualifier instead of `find -empty`. " +
					"Add `.` for regular files only: `*(.L0)`.",
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
		ID:       "ZC1343",
		Title:    "Use Zsh `*(m±N)` glob qualifier instead of `find -mtime N`",
		Severity: SeverityStyle,
		Description: "Zsh's `*(mN)`, `*(m+N)`, `*(m-N)` glob qualifiers match files by age in days " +
			"(exact / older / newer). For hours use `*(h±N)`, for minutes `*(M±N)`. " +
			"Same expressive power as `find -mtime`, no external process.",
		Check: checkZC1343,
	})
}

func checkZC1343(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "find" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-mtime" || val == "-mmin" || val == "-atime" || val == "-amin" ||
			val == "-ctime" || val == "-cmin" {
			return []Violation{{
				KataID: "ZC1343",
				Message: "Use Zsh glob qualifiers (`*(m±N)`, `*(M±N)`, `*(a±N)`, `*(c±N)`) instead of " +
					"`find -mtime`/`-mmin`/`-atime`/`-amin`/`-ctime`/`-cmin` for age predicates.",
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
		ID:       "ZC1344",
		Title:    "Use Zsh `*(L±Nk)` glob qualifier instead of `find -size`",
		Severity: SeverityStyle,
		Description: "Zsh's `*(LN)`, `*(L+N)`, `*(L-N)` match files by size in 512-byte blocks " +
			"(or bytes with a unit suffix: `k`, `m`, `p`). Same expressive power as " +
			"`find -size` without an external process.",
		Check: checkZC1344,
	})
}

func checkZC1344(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "find" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		if arg.String() == "-size" {
			return []Violation{{
				KataID: "ZC1344",
				Message: "Use Zsh `*(L±N[kmp])` glob qualifier instead of `find -size`. " +
					"Unit suffixes: `k` kilobytes, `m` megabytes, `p` 512-byte blocks.",
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
		ID:       "ZC1345",
		Title:    "Use Zsh `*(f:mode:)` glob qualifier instead of `find -perm`",
		Severity: SeverityStyle,
		Description: "Zsh's `*(f:mode:)` glob qualifier matches files by permission mode. " +
			"Use octal (`*(f:0755:)`) or symbolic (`*(f:u+x:)`) inside the colon-delimited form. " +
			"Avoids spawning `find` for permission filters.",
		Check: checkZC1345,
	})
}

func checkZC1345(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "find" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		if arg.String() == "-perm" {
			return []Violation{{
				KataID: "ZC1345",
				Message: "Use Zsh `*(f:mode:)` glob qualifier instead of `find -perm`. " +
					"Octal (`*(f:0755:)`) or symbolic (`*(f:u+x:)`) expressions are both supported.",
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
		ID:       "ZC1346",
		Title:    "Use Zsh `*(u:name:)` glob qualifier instead of `find -user`",
		Severity: SeverityStyle,
		Description: "Zsh's `*(u:name:)` and `*(u+uid)` glob qualifiers match files by owner " +
			"(name or numeric uid). The `*(U)` shorthand matches files owned by the current user. " +
			"Avoid `find -user` for the same selection.",
		Check: checkZC1346,
	})
}

func checkZC1346(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "find" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-user" || v == "-uid" || v == "-nouser" {
			return []Violation{{
				KataID: "ZC1346",
				Message: "Use Zsh `*(u:name:)` / `*(u+uid)` / `*(U)` glob qualifiers instead of " +
					"`find -user`/`-uid`/`-nouser`. Ownership predicates live entirely in the shell.",
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
		ID:       "ZC1347",
		Title:    "Use Zsh `*(g:name:)` glob qualifier instead of `find -group`",
		Severity: SeverityStyle,
		Description: "Zsh's `*(g:name:)` and `*(g+gid)` glob qualifiers match files by group " +
			"(name or numeric gid). The `*(G)` shorthand matches files in the current user's group. " +
			"Avoid `find -group`/`-gid`/`-nogroup` for the same selection.",
		Check: checkZC1347,
	})
}

func checkZC1347(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "find" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-group" || v == "-gid" || v == "-nogroup" {
			return []Violation{{
				KataID: "ZC1347",
				Message: "Use Zsh `*(g:name:)` / `*(g+gid)` / `*(G)` glob qualifiers instead of " +
					"`find -group`/`-gid`/`-nogroup`. Group predicates live entirely in the shell.",
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
		ID:       "ZC1348",
		Title:    "Use Zsh glob type qualifiers instead of `find -type`",
		Severity: SeverityStyle,
		Description: "Zsh glob qualifiers select node type directly: `*(/)` directories, `*(.)` " +
			"regular files, `*(@)` symlinks, `*(=)` sockets, `*(p)` named pipes, `*(*)` " +
			"executable regular files, `*(%)` char/block devices. Avoid `find -type X` for " +
			"the same selection.",
		Check: checkZC1348,
	})
}

func checkZC1348(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "find" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		if arg.String() == "-type" {
			return []Violation{{
				KataID: "ZC1348",
				Message: "Use Zsh glob type qualifiers (`*(/)`, `*(.)`, `*(@)`, `*(=)`, `*(p)`, `*(*)`, " +
					"`*(%)`) instead of `find -type`. No external process required.",
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
		ID:       "ZC1349",
		Title:    "Use `${#var}` instead of `expr length \"$var\"` for string length",
		Severity: SeverityStyle,
		Description: "Zsh (and POSIX) `${#var}` returns string length without spawning `expr`. " +
			"Use it wherever you would reach for `expr length` or `expr STRING : '.*'`.",
		Check: checkZC1349,
	})
}

func checkZC1349(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "expr" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		if arg.String() == "length" {
			return []Violation{{
				KataID: "ZC1349",
				Message: "Use `${#var}` instead of `expr length \"$var\"` for string length. " +
					"Parameter expansion avoids spawning an external process.",
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
		ID:       "ZC1350",
		Title:    "Use `${str:pos:len}` instead of `expr substr` for substring extraction",
		Severity: SeverityStyle,
		Description: "Zsh parameter expansion `${str:pos:len}` extracts a substring starting at " +
			"`pos` of length `len`. No external `expr` call, and the semantics are consistent " +
			"with `${str:pos}` (to end) and negative positions.",
		Check: checkZC1350,
	})
}

func checkZC1350(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "expr" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		if arg.String() == "substr" {
			return []Violation{{
				KataID: "ZC1350",
				Message: "Use `${str:pos:len}` instead of `expr substr` for substring extraction. " +
					"Parameter expansion avoids spawning an external process.",
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
		ID:       "ZC1351",
		Title:    "Use `[[ $str =~ pattern ]]` instead of `expr match` / `expr :` for regex",
		Severity: SeverityStyle,
		Description: "Zsh's `[[ $str =~ pattern ]]` evaluates regex natively and populates `$match` / " +
			"`$MATCH` / `$mbegin` / `$mend` arrays. Avoid shelling out to `expr match` or the " +
			"`expr STRING : REGEX` form.",
		Check: checkZC1351,
	})
}

func checkZC1351(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "expr" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "match" || v == "index" {
			return []Violation{{
				KataID: "ZC1351",
				Message: "Use `[[ $str =~ pattern ]]` with `$match` / `$MATCH` arrays instead of " +
					"`expr match`/`expr index`. Regex evaluation stays in the shell.",
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
		ID:       "ZC1352",
		Title:    "Avoid `xargs -I{}` — use a Zsh `for` loop for per-item substitution",
		Severity: SeverityStyle,
		Description: "`xargs -I{}` runs one command per item with `{}` substituted. A Zsh `for` " +
			"loop over the same input (`for x in ${(f)\"$(cmd)\"}`) is clearer and keeps state " +
			"in the current shell.",
		Check: checkZC1352,
	})
}

func checkZC1352(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "xargs" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		// -I, -I{}, -Irepl, --replace, --replace=STR
		if v == "-I" || v == "--replace" ||
			(len(v) > 2 && v[:2] == "-I") ||
			(len(v) > 9 && v[:10] == "--replace=") {
			return []Violation{{
				KataID: "ZC1352",
				Message: "Avoid `xargs -I{}` — iterate with `for x in ${(f)\"$(cmd)\"}; do ...; done` " +
					"in Zsh. A for loop avoids one-subprocess-per-item and keeps variables in scope.",
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
		ID:       "ZC1353",
		Title:    "Avoid `printf -v` — use `print -v` or command substitution in Zsh",
		Severity: SeverityStyle,
		Description: "`printf -v var fmt ...` is a Bash-ism. In Zsh use `print -v var -rf fmt ...` " +
			"or plain command substitution `var=$(printf fmt ...)`. `-v` is silently ignored by " +
			"POSIX printf, producing surprising bugs on portable scripts.",
		Check: checkZC1353,
	})
}

func checkZC1353(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "printf" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		if arg.String() == "-v" {
			return []Violation{{
				KataID: "ZC1353",
				Message: "Avoid `printf -v` in Zsh — use `print -v var -rf fmt ...` or " +
					"`var=$(printf fmt ...)`. `-v` is Bash-specific and ignored elsewhere.",
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
		ID:       "ZC1354",
		Title:    "Use `whence -w` instead of Bash-specific `type -t` for command classification",
		Severity: SeverityStyle,
		Description: "`type -t` returns the category (alias, keyword, function, builtin, file) " +
			"of a command in Bash. Zsh's `whence -w` produces `name: category` output with " +
			"the same information and without shelling out for the sub-field extraction.",
		Check: checkZC1354,
	})
}

func checkZC1354(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "type" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-t" || v == "-a" || v == "-P" {
			return []Violation{{
				KataID: "ZC1354",
				Message: "Use Zsh `whence -w` (category), `whence -a` (all), or `whence -p` (path) " +
					"instead of Bash-specific `type -t`/`-a`/`-P`.",
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
		ID:       "ZC1355",
		Title:    "Use `print -r` instead of `echo -E` for raw output",
		Severity: SeverityStyle,
		Description: "`echo -E` disables backslash interpretation, but the flag is Bash-ism and " +
			"ignored by POSIX `echo`. Zsh's `print -r` is the idiomatic raw-printer; combine " +
			"with `-n` (no newline), `-l` (one per line), `-u<fd>` (file descriptor), or `--` " +
			"(end of flags) as needed.",
		Check: checkZC1355,
		Fix:   fixZC1355,
	})
}

// fixZC1355 collapses `echo -E` into `print -r`. Span covers the
// command name, intervening whitespace, and the `-E` flag.
func fixZC1355(node ast.Node, v Violation, source []byte) []FixEdit {
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
	if i+2 > len(source) || source[i] != '-' || source[i+1] != 'E' {
		return nil
	}
	end := i + 2
	return []FixEdit{{
		Line:    v.Line,
		Column:  v.Column,
		Length:  end - nameOff,
		Replace: "print -r",
	}}
}

func checkZC1355(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "echo" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		if arg.String() == "-E" {
			return []Violation{{
				KataID: "ZC1355",
				Message: "Use `print -r` instead of `echo -E` for raw output. " +
					"`-E` is a Bash-ism and ignored by POSIX echo.",
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
		ID:       "ZC1356",
		Title:    "Use `read -A` instead of `read -a` for array read in Zsh",
		Severity: SeverityError,
		Description: "Zsh's `read` uses `-A` (uppercase A) to read into an array. Bash uses `-a` " +
			"(lowercase) for the same thing. In Zsh, `read -a` assigns a flag to a scalar " +
			"variable — not what Bash users expect. Use `-A` for portable-Zsh behavior.",
		Check: checkZC1356,
		Fix:   fixZC1356,
	})
}

// fixZC1356 rewrites the Bash-flavoured `read -a` flag to the
// uppercase `-A` that Zsh uses for array reads.
func fixZC1356(node ast.Node, _ Violation, source []byte) []FixEdit {
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

func checkZC1356(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "read" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		if arg.String() == "-a" {
			return []Violation{{
				KataID: "ZC1356",
				Message: "Use `read -A` (uppercase) in Zsh to read into an array. " +
					"`read -a` has different semantics in Zsh than in Bash.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1357",
		Title:    "Use Zsh `${(q)var}` instead of `printf '%q'` for shell-quoting",
		Severity: SeverityStyle,
		Description: "Bash's `printf '%q'` emits shell-quoted output. Zsh's `${(q)var}` parameter " +
			"flag does the same in-shell, with variants `${(qq)var}`, `${(qqq)var}`, `${(qqqq)var}` " +
			"for single-quote, double-quote, $'...', and POSIX ANSI-C styles respectively.",
		Check: checkZC1357,
	})
}

func checkZC1357(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "printf" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		val := arg.String()
		// Any format string containing %q (unescaped) triggers the kata.
		if strings.Contains(val, "%q") {
			return []Violation{{
				KataID: "ZC1357",
				Message: "Use Zsh `${(q)var}` for shell-quoting instead of `printf '%q'`. " +
					"Variants: `${(qq)}`, `${(qqq)}`, `${(qqqq)}` for different quote styles.",
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
		ID:       "ZC1358",
		Title:    "Use `${PWD:P}` instead of `pwd -P` for physical current directory",
		Severity: SeverityStyle,
		Description: "`pwd -P` resolves symlinks to the physical path. Zsh's `${PWD:P}` modifier " +
			"does the same without spawning the external — the `P` modifier returns the " +
			"canonical (absolute, symlink-resolved) form.",
		Check: checkZC1358,
	})
}

func checkZC1358(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "pwd" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-P" {
			return []Violation{{
				KataID: "ZC1358",
				Message: "Use `${PWD:P}` instead of `pwd -P` — the `P` modifier resolves symlinks " +
					"and returns the canonical path without spawning an external.",
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
		ID:       "ZC1359",
		Title:    "Avoid `id -Gn` — use Zsh `$groups` associative array",
		Severity: SeverityStyle,
		Description: "Zsh's `zsh/parameter` module exposes the `$groups` associative array mapping " +
			"group names to GIDs for the current process. Load with `zmodload zsh/parameter` " +
			"(often auto-loaded) and inspect `${(k)groups}` for names, avoiding an external " +
			"`id -Gn`/`groups` call.",
		Check: checkZC1359,
	})
}

func checkZC1359(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "id" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-Gn" || v == "-G" || v == "-gn" || v == "-g" {
			return []Violation{{
				KataID: "ZC1359",
				Message: "Avoid `id -Gn`/`-G`/`-gn`/`-g` — use Zsh `$groups` (names→gids assoc array) " +
					"or `$GID` for the primary group after `zmodload zsh/parameter`.",
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
		ID:       "ZC1360",
		Title:    "Use Zsh `*(OL)` glob qualifier instead of `ls -S` for size-ordered listing",
		Severity: SeverityStyle,
		Description: "Zsh glob qualifier `*(OL)` orders results by size (descending). `*(oL)` is " +
			"ascending. Combined with `[N]` subscript you get the N-th largest/smallest file " +
			"without `ls -S` and piping.",
		Check: checkZC1360,
	})
}

func checkZC1360(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "ls" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-S" || v == "-Sr" || v == "-lS" || v == "-lSr" {
			return []Violation{{
				KataID: "ZC1360",
				Message: "Use Zsh `*(OL)` (largest-first) or `*(oL)` (smallest-first) glob qualifier " +
					"instead of `ls -S`. No external process needed.",
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
		ID:       "ZC1361",
		Title:    "Avoid `awk 'NR==N'` — use Zsh array subscript on `${(f)...}`",
		Severity: SeverityStyle,
		Description: "Picking the N-th line with `awk 'NR==N'` spawns awk. Zsh can split file " +
			"contents on newlines with `${(f)\"$(<file)\"}` and index directly: `lines=(${(f)\"$(<f)\"}); print $lines[N]`.",
		Check: checkZC1361,
	})
}

func checkZC1361(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "awk" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		val := arg.String()
		// Look for NR== or NR == in the awk program
		if strings.Contains(val, "NR==") || strings.Contains(val, "NR ==") {
			return []Violation{{
				KataID: "ZC1361",
				Message: "Avoid `awk 'NR==N'` — split with `${(f)\"$(<file)\"}` in Zsh and index: " +
					"`lines=(${(f)\"$(<file)\"}); print $lines[N]`. No awk process needed.",
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
		ID:       "ZC1362",
		Title:    "Use `[[ -o option ]]` instead of `test -o option` for Zsh option checks",
		Severity: SeverityInfo,
		Description: "In Zsh, `[[ -o name ]]` tests whether a shell option is set. The `test` / `[` " +
			"builtin interprets `-o` as a logical OR, not an option-query — so `test -o foo` is " +
			"a syntax error or wrong behavior. Use the `[[ ... ]]` form for option tests.",
		Check: checkZC1362,
	})
}

func checkZC1362(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "test" && ident.Value != "[" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		if arg.String() == "-o" {
			return []Violation{{
				KataID: "ZC1362",
				Message: "Use `[[ -o option ]]` for option checks in Zsh — `test -o` means logical OR, " +
					"not option-query, producing wrong results.",
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
		ID:       "ZC1363",
		Title:    "Use Zsh `*(e:...:)` eval qualifier instead of `find -newer`/`-older`",
		Severity: SeverityStyle,
		Description: "Zsh's `*(e:expr:)` glob qualifier evaluates an arbitrary expression per match — " +
			"perfect for `-newer REF`-style predicates. Example: `*(e:'[[ $REPLY -nt reference ]]':)` " +
			"selects files newer than `reference`.",
		Check: checkZC1363,
	})
}

func checkZC1363(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "find" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-newer" || v == "-anewer" || v == "-cnewer" ||
			v == "-neweraa" || v == "-newercm" || v == "-newermt" {
			return []Violation{{
				KataID: "ZC1363",
				Message: "Use Zsh `*(e:'[[ $REPLY -nt REF ]]':)` eval glob qualifier instead of " +
					"`find -newer`/`-anewer`/`-cnewer`/`-newerXY`. `$REPLY` holds the current match.",
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
		ID:       "ZC1364",
		Title:    "Use Zsh `${var:pos:len}` instead of `cut -c` for character ranges",
		Severity: SeverityStyle,
		Description: "`cut -c N-M` extracts characters N through M from each line. Zsh's " +
			"`${var:pos:len}` (0-indexed position, length) does the same from a variable " +
			"without spawning `cut`.",
		Check: checkZC1364,
	})
}

func checkZC1364(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "cut" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		val := strings.TrimFunc(arg.String(), func(r rune) bool { return r == '\'' || r == '"' })
		if val == "-c" || val == "--characters" ||
			(len(val) > 2 && val[:2] == "-c") ||
			strings.HasPrefix(val, "--characters=") {
			return []Violation{{
				KataID: "ZC1364",
				Message: "Use Zsh `${var:pos:len}` for character ranges instead of `cut -c`. " +
					"Parameter expansion is in-shell and zero-indexed.",
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
		ID:       "ZC1365",
		Title:    "Use Zsh `zstat` module instead of `stat -c` for file metadata",
		Severity: SeverityStyle,
		Description: "Zsh's `zsh/stat` module (loaded with `zmodload zsh/stat` — the command is " +
			"named `zstat`) exposes every `stat(2)` field natively: mtime, size, owner, group, " +
			"mode, links, etc. Avoid external `stat -c '%...'` invocations.",
		Check: checkZC1365,
	})
}

func checkZC1365(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "stat" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-c" || v == "--format" || v == "--printf" {
			return []Violation{{
				KataID: "ZC1365",
				Message: "Use Zsh `zmodload zsh/stat; zstat -H meta file` for file metadata instead " +
					"of `stat -c '%...'`. The associative array `meta` exposes every stat field.",
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
		ID:       "ZC1366",
		Title:    "Use Zsh `limit` instead of POSIX `ulimit` for idiomatic resource queries",
		Severity: SeverityStyle,
		Description: "Zsh provides both `ulimit` (POSIX compatibility) and `limit` (Zsh native). " +
			"`limit` prints human-readable values (`cputime 10 seconds` vs `-t 10`) and accepts " +
			"`unlimited` as a value. Prefer `limit` for Zsh-idiomatic scripts; keep `ulimit` only " +
			"when the script must run under Bash as well.",
		Check: checkZC1366,
	})
}

func checkZC1366(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "ulimit" {
		return nil
	}

	return []Violation{{
		KataID: "ZC1366",
		Message: "Use Zsh `limit` (human-readable) or `limit -s` (stdout-only) instead of " +
			"POSIX `ulimit` for Zsh-native resource queries.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityStyle,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1367",
		Title:    "Use Zsh `strftime` instead of Bash `printf '%(fmt)T'`",
		Severity: SeverityStyle,
		Description: "Bash 4.2+ supports `printf '%(fmt)T\\n' seconds` to format a timestamp. Zsh's " +
			"`zsh/datetime` module provides `strftime` which is more readable and works " +
			"consistently across versions: `strftime '%Y-%m-%d' $EPOCHSECONDS`.",
		Check: checkZC1367,
	})
}

func checkZC1367(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "printf" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		val := arg.String()
		// Look for %(...)T format specifier
		if strings.Contains(val, ")T") && strings.Contains(val, "%(") {
			return []Violation{{
				KataID: "ZC1367",
				Message: "Use Zsh `strftime fmt seconds` (from `zsh/datetime`) instead of Bash " +
					"`printf '%(fmt)T' seconds`. Same formatting, more readable, no Bash-version gating.",
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
		ID:       "ZC1368",
		Title:    "Avoid `sh -c` / `bash -c` inside a Zsh script — inline or use a function",
		Severity: SeverityStyle,
		Description: "Invoking `sh -c` or `bash -c` inside a Zsh script spawns a second shell, " +
			"loses access to the parent script's functions, arrays, and associative arrays, and " +
			"re-interprets POSIX-only syntax. Inline the code as a function or use `zsh -c` when " +
			"a subshell is truly required.",
		Check: checkZC1368,
	})
}

func checkZC1368(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "sh" && ident.Value != "bash" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		if arg.String() == "-c" {
			return []Violation{{
				KataID: "ZC1368",
				Message: "Avoid `sh -c`/`bash -c` inside a Zsh script — inline the code as a " +
					"function to keep access to arrays, associative arrays, and Zsh features. " +
					"Use `zsh -c` only when a fresh shell is truly required.",
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
		ID:       "ZC1369",
		Title:    "Prefer Zsh `${(V)var}` over `od -c` for printable-visible character output",
		Severity: SeverityStyle,
		Description: "Zsh's `${(V)var}` parameter flag renders non-printable characters in " +
			"visible form (e.g. `\\n` for newline). For simple inspection of a variable's " +
			"contents, this avoids the `od -c` process entirely.",
		Check: checkZC1369,
	})
}

func checkZC1369(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "od" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-c" || v == "-C" {
			return []Violation{{
				KataID: "ZC1369",
				Message: "Use Zsh `${(V)var}` to see non-printable characters in a variable — " +
					"renders control chars as `\\n`, `\\t`, etc., without spawning `od`.",
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
		ID:       "ZC1370",
		Title:    "Prefer Zsh `repeat N { ... }` over `yes str | head -n N` for finite output",
		Severity: SeverityStyle,
		Description: "`yes` plus `head` is a common idiom for producing N copies of a line. " +
			"Zsh's `repeat N { print str }` does the same loop in-shell without spawning yes " +
			"or the pipe, and without the SIGPIPE handshake.",
		Check: checkZC1370,
	})
}

func checkZC1370(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "yes" {
		return nil
	}

	return []Violation{{
		KataID: "ZC1370",
		Message: "Prefer Zsh `repeat N { print str }` over `yes str | head -n N` for producing " +
			"N copies of a line. No external `yes` process, no pipe.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityStyle,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1371",
		Title:    "Use Zsh array `:t` modifier instead of `basename -a` for bulk path stripping",
		Severity: SeverityStyle,
		Description: "`basename -a a b c` returns the file name component of each path. Zsh's " +
			"`${array:t}` parameter modifier applies the same tail-component extraction to every " +
			"element of an array at once — no external process.",
		Check: checkZC1371,
	})
}

func checkZC1371(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "basename" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-a" || v == "--multiple" {
			return []Violation{{
				KataID: "ZC1371",
				Message: "Use Zsh `${paths:t}` on an array for bulk basename extraction instead of " +
					"`basename -a`. The `:t` modifier applies to every array element.",
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
		ID:       "ZC1372",
		Title:    "Use Zsh `zmv` autoload function instead of `rename`/`rename.ul`",
		Severity: SeverityStyle,
		Description: "Zsh's `zmv` (autoloaded via `autoload -Uz zmv`) batch-renames files using " +
			"glob patterns with capture groups. Safer than the various `rename`/`rename.ul`/`prename` " +
			"utilities (perl-based vs util-linux) and does not depend on which one is installed.",
		Check: checkZC1372,
	})
}

func checkZC1372(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "rename" && ident.Value != "rename.ul" && ident.Value != "prename" {
		return nil
	}

	return []Violation{{
		KataID: "ZC1372",
		Message: "Use Zsh `zmv` (autoload -Uz zmv) instead of `rename`/`rename.ul`/`prename`. " +
			"Glob-pattern renaming is handled in-shell with capture groups.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityStyle,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1373",
		Title:    "Use Zsh `${(0)var}` flag for NUL-split parsing instead of `env -0`",
		Severity: SeverityStyle,
		Description: "When reading NUL-terminated data (e.g. `/proc/*/environ`), Zsh's `${(0)var}` " +
			"parameter flag splits on NUL into an array natively. Avoid `env -0 | xargs -0 ...` " +
			"chains that require two additional processes.",
		Check: checkZC1373,
	})
}

func checkZC1373(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "env" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-0" || v == "--null" {
			return []Violation{{
				KataID: "ZC1373",
				Message: "Use Zsh `${(0)\"$(<file)\"}` to split NUL-terminated content in-shell. " +
					"`env -0` is usually followed by `xargs -0` or a read loop — both avoided.",
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

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1375",
		Title:    "Use `[[ -t fd ]]` instead of `tty -s` for tty-check",
		Severity: SeverityStyle,
		Description: "`tty -s` exits 0 if stdin is a terminal. Zsh's `[[ -t 0 ]]` (or `[[ -t 1 ]]` " +
			"for stdout, `[[ -t 2 ]]` for stderr) does the same check without spawning `tty`.",
		Check: checkZC1375,
	})
}

func checkZC1375(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "tty" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		if arg.String() == "-s" {
			return []Violation{{
				KataID: "ZC1375",
				Message: "Use `[[ -t 0 ]]` (stdin), `[[ -t 1 ]]` (stdout), or `[[ -t 2 ]]` (stderr) " +
					"instead of `tty -s`. In-shell file-descriptor test, no external process.",
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
		ID:       "ZC1376",
		Title:    "Avoid `BASH_XTRACEFD` — use Zsh `exec {fd}>file` + `setopt XTRACE`",
		Severity: SeverityWarning,
		Description: "Bash's `BASH_XTRACEFD` redirects `set -x` output to a file descriptor. Zsh " +
			"does not honor this variable; setting it is a silent no-op. To redirect trace output " +
			"in Zsh, open a dedicated fd with `exec {fd}>file` and redirect fd 2 through it: " +
			"`exec 2>&$fd; setopt XTRACE`.",
		Check: checkZC1376,
	})
}

func checkZC1376(node ast.Node) []Violation {
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
		if v == "$BASH_XTRACEFD" || v == "${BASH_XTRACEFD}" ||
			v == "BASH_XTRACEFD" {
			return []Violation{{
				KataID: "ZC1376",
				Message: "`BASH_XTRACEFD` is Bash-only. Zsh ignores it. Redirect trace output " +
					"with `exec {fd}>file; exec 2>&$fd; setopt XTRACE` instead.",
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
		ID:       "ZC1377",
		Title:    "Avoid `$BASH_ALIASES` — use Zsh `$aliases` associative array",
		Severity: SeverityWarning,
		Description: "Bash's `$BASH_ALIASES` is an associative array of alias→value mappings. Zsh " +
			"exposes the same information via `$aliases` (also an assoc array). `$BASH_ALIASES` " +
			"is unset in Zsh; reading it yields nothing.",
		Check: checkZC1377,
		Fix:   fixZC1377,
	})
}

// fixZC1377 renames every `BASH_ALIASES` token inside an echo / print /
// printf argument to `aliases`. Each occurrence becomes its own edit at
// the absolute source offset of that arg's token + the substring index;
// surrounding quoting and adjoining text stay byte-exact.
func fixZC1377(node ast.Node, _ Violation, source []byte) []FixEdit {
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
		if !strings.Contains(val, "BASH_ALIASES") {
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
			pos := strings.Index(val[idx:], "BASH_ALIASES")
			if pos < 0 {
				break
			}
			abs := off + idx + pos
			line, col := offsetLineColZC1377(source, abs)
			if line < 0 {
				break
			}
			edits = append(edits, FixEdit{
				Line:    line,
				Column:  col,
				Length:  len("BASH_ALIASES"),
				Replace: "aliases",
			})
			idx += pos + len("BASH_ALIASES")
		}
	}
	return edits
}

func offsetLineColZC1377(source []byte, offset int) (int, int) {
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

func checkZC1377(node ast.Node) []Violation {
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
		if strings.Contains(v, "BASH_ALIASES") {
			return []Violation{{
				KataID: "ZC1377",
				Message: "`$BASH_ALIASES` is Bash-only. In Zsh use `$aliases` (assoc array) — " +
					"same structure, e.g. `print -l ${(kv)aliases}`.",
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

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1379",
		Title:    "Avoid `$PROMPT_COMMAND` — use Zsh `precmd` function",
		Severity: SeverityWarning,
		Description: "Bash runs the command in `$PROMPT_COMMAND` before each prompt. Zsh does not " +
			"honor this variable; the equivalent is a function named `precmd` (or registered via " +
			"`add-zsh-hook precmd name`). Reading `$PROMPT_COMMAND` in Zsh is a no-op.",
		Check: checkZC1379,
	})
}

func checkZC1379(node ast.Node) []Violation {
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
		if strings.Contains(v, "PROMPT_COMMAND") {
			return []Violation{{
				KataID: "ZC1379",
				Message: "`PROMPT_COMMAND` is Bash-only. In Zsh define a `precmd` function or use " +
					"`autoload -Uz add-zsh-hook; add-zsh-hook precmd my_hook`.",
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

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1381",
		Title:    "Avoid `$COMP_WORDS`/`$COMP_CWORD` — Zsh uses `words`/`$CURRENT`",
		Severity: SeverityError,
		Description: "Bash programmable completion reads the partial command via `$COMP_WORDS` " +
			"(array of tokens) and `$COMP_CWORD` (index of cursor). Zsh's completion system " +
			"exposes the same via `words` (array) and `$CURRENT` (1-based cursor index). Using " +
			"the Bash names in Zsh completion functions produces empty expansions.",
		Check: checkZC1381,
		Fix:   fixZC1381,
	})
}

// fixZC1381 rewrites Bash completion variable names inside echo /
// print / printf args to their Zsh equivalents:
//
//	COMP_WORDS  → words
//	COMP_CWORD  → CURRENT
//	COMP_LINE   → BUFFER
//	COMP_POINT  → CURSOR
//
// Per-arg byte-anchored scan; one edit per match. Idempotent — a
// re-run sees the Zsh names, which the detector's substring guard
// won't match.
func fixZC1381(node ast.Node, _ Violation, source []byte) []FixEdit {
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
	mapping := []struct{ old, new string }{
		{"COMP_WORDS", "words"},
		{"COMP_CWORD", "CURRENT"},
		{"COMP_LINE", "BUFFER"},
		{"COMP_POINT", "CURSOR"},
	}
	var edits []FixEdit
	for _, arg := range cmd.Arguments {
		val := arg.String()
		tok := arg.TokenLiteralNode()
		off := LineColToByteOffset(source, tok.Line, tok.Column)
		if off < 0 || off+len(val) > len(source) {
			continue
		}
		if string(source[off:off+len(val)]) != val {
			continue
		}
		for _, m := range mapping {
			idx := 0
			for {
				pos := strings.Index(val[idx:], m.old)
				if pos < 0 {
					break
				}
				abs := off + idx + pos
				line, col := offsetLineColZC1381(source, abs)
				if line < 0 {
					break
				}
				edits = append(edits, FixEdit{
					Line:    line,
					Column:  col,
					Length:  len(m.old),
					Replace: m.new,
				})
				idx += pos + len(m.old)
			}
		}
	}
	return edits
}

func offsetLineColZC1381(source []byte, offset int) (int, int) {
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

func checkZC1381(node ast.Node) []Violation {
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
		if strings.Contains(v, "COMP_WORDS") || strings.Contains(v, "COMP_CWORD") ||
			strings.Contains(v, "COMP_LINE") || strings.Contains(v, "COMP_POINT") {
			return []Violation{{
				KataID: "ZC1381",
				Message: "Bash `$COMP_*` completion variables do not exist in Zsh. Use " +
					"`$words` (array of tokens), `$CURRENT` (cursor index), `$BUFFER`, or the " +
					"`_arguments`/`_values` helpers from `compsys`.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1382",
		Title:    "Avoid `$READLINE_LINE`/`$READLINE_POINT` — Zsh ZLE uses `$BUFFER`/`$CURSOR`",
		Severity: SeverityError,
		Description: "Bash readline exposes the current input line as `$READLINE_LINE` and cursor " +
			"offset as `$READLINE_POINT` inside `bind -x` handlers. Zsh's Line Editor (ZLE) uses " +
			"`$BUFFER` (line text) and `$CURSOR` (1-based column) inside widget functions. The " +
			"Bash names are unset in Zsh.",
		Check: checkZC1382,
		Fix:   fixZC1382,
	})
}

// fixZC1382 rewrites Bash readline variable names inside echo /
// print / printf args to their Zsh ZLE equivalents:
//
//	READLINE_LINE   → BUFFER
//	READLINE_POINT  → CURSOR
//	READLINE_MARK   → MARK
//
// Per-arg byte-anchored scan; one edit per match. Idempotent — a
// re-run sees the Zsh names, which the detector's substring guard
// won't match.
func fixZC1382(node ast.Node, _ Violation, source []byte) []FixEdit {
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
	mapping := []struct{ old, new string }{
		{"READLINE_LINE", "BUFFER"},
		{"READLINE_POINT", "CURSOR"},
		{"READLINE_MARK", "MARK"},
	}
	var edits []FixEdit
	for _, arg := range cmd.Arguments {
		val := arg.String()
		tok := arg.TokenLiteralNode()
		off := LineColToByteOffset(source, tok.Line, tok.Column)
		if off < 0 || off+len(val) > len(source) {
			continue
		}
		if string(source[off:off+len(val)]) != val {
			continue
		}
		for _, m := range mapping {
			idx := 0
			for {
				pos := strings.Index(val[idx:], m.old)
				if pos < 0 {
					break
				}
				abs := off + idx + pos
				line, col := offsetLineColZC1382(source, abs)
				if line < 0 {
					break
				}
				edits = append(edits, FixEdit{
					Line:    line,
					Column:  col,
					Length:  len(m.old),
					Replace: m.new,
				})
				idx += pos + len(m.old)
			}
		}
	}
	return edits
}

func offsetLineColZC1382(source []byte, offset int) (int, int) {
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

func checkZC1382(node ast.Node) []Violation {
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
		if strings.Contains(v, "READLINE_LINE") || strings.Contains(v, "READLINE_POINT") ||
			strings.Contains(v, "READLINE_MARK") {
			return []Violation{{
				KataID: "ZC1382",
				Message: "Bash `$READLINE_*` vars do not exist in Zsh. Inside ZLE widgets use " +
					"`$BUFFER`, `$CURSOR`, `$MARK`.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1383",
		Title:    "Avoid `$TIMEFORMAT` — Zsh uses `$TIMEFMT`",
		Severity: SeverityWarning,
		Description: "Bash's `$TIMEFORMAT` controls the output of the `time` builtin. Zsh uses a " +
			"shorter name, `$TIMEFMT`, for the same purpose. Setting `TIMEFORMAT` in a Zsh script " +
			"has no effect; the Zsh `time` builtin reads `$TIMEFMT`.",
		Check: checkZC1383,
		Fix:   fixZC1383,
	})
}

// fixZC1383 renames every `TIMEFORMAT` token inside an echo / print /
// printf / export argument to `TIMEFMT`. Each occurrence becomes its own
// edit at the absolute source offset of that arg's token + the substring
// index; surrounding quoting and adjoining text stay byte-exact.
func fixZC1383(node ast.Node, _ Violation, source []byte) []FixEdit {
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
		val := arg.String()
		if !strings.Contains(val, "TIMEFORMAT") {
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
			pos := strings.Index(val[idx:], "TIMEFORMAT")
			if pos < 0 {
				break
			}
			abs := off + idx + pos
			line, col := offsetLineColZC1383(source, abs)
			if line < 0 {
				break
			}
			edits = append(edits, FixEdit{
				Line:    line,
				Column:  col,
				Length:  len("TIMEFORMAT"),
				Replace: "TIMEFMT",
			})
			idx += pos + len("TIMEFORMAT")
		}
	}
	return edits
}

func offsetLineColZC1383(source []byte, offset int) (int, int) {
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

func checkZC1383(node ast.Node) []Violation {
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
		if strings.Contains(v, "TIMEFORMAT") {
			return []Violation{{
				KataID: "ZC1383",
				Message: "`$TIMEFORMAT` is Bash-only. Zsh reads `$TIMEFMT` (shorter name) for the " +
					"`time` builtin's output format.",
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
		ID:       "ZC1384",
		Title:    "Avoid `$EXECIGNORE` — Bash-only; Zsh uses completion-system ignore patterns",
		Severity: SeverityWarning,
		Description: "Bash's `$EXECIGNORE` excludes matching commands from PATH hashing. Zsh does " +
			"not honor this variable; use the compsys tag-based filters " +
			"(`zstyle ':completion:*' ignored-patterns ...`) for a similar effect on completion.",
		Check: checkZC1384,
	})
}

func checkZC1384(node ast.Node) []Violation {
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
		if strings.Contains(v, "EXECIGNORE") {
			return []Violation{{
				KataID: "ZC1384",
				Message: "`$EXECIGNORE` is Bash-only. For completion filtering in Zsh use " +
					"`zstyle ':completion:*' ignored-patterns 'pattern'`.",
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
		ID:       "ZC1385",
		Title:    "Avoid `$PS0` — Bash-only; Zsh uses `preexec` hook",
		Severity: SeverityWarning,
		Description: "Bash 4.4+ prints `$PS0` after reading a command and before executing it. Zsh " +
			"does not honor `$PS0`; the equivalent is a `preexec` function (or " +
			"`add-zsh-hook preexec funcname`) which receives the command line as `$1`.",
		Check: checkZC1385,
	})
}

func checkZC1385(node ast.Node) []Violation {
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
		if strings.Contains(v, "$PS0") || strings.Contains(v, "${PS0}") || v == "PS0" {
			return []Violation{{
				KataID: "ZC1385",
				Message: "`$PS0` is Bash-only. Zsh uses the `preexec` hook function for " +
					"pre-execution prompts.",
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
		ID:       "ZC1386",
		Title:    "Avoid `$FIGNORE` — Bash-only; Zsh uses compsys tag patterns",
		Severity: SeverityWarning,
		Description: "Bash's `$FIGNORE` hides filenames matching listed suffixes from completion. " +
			"Zsh does not honor this variable; use `zstyle ':completion:*' ignored-patterns '*.o *.pyc'` " +
			"or the file-patterns tag for equivalent filtering.",
		Check: checkZC1386,
	})
}

func checkZC1386(node ast.Node) []Violation {
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
		if strings.Contains(v, "FIGNORE") {
			return []Violation{{
				KataID: "ZC1386",
				Message: "`$FIGNORE` is Bash-only. In Zsh use " +
					"`zstyle ':completion:*' ignored-patterns '*.o *.pyc'` for completion filtering.",
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
		ID:       "ZC1387",
		Title:    "Avoid `$SHELLOPTS` — Zsh uses `$options` associative array",
		Severity: SeverityWarning,
		Description: "Bash's `$SHELLOPTS` is a colon-separated list of set options. Zsh exposes " +
			"the same information via the `$options` associative array (keys are option names, " +
			"values are `on`/`off`). `$SHELLOPTS` is unset in Zsh.",
		Check: checkZC1387,
	})
}

func checkZC1387(node ast.Node) []Violation {
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
		if strings.Contains(v, "SHELLOPTS") {
			return []Violation{{
				KataID: "ZC1387",
				Message: "`$SHELLOPTS` is Bash-only. In Zsh inspect `$options` (assoc array, " +
					"keys are option names) via `print -l ${(kv)options}`.",
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
		ID:       "ZC1388",
		Title:    "Use Zsh lowercase `$mailpath` array instead of colon-separated `$MAILPATH`",
		Severity: SeverityWarning,
		Description: "Bash uses `$MAILPATH` — a colon-separated string of mail files with " +
			"optional `?message` suffixes. Zsh uses lowercase `$mailpath` as an array (each " +
			"element: `file?message`), which is typed and parseable. Setting the uppercase " +
			"name in Zsh is ignored.",
		Check: checkZC1388,
	})
}

func checkZC1388(node ast.Node) []Violation {
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
		if strings.Contains(v, "$MAILPATH") || strings.Contains(v, "${MAILPATH}") ||
			strings.Contains(v, "MAILPATH=") {
			return []Violation{{
				KataID: "ZC1388",
				Message: "Use Zsh lowercase `$mailpath` (array) instead of Bash uppercase " +
					"`$MAILPATH` (colon-separated string).",
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
		ID:       "ZC1389",
		Title:    "Avoid `$HOSTFILE` — Bash-only; Zsh uses `$hosts` array",
		Severity: SeverityWarning,
		Description: "Bash reads `$HOSTFILE` to feed hostname completion. Zsh populates hostname " +
			"completion from the `$hosts` array (lowercase). Setting `$HOSTFILE` in Zsh is " +
			"ignored; extend `$hosts` instead.",
		Check: checkZC1389,
	})
}

func checkZC1389(node ast.Node) []Violation {
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
		if strings.Contains(v, "HOSTFILE") {
			return []Violation{{
				KataID: "ZC1389",
				Message: "`$HOSTFILE` is Bash-only. Zsh reads hostnames for completion from the " +
					"`$hosts` array (lowercase).",
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
		ID:       "ZC1390",
		Title:    "Avoid `$GROUPS[@]` — Zsh `$GROUPS` is a scalar, not an array",
		Severity: SeverityError,
		Description: "Bash's `$GROUPS` is an array of all group IDs the user belongs to, so " +
			"`${GROUPS[@]}` iterates them. In Zsh, `$GROUPS` is a scalar (primary GID). The " +
			"array of all group IDs is `$(groups)` output or `${(k)groups}` (if the " +
			"`zsh/parameter` module is loaded, `$groups` is an assoc array name→gid).",
		Check: checkZC1390,
	})
}

func checkZC1390(node ast.Node) []Violation {
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
		if strings.Contains(v, "GROUPS[") || strings.Contains(v, "${GROUPS[") {
			return []Violation{{
				KataID: "ZC1390",
				Message: "Zsh `$GROUPS` is a scalar (primary GID), not an array. For all group " +
					"IDs use `${(k)groups}` (after `zmodload zsh/parameter`).",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1391",
		Title:    "Avoid `[[ -v VAR ]]` for Bash set-check — use Zsh `(( ${+VAR} ))`",
		Severity: SeverityWarning,
		Description: "Bash 4.2+ supports `[[ -v VAR ]]` to test whether a variable is set. Zsh " +
			"`[[ -v VAR ]]` is parsed but not as the set-check — Zsh's canonical form is " +
			"`(( ${+VAR} ))` which evaluates to 1 when set and 0 when unset, working reliably " +
			"across Zsh versions.",
		Check: checkZC1391,
	})
}

func checkZC1391(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	// `[[` is its own node type in most AST designs, but bracket-test tokens
	// may come through as commands. We look for "-v" as a bare arg with a
	// following identifier in a context that looks like a bracket test.
	if ident.Value != "test" && ident.Value != "[" && ident.Value != "[[" {
		return nil
	}

	for i, arg := range cmd.Arguments {
		if arg.String() == "-v" && i+1 < len(cmd.Arguments) {
			next := cmd.Arguments[i+1].String()
			// Only flag if the "-v" is followed by an identifier (not a value to compare)
			if len(next) > 0 && !strings.Contains(next, "=") &&
				!strings.ContainsAny(next, "<>!/") {
				return []Violation{{
					KataID: "ZC1391",
					Message: "Use `(( ${+VAR} ))` for Zsh set-check — `-v` is a Bash 4.2+ " +
						"extension, not reliably portable to Zsh.",
					Line:   cmd.Token.Line,
					Column: cmd.Token.Column,
					Level:  SeverityWarning,
				}}
			}
		}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1392",
		Title:    "Avoid `$CHILD_MAX` — Bash-only; Zsh uses `limit` / `ulimit -u`",
		Severity: SeverityInfo,
		Description: "Bash's `$CHILD_MAX` reports the maximum number of exited child processes " +
			"Bash remembers. Zsh does not export this var. For current process limits use " +
			"`limit -s maxproc` or `ulimit -u` — but the exact Bash semantic is not mirrored.",
		Check: checkZC1392,
	})
}

func checkZC1392(node ast.Node) []Violation {
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
		if strings.Contains(v, "CHILD_MAX") {
			return []Violation{{
				KataID: "ZC1392",
				Message: "`$CHILD_MAX` is Bash-only. Zsh uses `limit -s maxproc` or `ulimit -u` " +
					"for process-count limits.",
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
		ID:       "ZC1393",
		Title:    "Avoid `$SRANDOM` — Bash 5.1+ only, read `/dev/urandom` in Zsh",
		Severity: SeverityWarning,
		Description: "Bash 5.1 added `$SRANDOM` as a cryptographically secure 32-bit random value. " +
			"Zsh does not have an equivalent variable. For secure random integers, read bytes " +
			"from `/dev/urandom` (e.g. `(( n = 0x$(od -N4 -An -tx1 /dev/urandom | tr -d ' ') ))`) " +
			"or use an external such as `openssl rand`.",
		Check: checkZC1393,
	})
}

func checkZC1393(node ast.Node) []Violation {
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
		if strings.Contains(v, "SRANDOM") {
			return []Violation{{
				KataID: "ZC1393",
				Message: "`$SRANDOM` is Bash 5.1+. In Zsh read `/dev/urandom` directly or use an " +
					"external (`openssl rand`) for secure random integers.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}

	return nil
}

// bashVarRE matches `$BASH` used as a standalone variable (not `$BASH_`).
var bashVarRE = regexp.MustCompile(`\$BASH(?:[^_A-Z]|$)`)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1394",
		Title:    "Avoid `$BASH` — Zsh uses `$ZSH_NAME` for the interpreter name",
		Severity: SeverityInfo,
		Description: "Bash's `$BASH` holds the path to the running Bash executable. Zsh's " +
			"equivalent is `$ZSH_NAME` (for the binary name) or `$0` (interactive shell). " +
			"Using `$BASH` in a Zsh script yields empty output.",
		Check: checkZC1394,
		Fix:   fixZC1394,
	})
}

// fixZC1394 renames every `$BASH` token (not part of a longer
// `$BASH_*` identifier) inside an echo / print / printf argument to
// `$ZSH_NAME`. Each occurrence becomes its own edit at the absolute
// source offset of that arg's token + the substring index; surrounding
// quoting, trailing punctuation, and adjoining text stay byte-exact.
func fixZC1394(node ast.Node, _ Violation, source []byte) []FixEdit {
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
		matches := bashVarRE.FindAllStringIndex(val, -1)
		if len(matches) == 0 {
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
		for _, m := range matches {
			// The regex spans `$BASH` plus one trailing byte (or end).
			// Rewrite only the `$BASH` prefix; leave the trailing byte
			// (the boundary char such as a space or quote) intact.
			abs := off + m[0]
			line, col := offsetLineColZC1394(source, abs)
			if line < 0 {
				continue
			}
			edits = append(edits, FixEdit{
				Line:    line,
				Column:  col,
				Length:  len("$BASH"),
				Replace: "$ZSH_NAME",
			})
		}
	}
	return edits
}

func offsetLineColZC1394(source []byte, offset int) (int, int) {
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

func checkZC1394(node ast.Node) []Violation {
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
		if bashVarRE.MatchString(v) {
			return []Violation{{
				KataID: "ZC1394",
				Message: "`$BASH` is Bash-only. Zsh exposes the interpreter name via `$ZSH_NAME` " +
					"and the executable path indirectly via `$0`.",
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
		ID:       "ZC1395",
		Title:    "Avoid `wait -n` — Bash 4.3+ only; Zsh `wait` on job IDs",
		Severity: SeverityWarning,
		Description: "Bash 4.3+ added `wait -n` (wait for any job to finish). Zsh's `wait` does " +
			"not accept `-n`; instead wait explicitly on job IDs or PIDs, or use `wait` with no " +
			"args (waits for all). For any-of semantics use `wait $pid1 $pid2; ...` in a loop.",
		Check: checkZC1395,
	})
}

func checkZC1395(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "wait" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		if arg.String() == "-n" {
			return []Violation{{
				KataID: "ZC1395",
				Message: "`wait -n` is Bash 4.3+. Zsh's `wait` waits on specific PIDs/jobs or " +
					"(bare `wait`) all jobs. For any-child semantics, loop over PIDs with " +
					"individual `wait $pid` calls.",
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
		ID:       "ZC1396",
		Title:    "Avoid `unset -n` — Bash nameref semantics not in Zsh",
		Severity: SeverityError,
		Description: "Bash's `unset -n NAME` unsets the nameref itself rather than the target " +
			"variable it points to. Zsh does not implement namerefs; `unset -n` flags as an " +
			"error or unsets something unintended. Use `unset -v` for variable unset and " +
			"`unset -f` for function unset explicitly.",
		Check: checkZC1396,
	})
}

func checkZC1396(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "unset" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		if arg.String() == "-n" {
			return []Violation{{
				KataID: "ZC1396",
				Message: "`unset -n` is a Bash nameref operation. Zsh does not honor it; use " +
					"`unset -v NAME` (variable) or `unset -f NAME` (function) explicitly.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1397",
		Title:    "Avoid `$COMP_TYPE`/`$COMP_KEY` — Bash completion globals, not in Zsh",
		Severity: SeverityError,
		Description: "Bash programmable completion exposes `$COMP_TYPE` (completion type) and " +
			"`$COMP_KEY` (completion key pressed). Zsh's compsys does not use these variables; " +
			"query completion context via `$compstate` assoc array or context keys from " +
			"`_arguments`/`_values` instead.",
		Check: checkZC1397,
	})
}

func checkZC1397(node ast.Node) []Violation {
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
		if strings.Contains(v, "COMP_TYPE") || strings.Contains(v, "COMP_KEY") ||
			strings.Contains(v, "COMP_WORDBREAKS") {
			return []Violation{{
				KataID: "ZC1397",
				Message: "Bash `$COMP_TYPE`/`$COMP_KEY`/`$COMP_WORDBREAKS` are not Zsh-native. " +
					"Use `$compstate` associative array for completion context in Zsh.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1398",
		Title:    "Avoid `$PROMPT_DIRTRIM` — use Zsh `%N~` prompt modifier",
		Severity: SeverityWarning,
		Description: "Bash's `$PROMPT_DIRTRIM` limits the number of directory components shown " +
			"in `\\w`. Zsh has no such variable; use the `%N~` prompt escape (N is component " +
			"count) or `%/` / `%~` with precmd adjustments for Zsh-native directory truncation.",
		Check: checkZC1398,
	})
}

func checkZC1398(node ast.Node) []Violation {
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
		if strings.Contains(v, "PROMPT_DIRTRIM") {
			return []Violation{{
				KataID: "ZC1398",
				Message: "`$PROMPT_DIRTRIM` is Bash-only. Use the Zsh prompt escape `%N~` " +
					"(N = number of path components to keep) for directory truncation.",
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
		ID:       "ZC1399",
		Title:    "Use Zsh `$signals` array instead of `kill -l` for signal enumeration",
		Severity: SeverityStyle,
		Description: "Zsh exposes the `$signals` array (from `zsh/parameter`) holding all signal " +
			"names indexed from 0. `print -l $signals` produces the same list as `kill -l` " +
			"without spawning an external process.",
		Check: checkZC1399,
	})
}

func checkZC1399(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "kill" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		if arg.String() == "-l" {
			return []Violation{{
				KataID: "ZC1399",
				Message: "Use Zsh `print -l $signals` (after `zmodload zsh/parameter`) instead " +
					"of `kill -l` for listing signal names.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityStyle,
			}}
		}
	}

	return nil
}
