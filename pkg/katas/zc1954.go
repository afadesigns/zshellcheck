package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1954",
		Title:    "Warn on `setfattr -n security.capability|security.selinux|security.ima` — bypasses `setcap`/`chcon`",
		Severity: SeverityWarning,
		Description: "`setfattr -n security.capability -v …` writes the raw file-capability xattr " +
			"that the kernel consults when a binary `execve()`s, bypassing the `setcap` " +
			"wrapper's validation and audit trail. Similarly, `security.selinux` replaces the " +
			"SELinux label without going through `chcon` / `semanage`, and `security.ima` " +
			"overwrites the IMA hash that integrity-measurement trusts. These attributes are " +
			"the raw kernel knobs behind purpose-built tools; script usage is almost always " +
			"wrong. Use `setcap`, `chcon`/`semanage fcontext`, and `evmctl` instead.",
		Check: checkZC1954,
	})
}

func checkZC1954(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "setfattr" {
		return nil
	}

	for i, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-n" && i+1 < len(cmd.Arguments) {
			name := cmd.Arguments[i+1].String()
			if zc1954SecurityAttr(name) {
				return zc1954Hit(cmd, name)
			}
		}
		if strings.HasPrefix(v, "-n") && len(v) > 2 {
			if zc1954SecurityAttr(v[2:]) {
				return zc1954Hit(cmd, v[2:])
			}
		}
		if strings.HasPrefix(v, "--name=") {
			if zc1954SecurityAttr(strings.TrimPrefix(v, "--name=")) {
				return zc1954Hit(cmd, strings.TrimPrefix(v, "--name="))
			}
		}
	}
	return nil
}

func zc1954SecurityAttr(name string) bool {
	switch {
	case name == "security.capability",
		name == "security.selinux",
		name == "security.ima",
		name == "security.evm":
		return true
	case strings.HasPrefix(name, "security.apparmor"):
		return true
	}
	return false
}

func zc1954Hit(cmd *ast.SimpleCommand, attr string) []Violation {
	return []Violation{{
		KataID: "ZC1954",
		Message: "`setfattr -n " + attr + "` writes the raw kernel xattr — bypasses " +
			"`setcap`/`chcon`/`evmctl` validation and audit. Use the purpose-built tool.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}
