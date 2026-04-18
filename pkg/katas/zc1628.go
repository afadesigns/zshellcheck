package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1628",
		Title:    "Warn on `insmod` / `modprobe -f` — loads modules bypassing blacklist / signature checks",
		Severity: SeverityWarning,
		Description: "`insmod PATH.ko` loads a kernel module from a file, skipping the depmod-" +
			"built dependency graph and the `/etc/modprobe.d/*.conf` blacklist. `modprobe " +
			"-f` instructs modprobe to ignore version-magic and kernel-mismatch checks. " +
			"Either path lets a module enter the kernel that the administrator explicitly " +
			"disabled, or one compiled against a different kernel — crash, privesc, or full " +
			"kernel compromise. Use plain `modprobe MODNAME` so the system's policy and " +
			"signature verification run.",
		Check: checkZC1628,
	})
}

func checkZC1628(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value == "insmod" {
		if len(cmd.Arguments) == 0 {
			return nil
		}
		return []Violation{{
			KataID: "ZC1628",
			Message: "`insmod` loads a kernel module bypassing depmod / blacklist — prefer " +
				"`modprobe MODNAME` so system policy and signature checks apply.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityWarning,
		}}
	}
	if ident.Value == "modprobe" {
		for _, arg := range cmd.Arguments {
			if arg.String() == "-f" {
				return []Violation{{
					KataID: "ZC1628",
					Message: "`modprobe -f` ignores version-magic and kernel-mismatch " +
						"checks — a mismatched module can crash or compromise the kernel. " +
						"Drop the flag and fix the underlying version mismatch.",
					Line:   cmd.Token.Line,
					Column: cmd.Token.Column,
					Level:  SeverityWarning,
				}}
			}
		}
	}
	return nil
}
