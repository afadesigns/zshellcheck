package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1949",
		Title:    "Error on `rmmod -f` / `rmmod --force` — bypasses refcount, can panic the kernel",
		Severity: SeverityError,
		Description: "`rmmod -f` asks the kernel to tear down a module even if its reference count " +
			"is non-zero. Any live `open(\"/dev/…\")`, mounted filesystem, or in-flight network " +
			"device driven by that module becomes a dangling pointer — the kernel oopses or " +
			"outright panics as soon as the next callback fires. The feature is compiled out " +
			"on most distros (`CONFIG_MODULE_FORCE_UNLOAD=n`), but when present it is strictly " +
			"a break-glass recovery tool. Stop the holders first (`lsof /dev/FOO`, `umount`, " +
			"`ip link set dev … down`), then use plain `rmmod` or `modprobe -r`.",
		Check: checkZC1949,
	})
}

func checkZC1949(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	// Parser caveat: `rmmod --force MOD` mangles the command name to `force`.
	if ident.Value == "force" {
		return zc1949Hit(cmd, "rmmod --force")
	}
	if ident.Value != "rmmod" {
		return nil
	}
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-f" || v == "--force" {
			return zc1949Hit(cmd, "rmmod "+v)
		}
	}
	return nil
}

func zc1949Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1949",
		Message: "`" + form + "` tears down a module even when its refcount is non-zero — " +
			"in-use drivers dangle, kernel oopses on the next callback. Stop holders first " +
			"(`lsof`/`umount`/`ip link down`), then `rmmod` without `-f`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}
