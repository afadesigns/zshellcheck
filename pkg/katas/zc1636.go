package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1636",
		Title:    "Warn on `virsh destroy DOMAIN` — force-stops VM (no graceful shutdown)",
		Severity: SeverityWarning,
		Description: "`virsh destroy DOM` is the libvirt equivalent of pulling the plug on a " +
			"running VM. The guest OS gets no chance to flush filesystems, close network " +
			"connections, or run its own shutdown services — data corruption risk on any " +
			"open file in the guest. For graceful shutdown use `virsh shutdown DOM` (ACPI " +
			"event), wait for completion, and only fall back to `destroy` for a genuinely " +
			"unresponsive guest. `virsh destroy --graceful DOM` attempts a timed graceful " +
			"first, then forces — that variant is not flagged.",
		Check: checkZC1636,
	})
}

func checkZC1636(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "virsh" {
		return nil
	}
	if len(cmd.Arguments) == 0 || cmd.Arguments[0].String() != "destroy" {
		return nil
	}

	for _, arg := range cmd.Arguments[1:] {
		if arg.String() == "--graceful" {
			return nil
		}
	}

	return []Violation{{
		KataID: "ZC1636",
		Message: "`virsh destroy` yanks power from the VM — filesystem corruption risk. Use " +
			"`virsh shutdown` for graceful stop, or `virsh destroy --graceful` as a timed " +
			"fallback.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}
