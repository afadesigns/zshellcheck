package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1496",
		Title:    "Error on reading `/dev/mem` / `/dev/kmem` / `/dev/port` — leaks physical memory",
		Severity: SeverityError,
		Description: "These device nodes map physical memory, kernel memory, and x86 I/O ports. " +
			"Reading them (with `strings`, `xxd`, `cat`, or `dd`) exposes kernel state, keys, " +
			"and any other live secret on the box. Modern kernels gate `/dev/mem` behind " +
			"`CONFIG_STRICT_DEVMEM` but most distros also carry `CAP_SYS_RAWIO` on installed " +
			"debugging tools, so the protection is fragile. If you really need a memory dump, " +
			"use `kdump` + `crash` on a proper crash-kernel image.",
		Check: checkZC1496,
	})
}

func checkZC1496(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "strings" && ident.Value != "xxd" && ident.Value != "cat" &&
		ident.Value != "dd" && ident.Value != "od" && ident.Value != "hexdump" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "/dev/mem" || v == "/dev/kmem" || v == "/dev/port" ||
			v == "if=/dev/mem" || v == "if=/dev/kmem" {
			return []Violation{{
				KataID: "ZC1496",
				Message: "Reading `" + v + "` leaks kernel / physical memory. Use kdump + " +
					"crash on a crash-kernel image if you need a dump.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}
	return nil
}
