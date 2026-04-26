// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1400",
		Title:    "Use Zsh `$CPUTYPE` for architecture detection instead of parsing `$HOSTTYPE`",
		Severity: SeverityInfo,
		Description: "Bash's `$HOSTTYPE` is a combined architecture/vendor/OS string (e.g. " +
			"`x86_64-pc-linux-gnu`). Zsh exposes the same as `$HOSTTYPE` but additionally splits " +
			"out `$CPUTYPE` (e.g. `x86_64`) for pure architecture queries — no `awk -F-` " +
			"needed to extract.",
		Check: checkZC1400,
	})
}

func checkZC1400(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	// Only fire when HOSTTYPE is being parsed/split (cut, awk, sed usage).
	switch ident.Value {
	case "cut", "awk", "sed":
	default:
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if strings.Contains(v, "HOSTTYPE") {
			return []Violation{{
				KataID: "ZC1400",
				Message: "Use Zsh `$CPUTYPE` for pure architecture instead of splitting " +
					"`$HOSTTYPE` with `cut`/`awk`/`sed`.",
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
		ID:       "ZC1401",
		Title:    "Prefer Zsh `$VENDOR` over parsing `$MACHTYPE` for vendor detection",
		Severity: SeverityInfo,
		Description: "Both Bash and Zsh expose `$MACHTYPE` (e.g. `x86_64-pc-linux-gnu`). Zsh " +
			"additionally pre-parses the vendor component into `$VENDOR` (e.g. `pc`, `apple`). " +
			"Avoid `cut -d- -f2 <<< $MACHTYPE` when `$VENDOR` is available directly.",
		Check: checkZC1401,
	})
}

func checkZC1401(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	switch ident.Value {
	case "cut", "awk", "sed":
	default:
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if strings.Contains(v, "MACHTYPE") {
			return []Violation{{
				KataID: "ZC1401",
				Message: "Use Zsh `$VENDOR` for vendor field instead of splitting `$MACHTYPE` " +
					"with `cut`/`awk`/`sed`.",
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
		ID:       "ZC1402",
		Title:    "Avoid `date -d @seconds` — use Zsh `strftime` for epoch formatting",
		Severity: SeverityStyle,
		Description: "`date -d @N -- '+fmt'` / `date --date=@N` converts epoch seconds to a " +
			"formatted date. Zsh's `zsh/datetime` module provides `strftime fmt N` directly " +
			"— a single builtin, no `date` spawn, and the `-d`/`@` form is GNU-specific " +
			"(not portable to BSD `date`).",
		Check: checkZC1402,
	})
}

func checkZC1402(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "date" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-d" || v == "--date" ||
			(len(v) > 6 && v[:6] == "--date=") {
			return []Violation{{
				KataID: "ZC1402",
				Message: "Use Zsh `strftime` (from `zsh/datetime`) instead of `date -d @N -- +fmt`. " +
					"The `-d`/`@` form is GNU-specific; `strftime` is portable Zsh.",
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
		ID:       "ZC1403",
		Title:    "Setting `$HISTFILESIZE` alone is incomplete in Zsh — pair with `$SAVEHIST`",
		Severity: SeverityWarning,
		Description: "Bash uses `$HISTSIZE` (in-memory) and `$HISTFILESIZE` (on disk). Zsh uses " +
			"`$HISTSIZE` (in-memory) and `$SAVEHIST` (on disk). Setting only `$HISTFILESIZE` in " +
			"Zsh has no effect on disk — `$SAVEHIST` must be set. Mixing both names leaves " +
			"disk-history behavior undefined.",
		Check: checkZC1403,
		Fix:   fixZC1403,
	})
}

// fixZC1403 rewrites `HISTFILESIZE` → `SAVEHIST` inside echo /
// print / printf / export args. Per-arg substring scan; one edit
// per match. Idempotent — a re-run sees `SAVEHIST`, which the
// detector's substring guard won't match.
func fixZC1403(node ast.Node, _ Violation, source []byte) []FixEdit {
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
	const oldName = "HISTFILESIZE"
	const newName = "SAVEHIST"
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
		idx := 0
		for {
			pos := strings.Index(val[idx:], oldName)
			if pos < 0 {
				break
			}
			abs := off + idx + pos
			line, col := offsetLineColZC1403(source, abs)
			if line < 0 {
				break
			}
			edits = append(edits, FixEdit{
				Line:    line,
				Column:  col,
				Length:  len(oldName),
				Replace: newName,
			})
			idx += pos + len(oldName)
		}
	}
	return edits
}

func offsetLineColZC1403(source []byte, offset int) (int, int) {
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

func checkZC1403(node ast.Node) []Violation {
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
		if strings.Contains(v, "HISTFILESIZE") {
			return []Violation{{
				KataID: "ZC1403",
				Message: "`$HISTFILESIZE` is Bash-only. Zsh uses `$SAVEHIST` for on-disk history " +
					"size. Setting `HISTFILESIZE` in Zsh has no effect.",
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
		ID:       "ZC1404",
		Title:    "Avoid `$BASH_CMDS` — Bash-specific hash-table mirror, use Zsh `$commands`",
		Severity: SeverityWarning,
		Description: "Bash's `$BASH_CMDS` associative array mirrors the hash-table of command " +
			"names→paths. Zsh exposes the same via `$commands` (assoc array from " +
			"`zsh/parameter`). `$BASH_CMDS` is unset in Zsh.",
		Check: checkZC1404,
		Fix:   fixZC1404,
	})
}

// fixZC1404 rewrites `BASH_CMDS` → `commands` inside echo / print /
// printf args. Per-arg substring scan; one edit per match.
// Idempotent — a re-run sees `commands`, which the detector's
// substring guard won't match.
func fixZC1404(node ast.Node, _ Violation, source []byte) []FixEdit {
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
	const oldName = "BASH_CMDS"
	const newName = "commands"
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
		idx := 0
		for {
			pos := strings.Index(val[idx:], oldName)
			if pos < 0 {
				break
			}
			abs := off + idx + pos
			line, col := offsetLineColZC1404(source, abs)
			if line < 0 {
				break
			}
			edits = append(edits, FixEdit{
				Line:    line,
				Column:  col,
				Length:  len(oldName),
				Replace: newName,
			})
			idx += pos + len(oldName)
		}
	}
	return edits
}

func offsetLineColZC1404(source []byte, offset int) (int, int) {
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

func checkZC1404(node ast.Node) []Violation {
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
		if strings.Contains(v, "BASH_CMDS") {
			return []Violation{{
				KataID: "ZC1404",
				Message: "`$BASH_CMDS` is Bash-only. In Zsh use `$commands` (assoc array, " +
					"names→paths) via `zsh/parameter`.",
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
		ID:       "ZC1405",
		Title:    "Avoid `env -u VAR cmd` — use Zsh `(unset VAR; cmd)` subshell",
		Severity: SeverityStyle,
		Description: "`env -u VAR cmd` unsets a variable for a single command. In Zsh the " +
			"idiomatic form is a subshell: `(unset VAR; cmd)` — no external `env` spawn, and " +
			"the unset is naturally scoped to the subshell.",
		Check: checkZC1405,
	})
}

func checkZC1405(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "env" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		if arg.String() == "-u" {
			return []Violation{{
				KataID: "ZC1405",
				Message: "Use `(unset VAR; cmd)` subshell instead of `env -u VAR cmd`. " +
					"In-shell scoping, no external process.",
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
		ID:       "ZC1406",
		Title:    "Prefer Zsh `zargs -P N` autoload over `xargs -P N` for parallel execution",
		Severity: SeverityStyle,
		Description: "Zsh provides `zargs` (loaded via `autoload -Uz zargs`) — a native equivalent " +
			"of `xargs` with parallel execution via `-P`. It keeps variables and functions in " +
			"scope (unlike xargs) and avoids the utility-quoting surprises of `xargs`.",
		Check: checkZC1406,
	})
}

func checkZC1406(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "xargs" {
		return nil
	}

	for i, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-P" || v == "--max-procs" ||
			(len(v) > 2 && v[:2] == "-P") {
			_ = i
			return []Violation{{
				KataID: "ZC1406",
				Message: "Consider `zargs -P N` (autoload -Uz zargs) instead of `xargs -P N`. " +
					"Parallel execution with Zsh functions in scope — no subshell-per-item.",
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
		ID:       "ZC1407",
		Title:    "Avoid `/dev/tcp/...` — use Zsh `zsh/net/tcp` module",
		Severity: SeverityError,
		Description: "`/dev/tcp/host/port` is a Bash-specific virtual-file interface for TCP " +
			"connections; Zsh does not implement it. For TCP in Zsh, load `zmodload zsh/net/tcp` " +
			"and use `ztcp host port` which exposes the connection as a regular file descriptor.",
		Check: checkZC1407,
	})
}

func checkZC1407(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	// Check all args for /dev/tcp or /dev/udp paths
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if strings.Contains(v, "/dev/tcp/") || strings.Contains(v, "/dev/udp/") {
			return []Violation{{
				KataID: "ZC1407",
				Message: "`/dev/tcp/...` and `/dev/udp/...` are Bash-only virtual files. In Zsh " +
					"load `zsh/net/tcp` and use `ztcp host port` / `ztcp -c $fd` for TCP.",
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
		ID:       "ZC1408",
		Title:    "Avoid `$BASH_FUNC_...%%` — Bash-specific exported-function envvar",
		Severity: SeverityError,
		Description: "Bash exports functions into environment variables named `BASH_FUNC_NAME%%`. " +
			"These are consumed only by other Bash shells. Zsh does not recognize the format " +
			"and will neither inherit the function nor clean these envvars.",
		Check: checkZC1408,
	})
}

func checkZC1408(node ast.Node) []Violation {
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
		if strings.Contains(v, "BASH_FUNC_") {
			return []Violation{{
				KataID: "ZC1408",
				Message: "`BASH_FUNC_*` exported-function envvars are Bash-only. Zsh does not " +
					"consume them; export function definitions via `autoload` + `$FPATH` instead.",
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
		ID:       "ZC1409",
		Title:    "Avoid `[ -N file ]` / `test -N file` — Bash-only, use Zsh `zstat` for mtime comparison",
		Severity: SeverityInfo,
		Description: "`[ -N file ]` and `test -N file` test whether a file has been modified since " +
			"last read (Bash extension). Zsh does not implement `-N`. Use the `zsh/stat` module " +
			"to compare `atime` and `mtime` explicitly: `zstat -H s file; (( s[mtime] > s[atime] ))`.",
		Check: checkZC1409,
	})
}

func checkZC1409(node ast.Node) []Violation {
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
		if arg.String() == "-N" {
			return []Violation{{
				KataID: "ZC1409",
				Message: "`test -N file` (modified-since-read) is a Bash extension. In Zsh use " +
					"`zmodload zsh/stat; zstat -H s file; (( s[mtime] > s[atime] ))`.",
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
		ID:       "ZC1410",
		Title:    "Avoid `compopt` — Bash programmable-completion modifier, not in Zsh",
		Severity: SeverityError,
		Description: "`compopt` tweaks Bash programmable-completion options for the current " +
			"completion. Zsh's compsys does not implement `compopt`; completion options are set " +
			"via `zstyle` / completion-function context instead.",
		Check: checkZC1410,
	})
}

func checkZC1410(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "compopt" {
		return nil
	}

	return []Violation{{
		KataID: "ZC1410",
		Message: "`compopt` is a Bash-only completion builtin. Zsh compsys uses `zstyle` " +
			"(e.g. `zstyle ':completion:*' menu select`) for equivalent tuning.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1411",
		Title:    "Use Zsh `disable` instead of Bash `enable -n` to hide builtins",
		Severity: SeverityStyle,
		Description: "Bash's `enable -n name` disables a builtin so that the external of the same " +
			"name is used. Zsh provides a dedicated `disable` builtin: `disable name` achieves " +
			"the same in one verb. Re-enable later with `enable name`.",
		Check: checkZC1411,
		Fix:   fixZC1411,
	})
}

// fixZC1411 collapses `enable -n NAME` into `disable NAME`. The span
// covers the `enable` command name and the `-n` flag in a single edit;
// trailing builtin name(s) stay in place.
func fixZC1411(node ast.Node, v Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "enable" {
		return nil
	}
	var dashN ast.Expression
	for _, arg := range cmd.Arguments {
		if arg.String() == "-n" {
			dashN = arg
			break
		}
	}
	if dashN == nil {
		return nil
	}
	nameOff := LineColToByteOffset(source, v.Line, v.Column)
	if nameOff < 0 || nameOff+len("enable") > len(source) {
		return nil
	}
	if string(source[nameOff:nameOff+len("enable")]) != "enable" {
		return nil
	}
	dashTok := dashN.TokenLiteralNode()
	dashOff := LineColToByteOffset(source, dashTok.Line, dashTok.Column)
	if dashOff < 0 || dashOff+2 > len(source) {
		return nil
	}
	if string(source[dashOff:dashOff+2]) != "-n" {
		return nil
	}
	return []FixEdit{{
		Line:    v.Line,
		Column:  v.Column,
		Length:  dashOff + 2 - nameOff,
		Replace: "disable",
	}}
}

func checkZC1411(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "enable" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		if arg.String() == "-n" {
			return []Violation{{
				KataID: "ZC1411",
				Message: "Use Zsh `disable name` instead of `enable -n name`. Zsh has a " +
					"dedicated `disable` builtin that reads clearer.",
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
		ID:       "ZC1412",
		Title:    "Avoid `$COMPREPLY` — Bash completion output, use Zsh `compadd`",
		Severity: SeverityError,
		Description: "Bash completion functions populate the `$COMPREPLY` array to declare " +
			"candidates. Zsh's compsys uses the `compadd` builtin: `compadd -- foo bar baz`. " +
			"Setting `$COMPREPLY` in a Zsh completion does nothing.",
		Check: checkZC1412,
	})
}

func checkZC1412(node ast.Node) []Violation {
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
		if strings.Contains(v, "COMPREPLY") {
			return []Violation{{
				KataID: "ZC1412",
				Message: "`$COMPREPLY` is a Bash-only completion output array. In Zsh compsys " +
					"use `compadd -- candidate1 candidate2`.",
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
		ID:       "ZC1413",
		Title:    "Use Zsh `whence -p cmd` instead of `hash -t cmd` for resolved path",
		Severity: SeverityStyle,
		Description: "Bash's `hash -t cmd` prints the hashed path for `cmd` (or fails if not " +
			"hashed). Zsh's `whence -p cmd` prints the PATH-resolved absolute path, whether " +
			"hashed or not — more reliable and the native Zsh idiom.",
		Check: checkZC1413,
		Fix:   fixZC1413,
	})
}

// fixZC1413 rewrites `hash -t cmd` to `whence -p cmd`. Two edits per
// fire: the command name and the `-t` flag. Idempotent — a re-run
// sees `whence`, not `hash`, so the detector won't fire. Defensive
// byte-match guards on both edits refuse to insert if the source
// at the offset doesn't match.
func fixZC1413(node ast.Node, v Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "hash" {
		return nil
	}
	var dashT ast.Node
	for _, arg := range cmd.Arguments {
		if arg.String() == "-t" {
			dashT = arg
			break
		}
	}
	if dashT == nil {
		return nil
	}
	cmdOff := LineColToByteOffset(source, v.Line, v.Column)
	if cmdOff < 0 || cmdOff+len("hash") > len(source) {
		return nil
	}
	if string(source[cmdOff:cmdOff+len("hash")]) != "hash" {
		return nil
	}
	tTok := dashT.TokenLiteralNode()
	tOff := LineColToByteOffset(source, tTok.Line, tTok.Column)
	if tOff < 0 || tOff+len("-t") > len(source) {
		return nil
	}
	if string(source[tOff:tOff+len("-t")]) != "-t" {
		return nil
	}
	return []FixEdit{
		{Line: v.Line, Column: v.Column, Length: len("hash"), Replace: "whence"},
		{Line: tTok.Line, Column: tTok.Column, Length: len("-t"), Replace: "-p"},
	}
}

func checkZC1413(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "hash" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		if arg.String() == "-t" {
			return []Violation{{
				KataID: "ZC1413",
				Message: "Use `whence -p cmd` (Zsh) instead of `hash -t cmd`. " +
					"`whence -p` always returns the absolute path, regardless of hash state.",
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
		ID:       "ZC1414",
		Title:    "Beware `hash -d` — Bash deletes from hash table, Zsh defines named directory",
		Severity: SeverityError,
		Description: "The `-d` flag has opposite meanings across shells: Bash `hash -d NAME` " +
			"removes `NAME` from the command-hash table. Zsh `hash -d NAME=PATH` **defines** a " +
			"named directory (`~NAME` expansion). A Bash script ported to Zsh breaks silently " +
			"when `hash -d ls` is interpreted as defining `~ls`.",
		Check: checkZC1414,
	})
}

func checkZC1414(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "hash" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		if arg.String() == "-d" {
			return []Violation{{
				KataID: "ZC1414",
				Message: "`hash -d` has opposite semantics in Bash (delete) vs Zsh (define " +
					"named directory). Use `unhash cmd` for Zsh command-hash removal, or " +
					"`hash -d NAME=/path` for named-directory definition.",
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
		ID:       "ZC1415",
		Title:    "Prefer Zsh `TRAPZERR` function over `trap 'cmd' ERR`",
		Severity: SeverityInfo,
		Description: "Both Bash and Zsh accept `trap 'cmd' ERR`, but Zsh's idiomatic form is the " +
			"named function `TRAPZERR`: `TRAPZERR() { echo \"err at $LINENO\"; }`. The named " +
			"function receives `$1` = signal and is easier to compose than an inline string.",
		Check: checkZC1415,
	})
}

func checkZC1415(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "trap" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "ERR" || v == "ZERR" {
			return []Violation{{
				KataID: "ZC1415",
				Message: "Prefer Zsh `TRAPZERR() { ... }` function over `trap 'cmd' ERR`. " +
					"The named-function form is more idiomatic and composable in Zsh.",
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
		ID:       "ZC1416",
		Title:    "Prefer Zsh `preexec` hook over `trap 'cmd' DEBUG`",
		Severity: SeverityWarning,
		Description: "Bash's `trap 'cmd' DEBUG` runs `cmd` before each simple command. Zsh's " +
			"equivalent is the `preexec` function (or `add-zsh-hook preexec name`) which " +
			"receives the about-to-execute command line as `$1`, `$2`, `$3`. The DEBUG trap " +
			"is not fired in Zsh the way it is in Bash — use preexec for portability.",
		Check: checkZC1416,
	})
}

func checkZC1416(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "trap" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		if arg.String() == "DEBUG" {
			return []Violation{{
				KataID: "ZC1416",
				Message: "Use Zsh `preexec() { ... }` (or `add-zsh-hook preexec`) instead of " +
					"`trap 'cmd' DEBUG`. Zsh's DEBUG trap does not fire the same way as Bash's.",
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
		ID:       "ZC1417",
		Title:    "Prefer Zsh `TRAPRETURN` function over `trap 'cmd' RETURN`",
		Severity: SeverityInfo,
		Description: "Bash's `trap 'cmd' RETURN` runs `cmd` when a function returns. Zsh accepts " +
			"the `RETURN` signal name but the idiomatic form is a function named `TRAPRETURN`: " +
			"`TRAPRETURN() { print \"returning $?\"; }`.",
		Check: checkZC1417,
	})
}

func checkZC1417(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "trap" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		if arg.String() == "RETURN" {
			return []Violation{{
				KataID: "ZC1417",
				Message: "Prefer Zsh `TRAPRETURN() { ... }` function over `trap 'cmd' RETURN`. " +
					"Named-function form is more idiomatic in Zsh.",
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
		ID:       "ZC1418",
		Title:    "Use Zsh `limit -h`/`-s` instead of `ulimit -H`/`-S` for hard/soft limits",
		Severity: SeverityStyle,
		Description: "Bash's `ulimit` uses uppercase `-H` (hard) and `-S` (soft). Zsh's native " +
			"`limit` builtin uses lowercase `-h` and `-s` for the same. The Zsh form is easier " +
			"to remember and produces human-readable output.",
		Check: checkZC1418,
	})
}

func checkZC1418(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "ulimit" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-H" || v == "-S" || v == "-HS" || v == "-SH" {
			return []Violation{{
				KataID: "ZC1418",
				Message: "Use Zsh `limit -h` (hard) / `limit -s` (soft) instead of " +
					"`ulimit -H`/`-S`. Zsh's `limit` builtin is more human-readable.",
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
		ID:       "ZC1419",
		Title:    "Avoid `chmod 777` — grants world-writable access",
		Severity: SeverityWarning,
		Description: "Mode 777 (or 0777) grants read/write/execute to owner, group, and world. " +
			"Files become world-writable, which on a multi-user system or inside a container " +
			"with mapped UIDs is almost always wrong. Use 755 for executables, 644 for regular " +
			"files, 700 for private directories, or `umask`-aware helpers.",
		Check: checkZC1419,
	})
}

func checkZC1419(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "chmod" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "777" || v == "0777" || v == "a+rwx" || v == "ugo+rwx" {
			return []Violation{{
				KataID: "ZC1419",
				Message: "Avoid `chmod 777`/`a+rwx` — grants world-writable access. Prefer " +
					"restrictive modes (755, 644, 700, 600) matched to the actual file purpose.",
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
		ID:       "ZC1420",
		Title:    "Avoid `chmod +s` / `chmod u+s` — setuid/setgid is a security risk",
		Severity: SeverityWarning,
		Description: "Setuid (mode bit 4000) and setgid (2000) cause the program to run with the " +
			"file-owner's (or group's) privileges, not the caller's. Any bug in such a program " +
			"is a privilege-escalation vector. Reserve setuid for audited, minimal binaries; " +
			"prefer sudo + policy, capabilities, or containers for less-trusted tooling.",
		Check: checkZC1420,
	})
}

func checkZC1420(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "chmod" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := strings.Trim(arg.String(), "'\"")
		// Symbolic setuid/setgid
		if v == "+s" || v == "u+s" || v == "g+s" || v == "+st" || v == "u+st" {
			return []Violation{{
				KataID: "ZC1420",
				Message: "`chmod +s` / `u+s` / `g+s` sets setuid/setgid — privilege-escalation risk. " +
					"Prefer sudo policy, capabilities, or containerization.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
		// Numeric modes starting with 4xxx or 2xxx (setuid/setgid)
		if len(v) == 4 && (v[0] == '4' || v[0] == '2' || v[0] == '6') &&
			v[1] >= '0' && v[1] <= '7' && v[2] >= '0' && v[2] <= '7' && v[3] >= '0' && v[3] <= '7' {
			return []Violation{{
				KataID:  "ZC1420",
				Message: "Numeric mode with leading 4/2/6 sets setuid/setgid — privilege-escalation risk.",
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
		ID:       "ZC1421",
		Title:    "Avoid `chpasswd` / `passwd --stdin` — plaintext passwords in process tree",
		Severity: SeverityError,
		Description: "Passing passwords on stdin to `chpasswd` or `passwd --stdin` exposes the " +
			"plaintext in the process command line or pipeline — visible to `ps`, logs, and " +
			"environment. Use encrypted-hash input (`chpasswd -e`), `usermod -p` with a hash, " +
			"or an IaC tool that handles credentials outside the process tree.",
		Check: checkZC1421,
	})
}

func checkZC1421(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "chpasswd" {
		return nil
	}

	// Any chpasswd invocation without -e (encrypted) is risky.
	hasEncrypted := false
	for _, arg := range cmd.Arguments {
		if arg.String() == "-e" || arg.String() == "--encrypted" {
			hasEncrypted = true
		}
	}
	if !hasEncrypted {
		return []Violation{{
			KataID: "ZC1421",
			Message: "`chpasswd` without `-e`/`--encrypted` accepts plaintext passwords — avoid " +
				"piping cleartext credentials into the process tree. Use a password hash (`-e`) " +
				"or a credentials store.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityError,
		}}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1422",
		Title:    "Avoid `sudo -S` — reads password from stdin, exposes plaintext",
		Severity: SeverityError,
		Description: "`sudo -S` reads the password from stdin, enabling `echo $PW | sudo -S cmd` " +
			"patterns that place the plaintext password in the process tree and shell history. " +
			"Prefer `sudo -A` with a graphical askpass, `NOPASSWD:` in sudoers for specific " +
			"commands, or `pkexec` for policy-based privilege elevation.",
		Check: checkZC1422,
	})
}

func checkZC1422(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "sudo" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		if arg.String() == "-S" {
			return []Violation{{
				KataID: "ZC1422",
				Message: "`sudo -S` enables password-via-stdin. Avoid piping plaintext " +
					"credentials. Use `sudo -A` (askpass), `NOPASSWD:` in sudoers, or `pkexec`.",
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
		ID:       "ZC1423",
		Title:    "Dangerous: `iptables -F` / `nft flush ruleset` — drops all firewall rules",
		Severity: SeverityWarning,
		Description: "Flushing the firewall ruleset removes every existing rule, typically " +
			"reverting to the default policy. On a remote machine with policy=DROP, this locks " +
			"you out. Save existing rules first (`iptables-save > backup`) and consider " +
			"`iptables-apply` with a rollback timer.",
		Check: checkZC1423,
	})
}

func checkZC1423(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	switch ident.Value {
	case "iptables", "ip6tables":
		for _, arg := range cmd.Arguments {
			if arg.String() == "-F" || arg.String() == "--flush" {
				return []Violation{{
					KataID: "ZC1423",
					Message: "Flushing firewall rules with `-F` removes every rule — risk of " +
						"locking yourself out of remote hosts. Save + use rollback mechanism.",
					Line:   cmd.Token.Line,
					Column: cmd.Token.Column,
					Level:  SeverityWarning,
				}}
			}
		}
	case "nft":
		for _, arg := range cmd.Arguments {
			if arg.String() == "flush" {
				return []Violation{{
					KataID: "ZC1423",
					Message: "`nft flush ruleset` clears every firewall table — risk of locking " +
						"yourself out of remote hosts. Save + use rollback mechanism.",
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
		ID:       "ZC1424",
		Title:    "Dangerous: `mkfs.*` / `mkfs -t` — formats a filesystem, destroys data",
		Severity: SeverityError,
		Description: "`mkfs.ext4 /dev/sda1`, `mkfs.xfs /dev/...`, `mkfs -t ...` all destroy the " +
			"existing filesystem on the target device. A typo on the target path reformats the " +
			"wrong disk. Validate the device path, use `blkid` / `lsblk` first, and consider a " +
			"confirmation prompt.",
		Check: checkZC1424,
	})
}

func checkZC1424(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	// mkfs, mkfs.ext4, mkfs.xfs, mkfs.btrfs, mkfs.vfat, etc.
	name := ident.Value
	if name == "mkfs" || strings.HasPrefix(name, "mkfs.") || name == "mke2fs" ||
		name == "mkswap" || name == "wipefs" {
		return []Violation{{
			KataID: "ZC1424",
			Message: "`" + name + "` formats / wipes a device — destroys data. Validate the " +
				"target with `lsblk` / `blkid` first, and consider an interactive confirmation.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityError,
		}}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1425",
		Title:    "`shutdown` / `reboot` / `halt` / `poweroff` — confirm before scripting",
		Severity: SeverityWarning,
		Description: "Scripts that invoke `shutdown`, `reboot`, `halt`, `poweroff`, or " +
			"`systemctl poweroff` take down the system. Unattended invocation in automation is " +
			"often wrong (e.g. leftover test step). Prefer `systemctl isolate rescue.target` for " +
			"controlled scenarios, and require explicit confirmation for interactive scripts.",
		Check: checkZC1425,
	})
}

func checkZC1425(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	switch ident.Value {
	case "shutdown", "reboot", "halt", "poweroff":
		return []Violation{{
			KataID: "ZC1425",
			Message: "`" + ident.Value + "` takes down the system. In scripts, confirm the " +
				"caller really wants this (interactive prompt, feature flag, or CI guard).",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityWarning,
		}}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1426",
		Title:    "Avoid `git clone http://` — unencrypted transport, use `https://` or `git://`+verify",
		Severity: SeverityWarning,
		Description: "`git clone http://...` transfers repository content unencrypted and " +
			"unauthenticated — susceptible to MITM insertion of malicious commits. Use " +
			"`https://` for authenticated hosts (GitHub, GitLab) or SSH (`git@host:path`) with " +
			"verified host keys. Plain `http://` has no integrity guarantee.",
		Check: checkZC1426,
	})
}

func checkZC1426(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "git" {
		return nil
	}

	isClone := false
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "clone" {
			isClone = true
			continue
		}
		if isClone && strings.HasPrefix(v, "http://") {
			return []Violation{{
				KataID: "ZC1426",
				Message: "`git clone http://` is unencrypted/unauthenticated. Use `https://` " +
					"or SSH with verified host keys.",
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
		ID:       "ZC1427",
		Title:    "Dangerous: `nc -e` / `ncat -e` — spawns arbitrary command on network connect",
		Severity: SeverityError,
		Description: "`nc -e cmd` and `ncat --exec cmd` pipe the network socket to an arbitrary " +
			"command. Incoming connections get a shell or any command you specify — the " +
			"classic reverse-shell pattern. Many distros ship `nc` compiled without `-e` for " +
			"this reason. Remove `-e` from scripts except in audited, restricted contexts.",
		Check: checkZC1427,
	})
}

func checkZC1427(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "nc" && ident.Value != "ncat" && ident.Value != "netcat" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-e" || v == "-c" {
			return []Violation{{
				KataID: "ZC1427",
				Message: "`" + ident.Value + " " + v + "` spawns an arbitrary command for " +
					"each connection — reverse-shell territory. Remove from scripts unless " +
					"audited and restricted.",
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
		ID:       "ZC1428",
		Title:    "Avoid `curl -u user:pass` — credentials visible in process list",
		Severity: SeverityError,
		Description: "`curl -u user:password` places the credentials in the command line, where " +
			"they show up in `ps`, `/proc/*/cmdline`, shell history, and most audit logs. Use " +
			"`-u user:` with an interactive password prompt, `--netrc`/`--netrc-file` for " +
			"persistent credentials, or a credentials manager.",
		Check: checkZC1428,
	})
}

func checkZC1428(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "curl" && ident.Value != "wget" {
		return nil
	}

	var sawU bool
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-u" {
			sawU = true
			continue
		}
		// Next arg after -u containing ':' signals user:pass literal
		if sawU {
			sawU = false
			for i := 0; i < len(v); i++ {
				if v[i] == ':' && i+1 < len(v) && v[i+1] != '\x00' {
					return []Violation{{
						KataID: "ZC1428",
						Message: "`curl -u user:pass` leaks credentials into the process list. " +
							"Use `-u user:` (prompt), `--netrc`, or a credentials manager.",
						Line:   cmd.Token.Line,
						Column: cmd.Token.Column,
						Level:  SeverityError,
					}}
				}
			}
		}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1429",
		Title:    "Avoid `umount -f` / `-l` — force/lazy unmount masks real issues",
		Severity: SeverityWarning,
		Description: "`umount -f` forces the unmount even if the FS is busy; `-l` (lazy) " +
			"detaches immediately but keeps the FS in-use. Both can leave stale file handles " +
			"and data loss. Fix the underlying 'target busy' (use `lsof` / `fuser -m` to find " +
			"users) instead of forcing.",
		Check: checkZC1429,
	})
}

func checkZC1429(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "umount" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-f" || v == "-l" || v == "-fl" || v == "-lf" {
			return []Violation{{
				KataID: "ZC1429",
				Message: "`umount -f`/`-l` force/lazy unmount masks the underlying 'busy' error. " +
					"Find open files with `lsof` / `fuser -m` and close them properly.",
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
		ID:       "ZC1430",
		Title:    "Prefer Zsh `zsh/sched` module over `at now` / `batch` for in-shell scheduling",
		Severity: SeverityStyle,
		Description: "`at`/`batch` schedule commands via the atd daemon — requires daemon " +
			"running, leaves a spool-file audit trail, and runs in a fresh environment. For " +
			"in-shell scheduling the Zsh `zsh/sched` module (`sched +1:00 cmd`) runs the " +
			"command from the current shell without the daemon dependency.",
		Check: checkZC1430,
	})
}

func checkZC1430(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "at" && ident.Value != "batch" {
		return nil
	}

	return []Violation{{
		KataID: "ZC1430",
		Message: "Prefer Zsh `sched` (from `zsh/sched`) for in-shell scheduling instead of " +
			"`at`/`batch`. No daemon dependency, runs in the current shell's environment.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityStyle,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1431",
		Title:    "Dangerous: `crontab -r` — removes all the user's cron jobs without confirmation",
		Severity: SeverityWarning,
		Description: "`crontab -r` deletes the entire crontab for the current user (or the target " +
			"user with `-u`). There is no `.bak` left behind, no `-i` prompt by default on most " +
			"platforms. Back up first with `crontab -l > /tmp/cron.bak`, then use `crontab -ir` " +
			"(interactive) to require confirmation.",
		Check: checkZC1431,
	})
}

func checkZC1431(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "crontab" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-r" || v == "-ur" || v == "-ru" {
			return []Violation{{
				KataID: "ZC1431",
				Message: "`crontab -r` removes all cron jobs with no backup. Save first " +
					"(`crontab -l > cron.bak`) and use `crontab -ir` for interactive confirmation.",
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
		ID:       "ZC1432",
		Title:    "Dangerous: `passwd -d user` — deletes the password, leaving the account passwordless",
		Severity: SeverityError,
		Description: "`passwd -d user` removes the password entirely, making the account usable " +
			"without any password (depending on PAM config). This is almost never what you want — " +
			"use `passwd -l user` to lock the account, or `usermod -L` + delete the ssh keys to " +
			"fully disable login.",
		Check: checkZC1432,
	})
}

func checkZC1432(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "passwd" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		if arg.String() == "-d" {
			return []Violation{{
				KataID: "ZC1432",
				Message: "`passwd -d user` deletes the password — account becomes passwordless. " +
					"Use `passwd -l user` to lock, or `usermod -L` + delete SSH keys to disable login.",
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
		ID:       "ZC1433",
		Title:    "Caution with `userdel -f` / `-r` — removes home directory and kills processes",
		Severity: SeverityWarning,
		Description: "`userdel -f` proceeds even when the user is logged in or has running " +
			"processes, potentially killing unsaved work. `-r` additionally deletes the home " +
			"directory and mail spool. Combined (`-rf`) these are destructive and often " +
			"misused for 'clean up a user' without warning. Verify no active sessions first.",
		Check: checkZC1433,
	})
}

func checkZC1433(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "userdel" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-f" || v == "-r" || v == "-rf" || v == "-fr" ||
			v == "--force" || v == "--remove" {
			return []Violation{{
				KataID: "ZC1433",
				Message: "`userdel -f`/`-r` forcibly removes user (kills processes, deletes home). " +
					"Check for active sessions first with `who -u` / `loginctl list-users`.",
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
		ID:       "ZC1434",
		Title:    "Warn on `swapoff -a` — disables all swap, can OOM-kill",
		Severity: SeverityWarning,
		Description: "`swapoff -a` disables every active swap. On a memory-constrained host " +
			"this pushes data back into RAM, potentially triggering OOM-killer. Prefer " +
			"disabling specific devices/files (`swapoff /swapfile`) and verify memory headroom " +
			"with `free -m` first.",
		Check: checkZC1434,
	})
}

func checkZC1434(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "swapoff" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		if arg.String() == "-a" || arg.String() == "--all" {
			return []Violation{{
				KataID: "ZC1434",
				Message: "`swapoff -a` disables ALL swap areas — risks OOM on memory-constrained " +
					"hosts. Disable specific swaps (`swapoff /swapfile`) after checking `free -m`.",
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
		ID:       "ZC1435",
		Title:    "Avoid `killall -9` / `killall -KILL` — force-kill by process name",
		Severity: SeverityWarning,
		Description: "`killall -9 name` sends SIGKILL to every process matching `name` — in " +
			"multi-user or containerized environments, this can hit unrelated processes that " +
			"happen to share the name. Prefer `killall -TERM` first (graceful), or kill by PID " +
			"after locating with `pgrep` / `pidof`.",
		Check: checkZC1435,
	})
}

func checkZC1435(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "killall" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-9" || v == "-KILL" || v == "-s" {
			return []Violation{{
				KataID: "ZC1435",
				Message: "`killall -9 name` force-kills every matching process, including " +
					"unrelated instances on multi-user or containerized hosts. Start with -TERM, " +
					"or kill by PID after `pgrep`/`pidof`.",
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
		ID:       "ZC1436",
		Title:    "`sysctl -w` is ephemeral — persist in `/etc/sysctl.d/*.conf` for surviving reboots",
		Severity: SeverityInfo,
		Description: "`sysctl -w key=value` sets a kernel parameter until the next reboot. For " +
			"configuration that must survive reboots, write a file in `/etc/sysctl.d/` and apply " +
			"with `sysctl --system`. Using only `-w` in provisioning scripts creates silent drift.",
		Check: checkZC1436,
	})
}

func checkZC1436(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "sysctl" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		if arg.String() == "-w" || arg.String() == "--write" {
			return []Violation{{
				KataID: "ZC1436",
				Message: "`sysctl -w` setting is lost on reboot. Persist in `/etc/sysctl.d/*.conf` " +
					"and reload with `sysctl --system`.",
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
		ID:       "ZC1437",
		Title:    "`dmesg -c` / `-C` clears the kernel ring buffer — destroys evidence",
		Severity: SeverityWarning,
		Description: "`dmesg -c` prints the ring buffer and then **clears** it. `dmesg -C` clears " +
			"without printing. Any later debugging loses the earlier messages. Prefer plain " +
			"`dmesg` for read-only inspection, or `journalctl -k` with a time filter.",
		Check: checkZC1437,
	})
}

func checkZC1437(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "dmesg" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-c" || v == "-C" || v == "--clear" || v == "--read-clear" {
			return []Violation{{
				KataID: "ZC1437",
				Message: "`dmesg -c`/`-C` clears the kernel ring buffer — subsequent debugging " +
					"loses earlier messages. Use plain `dmesg` or `journalctl -k`.",
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
		ID:       "ZC1438",
		Title:    "`systemctl mask` permanently prevents service start — document the unmask path",
		Severity: SeverityWarning,
		Description: "`systemctl mask unit` symlinks the unit to `/dev/null`, preventing any " +
			"start (manual, dependency, or at boot). Even `systemctl start` fails with 'Unit is " +
			"masked.'. The reverse `systemctl unmask` is easy to forget. Document the unmask in " +
			"provisioning scripts or use `disable` (which still allows manual start).",
		Check: checkZC1438,
	})
}

func checkZC1438(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "systemctl" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		if arg.String() == "mask" {
			return []Violation{{
				KataID: "ZC1438",
				Message: "`systemctl mask` permanently blocks service start. If this is a " +
					"policy choice, document the `unmask` path. For a softer block, use `disable`.",
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
		ID:       "ZC1439",
		Title:    "Enabling IP forwarding in a script — document firewall posture",
		Severity: SeverityWarning,
		Description: "Setting `net.ipv4.ip_forward=1` (or `-w`-ing a sysctl to the same effect) " +
			"turns the host into a router. Without matching iptables/nftables rules this can " +
			"silently expose services between interfaces. If intentional (VPN, container host, " +
			"NAT gateway), pair with explicit firewall rules and persist via `/etc/sysctl.d/`.",
		Check: checkZC1439,
	})
}

func checkZC1439(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "sysctl" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if strings.Contains(v, "ip_forward=1") ||
			strings.Contains(v, "forwarding=1") ||
			strings.Contains(v, "ip_forward =1") {
			return []Violation{{
				KataID: "ZC1439",
				Message: "Enabling `ip_forward` turns the host into a router. Verify firewall " +
					"posture (iptables/nftables) and persist the setting in `/etc/sysctl.d/`.",
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
		ID:       "ZC1440",
		Title:    "`usermod -G group user` replaces supplementary groups — use `-aG` to append",
		Severity: SeverityWarning,
		Description: "`usermod -G group user` overwrites the user's supplementary group list — " +
			"any prior group memberships are removed. Users commonly add themselves to `docker` " +
			"or `wheel` via `-G` and inadvertently lose `sudo`/`audio`/other memberships. Always " +
			"pair with `-a` (`-aG`) to append instead of replace.",
		Check: checkZC1440,
	})
}

func checkZC1440(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "usermod" {
		return nil
	}

	hasG := false
	hasA := false
	for _, arg := range cmd.Arguments {
		v := arg.String()
		switch v {
		case "-G", "--groups":
			hasG = true
		case "-a", "--append":
			hasA = true
		case "-aG", "-Ga":
			return nil // safe combined flag
		}
	}
	if hasG && !hasA {
		return []Violation{{
			KataID: "ZC1440",
			Message: "`usermod -G` without `-a` overwrites supplementary groups. Use `-aG` to " +
				"append — existing memberships are preserved.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityWarning,
		}}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1441",
		Title:    "Warn on `docker system prune -af` / `-a --force` (or similar podman/k8s)",
		Severity: SeverityWarning,
		Description: "`docker system prune -af` deletes every unused image, container, network, " +
			"and (with `--volumes`) volume. On shared CI runners or build hosts this obliterates " +
			"cached layers and slows future builds. Scope prunes with `--filter \"until=168h\"` " +
			"or target one resource type at a time.",
		Check: checkZC1441,
	})
}

func checkZC1441(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "docker" && ident.Value != "podman" {
		return nil
	}

	seenPrune := false
	seenA := false
	seenF := false
	seenVolumes := false
	for _, arg := range cmd.Arguments {
		v := arg.String()
		switch v {
		case "prune":
			seenPrune = true
		case "-a", "--all":
			seenA = true
		case "-f", "--force":
			seenF = true
		case "-af", "-fa":
			seenA = true
			seenF = true
		case "--volumes":
			seenVolumes = true
		}
	}
	// `--volumes` is the stricter superset handled by ZC1545; avoid double-firing.
	if seenPrune && seenA && seenF && !seenVolumes {
		return []Violation{{
			KataID: "ZC1441",
			Message: "`docker prune -af` / `-a --force` deletes all unused resources without " +
				"prompt. Scope with `--filter` or target one resource type.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityWarning,
		}}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1442",
		Title:    "Dangerous: `kubectl delete --all` / `--all-namespaces` deletes cluster resources",
		Severity: SeverityError,
		Description: "`kubectl delete --all pods` (in the current namespace) or " +
			"`-A`/`--all-namespaces` scopes delete operations across the whole cluster. A typo " +
			"on the resource type can wipe deployments, services, secrets, or even CRDs. " +
			"Always use `--dry-run=client` first, then apply with `-n` explicit namespace.",
		Check: checkZC1442,
	})
}

func checkZC1442(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "kubectl" && ident.Value != "oc" {
		return nil
	}

	hasDelete := false
	hasAll := false
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "delete" {
			hasDelete = true
		}
		if v == "--all" || v == "-A" || v == "--all-namespaces" {
			hasAll = true
		}
	}
	if hasDelete && hasAll {
		return []Violation{{
			KataID: "ZC1442",
			Message: "`kubectl delete --all` (or `-A`) deletes resources cluster-wide. Dry-run " +
				"with `--dry-run=client -o yaml` first, and scope with `-n` namespace.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityError,
		}}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1443",
		Title:    "Dangerous: `terraform destroy` / `apply -destroy` without `-target`",
		Severity: SeverityWarning,
		Description: "`terraform destroy` (or `terraform apply -destroy`) without a `-target` " +
			"removes every resource in state — entire environments, databases, volumes, DNS, " +
			"everything. Always prefer targeted destroy or scope via workspaces. Consider " +
			"guarding state-destroying commands behind an interactive confirmation.",
		Check: checkZC1443,
	})
}

func checkZC1443(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || (ident.Value != "terraform" && ident.Value != "tofu") {
		return nil
	}

	hasDestroy := false
	hasTarget := false
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "destroy" || v == "-destroy" {
			hasDestroy = true
		}
		if strings.HasPrefix(v, "-target") || v == "-target" {
			hasTarget = true
		}
	}
	if hasDestroy && !hasTarget {
		return []Violation{{
			KataID: "ZC1443",
			Message: "`terraform destroy` without `-target` removes every resource in state. " +
				"Scope with `-target=...` or gate behind interactive confirmation.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityWarning,
		}}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1444",
		Title:    "Dangerous: `redis-cli FLUSHALL` / `FLUSHDB` — wipes Redis data",
		Severity: SeverityError,
		Description: "`FLUSHALL` deletes every key in every database; `FLUSHDB` clears the current " +
			"DB. Running against production is usually catastrophic. Either rename the command " +
			"in `redis.conf` (`rename-command FLUSHALL \"\"`) or require an explicit " +
			"confirmation in scripts.",
		Check: checkZC1444,
	})
}

func checkZC1444(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "redis-cli" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := strings.ToUpper(arg.String())
		if v == "FLUSHALL" || v == "FLUSHDB" {
			return []Violation{{
				KataID: "ZC1444",
				Message: "`redis-cli FLUSHALL`/`FLUSHDB` wipes Redis data. Disable via " +
					"`rename-command` in redis.conf on production, or require explicit confirmation.",
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
		ID:       "ZC1445",
		Title:    "Dangerous: `dropdb` / `mysqladmin drop` — deletes a database",
		Severity: SeverityError,
		Description: "`dropdb NAME` removes a PostgreSQL database including all data and " +
			"schemas. `mysqladmin drop NAME` does the same for MySQL. Always `pg_dump` / " +
			"`mysqldump` first and consider requiring `-i`/`-y`-less forms so operators must " +
			"type confirmation.",
		Check: checkZC1445,
	})
}

func checkZC1445(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	switch ident.Value {
	case "dropdb":
		return []Violation{{
			KataID:  "ZC1445",
			Message: "`dropdb` removes a PostgreSQL database. Verify target and backup first (`pg_dump`).",
			Line:    cmd.Token.Line,
			Column:  cmd.Token.Column,
			Level:   SeverityError,
		}}
	case "mysqladmin":
		for _, arg := range cmd.Arguments {
			if arg.String() == "drop" {
				return []Violation{{
					KataID:  "ZC1445",
					Message: "`mysqladmin drop` removes a MySQL database. Verify target and backup first (`mysqldump`).",
					Line:    cmd.Token.Line,
					Column:  cmd.Token.Column,
					Level:   SeverityError,
				}}
			}
		}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1446",
		Title:    "Dangerous: `aws s3 rm --recursive` / `s3 rb --force` — bulk S3 deletion",
		Severity: SeverityError,
		Description: "`aws s3 rm s3://bucket/prefix --recursive` deletes every key under the " +
			"prefix. `aws s3 rb --force` deletes the bucket along with its contents. Combine " +
			"with a wrong prefix or bucket name and data loss is total. Enable versioning on " +
			"production buckets and use `aws s3api list-object-versions` before bulk removals.",
		Check: checkZC1446,
	})
}

func checkZC1446(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "aws" {
		return nil
	}

	var seenS3 bool
	var seenRm bool
	var seenRb bool
	var seenRecursive bool
	var seenForce bool
	for _, arg := range cmd.Arguments {
		v := arg.String()
		switch v {
		case "s3":
			seenS3 = true
		case "rm":
			seenRm = true
		case "rb":
			seenRb = true
		case "--recursive":
			seenRecursive = true
		case "--force":
			seenForce = true
		}
	}
	if seenS3 && ((seenRm && seenRecursive) || (seenRb && seenForce)) {
		return []Violation{{
			KataID: "ZC1446",
			Message: "`aws s3 rm --recursive` / `s3 rb --force` mass-deletes objects/buckets. " +
				"Enable versioning and dry-run with `--dryrun`.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityError,
		}}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1447",
		Title:    "Avoid deprecated `ifconfig` / `netstat` — prefer `ip` / `ss`",
		Severity: SeverityStyle,
		Description: "On modern Linux, `ifconfig` and `netstat` (from net-tools) are deprecated " +
			"in favor of the iproute2 suite: `ip addr`, `ip link`, `ip route`, `ss`. net-tools " +
			"is not installed by default on many distros (Alpine, Fedora Cloud, minimal images), " +
			"so scripts break. Use iproute2 commands for portability.",
		Check: checkZC1447,
	})
}

func checkZC1447(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	switch ident.Value {
	case "ifconfig":
		return []Violation{{
			KataID:  "ZC1447",
			Message: "`ifconfig` is deprecated. Use `ip addr` / `ip link` / `ip route` from iproute2.",
			Line:    cmd.Token.Line,
			Column:  cmd.Token.Column,
			Level:   SeverityStyle,
		}}
	case "netstat":
		return []Violation{{
			KataID:  "ZC1447",
			Message: "`netstat` is deprecated. Use `ss` from iproute2 (same flags, faster output).",
			Line:    cmd.Token.Line,
			Column:  cmd.Token.Column,
			Level:   SeverityStyle,
		}}
	case "route":
		return []Violation{{
			KataID:  "ZC1447",
			Message: "`route` is deprecated. Use `ip route`.",
			Line:    cmd.Token.Line,
			Column:  cmd.Token.Column,
			Level:   SeverityStyle,
		}}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1448",
		Title:    "`apt-get install` / `apt install` without `-y` hangs in non-interactive scripts",
		Severity: SeverityWarning,
		Description: "In provisioning scripts, `apt-get install foo` (no `-y`) waits for " +
			"interactive confirmation and stalls CI/Dockerfiles indefinitely. Always pass `-y` " +
			"(or `--yes`), and for unattended upgrades also set " +
			"`DEBIAN_FRONTEND=noninteractive` in the environment.",
		Check: checkZC1448,
		Fix:   fixZC1448,
	})
}

// fixZC1448 inserts ` -y` after the `apt` command name so install /
// upgrade / dist-upgrade / full-upgrade run without interactive
// confirmation. Only fires for plain `apt` — for `apt-get` the legacy
// ZC1213 fix already handles the rewrite, and emitting a duplicate
// zero-length insert here would yield ` -y -y` after both edits apply.
func fixZC1448(node ast.Node, v Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "apt" {
		return nil
	}
	nameOff := LineColToByteOffset(source, v.Line, v.Column)
	if nameOff < 0 {
		return nil
	}
	nameLen := IdentLenAt(source, nameOff)
	if nameLen != len("apt") {
		return nil
	}
	insertAt := nameOff + nameLen
	insLine, insCol := offsetLineColZC1448(source, insertAt)
	if insLine < 0 {
		return nil
	}
	return []FixEdit{{
		Line:    insLine,
		Column:  insCol,
		Length:  0,
		Replace: " -y",
	}}
}

func offsetLineColZC1448(source []byte, offset int) (int, int) {
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

func checkZC1448(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "apt-get" && ident.Value != "apt" {
		return nil
	}

	hasInstall := false
	hasYes := false
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "install" || v == "upgrade" || v == "dist-upgrade" || v == "full-upgrade" {
			hasInstall = true
		}
		if v == "-y" || v == "--yes" || v == "--assume-yes" {
			hasYes = true
		}
	}
	if hasInstall && !hasYes {
		return []Violation{{
			KataID: "ZC1448",
			Message: "`apt-get install`/`apt install` without `-y` hangs on the interactive " +
				"prompt in scripts. Add `-y` and set DEBIAN_FRONTEND=noninteractive.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityWarning,
		}}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1449",
		Title:    "`dnf`/`yum` install without `-y` hangs in non-interactive scripts",
		Severity: SeverityWarning,
		Description: "In CI/Dockerfiles, `dnf install pkg` or `yum install pkg` prompts for " +
			"confirmation and stalls. Always pass `-y` (or `--assumeyes`) for unattended runs. " +
			"Also consider `--nodocs` and `--setopt=install_weak_deps=False` for slim images.",
		Check: checkZC1449,
	})
}

func checkZC1449(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "dnf" && ident.Value != "yum" && ident.Value != "microdnf" {
		return nil
	}

	hasInstall := false
	hasYes := false
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "install" || v == "upgrade" || v == "update" || v == "remove" {
			hasInstall = true
		}
		if v == "-y" || v == "--assumeyes" {
			hasYes = true
		}
	}
	if hasInstall && !hasYes {
		return []Violation{{
			KataID: "ZC1449",
			Message: "`" + ident.Value + "` without `-y` hangs on confirmation. Add `-y` for " +
				"unattended runs.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityWarning,
		}}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1450",
		Title:    "`pacman -S` / `zypper install` without non-interactive flag hangs in scripts",
		Severity: SeverityWarning,
		Description: "Arch's `pacman -S` waits on confirmation unless `--noconfirm` is passed. " +
			"SUSE's `zypper install` needs `--non-interactive` (or `-n`). Both stall CI pipelines " +
			"and Dockerfiles without these flags.",
		Check: checkZC1450,
	})
}

func checkZC1450(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	switch ident.Value {
	case "pacman":
		hasInstall := false
		hasNoConfirm := false
		for _, arg := range cmd.Arguments {
			v := arg.String()
			if strings.HasPrefix(v, "-S") && !strings.HasPrefix(v, "-Ss") && !strings.HasPrefix(v, "-Si") {
				hasInstall = true
			}
			if v == "--noconfirm" {
				hasNoConfirm = true
			}
		}
		if hasInstall && !hasNoConfirm {
			return []Violation{{
				KataID:  "ZC1450",
				Message: "`pacman -S` without `--noconfirm` hangs in scripts.",
				Line:    cmd.Token.Line,
				Column:  cmd.Token.Column,
				Level:   SeverityWarning,
			}}
		}
	case "zypper":
		hasInstall := false
		hasN := false
		for _, arg := range cmd.Arguments {
			v := arg.String()
			if v == "install" || v == "in" || v == "update" || v == "up" {
				hasInstall = true
			}
			if v == "-n" || v == "--non-interactive" {
				hasN = true
			}
		}
		if hasInstall && !hasN {
			return []Violation{{
				KataID:  "ZC1450",
				Message: "`zypper install` without `--non-interactive` (`-n`) hangs in scripts.",
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
		ID:       "ZC1451",
		Title:    "Avoid `pip install` without `--user` or virtualenv",
		Severity: SeverityWarning,
		Description: "`pip install pkg` (no `--user`, no active venv) targets the system Python, " +
			"potentially breaking system tools or requiring sudo. On modern Linux this now fails " +
			"with PEP 668 `externally-managed-environment`. Always use a virtualenv (`python -m " +
			"venv`, `uv`, `poetry`) or `--user` for scoped installs.",
		Check: checkZC1451,
	})
}

func checkZC1451(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "pip" && ident.Value != "pip3" {
		return nil
	}

	hasInstall := false
	hasUser := false
	hasBreakSystem := false
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "install" {
			hasInstall = true
		}
		if v == "--user" {
			hasUser = true
		}
		if v == "--break-system-packages" {
			hasBreakSystem = true
		}
	}
	if hasInstall && !hasUser && !hasBreakSystem {
		return []Violation{{
			KataID: "ZC1451",
			Message: "`pip install` without `--user` or an active venv targets system Python. " +
				"Use `python -m venv` / `uv` / `--user` for scoped installs.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityWarning,
		}}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1452",
		Title:    "Avoid `npm install -g` — global installs need root, break under multiple Node versions",
		Severity: SeverityStyle,
		Description: "`npm install -g` places packages in a system-wide prefix (typically " +
			"`/usr/local`). That requires sudo, conflicts with Node version managers (nvm, " +
			"asdf, volta), and is rarely what you want in a project. Prefer project-local " +
			"installs (`npm i`), or `pnpm dlx`/`npx` for one-off tools.",
		Check: checkZC1452,
	})
}

func checkZC1452(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "npm" && ident.Value != "yarn" && ident.Value != "pnpm" {
		return nil
	}

	hasInstall := false
	hasGlobal := false
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "install" || v == "i" || v == "add" {
			hasInstall = true
		}
		if v == "-g" || v == "--global" {
			hasGlobal = true
		}
	}
	if hasInstall && hasGlobal {
		return []Violation{{
			KataID: "ZC1452",
			Message: "`" + ident.Value + " install -g` installs system-wide. Prefer project-local " +
				"install or `npx`/`pnpm dlx` for one-off tools.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityStyle,
		}}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1453",
		Title:    "Avoid `sudo pip` / `sudo npm` / `sudo gem` — language package managers as root",
		Severity: SeverityWarning,
		Description: "Running a language package manager as root installs third-party code with " +
			"full privileges, may overwrite distro-managed libs, and can execute arbitrary " +
			"install-time hooks as root. Use `--user`, a virtualenv/venv, or a version manager " +
			"(nvm, pyenv, rbenv) instead.",
		Check: checkZC1453,
	})
}

func checkZC1453(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "sudo" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		switch v {
		case "pip", "pip3", "npm", "yarn", "pnpm", "gem", "cpan", "luarocks":
			return []Violation{{
				KataID: "ZC1453",
				Message: "`sudo " + v + "` runs a language package manager as root. Prefer " +
					"`--user`, a virtualenv/venv, or a version manager (nvm/pyenv/rbenv).",
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
		ID:       "ZC1454",
		Title:    "Avoid `docker/podman run --privileged` — disables most container isolation",
		Severity: SeverityError,
		Description: "`--privileged` disables the seccomp profile, grants all Linux capabilities, " +
			"and lets the container access all host devices. It is effectively equivalent to " +
			"running the process as host root. Add specific capabilities with `--cap-add` and " +
			"bind-mount specific devices with `--device` instead.",
		Check: checkZC1454,
	})
}

func checkZC1454(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "docker" && ident.Value != "podman" {
		return nil
	}

	hasRun := false
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "run" || v == "exec" || v == "create" {
			hasRun = true
		}
		if hasRun && v == "--privileged" {
			return []Violation{{
				KataID: "ZC1454",
				Message: "`--privileged` disables container isolation — effectively host root. " +
					"Use `--cap-add` + `--device` for narrow permissions.",
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
		ID:       "ZC1455",
		Title:    "Avoid `docker run --net=host` / `--network=host` — disables network isolation",
		Severity: SeverityWarning,
		Description: "Host networking gives the container direct access to the host's network " +
			"stack, including localhost services. A vulnerable container can reach services " +
			"meant to be local-only. Use `-p hostport:containerport` for specific publishes and " +
			"dedicated networks for inter-container traffic.",
		Check: checkZC1455,
	})
}

func checkZC1455(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "docker" && ident.Value != "podman" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "--net=host" || v == "--network=host" || v == "-net=host" {
			return []Violation{{
				KataID: "ZC1455",
				Message: "`--net=host` / `--network=host` lets the container reach host-local " +
					"services. Use `-p` for explicit publishes or dedicated container networks.",
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
		ID:       "ZC1456",
		Title:    "Avoid `docker run -v /:...` — bind-mounts host root into container",
		Severity: SeverityError,
		Description: "Mounting `/` (host root) into a container gives the container read/write " +
			"access to the entire host filesystem — a trivial container escape. Mount only the " +
			"specific host paths the container needs, using `:ro` for read-only where possible.",
		Check: checkZC1456,
	})
}

func checkZC1456(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "docker" && ident.Value != "podman" {
		return nil
	}

	var prevV bool
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if prevV {
			prevV = false
			// v should be host:container or host:container:opts
			if strings.HasPrefix(v, "/:") || v == "/" {
				return []Violation{{
					KataID: "ZC1456",
					Message: "`-v /:...` mounts the host root into the container — trivial " +
						"container escape. Scope to specific paths.",
					Line:   cmd.Token.Line,
					Column: cmd.Token.Column,
					Level:  SeverityError,
				}}
			}
		}
		if v == "-v" || v == "--volume" {
			prevV = true
		}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1457",
		Title:    "Warn on bind-mount of `/var/run/docker.sock` — container escape vector",
		Severity: SeverityWarning,
		Description: "Mounting `/var/run/docker.sock` into a container lets the container start " +
			"any privileged container, mount host filesystems, and effectively gain root on the " +
			"host. Reserve this for trusted CI/tooling images; for general workloads use " +
			"rootless containers or a dedicated orchestrator API.",
		Check: checkZC1457,
	})
}

func checkZC1457(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "docker" && ident.Value != "podman" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if strings.Contains(v, "docker.sock") || strings.Contains(v, "/var/run/docker") {
			return []Violation{{
				KataID: "ZC1457",
				Message: "Mounting `/var/run/docker.sock` gives the container effective root on " +
					"the host. Reserve for trusted tooling.",
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
		ID:       "ZC1458",
		Title:    "Warn on explicit `docker run --user root` / `--user 0`",
		Severity: SeverityWarning,
		Description: "Running as UID 0 inside a container means a break-out bug leaves the " +
			"attacker as root on the host (absent user namespaces). Build images with a " +
			"non-root `USER` directive and avoid overriding to root at runtime.",
		Check: checkZC1458,
	})
}

func checkZC1458(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "docker" && ident.Value != "podman" {
		return nil
	}

	var prevU bool
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if prevU {
			prevU = false
			if v == "root" || v == "0" || strings.HasPrefix(v, "0:") || strings.HasPrefix(v, "root:") {
				return []Violation{{
					KataID: "ZC1458",
					Message: "Explicit root UID inside a container lets container-escape bugs " +
						"become host root. Use a non-root USER in the image.",
					Line:   cmd.Token.Line,
					Column: cmd.Token.Column,
					Level:  SeverityWarning,
				}}
			}
		}
		if v == "-u" || v == "--user" {
			prevU = true
		}
	}

	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1459",
		Title:    "Warn on `docker run --cap-add=SYS_ADMIN` / other dangerous capabilities",
		Severity: SeverityWarning,
		Description: "Granting `SYS_ADMIN`, `SYS_PTRACE`, `SYS_MODULE`, `NET_ADMIN`, or `ALL` " +
			"capabilities effectively disables the container's security boundary — most " +
			"container escapes rely on exactly these. Drop all capabilities and add back only " +
			"the specific ones the workload needs (usually none).",
		Check: checkZC1459,
	})
}

var dangerousCaps = map[string]struct{}{
	"SYS_ADMIN":       {},
	"SYS_PTRACE":      {},
	"SYS_MODULE":      {},
	"SYS_RAWIO":       {},
	"NET_ADMIN":       {},
	"DAC_READ_SEARCH": {},
	"ALL":             {},
}

func checkZC1459(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "docker" && ident.Value != "podman" {
		return nil
	}

	var prevCap bool
	for _, arg := range cmd.Arguments {
		v := arg.String()

		// Split form: --cap-add=VALUE
		if strings.HasPrefix(v, "--cap-add=") {
			val := strings.TrimPrefix(v, "--cap-add=")
			if _, bad := dangerousCaps[strings.ToUpper(val)]; bad {
				return violateZC1459(cmd)
			}
			continue
		}

		// Space form: --cap-add VALUE
		if prevCap {
			prevCap = false
			if _, bad := dangerousCaps[strings.ToUpper(v)]; bad {
				return violateZC1459(cmd)
			}
		}
		if v == "--cap-add" {
			prevCap = true
		}
	}

	return nil
}

func violateZC1459(cmd *ast.SimpleCommand) []Violation {
	return []Violation{{
		KataID: "ZC1459",
		Message: "Dangerous Linux capability granted — breaks the container's security " +
			"boundary. Prefer `--cap-drop=ALL` and add back only minimum needed.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1460",
		Title:    "Warn on `docker run --security-opt seccomp=unconfined` / `apparmor=unconfined`",
		Severity: SeverityWarning,
		Description: "Disabling seccomp or AppArmor removes the syscall / MAC filter that blocks " +
			"most container escape exploits. Only disable these in a known-safe development " +
			"context; production workloads should keep the default profile or ship a stricter " +
			"custom profile.",
		Check: checkZC1460,
	})
}

func checkZC1460(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "docker" && ident.Value != "podman" {
		return nil
	}

	var prevOpt bool
	for _, arg := range cmd.Arguments {
		v := arg.String()

		if strings.HasPrefix(v, "--security-opt=") {
			val := strings.TrimPrefix(v, "--security-opt=")
			if isUnconfined(val) {
				return violateZC1460(cmd)
			}
			continue
		}

		if prevOpt {
			prevOpt = false
			if isUnconfined(v) {
				return violateZC1460(cmd)
			}
		}
		if v == "--security-opt" {
			prevOpt = true
		}
	}

	return nil
}

func isUnconfined(v string) bool {
	return v == "seccomp=unconfined" ||
		v == "apparmor=unconfined" ||
		v == "seccomp:unconfined" ||
		v == "apparmor:unconfined"
}

func violateZC1460(cmd *ast.SimpleCommand) []Violation {
	return []Violation{{
		KataID: "ZC1460",
		Message: "Disabling seccomp or AppArmor removes the main syscall/MAC filter that " +
			"blocks container escapes. Keep the default profile.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1461",
		Title:    "Avoid `docker run --pid=host` — shares host PID namespace with the container",
		Severity: SeverityWarning,
		Description: "`--pid=host` lets the container see every host process and send signals to " +
			"them, including sending SIGKILL to init-managed daemons or attaching a debugger to " +
			"host-side processes. Use only for diagnostic tools (e.g. strace/perf containers) and " +
			"never for general workloads.",
		Check: checkZC1461,
	})
}

func checkZC1461(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "docker" && ident.Value != "podman" {
		return nil
	}

	var prevPid bool
	for _, arg := range cmd.Arguments {
		v := arg.String()

		if v == "--pid=host" {
			return violateZC1461(cmd)
		}
		if prevPid {
			prevPid = false
			if v == "host" {
				return violateZC1461(cmd)
			}
		}
		if v == "--pid" {
			prevPid = true
		}
	}

	return nil
}

func violateZC1461(cmd *ast.SimpleCommand) []Violation {
	return []Violation{{
		KataID: "ZC1461",
		Message: "`--pid=host` shares the host PID namespace — container can signal and " +
			"inspect every host process. Avoid outside debug tools.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1462",
		Title:    "Avoid `docker run --ipc=host` — shares host IPC namespace (/dev/shm, SysV IPC)",
		Severity: SeverityWarning,
		Description: "`--ipc=host` makes the container share `/dev/shm` and the SysV IPC keyspace " +
			"with the host. Any process on the host can read/write the container's shared memory " +
			"(and vice-versa), making side-channel and data-theft attacks trivial. Use the default " +
			"private IPC namespace unless two containers explicitly need to share IPC.",
		Check: checkZC1462,
	})
}

func checkZC1462(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "docker" && ident.Value != "podman" {
		return nil
	}

	var prevIpc bool
	for _, arg := range cmd.Arguments {
		v := arg.String()

		if v == "--ipc=host" {
			return violateZC1462(cmd)
		}
		if prevIpc {
			prevIpc = false
			if v == "host" {
				return violateZC1462(cmd)
			}
		}
		if v == "--ipc" {
			prevIpc = true
		}
	}

	return nil
}

func violateZC1462(cmd *ast.SimpleCommand) []Violation {
	return []Violation{{
		KataID: "ZC1462",
		Message: "`--ipc=host` shares host shared memory and SysV IPC with the container — " +
			"trivial data theft and side-channel vector. Use the default private IPC.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1463",
		Title:    "Avoid `docker run --userns=host` — disables user-namespace remapping",
		Severity: SeverityWarning,
		Description: "`--userns=host` turns off the user-namespace remap, meaning UID 0 in the " +
			"container maps to UID 0 on the host. Combined with any of the `--cap-add`, " +
			"`--privileged`, or bind-mount footguns, this becomes a direct host-root escalation. " +
			"Leave the default (container-side remap) enabled.",
		Check: checkZC1463,
	})
}

func checkZC1463(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "docker" && ident.Value != "podman" {
		return nil
	}

	var prevNs bool
	for _, arg := range cmd.Arguments {
		v := arg.String()

		if v == "--userns=host" {
			return violateZC1463(cmd)
		}
		if prevNs {
			prevNs = false
			if v == "host" {
				return violateZC1463(cmd)
			}
		}
		if v == "--userns" {
			prevNs = true
		}
	}

	return nil
}

func violateZC1463(cmd *ast.SimpleCommand) []Violation {
	return []Violation{{
		KataID: "ZC1463",
		Message: "`--userns=host` disables user-namespace remap — UID 0 in the container == " +
			"UID 0 on the host. Leave the default remap on.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1464",
		Title:    "Warn on `iptables -F` / `-P INPUT ACCEPT` — flushes or opens the host firewall",
		Severity: SeverityWarning,
		Description: "Flushing all rules (`-F`) or setting the default INPUT/FORWARD policy to " +
			"ACCEPT leaves the host with no network filter. This is rarely correct outside a " +
			"first-boot provisioning script, and is a frequent post-compromise persistence step. " +
			"Use `iptables-save`/`iptables-restore` for atomic reloads and keep a default-drop " +
			"policy on all hook chains.",
		Check: checkZC1464,
	})
}

var (
	zc1464FirewallNames = map[string]struct{}{"iptables": {}, "ip6tables": {}, "nft": {}}
	zc1464FlushFlags    = map[string]struct{}{"-F": {}, "--flush": {}}
	zc1464PolicyFlags   = map[string]struct{}{"-P": {}, "--policy": {}}
	zc1464OpenChains    = map[string]struct{}{"INPUT": {}, "FORWARD": {}}
)

func checkZC1464(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	if _, hit := zc1464FirewallNames[CommandIdentifier(cmd)]; !hit {
		return nil
	}
	args := zc1464StringArgs(cmd)
	if zc1464FlushHit(args) {
		return violateZC1464(cmd, "flushing all firewall rules")
	}
	if chain := zc1464AcceptPolicy(args); chain != "" {
		return violateZC1464(cmd, "default-ACCEPT policy on "+chain+" chain")
	}
	return nil
}

func zc1464StringArgs(cmd *ast.SimpleCommand) []string {
	out := make([]string, 0, len(cmd.Arguments))
	for _, arg := range cmd.Arguments {
		out = append(out, arg.String())
	}
	return out
}

func zc1464FlushHit(args []string) bool {
	for _, a := range args {
		if _, hit := zc1464FlushFlags[a]; hit {
			return true
		}
	}
	return false
}

func zc1464AcceptPolicy(args []string) string {
	for i, a := range args {
		if _, hit := zc1464PolicyFlags[a]; !hit {
			continue
		}
		if i+2 >= len(args) || args[i+2] != "ACCEPT" {
			continue
		}
		if _, hit := zc1464OpenChains[args[i+1]]; hit {
			return args[i+1]
		}
	}
	return ""
}

func violateZC1464(cmd *ast.SimpleCommand, what string) []Violation {
	return []Violation{{
		KataID:  "ZC1464",
		Message: "Firewall hardening weakened (" + what + "). Keep default-drop and use atomic reload.",
		Line:    cmd.Token.Line,
		Column:  cmd.Token.Column,
		Level:   SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1465",
		Title:    "Warn on `setenforce 0` — disables SELinux enforcement",
		Severity: SeverityWarning,
		Description: "`setenforce 0` switches SELinux to permissive mode, silencing every policy " +
			"decision into an audit log line instead of a deny. It is the textbook post-" +
			"compromise persistence step and also a common \"fix\" that papers over an actual " +
			"policy bug. Address the specific AVC with `audit2allow` instead, and leave " +
			"`setenforce 1` (enforcing) in production.",
		Check: checkZC1465,
	})
}

func checkZC1465(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "setenforce" {
		return nil
	}

	if len(cmd.Arguments) == 0 {
		return nil
	}
	v := cmd.Arguments[0].String()
	if v == "0" || v == "Permissive" || v == "permissive" {
		return []Violation{{
			KataID: "ZC1465",
			Message: "`setenforce 0` disables SELinux enforcement host-wide. Fix the AVC with " +
				"`audit2allow` instead and keep enforcing mode on.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityWarning,
		}}
	}
	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1466",
		Title:    "Warn on disabling the host firewall (`ufw disable` / `systemctl stop firewalld`)",
		Severity: SeverityWarning,
		Description: "Disabling the host firewall leaves every listening port reachable from " +
			"every network the host is on. This is a common \"just make it work\" shortcut that " +
			"has shipped to production more than once. Keep the firewall running and open the " +
			"specific port with `ufw allow <port>` / `firewall-cmd --add-port=<port>/tcp`.",
		Check: checkZC1466,
	})
}

func checkZC1466(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	if ident.Value == "ufw" && len(cmd.Arguments) >= 1 {
		if cmd.Arguments[0].String() == "disable" {
			return violateZC1466(cmd, "ufw disable")
		}
	}

	if ident.Value == "systemctl" && len(cmd.Arguments) >= 2 {
		verb := cmd.Arguments[0].String()
		unit := cmd.Arguments[1].String()
		if (verb == "stop" || verb == "disable" || verb == "mask") &&
			(unit == "firewalld" || unit == "firewalld.service" ||
				unit == "ufw" || unit == "ufw.service" ||
				unit == "nftables" || unit == "nftables.service" ||
				unit == "iptables" || unit == "iptables.service") {
			return violateZC1466(cmd, "systemctl "+verb+" "+unit)
		}
	}

	return nil
}

func violateZC1466(cmd *ast.SimpleCommand, what string) []Violation {
	return []Violation{{
		KataID:  "ZC1466",
		Message: "Host firewall disabled (" + what + "). Keep it on and open specific ports.",
		Line:    cmd.Token.Line,
		Column:  cmd.Token.Column,
		Level:   SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1467",
		Title:    "Warn on `sysctl -w kernel.core_pattern=|...` / `kernel.modprobe=...` (kernel hijack)",
		Severity: SeverityError,
		Description: "Writing `kernel.core_pattern` to a pipe handler or `kernel.modprobe` to a " +
			"user-writable path is a textbook privilege-escalation trick: the next crashing " +
			"setuid process (or the next auto-load of an absent module) executes the supplied " +
			"binary as root. Keep `core_pattern` set to `core` or `systemd-coredump` and leave " +
			"`kernel.modprobe` at the distro default (`/sbin/modprobe`).",
		Check: checkZC1467,
	})
}

func checkZC1467(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "sysctl" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := stripOuterQuotes(arg.String())
		// Accept both `key=value` and `-w key=value` — `-w` shows up as its own arg.
		k, val, ok := strings.Cut(v, "=")
		if !ok {
			continue
		}
		k = strings.TrimSpace(k)
		val = stripOuterQuotes(val)
		if k == "kernel.core_pattern" && strings.HasPrefix(val, "|") {
			return zc1467Violation(cmd, "kernel.core_pattern pipe handler")
		}
		if k == "kernel.modprobe" && val != "" && val != "/sbin/modprobe" {
			return zc1467Violation(cmd, "kernel.modprobe override")
		}
	}
	return nil
}

func stripOuterQuotes(s string) string {
	if len(s) >= 2 {
		first, last := s[0], s[len(s)-1]
		if (first == '\'' && last == '\'') || (first == '"' && last == '"') {
			return s[1 : len(s)-1]
		}
	}
	return s
}

func zc1467Violation(cmd *ast.SimpleCommand, what string) []Violation {
	return []Violation{{
		KataID: "ZC1467",
		Message: "Kernel hijack vector (" + what + ") — next crash / module load runs " +
			"attacker-supplied binary as root.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1468",
		Title:    "Error on apt `--allow-unauthenticated` / `--force-yes` — installs unsigned packages",
		Severity: SeverityError,
		Description: "`--allow-unauthenticated` and the deprecated `--force-yes` disable APT's " +
			"package-signature verification, turning any MITM or typo-squat into arbitrary " +
			"code execution as root. Always sign internal packages and leave verification on.",
		Check: checkZC1468,
	})
}

func checkZC1468(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "apt" && ident.Value != "apt-get" && ident.Value != "aptitude" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "--allow-unauthenticated" ||
			v == "--force-yes" ||
			v == "--allow-downgrades" ||
			v == "--allow-remove-essential" ||
			v == "--allow-change-held-packages" {
			return []Violation{{
				KataID: "ZC1468",
				Message: "APT installing unsigned or override-policy packages (" + v + ") — " +
					"disables signature verification, MITM-to-root trivial.",
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
		ID:       "ZC1469",
		Title:    "Error on `dnf/yum --nogpgcheck` or `rpm --nosignature` (unsigned RPM install)",
		Severity: SeverityError,
		Description: "`--nogpgcheck` / `--nosignature` / `--nodigest` disable RPM package " +
			"signature and digest verification. This turns every mirror, cache, or MITM into a " +
			"direct root compromise. Always keep GPG/signature checking on; sign internal repos " +
			"with your own key.",
		Check: checkZC1469,
	})
}

func checkZC1469(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	switch ident.Value {
	case "dnf", "yum", "microdnf", "zypper":
		for _, arg := range cmd.Arguments {
			v := arg.String()
			if v == "--nogpgcheck" || v == "--no-gpg-checks" {
				return zc1469Violation(cmd, ident.Value+" "+v)
			}
		}
	case "rpm", "rpmbuild":
		for _, arg := range cmd.Arguments {
			v := arg.String()
			if v == "--nosignature" || v == "--nodigest" || v == "--nofiledigest" ||
				v == "--noverify" || v == "--nochecksum" {
				return zc1469Violation(cmd, ident.Value+" "+v)
			}
		}
	}
	return nil
}

func zc1469Violation(cmd *ast.SimpleCommand, what string) []Violation {
	return []Violation{{
		KataID: "ZC1469",
		Message: "Package signature verification disabled (" + what + ") — any mirror / MITM " +
			"becomes immediate root.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1470",
		Title:    "Error on `git config http.sslVerify false` / `git -c http.sslVerify=false`",
		Severity: SeverityError,
		Description: "Disabling `http.sslVerify` in git means every subsequent fetch / clone " +
			"accepts any TLS certificate — MITM trivially replaces the tree you are cloning with " +
			"attacker-controlled code. Fix the broken CA instead: install the certificate, " +
			"point at the right store with `GIT_SSL_CAINFO`, or use an SSH transport.",
		Check: checkZC1470,
	})
}

func checkZC1470(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "git" {
		return nil
	}

	args := make([]string, 0, len(cmd.Arguments))
	for _, a := range cmd.Arguments {
		args = append(args, a.String())
	}

	// Form A: `git config [--global|--local|--system] http.sslVerify false`
	for i := 0; i < len(args); i++ {
		if args[i] == "config" {
			// Walk following args, skip scope flag if present.
			j := i + 1
			for j < len(args) && strings.HasPrefix(args[j], "--") && args[j] != "--" {
				j++
			}
			if j+1 < len(args) && strings.EqualFold(args[j], "http.sslVerify") {
				if strings.EqualFold(args[j+1], "false") || args[j+1] == "0" {
					return zc1470Violation(cmd)
				}
			}
		}
	}

	// Form B: `git -c http.sslVerify=false ...`
	for i := 0; i < len(args); i++ {
		if args[i] == "-c" && i+1 < len(args) {
			kv := args[i+1]
			if k, v, ok := strings.Cut(kv, "="); ok {
				if strings.EqualFold(k, "http.sslVerify") &&
					(strings.EqualFold(v, "false") || v == "0") {
					return zc1470Violation(cmd)
				}
			}
		}
	}

	return nil
}

func zc1470Violation(cmd *ast.SimpleCommand) []Violation {
	return []Violation{{
		KataID: "ZC1470",
		Message: "`http.sslVerify=false` disables TLS verification — any MITM swaps the " +
			"clone for attacker code. Fix the CA instead.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1471",
		Title:    "Error on `kubectl/helm --insecure-skip-tls-verify` (cluster MITM)",
		Severity: SeverityError,
		Description: "`--insecure-skip-tls-verify` tells kubectl / helm to accept any certificate " +
			"from the API server. Against a production cluster, this hands every secret and " +
			"admission payload to a MITM. Fix the trust chain: point `--certificate-authority` " +
			"at the right CA bundle, or restore `KUBECONFIG` with the cluster's embedded CA.",
		Check: checkZC1471,
	})
}

func checkZC1471(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "kubectl" && ident.Value != "helm" && ident.Value != "oc" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "--insecure-skip-tls-verify" ||
			v == "--insecure-skip-tls-verify=true" ||
			v == "--kube-insecure-skip-tls-verify" {
			return []Violation{{
				KataID: "ZC1471",
				Message: "`--insecure-skip-tls-verify` turns off API-server certificate " +
					"verification — MITM steals every secret. Fix the CA bundle instead.",
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
		ID:       "ZC1472",
		Title:    "Error on `aws s3 --acl public-read` / `public-read-write` (public bucket)",
		Severity: SeverityError,
		Description: "Using the `public-read` or `public-read-write` canned ACL when uploading, " +
			"syncing, or setting a bucket policy makes the object (and often the bucket) readable " +
			"by anyone on the internet. Prefer bucket policies scoped to specific principals, or " +
			"CloudFront with Origin Access Identity if you truly need public read.",
		Check: checkZC1472,
	})
}

func checkZC1472(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "aws" {
		return nil
	}

	// Must see `s3` or `s3api` service argument anywhere before `--acl`.
	var sawService bool
	var prevAcl bool
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "s3" || v == "s3api" {
			sawService = true
		}
		if !sawService {
			continue
		}
		if prevAcl {
			prevAcl = false
			if v == "public-read" || v == "public-read-write" {
				return zc1472Violation(cmd, v)
			}
		}
		if v == "--acl" {
			prevAcl = true
		}
		if v == "--acl=public-read" || v == "--acl=public-read-write" {
			return zc1472Violation(cmd, v[len("--acl="):])
		}
	}
	return nil
}

func zc1472Violation(cmd *ast.SimpleCommand, acl string) []Violation {
	return []Violation{{
		KataID:  "ZC1472",
		Message: "Canned ACL `" + acl + "` makes the object (often the bucket) world-readable. Use a scoped bucket policy instead.",
		Line:    cmd.Token.Line,
		Column:  cmd.Token.Column,
		Level:   SeverityError,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1473",
		Title:    "Warn on `openssl req ... -nodes` / `genrsa` without passphrase — unencrypted private key",
		Severity: SeverityWarning,
		Description: "`-nodes` tells OpenSSL not to encrypt the private key that is written to " +
			"disk. The file ends up at whatever filesystem permissions the umask dictates, and " +
			"any subsequent backup / container image / rsync picks up a usable key with no " +
			"passphrase. Use `-aes256` / `-aes-256-cbc` and keep the passphrase in a secrets " +
			"store, or rely on a hardware-backed key via PKCS#11 / TPM.",
		Check: checkZC1473,
	})
}

func checkZC1473(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "openssl" {
		return nil
	}

	if len(cmd.Arguments) == 0 {
		return nil
	}
	sub := cmd.Arguments[0].String()
	// Only flag on subcommands that actually produce a private key file.
	if sub != "req" && sub != "genrsa" && sub != "genpkey" && sub != "ecparam" &&
		sub != "dsaparam" && sub != "pkcs12" {
		return nil
	}

	for _, arg := range cmd.Arguments[1:] {
		v := arg.String()
		if v == "-nodes" || v == "-noenc" {
			return []Violation{{
				KataID: "ZC1473",
				Message: "`" + v + "` writes the private key to disk unencrypted. Use `-aes256` " +
					"(or an HSM/TPM) and keep the passphrase in a secrets store.",
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
		ID:       "ZC1474",
		Title:    "Warn on `ssh-keygen -N \"\"` — generates passwordless SSH key",
		Severity: SeverityWarning,
		Description: "Generating an SSH key with an empty passphrase (`-N \"\"`) leaves the key " +
			"usable by anything that can read the file. Combined with a weak umask or a backup " +
			"that follows the file, this is a common lateral-movement vector. Use a real " +
			"passphrase, or delegate key storage to `ssh-agent` / a hardware token.",
		Check: checkZC1474,
	})
}

func checkZC1474(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "ssh-keygen" {
		return nil
	}

	var prevN bool
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if prevN {
			prevN = false
			if v == `""` || v == `''` || v == "" {
				return []Violation{{
					KataID:  "ZC1474",
					Message: "`ssh-keygen -N \"\"` generates a passwordless key — anything that reads the file can use it. Use a passphrase or ssh-agent/HSM.",
					Line:    cmd.Token.Line,
					Column:  cmd.Token.Column,
					Level:   SeverityWarning,
				}}
			}
		}
		if v == "-N" {
			prevN = true
		}
	}
	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1475",
		Title:    "Warn on `setcap` granting dangerous capabilities to a binary (privesc)",
		Severity: SeverityWarning,
		Description: "Adding CAP_SYS_ADMIN, CAP_DAC_OVERRIDE, CAP_DAC_READ_SEARCH, CAP_SYS_PTRACE, " +
			"or CAP_SETUID to a binary lets any user who can execute it perform operations " +
			"roughly equivalent to root — read any file, change any UID, attach ptrace to root " +
			"processes. Scope the capability as narrowly as possible (e.g. CAP_NET_BIND_SERVICE) " +
			"or run the binary under a dedicated service user with a systemd unit.",
		Check: checkZC1475,
	})
}

var setcapDangerous = []string{
	"cap_sys_admin",
	"cap_dac_override",
	"cap_dac_read_search",
	"cap_sys_ptrace",
	"cap_sys_module",
	"cap_setuid",
	"cap_setgid",
	"cap_chown",
}

func checkZC1475(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "setcap" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := strings.ToLower(stripOuterQuotes(arg.String()))
		for _, cap := range setcapDangerous {
			if strings.Contains(v, cap) {
				return []Violation{{
					KataID: "ZC1475",
					Message: "`setcap` granting dangerous capability `" + cap + "` makes the " +
						"binary a privesc vector for any executing user. Scope narrower or use a " +
						"dedicated service user.",
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
		ID:       "ZC1476",
		Title:    "Warn on `apt-key add` — deprecated, trusts every repo system-wide",
		Severity: SeverityWarning,
		Description: "`apt-key` was deprecated in APT 2.2 and removed from `apt` 2.5. Keys added " +
			"with `apt-key add` end up in a global keyring that signs every repo on the system, " +
			"so a typo-squatted third-party PPA can ship updates for `apt`, `libc6`, or " +
			"`openssh-server`. Store the key in `/etc/apt/keyrings/<vendor>.gpg` and scope it in " +
			"`signed-by=` on the specific `sources.list.d` entry.",
		Check: checkZC1476,
	})
}

func checkZC1476(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "apt-key" {
		return nil
	}

	if len(cmd.Arguments) == 0 {
		return nil
	}
	sub := cmd.Arguments[0].String()
	if sub == "add" || sub == "adv" {
		return []Violation{{
			KataID: "ZC1476",
			Message: "`apt-key " + sub + "` adds to a global keyring that signs every repo. " +
				"Use `/etc/apt/keyrings/<vendor>.gpg` + `signed-by=` instead.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityWarning,
		}}
	}
	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1477",
		Title:    "Warn on `printf \"$var\"` — variable in format-string position (printf-fmt attack)",
		Severity: SeverityWarning,
		Description: "The first argument to `printf` is a format string. Interpolating a shell " +
			"variable into it means any `%` sequence inside the variable is interpreted as a " +
			"format specifier — at best producing garbage, at worst crashing with " +
			"`%s`-out-of-bounds reads or writing attacker-controlled data with `%n`. Always " +
			"use a literal format string: `printf '%s\\n' \"$var\"`.",
		Check: checkZC1477,
	})
}

func checkZC1477(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "printf" {
		return nil
	}

	if len(cmd.Arguments) == 0 {
		return nil
	}
	first := cmd.Arguments[0].String()
	raw := stripOuterQuotes(first)

	// Single-quoted strings don't interpolate; treat them as safe even if `$` is present.
	if strings.HasPrefix(first, "'") && strings.HasSuffix(first, "'") {
		return nil
	}

	// Look for an unescaped `$` (variable, command substitution, or arithmetic).
	for i := 0; i < len(raw); i++ {
		if raw[i] == '\\' {
			i++
			continue
		}
		if raw[i] == '$' {
			return []Violation{{
				KataID: "ZC1477",
				Message: "`printf` format string contains a variable — `%` inside `$var` is " +
					"reparsed as a format specifier. Use `printf '%s' \"$var\"` instead.",
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
		ID:       "ZC1478",
		Title:    "Avoid `mktemp -u` — returns a name without creating the file (TOCTOU)",
		Severity: SeverityWarning,
		Description: "`mktemp -u` allocates a unique name but does not create the file, leaving " +
			"a classic time-of-check to time-of-use race: a second process (possibly attacker- " +
			"controlled on a multi-user host or shared CI runner) can claim the name before you " +
			"redirect into it. Drop `-u` and operate on the file `mktemp` creates for you, or " +
			"use `mktemp -d` if you need a directory path.",
		Check: checkZC1478,
	})
}

func checkZC1478(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "mktemp" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-u" || v == "--dry-run" {
			return []Violation{{
				KataID: "ZC1478",
				Message: "`mktemp -u` returns a unique name but does not create the file — " +
					"TOCTOU race. Let `mktemp` create the file (or use `-d` for a directory).",
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
		ID:       "ZC1479",
		Title:    "Error on `ssh/scp -o StrictHostKeyChecking=no` / `UserKnownHostsFile=/dev/null`",
		Severity: SeverityError,
		Description: "Setting `StrictHostKeyChecking=no` or pointing `UserKnownHostsFile` at " +
			"`/dev/null` makes the client accept any server key on the first (and every) " +
			"connection, stripping the protection against MITM that SSH is designed to provide. " +
			"For ephemeral CI targets, pin the host key in `known_hosts` with `ssh-keyscan` and " +
			"verify the fingerprint out of band, or use `StrictHostKeyChecking=accept-new` at " +
			"most.",
		Check: checkZC1479,
	})
}

func checkZC1479(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "ssh" && ident.Value != "scp" && ident.Value != "sftp" {
		return nil
	}

	check := func(spec string) []Violation {
		s := strings.TrimSpace(strings.ToLower(spec))
		if s == "stricthostkeychecking=no" {
			return zc1479Violation(cmd, "StrictHostKeyChecking=no")
		}
		if s == "userknownhostsfile=/dev/null" {
			return zc1479Violation(cmd, "UserKnownHostsFile=/dev/null")
		}
		return nil
	}

	var prevO bool
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if prevO {
			prevO = false
			if res := check(v); res != nil {
				return res
			}
		}
		if v == "-o" {
			prevO = true
			continue
		}
		if strings.HasPrefix(v, "-o") {
			if res := check(v[2:]); res != nil {
				return res
			}
		}
	}
	return nil
}

func zc1479Violation(cmd *ast.SimpleCommand, what string) []Violation {
	return []Violation{{
		KataID:  "ZC1479",
		Message: "`" + what + "` disables SSH host-key verification — first MITM owns the connection. Pin the fingerprint in known_hosts instead.",
		Line:    cmd.Token.Line,
		Column:  cmd.Token.Column,
		Level:   SeverityError,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1480",
		Title:    "Warn on `terraform apply -auto-approve` / `destroy -auto-approve` in scripts",
		Severity: SeverityWarning,
		Description: "Running `terraform apply -auto-approve` or `destroy -auto-approve` from a " +
			"shell script skips the plan-review step that exists to catch schema drift, " +
			"accidental `-replace`, and resources being deleted. Fine for throwaway CI against " +
			"a PR environment, but dangerous against shared state. Prefer running `plan` + " +
			"`apply` with an out-file and human approval, or scope the auto-apply to specific " +
			"branches/environments.",
		Check: checkZC1480,
	})
}

func checkZC1480(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "terraform" && ident.Value != "terragrunt" && ident.Value != "tofu" {
		return nil
	}

	if len(cmd.Arguments) == 0 {
		return nil
	}
	sub := cmd.Arguments[0].String()
	if sub != "apply" && sub != "destroy" {
		return nil
	}

	for _, arg := range cmd.Arguments[1:] {
		v := arg.String()
		if v == "-auto-approve" || v == "--auto-approve" ||
			v == "-auto-approve=true" || v == "--auto-approve=true" {
			return []Violation{{
				KataID: "ZC1480",
				Message: "`" + ident.Value + " " + sub + " " + v + "` skips plan review. " +
					"Gate behind a branch/env check or require manual approval.",
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
		ID:       "ZC1481",
		Title:    "Warn on `unset HISTFILE` / `export HISTFILE=/dev/null` — disables shell history",
		Severity: SeverityWarning,
		Description: "Disabling shell history (`unset HISTFILE`, `HISTFILE=/dev/null`, " +
			"`HISTSIZE=0`) is a classic stepping stone for hiding post-compromise activity. " +
			"Legitimate scripts almost never need this — if you are pasting a secret on the " +
			"command line, use `HISTCONTROL=ignorespace` and prefix the line with a space, or " +
			"read the value from a file / stdin.",
		Check: checkZC1481,
	})
}

var (
	zc1481UnsetVars  = map[string]struct{}{"HISTFILE": {}, "HISTSIZE": {}, "SAVEHIST": {}, "HISTCMD": {}}
	zc1481EmptyHist  = map[string]struct{}{"": {}, "/dev/null": {}, "''": {}, `""`: {}}
	zc1481ZeroAssign = map[string]struct{}{"HISTSIZE=0": {}, "SAVEHIST=0": {}}
)

func checkZC1481(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	switch CommandIdentifier(cmd) {
	case "unset":
		if hit := zc1481UnsetHit(cmd); hit != "" {
			return zc1481Violation(cmd, "unset "+hit)
		}
	case "export", "typeset":
		if hit := zc1481AssignHit(cmd); hit != "" {
			return zc1481Violation(cmd, hit)
		}
	}
	return nil
}

func zc1481UnsetHit(cmd *ast.SimpleCommand) string {
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if _, hit := zc1481UnsetVars[v]; hit {
			return v
		}
	}
	return ""
}

func zc1481AssignHit(cmd *ast.SimpleCommand) string {
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if strings.HasPrefix(v, "HISTFILE=") {
			val := strings.TrimPrefix(v, "HISTFILE=")
			if _, hit := zc1481EmptyHist[val]; hit {
				return v
			}
		}
		if _, hit := zc1481ZeroAssign[v]; hit {
			return v
		}
	}
	return ""
}

func zc1481Violation(cmd *ast.SimpleCommand, what string) []Violation {
	return []Violation{{
		KataID: "ZC1481",
		Message: "`" + what + "` disables shell history — textbook post-compromise tactic. " +
			"Legitimate alternative: `HISTCONTROL=ignorespace` plus leading-space prefix.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1482",
		Title:    "Error on `docker login -p` / `--password=` — credential in process list",
		Severity: SeverityError,
		Description: "Passing the registry password on the command line puts it in the output of " +
			"`ps`, `/proc/<pid>/cmdline`, and the shell history. On a shared CI runner or a host " +
			"with unprivileged users, that is an immediate leak. Use `--password-stdin` and " +
			"pipe the secret in from `cat /run/secrets/foo` or a credential helper.",
		Check: checkZC1482,
	})
}

func checkZC1482(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "docker" && ident.Value != "podman" && ident.Value != "buildah" &&
		ident.Value != "skopeo" && ident.Value != "helm" {
		return nil
	}

	// Must see `login` subcommand anywhere.
	var sawLogin bool
	var prevP bool
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "login" {
			sawLogin = true
			continue
		}
		if !sawLogin {
			continue
		}
		if prevP {
			return zc1482Violation(cmd, "-p "+v)
		}
		if v == "-p" {
			prevP = true
			continue
		}
		if strings.HasPrefix(v, "--password=") {
			return zc1482Violation(cmd, v)
		}
	}
	return nil
}

func zc1482Violation(cmd *ast.SimpleCommand, what string) []Violation {
	return []Violation{{
		KataID: "ZC1482",
		Message: "`" + what + "` puts the password in ps / /proc / history. Use " +
			"`--password-stdin` piped from a secrets file or credential helper.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1483",
		Title:    "Warn on `pip install --break-system-packages` — bypasses PEP 668 externally-managed guard",
		Severity: SeverityWarning,
		Description: "`--break-system-packages` tells pip to ignore the distro's PEP 668 marker " +
			"and install into `/usr/lib/python*`, overwriting files the package manager owns. " +
			"The next `apt`/`dnf` upgrade clobbers or gets clobbered by the pip-installed " +
			"version, and you now have two sources of truth for Python dependencies. Install " +
			"into a virtualenv (`python -m venv`), use `pipx` for application scripts, or use " +
			"`uv` / `poetry` for project dependencies.",
		Check: checkZC1483,
	})
}

func checkZC1483(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "pip" && ident.Value != "pip3" && ident.Value != "pipx" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "--break-system-packages" {
			return []Violation{{
				KataID: "ZC1483",
				Message: "`--break-system-packages` installs into distro-managed paths and " +
					"collides with apt/dnf. Use a venv, pipx, or uv/poetry instead.",
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
		ID:       "ZC1484",
		Title:    "Error on `npm/yarn/pnpm config set strict-ssl false` — disables registry TLS verification",
		Severity: SeverityError,
		Description: "Turning off `strict-ssl` for npm, yarn, or pnpm makes the client accept any " +
			"TLS certificate from the registry — a MITM (corporate proxy, compromised WiFi, rogue " +
			"BGP) can substitute any package, including new versions of `react` or `lodash`. If " +
			"the registry uses a private CA, point `cafile` / `NODE_EXTRA_CA_CERTS` at the right " +
			"bundle instead.",
		Check: checkZC1484,
	})
}

func checkZC1484(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "npm" && ident.Value != "yarn" && ident.Value != "pnpm" &&
		ident.Value != "bun" {
		return nil
	}

	args := make([]string, 0, len(cmd.Arguments))
	for _, a := range cmd.Arguments {
		args = append(args, a.String())
	}

	// Form A: `npm config set strict-ssl false`
	for i := 0; i+2 < len(args); i++ {
		if args[i] == "config" && args[i+1] == "set" {
			j := i + 2
			// Skip optional scope flags like `--global` / `-g`.
			for j < len(args) && strings.HasPrefix(args[j], "-") {
				j++
			}
			if j+1 < len(args) && args[j] == "strict-ssl" {
				val := strings.ToLower(args[j+1])
				if val == "false" || val == "0" || val == "no" {
					return zc1484Violation(cmd)
				}
			}
			if j < len(args) && strings.HasPrefix(strings.ToLower(args[j]), "strict-ssl=") {
				val := strings.ToLower(strings.TrimPrefix(args[j], "strict-ssl="))
				if val == "false" || val == "0" || val == "no" {
					return zc1484Violation(cmd)
				}
			}
		}
	}

	// Form B: `npm install --strict-ssl=false` (one-shot)
	for _, v := range args {
		if strings.EqualFold(v, "--strict-ssl=false") ||
			strings.EqualFold(v, "--no-strict-ssl") {
			return zc1484Violation(cmd)
		}
	}
	return nil
}

func zc1484Violation(cmd *ast.SimpleCommand) []Violation {
	return []Violation{{
		KataID: "ZC1484",
		Message: "`strict-ssl=false` disables npm/yarn/pnpm registry TLS verification — any " +
			"MITM swaps packages. Point `cafile` at the right CA bundle instead.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1485",
		Title:    "Warn on `openssl s_client -ssl3 / -tls1 / -tls1_1` — legacy TLS",
		Severity: SeverityWarning,
		Description: "Forcing SSLv3, TLSv1.0, or TLSv1.1 connects with protocols that have known " +
			"downgrade and bit-flip attacks (POODLE, BEAST). These are disabled by default in " +
			"every maintained OpenSSL build. If the remote only speaks an old protocol, the " +
			"right fix is to update the remote, not downgrade your client.",
		Check: checkZC1485,
	})
}

func checkZC1485(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "openssl" {
		return nil
	}

	if len(cmd.Arguments) == 0 {
		return nil
	}
	sub := cmd.Arguments[0].String()
	if sub != "s_client" && sub != "s_server" {
		return nil
	}

	for _, arg := range cmd.Arguments[1:] {
		v := arg.String()
		if v == "-ssl2" || v == "-ssl3" || v == "-tls1" || v == "-tls1_1" ||
			v == "-no_tls1_2" || v == "-no_tls1_3" {
			return []Violation{{
				KataID: "ZC1485",
				Message: "`openssl " + sub + " " + v + "` forces a legacy / disabled TLS " +
					"version (downgrade-attack surface). Update the remote instead.",
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
		ID:       "ZC1486",
		Title:    "Warn on `curl -2` / `-3` — forces broken SSLv2 / SSLv3",
		Severity: SeverityWarning,
		Description: "`curl -2` (SSLv2) and `-3` (SSLv3) force protocols that are removed from " +
			"every current TLS library. `-2` matches no working server; `-3` leaves you open to " +
			"POODLE. If the remote really needs an old protocol the fix is on the server, not " +
			"the client.",
		Check: checkZC1486,
	})
}

func checkZC1486(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "curl" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-2" || v == "-3" {
			return []Violation{{
				KataID: "ZC1486",
				Message: "`curl " + v + "` forces SSLv2/SSLv3 — removed from modern TLS " +
					"libraries and subject to POODLE. Fix the server instead.",
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
		ID:       "ZC1487",
		Title:    "Warn on `history -c` — clears shell history (and is a Bash-ism under Zsh)",
		Severity: SeverityWarning,
		Description: "`history -c` clears the in-memory history buffer in Bash. It is a standard " +
			"post-compromise anti-forensics step. It is also a Bash-ism: in Zsh, `history` " +
			"takes completely different arguments, so a copy-pasted `history -c` silently no-ops " +
			"and leaves the author thinking history was cleared when it was not. If you really " +
			"need to rotate history in a Zsh script, unset `HISTFILE` before the sensitive " +
			"block or redirect to `/dev/null` explicitly.",
		Check: checkZC1487,
	})
}

func checkZC1487(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "history" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-c" || v == "-d" {
			return []Violation{{
				KataID: "ZC1487",
				Message: "`history " + v + "` is a Bash-ism for clearing history — does " +
					"nothing in Zsh and is a classic post-compromise tactic elsewhere.",
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
		ID:       "ZC1488",
		Title:    "Warn on `ssh -R 0.0.0.0:...` / `*:...` — reverse tunnel bound to all interfaces",
		Severity: SeverityWarning,
		Description: "The default for `ssh -R` binds the remote listener to `localhost`. Pointing " +
			"it at `0.0.0.0` or `*` (or an explicit public IP) exposes the forwarded port to the " +
			"whole network, including anything else that has reached the jump host. For " +
			"persistent ops tunnels, pin the bind address to a specific private interface and " +
			"require `GatewayPorts clientspecified` server-side.",
		Check: checkZC1488,
	})
}

func checkZC1488(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "ssh" && ident.Value != "autossh" {
		return nil
	}

	var prevForward bool
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if prevForward {
			prevForward = false
			if strings.HasPrefix(v, "0.0.0.0:") || strings.HasPrefix(v, "*:") ||
				strings.HasPrefix(v, "::") {
				return []Violation{{
					KataID: "ZC1488",
					Message: "SSH reverse tunnel bound to all interfaces (`" + v + "`) — " +
						"forwarded port reachable from any network. Bind to a specific IP.",
					Line:   cmd.Token.Line,
					Column: cmd.Token.Column,
					Level:  SeverityWarning,
				}}
			}
		}
		if v == "-R" || v == "-L" || v == "-D" {
			prevForward = true
		}
	}
	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1489",
		Title:    "Error on `nc -e` / `ncat -e` — classic reverse-shell invocation",
		Severity: SeverityError,
		Description: "`nc -e <shell>` and `ncat -e <shell>` pipe a shell to a network socket. " +
			"This is the canonical reverse-shell payload. Most distro builds of `nc` have " +
			"`-e` disabled for precisely this reason, so seeing it in a script is either an " +
			"attacker backdoor or a deployment time bomb waiting on a different packaging " +
			"of netcat. If you need a bidirectional pipe, use `socat TCP:... EXEC:...,pty` " +
			"with an explicit authorization check and document the use.",
		Check: checkZC1489,
	})
}

func checkZC1489(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "nc" && ident.Value != "ncat" && ident.Value != "netcat" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-e" || v == "-c" {
			return []Violation{{
				KataID: "ZC1489",
				Message: "`" + ident.Value + " " + v + "` is the classic reverse-shell flag. " +
					"Use socat with explicit PTY + authorization instead.",
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
		ID:       "ZC1490",
		Title:    "Error on `socat ... EXEC:<shell>` / `SYSTEM:<shell>` — socat reverse-shell pattern",
		Severity: SeverityError,
		Description: "The `EXEC:` and `SYSTEM:` socat address types spawn a subprocess connected " +
			"to the other socat endpoint. Paired with `TCP:` or `TCP-LISTEN:`, they form the " +
			"second-most-common reverse/bind shell payload after `nc -e`. Legitimate uses exist " +
			"(test harnesses, pty brokers) but should be gated behind explicit authorization " +
			"and a non-shell command. Scan hits are worth a look.",
		Check: checkZC1490,
	})
}

func checkZC1490(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "socat" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		low := strings.ToLower(v)
		// EXEC:/bin/bash or "EXEC:\"/bin/sh -i\",pty,stderr"
		if strings.Contains(low, "exec:/bin/bash") ||
			strings.Contains(low, "exec:/bin/sh") ||
			strings.Contains(low, "exec:/bin/zsh") ||
			strings.Contains(low, "exec:\"/bin/bash") ||
			strings.Contains(low, "exec:\"/bin/sh") ||
			strings.Contains(low, "system:/bin/bash") ||
			strings.Contains(low, "system:/bin/sh") {
			return []Violation{{
				KataID: "ZC1490",
				Message: "`socat` pointed at a shell via `EXEC:` / `SYSTEM:` — matches the " +
					"classic reverse/bind-shell pattern. Gate behind explicit authorization.",
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
		ID:       "ZC1491",
		Title:    "Warn on `export LD_PRELOAD=...` / `LD_LIBRARY_PATH=...` — library injection",
		Severity: SeverityWarning,
		Description: "Setting `LD_PRELOAD` in a script forces every subsequent dynamically-linked " +
			"command to load the specified shared object first, a classic post-compromise " +
			"privesc and persistence technique. Setting `LD_LIBRARY_PATH` to a writable path is " +
			"a gentler variant of the same class. Legitimate uses exist (perf profiling, " +
			"asan instrumentation) but should be scoped to a single invocation and the path " +
			"pinned to a read-only location.",
		Check: checkZC1491,
	})
}

func checkZC1491(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "export" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if strings.HasPrefix(v, "LD_PRELOAD=") && len(v) > len("LD_PRELOAD=") {
			return zc1491Violation(cmd, "LD_PRELOAD")
		}
		if strings.HasPrefix(v, "LD_LIBRARY_PATH=") && len(v) > len("LD_LIBRARY_PATH=") {
			return zc1491Violation(cmd, "LD_LIBRARY_PATH")
		}
		if strings.HasPrefix(v, "LD_AUDIT=") && len(v) > len("LD_AUDIT=") {
			return zc1491Violation(cmd, "LD_AUDIT")
		}
	}
	return nil
}

func zc1491Violation(cmd *ast.SimpleCommand, varName string) []Violation {
	return []Violation{{
		KataID: "ZC1491",
		Message: "`export " + varName + "=...` forces every subsequent binary to load a custom " +
			"library — classic privesc/persistence. Scope to a single invocation if needed.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1492",
		Title:    "Style: `at` / `batch` for deferred execution — prefer systemd timers for auditability",
		Severity: SeverityStyle,
		Description: "`at` and `batch` schedule one-shot deferred jobs via `atd`. The job payload " +
			"lands in `/var/spool/at*/` with no unit file or dependency graph, which makes it " +
			"harder to review in fleet audits, easier to miss in a compromise triage, and one of " +
			"the less-watched places adversaries stash persistence. Prefer `systemd-run " +
			"--on-calendar=` or a proper `.timer` unit with a corresponding `.service`.",
		Check: checkZC1492,
	})
}

func checkZC1492(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "at" {
		return nil
	}

	// Skip list/remove/query forms: -l, -r, -c, -q
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-l" || v == "-r" || v == "-c" || v == "-q" ||
			v == "--list" || v == "--remove" {
			return nil
		}
	}

	// Any remaining `at <time>` or `at -f <script> <time>` invocation.
	if len(cmd.Arguments) == 0 {
		return nil
	}
	return []Violation{{
		KataID: "ZC1492",
		Message: "`at` schedules via atd with no unit file — harder to audit. Prefer " +
			"`systemd-run --on-calendar=` or a `.timer` unit.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityStyle,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1493",
		Title:    "Warn on `watch -n 0` — zero-interval watch spins CPU",
		Severity: SeverityWarning,
		Description: "`watch -n 0` (or `-n 0.0` / `-n .0`) tells `watch` to re-run the command " +
			"with no delay, which immediately pins a core to 100% and usually saturates the " +
			"terminal emulator too. Pick a realistic interval (`-n 1`, `-n 2`, `-n 0.5`) — or " +
			"if you truly want tight polling, use a dedicated event API (`inotifywait`, " +
			"`systemd.path` unit, `journalctl -f`).",
		Check: checkZC1493,
	})
}

func checkZC1493(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "watch" {
		return nil
	}

	var prevN bool
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if prevN {
			prevN = false
			if v == "0" || v == "0.0" || v == "0.00" || v == ".0" || v == ".00" {
				return zc1493Violation(cmd, v)
			}
		}
		if v == "-n" || v == "--interval" {
			prevN = true
			continue
		}
		if v == "-n0" || v == "-n0.0" || v == "--interval=0" || v == "--interval=0.0" {
			return zc1493Violation(cmd, v)
		}
	}
	return nil
}

func zc1493Violation(cmd *ast.SimpleCommand, what string) []Violation {
	return []Violation{{
		KataID: "ZC1493",
		Message: "`watch -n " + what + "` pins a core at 100% and saturates the terminal. " +
			"Use a realistic interval or an event-driven API (inotifywait / journalctl -f).",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1494",
		Title:    "Warn on `tcpdump -w <file>` without `-Z <user>` — capture file owned by root",
		Severity: SeverityWarning,
		Description: "`tcpdump` needs root (or CAP_NET_RAW) to open the raw socket, but once the " +
			"socket is open it should drop privileges with `-Z <user>` before writing the pcap. " +
			"Without `-Z`, the capture file is owned by root, any bpf filter bug is exercised " +
			"with root privileges, and on a shared host the pcap can land with permissions that " +
			"leak sensitive traffic to other users. Pair `-w` with `-Z tcpdump` (or a dedicated " +
			"capture user).",
		Check: checkZC1494,
	})
}

func checkZC1494(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "tcpdump" {
		return nil
	}

	hasW := false
	hasZ := false
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-w" || v == "--write-file" {
			hasW = true
		}
		if v == "-Z" || v == "--relinquish-privileges" {
			hasZ = true
		}
	}
	if !hasW || hasZ {
		return nil
	}
	return []Violation{{
		KataID: "ZC1494",
		Message: "`tcpdump -w` without `-Z <user>` writes the pcap as root and never drops " +
			"privileges. Add `-Z tcpdump` (or a dedicated capture user).",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1495",
		Title:    "Warn on `ulimit -c unlimited` — enables core dumps from setuid binaries",
		Severity: SeverityWarning,
		Description: "`ulimit -c unlimited` enables unbounded core dumps for the current shell " +
			"and its children. On a system with `fs.suid_dumpable=2` and a world-readable " +
			"coredump directory, a setuid process that segfaults leaks its memory into a file " +
			"any user can read — Dirty COW-class keys, TLS session material, kerberos tickets. " +
			"Leave core dumps at the distro default (usually 0) and use systemd-coredump with " +
			"access controls if you genuinely need post-mortems.",
		Check: checkZC1495,
	})
}

func checkZC1495(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "ulimit" {
		return nil
	}

	var coreFlag bool
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-c" {
			coreFlag = true
			continue
		}
		if coreFlag && v == "unlimited" {
			return []Violation{{
				KataID: "ZC1495",
				Message: "`ulimit -c unlimited` exposes setuid-process memory via core dumps. " +
					"Leave the distro default and use systemd-coredump if you need post-mortems.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
		coreFlag = false
	}
	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1496",
		Title:    "Error on reading `/dev/mem` / `/dev/kmem` / `/dev/port` — leaks physical memory",
		Severity: SeverityError,
		Description: "These device nodes map physical memory, kernel memory, and x86 I/O ports. " +
			"Reading them (with `strings`, `xxd`, `cat`, or `dd`) exposes kernel state, keys, " +
			"and any other live secret on the box. Modern kernels gate `/dev/mem` behind " +
			"`CONFIG_STRICT_DEVMEM` but most distros also carry `CAP_SYS_RAWIO` on installed " +
			"debugging tools, so the protection is fragile. If you really need a memory dump, " +
			"use `kdump` + `crash` on a proper crash-kernel image.",
		Check: checkZC1496,
	})
}

func checkZC1496(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "strings" && ident.Value != "xxd" && ident.Value != "cat" &&
		ident.Value != "dd" && ident.Value != "od" && ident.Value != "hexdump" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "/dev/mem" || v == "/dev/kmem" || v == "/dev/port" ||
			v == "if=/dev/mem" || v == "if=/dev/kmem" {
			return []Violation{{
				KataID: "ZC1496",
				Message: "Reading `" + v + "` leaks kernel / physical memory. Use kdump + " +
					"crash on a crash-kernel image if you need a dump.",
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
		ID:       "ZC1497",
		Title:    "Error on `useradd -u 0` / `usermod -u 0` — creates a second root account",
		Severity: SeverityError,
		Description: "Creating a user with UID 0 makes them a second root — indistinguishable " +
			"from `root` for every access decision, but hiding behind a non-obvious username " +
			"(`backup`, `service`, `svc-updater`). This is a textbook persistence technique. " +
			"If you need privileged but auditable operations, grant sudo rules tied to a " +
			"specific non-0 UID and log via sudo's session plugin.",
		Check: checkZC1497,
	})
}

func checkZC1497(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "useradd" && ident.Value != "usermod" && ident.Value != "adduser" {
		return nil
	}

	var prevU bool
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if prevU {
			prevU = false
			if v == "0" {
				return zc1497Violation(cmd)
			}
		}
		if v == "-u" || v == "--uid" {
			prevU = true
		}
		if v == "-u0" || v == "--uid=0" {
			return zc1497Violation(cmd)
		}
	}
	return nil
}

func zc1497Violation(cmd *ast.SimpleCommand) []Violation {
	return []Violation{{
		KataID: "ZC1497",
		Message: "Creating a user with UID 0 produces a second root account — classic " +
			"persistence technique. Use sudo rules tied to a non-0 UID instead.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1498",
		Title:    "Warn on `mount -o remount,rw /` — makes read-only root filesystem writable",
		Severity: SeverityWarning,
		Description: "Remounting the root filesystem read-write is either an intentional config " +
			"change that belongs in `/etc/fstab` (in which case this script is the wrong place) " +
			"or a post-compromise step for persisting changes on an immutable / verity-backed " +
			"root. On distros that ship with RO root (Fedora Silverblue, Chrome OS, appliance " +
			"images) this also breaks rollback guarantees. Use `systemd-sysext` or " +
			"`ostree admin deploy` for legitimate modifications.",
		Check: checkZC1498,
	})
}

var zc1498Roots = map[string]struct{}{"/": {}, "/root": {}, "/boot": {}}

func checkZC1498(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok || CommandIdentifier(cmd) != "mount" {
		return nil
	}
	args := zc1464StringArgs(cmd)
	hasRemount, hasRW := zc1498RemountFlags(args)
	target := zc1498SystemTarget(args)
	if !hasRemount || !hasRW || target == "" {
		return nil
	}
	return []Violation{{
		KataID: "ZC1498",
		Message: "`mount -o remount,rw " + target + "` makes a read-only system path " +
			"writable — use ostree / systemd-sysext or fix /etc/fstab.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func zc1498RemountFlags(args []string) (hasRemount, hasRW bool) {
	for i, a := range args {
		if a != "-o" || i+1 >= len(args) {
			continue
		}
		for _, o := range strings.Split(args[i+1], ",") {
			switch o {
			case "remount":
				hasRemount = true
			case "rw":
				hasRW = true
			}
		}
	}
	return
}

func zc1498SystemTarget(args []string) string {
	for _, a := range args {
		if _, hit := zc1498Roots[a]; hit && !strings.HasPrefix(a, "-") {
			return a
		}
	}
	return ""
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1499",
		Title:    "Style: `docker pull <image>` / `:latest` — unpinned image tag",
		Severity: SeverityStyle,
		Description: "Pulling without a tag defaults to `:latest`, which is a moving label. That " +
			"breaks CI reproducibility (yesterday's build passed, today's fails for no reason " +
			"the author changed) and reintroduces supply-chain surface every pull. Pin to a " +
			"specific tag for convenience or to an immutable `@sha256:` digest for production.",
		Check: checkZC1499,
	})
}

func checkZC1499(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "docker" && ident.Value != "podman" && ident.Value != "nerdctl" {
		return nil
	}

	if len(cmd.Arguments) < 2 {
		return nil
	}
	sub := cmd.Arguments[0].String()
	if sub != "pull" && sub != "run" {
		return nil
	}

	// Find the image reference — for `pull` it is the first non-flag arg; for
	// `run` it is also typically the first non-flag non-option arg. Bail as
	// soon as we hit an image-looking token.
	for _, arg := range cmd.Arguments[1:] {
		v := arg.String()
		if strings.HasPrefix(v, "-") || strings.Contains(v, "=") {
			continue
		}
		// Digest pin — explicit immutable reference.
		if strings.Contains(v, "@sha256:") || strings.Contains(v, "@sha512:") {
			return nil
		}
		// Has a tag — flag only if tag is `latest`.
		if colon := strings.LastIndex(v, ":"); colon != -1 {
			tag := v[colon+1:]
			// Skip port-looking refs like localhost:5000 — not a tag.
			if strings.Contains(tag, "/") {
				// registry-with-port, no tag present → unpinned
				return zc1499Violation(cmd, v)
			}
			if tag == "latest" {
				return zc1499Violation(cmd, v)
			}
			return nil
		}
		// No colon, no @ — bare image means implicit `:latest`.
		return zc1499Violation(cmd, v)
	}
	return nil
}

func zc1499Violation(cmd *ast.SimpleCommand, ref string) []Violation {
	return []Violation{{
		KataID: "ZC1499",
		Message: "`" + ref + "` is unpinned (implicit `:latest`). Pin to a specific tag or " +
			"an immutable `@sha256:` digest for reproducibility.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityStyle,
	}}
}
