// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strconv"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1671",
		Title:    "Error on `install -m 777` / `mkdir -m 777` — creates world-writable target",
		Severity: SeverityError,
		Description: "`install -m MODE` / `mkdir -m MODE` applies MODE atomically at file or " +
			"directory creation, so the world-writable window from a later `chmod 777` is " +
			"not even needed — the path is wide-open from the moment it exists. Any local " +
			"user can swap binaries under `/usr/local/bin`, write shell-completion hooks " +
			"into `/etc/bash_completion.d`, or turn a shared directory into an LPE staging " +
			"ground. Drop the world-write bit: `0755` for binaries, `0644` for files, `2770` " +
			"with `chgrp` for shared directories.",
		Check: checkZC1671,
	})
}

func checkZC1671(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "install" && ident.Value != "mkdir" {
		return nil
	}

	for i, arg := range cmd.Arguments {
		if arg.String() != "-m" {
			continue
		}
		if i+1 >= len(cmd.Arguments) {
			continue
		}
		mode := cmd.Arguments[i+1].String()
		if !zc1671WorldWritable(mode) {
			continue
		}
		return []Violation{{
			KataID: "ZC1671",
			Message: "`" + ident.Value + " -m " + mode + "` creates a world-writable " +
				"target — drop the world-write bit (e.g. `0755` / `0644`).",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityError,
		}}
	}
	return nil
}

// zc1671WorldWritable returns true if MODE has the world-write (o+w) bit set.
// Users spell modes in octal. If the literal parses as octal, trust that
// reading. Otherwise (a digit 8/9 appears — that only happens because the
// parser normalized a leading-zero octal like `0666` to decimal `438`), parse
// as decimal and still check the o+w bit.
func zc1671WorldWritable(mode string) bool {
	if n, err := strconv.ParseInt(mode, 8, 32); err == nil {
		return n > 0 && n&0o002 != 0
	}
	if n, err := strconv.ParseInt(mode, 10, 32); err == nil && n > 0 && n&0o002 != 0 {
		return true
	}
	return false
}
