package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

var zc1924Tools = map[string]bool{
	"virt-cat":       true,
	"virt-copy-out":  true,
	"virt-tar-out":   true,
	"virt-edit":      true,
	"virt-copy-in":   true,
	"virt-tar-in":    true,
	"guestfish":      true,
	"guestmount":     true,
	"virt-customize": true,
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1924",
		Title:    "Warn on `virt-cat` / `virt-copy-out` / `guestfish` / `guestmount` — reads guest disk from host",
		Severity: SeverityWarning,
		Description: "libguestfs tools (`virt-cat`, `virt-copy-out`, `virt-tar-out`, `virt-edit`, " +
			"`virt-customize`, `guestfish`, `guestmount`) open a VM's disk image directly from " +
			"the hypervisor and read or mutate its contents without going through the guest " +
			"OS. That bypasses every in-guest permission, audit, and LUKS keyslot the VM was " +
			"using, and — if the VM is live — risks filesystem corruption because two writers " +
			"are now mounted on the same image. Snapshot the disk first, work on the clone, " +
			"and prefer in-guest `ssh`/`scp`/`ansible` for anything that does not need " +
			"out-of-band recovery.",
		Check: checkZC1924,
	})
}

func checkZC1924(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if !zc1924Tools[ident.Value] {
		return nil
	}
	return []Violation{{
		KataID: "ZC1924",
		Message: "`" + ident.Value + "` reads/writes the VM disk directly from the host — " +
			"bypasses in-guest permissions, audit, and LUKS; a live VM risks corruption " +
			"from double-mount. Snapshot first, work on the clone, prefer in-guest tooling.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}
