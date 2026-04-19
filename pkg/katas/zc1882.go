package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1882",
		Title:    "Warn on `sudo -s` / `sudo su` / `sudo bash` — spawns an interactive root shell from a script",
		Severity: SeverityWarning,
		Description: "`sudo -s`, `sudo -i`, `sudo su [-]`, and `sudo bash` (or `zsh`/`sh`/`ksh`) " +
			"with no trailing command hand you an interactive root shell. That is fine " +
			"at a prompt, but in a non-interactive script the shell either hangs " +
			"waiting for stdin or drains stdin into root's shell as if those lines were " +
			"the shell's commands — neither is what the script author meant. Pass the " +
			"actual command to sudo (`sudo /usr/local/bin/provision.sh`) so the " +
			"elevation is scoped and audit logs capture the real work.",
		Check: checkZC1882,
	})
}

func checkZC1882(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "sudo" {
		return nil
	}
	args := cmd.Arguments
	if len(args) == 0 {
		return nil
	}

	first := args[0].String()
	rest := args[1:]

	// sudo -s / sudo -i with no trailing positional command.
	if (first == "-s" || first == "-i") && !zc1882HasPositional(rest) {
		return zc1882Hit(cmd, "sudo "+first)
	}

	// sudo su, sudo su -, sudo su -l, sudo su --login (no -c).
	if first == "su" {
		if !zc1882HasArg(rest, "-c", "--command") {
			return zc1882Hit(cmd, "sudo su")
		}
	}

	// sudo bash / sudo zsh / sudo sh / sudo ksh without -c.
	switch first {
	case "bash", "zsh", "sh", "ksh", "dash", "ash":
		if !zc1882HasArg(rest, "-c") {
			return zc1882Hit(cmd, "sudo "+first)
		}
	}
	return nil
}

func zc1882HasPositional(args []ast.Expression) bool {
	for _, a := range args {
		v := a.String()
		if v == "" || v[0] == '-' {
			continue
		}
		return true
	}
	return false
}

func zc1882HasArg(args []ast.Expression, names ...string) bool {
	for _, a := range args {
		v := a.String()
		for _, n := range names {
			if v == n {
				return true
			}
		}
	}
	return false
}

func zc1882Hit(cmd *ast.SimpleCommand, where string) []Violation {
	return []Violation{{
		KataID: "ZC1882",
		Message: "`" + where + "` spawns an interactive root shell — in a script " +
			"either hangs on stdin or drains the rest of the file into root's " +
			"shell. Pass the command to sudo: `sudo /path/to/cmd arg …`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}
