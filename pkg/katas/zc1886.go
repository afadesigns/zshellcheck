package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1886",
		Title:    "Error on `tee/cp/mv/install/dd` writing system shell-init files — persistent privesc surface",
		Severity: SeverityError,
		Description: "`/etc/profile`, `/etc/bash.bashrc`, `/etc/zshrc`, `/etc/zsh/zshenv`, " +
			"`/etc/environment`, and every drop-in under `/etc/profile.d/` are sourced " +
			"by every interactive shell (and `/etc/zshenv` by every Zsh invocation). A " +
			"script that `tee`s, `cp`s, `mv`s, or `dd`s arbitrary content into any of " +
			"those paths becomes a persistent foothold — the next root login runs the " +
			"injected code. These files belong to the packaging system; hand-edit " +
			"carefully, stage a temp file, validate it with a dry-run login, and move " +
			"it into place with an atomic `install -m 644`.",
		Check: checkZC1886,
	})
}

var zc1886SensitivePaths = []string{
	"/etc/profile",
	"/etc/bash.bashrc",
	"/etc/bashrc",
	"/etc/zshrc",
	"/etc/zshenv",
	"/etc/zsh/zshrc",
	"/etc/zsh/zshenv",
	"/etc/zsh/zprofile",
	"/etc/zsh/zlogin",
	"/etc/zprofile",
	"/etc/zlogin",
	"/etc/environment",
}

func checkZC1886(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	switch ident.Value {
	case "tee", "cp", "mv", "install", "dd":
	default:
		return nil
	}
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if zc1886IsSensitivePath(v) {
			return []Violation{{
				KataID: "ZC1886",
				Message: "`" + ident.Value + " ... " + v + "` writes a shell-init " +
					"file sourced by every interactive shell — persistent " +
					"foothold for the next root login. Stage a temp file, " +
					"validate, and `install -m 644` atomically.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}
	return nil
}

func zc1886IsSensitivePath(v string) bool {
	trimmed := strings.Trim(v, "\"'")
	for _, p := range zc1886SensitivePaths {
		if trimmed == p {
			return true
		}
	}
	return strings.HasPrefix(trimmed, "/etc/profile.d/")
}
