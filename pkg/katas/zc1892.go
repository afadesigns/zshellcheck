// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strconv"
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1892",
		Title:    "Error on `install -m 4755|6755|2755` — sets setuid/setgid bit at install time",
		Severity: SeverityError,
		Description: "`install -m <mode>` with the setuid (`4xxx`), setgid (`2xxx`), or combined " +
			"(`6xxx`) octal prefix creates the target with those special bits set, which " +
			"turns every execution into a privilege-elevation vector. An uninspected " +
			"binary installed this way — especially from a build script or package " +
			"post-install — becomes a persistent local-privesc primitive if the binary " +
			"is writable, has command-injection, or links against attacker-influenced " +
			"libraries. Drop the setuid/setgid bits from the mode (`install -m 0755`) and " +
			"grant the narrow capability the program actually needs with `setcap " +
			"cap_net_bind_service+ep`; audit the remaining setuid binaries with " +
			"`find / -perm -4000`.",
		Check: checkZC1892,
	})
}

func checkZC1892(node ast.Node) []Violation {
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
	args := cmd.Arguments
	for i, arg := range args {
		v := arg.String()
		var mode string
		switch {
		case v == "-m" && i+1 < len(args):
			mode = args[i+1].String()
		case v == "--mode" && i+1 < len(args):
			mode = args[i+1].String()
		case strings.HasPrefix(v, "-m") && v != "-m":
			mode = strings.TrimPrefix(v, "-m")
		case strings.HasPrefix(v, "--mode="):
			mode = strings.TrimPrefix(v, "--mode=")
		default:
			continue
		}
		if zc1892HasSetuidBits(mode) {
			return []Violation{{
				KataID: "ZC1892",
				Message: "`" + ident.Value + " -m " + mode + "` sets setuid/setgid " +
					"bits at install time — every execution becomes a privesc " +
					"vector. Use `0755` and grant narrow caps with `setcap` instead.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}
	return nil
}

func zc1892HasSetuidBits(mode string) bool {
	mode = strings.Trim(mode, "\"'")
	if mode == "" {
		return false
	}
	for _, r := range mode {
		if r < '0' || r > '9' {
			return false
		}
	}
	var n int64
	if strings.ContainsAny(mode, "89") {
		n, _ = strconv.ParseInt(mode, 10, 32)
	} else {
		n, _ = strconv.ParseInt(mode, 8, 32)
	}
	// setuid (0o4000), setgid (0o2000)
	return (n & 0o6000) != 0
}
