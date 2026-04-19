package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1826",
		Title:    "Warn on `install -m 4xxx/2xxx/6xxx` — drops a setuid / setgid binary in one step",
		Severity: SeverityWarning,
		Description: "`install -m MODE SRC DEST` applies MODE atomically at copy time. A four-digit " +
			"mode whose leading digit is `4` (setuid), `2` (setgid), or `6` (both) places a " +
			"setuid / setgid binary into the destination path in a single operation — no " +
			"intermediate `chmod` step where a privilege-tripwire would fire, no time window " +
			"where the file exists without the special bit. If DEST is on `$PATH` (`/usr/" +
			"local/bin`, `/usr/bin`), every user can invoke the elevated binary. Only install " +
			"setuid / setgid binaries from trusted builds you have reviewed, and prefer " +
			"narrow capabilities (`setcap cap_net_bind_service+ep`) over broad setuid.",
		Check: checkZC1826,
	})
}

func checkZC1826(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "install" {
		return nil
	}
	for i, arg := range cmd.Arguments {
		v := arg.String()
		var mode string
		if v == "-m" || v == "--mode" {
			if i+1 < len(cmd.Arguments) {
				mode = cmd.Arguments[i+1].String()
			}
		} else if strings.HasPrefix(v, "-m") && len(v) > 2 {
			mode = v[2:]
		}
		if mode == "" {
			continue
		}
		mode = strings.TrimSpace(mode)
		if len(mode) == 4 && (mode[0] == '4' || mode[0] == '2' || mode[0] == '6') {
			return []Violation{{
				KataID: "ZC1826",
				Message: "`install -m " + mode + "` drops a setuid/setgid binary in one " +
					"step. If DEST is on `$PATH`, every local user can invoke the " +
					"elevated binary. Only install trusted builds, and prefer narrow " +
					"`setcap` capabilities over setuid.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}
