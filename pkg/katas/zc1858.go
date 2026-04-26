// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1858",
		Title:    "Error on `ssh -c 3des-cbc|arcfour|blowfish-cbc` — weak cipher forced on the tunnel",
		Severity: SeverityError,
		Description: "OpenSSH disables legacy ciphers by default; a script that explicitly forces " +
			"one with `-c 3des-cbc`, `-c arcfour`, `-c blowfish-cbc`, or a matching entry " +
			"in `-o Ciphers=...` downgrades the tunnel to an algorithm with known plaintext " +
			"recovery, IV-reuse, or birthday-bound attacks. Typically this is done to reach " +
			"an old appliance — but it drags every other session on the same invocation " +
			"down with it. Leave cipher selection to OpenSSH's default; if a legacy device " +
			"absolutely requires a weak cipher, isolate it in a `Host ...` block in " +
			"`~/.ssh/config` with explicit `HostKeyAlgorithms` and keep the rest of the " +
			"fleet on strong defaults.",
		Check: checkZC1858,
	})
}

var zc1858Weak = []string{
	"3des-cbc",
	"arcfour",
	"arcfour128",
	"arcfour256",
	"blowfish-cbc",
	"cast128-cbc",
	"des-cbc",
	"rijndael-cbc@lysator.liu.se",
}

func checkZC1858(node ast.Node) []Violation {
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
	args := cmd.Arguments
	for i, arg := range args {
		v := arg.String()
		var candidate string
		switch {
		case v == "-c" && i+1 < len(args):
			candidate = args[i+1].String()
		case v == "-o" && i+1 < len(args):
			candidate = zc1858ExtractCiphers(args[i+1].String())
		case strings.HasPrefix(v, "-o"):
			candidate = zc1858ExtractCiphers(strings.TrimPrefix(v, "-o"))
		}
		if candidate == "" {
			continue
		}
		if weak := zc1858FirstWeakCipher(candidate); weak != "" {
			return []Violation{{
				KataID: "ZC1858",
				Message: "`" + ident.Value + " ... " + weak + "` forces a weak cipher " +
					"with known plaintext-recovery / IV-reuse attacks. Leave " +
					"cipher selection to OpenSSH defaults; if a legacy peer needs " +
					"it, scope inside a `Host` block in `~/.ssh/config`.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}
	return nil
}

func zc1858ExtractCiphers(kv string) string {
	kv = strings.TrimSpace(strings.Trim(kv, "\"'"))
	lower := strings.ToLower(kv)
	if !strings.HasPrefix(lower, "ciphers") {
		return ""
	}
	rest := strings.TrimSpace(strings.TrimPrefix(lower, "ciphers"))
	if !strings.HasPrefix(rest, "=") {
		return ""
	}
	return strings.TrimSpace(strings.TrimPrefix(rest, "="))
}

func zc1858FirstWeakCipher(list string) string {
	lower := strings.ToLower(list)
	for _, entry := range strings.Split(lower, ",") {
		entry = strings.TrimSpace(strings.Trim(entry, "+^-"))
		for _, weak := range zc1858Weak {
			if entry == weak {
				return entry
			}
		}
	}
	return ""
}
