// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1747",
		Title:    "Error on `npm/yarn/pnpm --registry http://...` — plaintext registry allows MITM",
		Severity: SeverityError,
		Description: "`npm install --registry http://...`, `pnpm --registry http://...`, and " +
			"`yarn config set registry http://...` configure a plaintext HTTP package " +
			"registry. Any network-position attacker (open Wi-Fi, hostile transit, MITM " +
			"proxy) can replace tarball metadata or content in flight; npm install-time " +
			"`postinstall` scripts then execute the swapped code on the build host. Switch " +
			"the registry URL to `https://` (or terminate TLS at the internal mirror) and " +
			"pair it with a lockfile to pin tarball integrity hashes.",
		Check: checkZC1747,
	})
}

func checkZC1747(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	switch ident.Value {
	case "npm", "pnpm":
		return zc1747ScanFlags(cmd, ident.Value)
	case "yarn":
		if hit := zc1747ScanFlags(cmd, "yarn"); hit != nil {
			return hit
		}
		return zc1747YarnConfigSet(cmd)
	}
	return nil
}

func zc1747ScanFlags(cmd *ast.SimpleCommand, tool string) []Violation {
	prevRegistry := false
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if prevRegistry {
			if strings.HasPrefix(v, "http://") {
				return zc1747Hit(cmd, tool, "--registry "+v)
			}
			prevRegistry = false
			continue
		}
		switch {
		case v == "--registry":
			prevRegistry = true
		case strings.HasPrefix(v, "--registry="):
			if strings.HasPrefix(strings.TrimPrefix(v, "--registry="), "http://") {
				return zc1747Hit(cmd, tool, v)
			}
		}
	}
	return nil
}

func zc1747YarnConfigSet(cmd *ast.SimpleCommand) []Violation {
	if len(cmd.Arguments) < 4 {
		return nil
	}
	if cmd.Arguments[0].String() != "config" || cmd.Arguments[1].String() != "set" {
		return nil
	}
	if cmd.Arguments[2].String() != "registry" {
		return nil
	}
	url := cmd.Arguments[3].String()
	if strings.HasPrefix(url, "http://") {
		return zc1747Hit(cmd, "yarn", "config set registry "+url)
	}
	return nil
}

func zc1747Hit(cmd *ast.SimpleCommand, tool, what string) []Violation {
	return []Violation{{
		KataID: "ZC1747",
		Message: "`" + tool + " " + what + "` uses plaintext HTTP for the package registry — " +
			"any MITM swaps tarballs and runs install-time `postinstall` code. Use " +
			"`https://`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}
