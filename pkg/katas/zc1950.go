// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1950",
		Title:    "Error on `tune2fs -O ^has_journal` / `-m 0` — removes journal or root reserve",
		Severity: SeverityError,
		Description: "`tune2fs -O ^has_journal $DEV` strips the ext3/4 journal from the " +
			"filesystem. Crash recovery drops from \"replay the journal\" to \"scan the whole " +
			"block device with `fsck -y`\", which frequently truncates partially-written files. " +
			"`tune2fs -m 0 $DEV` takes the reserved-for-root space down to zero; when the " +
			"filesystem fills up there is no headroom for `journald`, `apt`, or even a root " +
			"shell to clean up — recovery needs rescue media. Keep the journal on and leave " +
			"`-m` at the distro default (5% is overkill on large disks, but `-m 1` is still " +
			"safe).",
		Check: checkZC1950,
	})
}

func checkZC1950(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "tune2fs" {
		return nil
	}

	for i, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-O" && i+1 < len(cmd.Arguments) {
			spec := cmd.Arguments[i+1].String()
			if zc1950StripsJournal(spec) {
				return zc1950Hit(cmd, "-O "+spec, "strips the journal — crash recovery needs a full `fsck -y` and may truncate files")
			}
		}
		if v == "-m" && i+1 < len(cmd.Arguments) {
			if cmd.Arguments[i+1].String() == "0" {
				return zc1950Hit(cmd, "-m 0", "zeroes the root reserve — a full fs leaves no headroom for `journald`/`apt`/root shells")
			}
		}
	}
	return nil
}

func zc1950StripsJournal(spec string) bool {
	for _, part := range strings.Split(spec, ",") {
		part = strings.TrimSpace(part)
		if part == "^has_journal" {
			return true
		}
	}
	return false
}

func zc1950Hit(cmd *ast.SimpleCommand, form, why string) []Violation {
	return []Violation{{
		KataID:  "ZC1950",
		Message: "`tune2fs " + form + "` " + why + ". Keep the default.",
		Line:    cmd.Token.Line,
		Column:  cmd.Token.Column,
		Level:   SeverityError,
	}}
}
