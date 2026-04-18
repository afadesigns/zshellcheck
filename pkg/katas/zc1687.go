package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1687",
		Title:    "Warn on `snap install --classic` / `--devmode` — weakens snap confinement",
		Severity: SeverityWarning,
		Description: "`snap install --classic` drops the AppArmor / cgroup / seccomp sandbox " +
			"entirely — the snap behaves like a normal Debian package with full system " +
			"access. `--devmode` keeps the sandbox wired up but logs violations instead of " +
			"blocking them. Both modes are documented escape hatches for snaps that cannot " +
			"yet fit the strict confinement (IDEs, compilers, some network tooling), but in " +
			"provisioning scripts they usually mean \"I could not be bothered to pick a " +
			"strict snap.\" Find a strict alternative, or install from the distro repository " +
			"with proper AppArmor profiles; if `--classic` is truly required, document the " +
			"specific snap and the interface that needed elevation.",
		Check: checkZC1687,
	})
}

func checkZC1687(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "snap" {
		return nil
	}
	if len(cmd.Arguments) == 0 || cmd.Arguments[0].String() != "install" {
		return nil
	}

	for _, arg := range cmd.Arguments[1:] {
		v := arg.String()
		if v == "--classic" {
			return zc1687Hit(cmd, "--classic", "drops AppArmor / cgroup / seccomp sandbox")
		}
		if v == "--devmode" {
			return zc1687Hit(cmd, "--devmode", "logs confinement violations instead of blocking")
		}
	}
	return nil
}

func zc1687Hit(cmd *ast.SimpleCommand, flag, what string) []Violation {
	return []Violation{{
		KataID: "ZC1687",
		Message: "`snap install " + flag + "` " + what + " — find a strict snap or a " +
			"distro-package alternative, or document why this specific snap needs it.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}
