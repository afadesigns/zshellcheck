package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1909",
		Title:    "Warn on `kexec -l` / `-e` — jumps to an alternate kernel, bypasses bootloader and Secure Boot",
		Severity: SeverityWarning,
		Description: "`kexec -l /path/to/vmlinuz …` stages a second kernel image, and `kexec -e` " +
			"(or `kexec -f`) then transfers control to it without going through the firmware, " +
			"GRUB, or shim. On a Secure-Boot system the staged kernel is never verified against " +
			"the enrolled MOK/PK — an attacker who lands a root exec can boot a hostile kernel " +
			"while leaving /boot untouched. Reserve `kexec` for the live-patching / crash-dump " +
			"workflow it was designed for, gate the call behind `sudo` + audit, and prefer " +
			"`systemctl kexec` or a normal reboot when possible.",
		Check: checkZC1909,
	})
}

func checkZC1909(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	if ident.Value == "load" || ident.Value == "exec" || ident.Value == "unload" {
		// Parser caveat: `kexec --load X` mangles to name=`load`.
		return zc1909Hit(cmd, "kexec --"+ident.Value)
	}
	if ident.Value != "kexec" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		switch v {
		case "-l", "-e", "-f", "-p":
			return zc1909Hit(cmd, "kexec "+v)
		case "--load", "--exec", "--force", "--load-panic":
			return zc1909Hit(cmd, "kexec "+v)
		}
	}
	return nil
}

func zc1909Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1909",
		Message: "`" + form + "` stages or jumps to a kernel without firmware / " +
			"bootloader verification — Secure Boot never checks the signature. Gate behind " +
			"`sudo` + audit and prefer `systemctl kexec` or a real reboot.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}
