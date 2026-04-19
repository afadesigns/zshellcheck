package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

var zc1780FsProtections = map[string]string{
	"fs.protected_symlinks=0":  "symlink follow protection in sticky dirs",
	"fs.protected_hardlinks=0": "hardlink creation protection in sticky dirs",
	"fs.protected_fifos=0":     "FIFO open protection in sticky dirs",
	"fs.protected_regular=0":   "regular-file open protection in sticky dirs",
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1780",
		Title:    "Warn on `sysctl -w fs.protected_symlinks=0|protected_hardlinks=0|…` — TOCTOU guard disabled",
		Severity: SeverityWarning,
		Description: "The `fs.protected_*` sysctls close a classic race: in a sticky directory " +
			"(`/tmp`, `/var/tmp`, `/dev/shm`), a non-owner cannot follow a symlink, create a " +
			"hardlink to a file they don't own, or open a FIFO / regular file they didn't " +
			"create. Those four gates block the shape of attack where a privileged program " +
			"predictably opens a `/tmp/NAME` that an attacker has already placed as a " +
			"symlink to `/etc/shadow`. Setting any of them to `0` re-enables the race across " +
			"the whole host. Leave the defaults (`1` / `2`) in place; if a specific " +
			"application legitimately needs the old behavior, run it in a mount namespace " +
			"where `/tmp` is not sticky-shared.",
		Check: checkZC1780,
	})
}

func checkZC1780(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "sysctl" {
		return nil
	}
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if note, ok := zc1780FsProtections[v]; ok {
			return []Violation{{
				KataID: "ZC1780",
				Message: "`sysctl " + v + "` disables " + note + " — re-opens a TOCTOU " +
					"race in sticky dirs. Leave the default unless you have a specific " +
					"reason; otherwise scope the change to a mount namespace.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}
