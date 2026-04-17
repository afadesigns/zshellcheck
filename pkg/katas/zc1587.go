package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1587",
		Title:    "Warn on `modprobe -r` / `rmmod` from scripts — unloading active kernel modules",
		Severity: SeverityWarning,
		Description: "Unloading a kernel module that is in use — `nvme` (storage), `nvidia` " +
			"(GPU), `e1000`/`ixgbe` (network), `kvm` (virt) — instantly takes the backing " +
			"subsystem offline. On a remote host the script loses its storage or network " +
			"mid-run. Reserve `modprobe -r` / `rmmod` for console maintenance, and consider " +
			"`systemctl stop <unit>` if you are trying to stop a service.",
		Check: checkZC1587,
	})
}

func checkZC1587(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	if ident.Value == "modprobe" {
		for _, arg := range cmd.Arguments {
			v := arg.String()
			if v == "-r" || v == "--remove" {
				return []Violation{{
					KataID: "ZC1587",
					Message: "`modprobe -r` unloads an in-use module — the backing subsystem " +
						"goes offline. Use `systemctl stop` if you meant to stop a service.",
					Line:   cmd.Token.Line,
					Column: cmd.Token.Column,
					Level:  SeverityWarning,
				}}
			}
		}
	}
	if ident.Value == "rmmod" && len(cmd.Arguments) > 0 {
		return []Violation{{
			KataID: "ZC1587",
			Message: "`rmmod` unloads a kernel module — the backing subsystem goes offline. " +
				"Use `systemctl stop` if you meant to stop a service.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityWarning,
		}}
	}
	return nil
}
