package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

var zc1619NetworkFS = map[string]bool{
	"nfs": true, "nfs4": true,
	"cifs": true, "smbfs": true, "smb3": true,
	"sshfs": true,
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1619",
		Title:    "Warn on `mount -t nfs/cifs/smb/sshfs` missing `nosuid` or `nodev`",
		Severity: SeverityWarning,
		Description: "Network filesystems present files whose mode bits are controlled by a " +
			"remote server. Without `nosuid` in the mount options, a compromised or hostile " +
			"server can plant a setuid-root binary on the share; the client kernel honors the " +
			"suid bit and the binary runs as root on the mounting host. Without `nodev`, the " +
			"server can plant device nodes the kernel treats as real. Always mount network " +
			"shares with `nosuid,nodev`; add `noexec` unless the export is intended to hold " +
			"executables.",
		Check: checkZC1619,
	})
}

func checkZC1619(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "mount" {
		return nil
	}

	var fsType, opts string
	for i, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-t" && i+1 < len(cmd.Arguments) {
			fsType = cmd.Arguments[i+1].String()
		}
		if v == "-o" && i+1 < len(cmd.Arguments) {
			opts = cmd.Arguments[i+1].String()
		}
	}

	if !zc1619NetworkFS[fsType] {
		return nil
	}
	if strings.Contains(opts, "nosuid") && strings.Contains(opts, "nodev") {
		return nil
	}

	missing := []string{}
	if !strings.Contains(opts, "nosuid") {
		missing = append(missing, "nosuid")
	}
	if !strings.Contains(opts, "nodev") {
		missing = append(missing, "nodev")
	}
	return []Violation{{
		KataID: "ZC1619",
		Message: "`mount -t " + fsType + "` without " + strings.Join(missing, ",") +
			" — a hostile server can plant setuid binaries or device nodes that the " +
			"client kernel honors. Add `nosuid,nodev` to the `-o` options.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}
