package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1681",
		Title:    "Error on `tar -P` / `--absolute-names` — archive absolute paths, can overwrite host files",
		Severity: SeverityError,
		Description: "By default GNU tar strips the leading `/` from archive member paths so " +
			"that extraction stays under the current directory. `-P` (or the long form " +
			"`--absolute-names`) disables that strip: `tar -xPf evil.tar` happily writes to " +
			"`/etc/cron.d/evil`, `/usr/local/bin/sshd`, or any other absolute path the " +
			"archive mentions. Archives from untrusted sources should never be unpacked " +
			"with `-P`. Drop the flag, extract with `-C <scratch-dir>`, audit the tree, " +
			"then copy files into place with `install` or `cp`.",
		Check: checkZC1681,
	})
}

func checkZC1681(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	switch ident.Value {
	case "tar", "bsdtar", "gtar":
	case "absolute-names":
		// `tar --absolute-names …` parses with `tar` consumed — the trailing
		// name alone is unambiguous evidence of the flag.
		return zc1681Hit(cmd, "--absolute-names")
	default:
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-P" || v == "--absolute-names" {
			return zc1681Hit(cmd, v)
		}
		if !strings.HasPrefix(v, "-") || strings.HasPrefix(v, "--") {
			continue
		}
		body := strings.TrimPrefix(v, "-")
		if strings.ContainsRune(body, 'P') {
			return zc1681Hit(cmd, v)
		}
	}
	return nil
}

func zc1681Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1681",
		Message: "`tar " + form + "` keeps absolute paths during extraction — an " +
			"untrusted archive can overwrite `/etc/cron.d`, `/usr/local/bin`, etc. Drop " +
			"the flag and extract with `-C <scratch-dir>` instead.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}
