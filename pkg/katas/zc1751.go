package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1751",
		Title:    "Error on `rpm/dnf/yum remove --nodeps` — bypasses dependency check, breaks dependents",
		Severity: SeverityError,
		Description: "`rpm -e --nodeps PKG` (also `dnf remove --nodeps`, `yum remove --nodeps`, " +
			"`zypper remove --force`) removes the package while skipping the dependency " +
			"solver. Anything transitively depending on the target immediately breaks — " +
			"`libc`, `openssl`, `systemd` units, even `dnf` itself can get pulled out, " +
			"leaving the host unbootable or unpackageable. Resolve the dependency conflict " +
			"explicitly (`dnf swap`, `rpm -e --rebuilddb` never, pin the conflicting package) " +
			"instead of bypassing the check.",
		Check: checkZC1751,
	})
}

func checkZC1751(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	var tool string
	var verbOK bool
	switch ident.Value {
	case "rpm":
		tool = "rpm"
		for _, arg := range cmd.Arguments {
			v := arg.String()
			if v == "-e" || v == "--erase" {
				verbOK = true
				break
			}
		}
	case "dnf", "yum", "microdnf":
		tool = ident.Value
		if len(cmd.Arguments) == 0 {
			return nil
		}
		sub := cmd.Arguments[0].String()
		if sub == "remove" || sub == "erase" {
			verbOK = true
		}
	case "zypper":
		tool = "zypper"
		if len(cmd.Arguments) == 0 {
			return nil
		}
		sub := cmd.Arguments[0].String()
		if sub == "remove" || sub == "rm" {
			verbOK = true
		}
	default:
		return nil
	}
	if !verbOK {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "--nodeps" || v == "--no-deps" {
			return []Violation{{
				KataID: "ZC1751",
				Message: "`" + tool + " ... " + v + "` removes the package without the " +
					"dependency solver — dependents break (libc, openssl, systemd " +
					"units). Resolve the conflict explicitly instead of bypassing.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}
	return nil
}
