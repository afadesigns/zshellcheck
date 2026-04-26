// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC2000",
		Title:    "Error on `kubectl taint nodes $NODE key=value:NoExecute` — evicts every non-tolerating pod off the node",
		Severity: SeverityError,
		Description: "A `NoExecute` taint kicks every existing pod off the node unless the pod " +
			"spec explicitly tolerates it. Draining one node during a rolling upgrade " +
			"is one thing; a script that types the taint wrong (typoed " +
			"toleration value, applying to `--all` nodes, or iterating a node list " +
			"without a pause) can empty a whole cluster in seconds and trigger " +
			"cascade reschedules that overwhelm the scheduler. Prefer `kubectl drain " +
			"$NODE` (which respects PodDisruptionBudget and runs PreStop hooks) or a " +
			"`NoSchedule` taint for gentle drain; reserve `NoExecute` for genuine " +
			"incident response with a runbook and a safety countdown.",
		Check: checkZC2000,
	})
}

func checkZC2000(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "kubectl" {
		return nil
	}
	if len(cmd.Arguments) < 2 || cmd.Arguments[0].String() != "taint" {
		return nil
	}
	for _, arg := range cmd.Arguments[1:] {
		v := arg.String()
		if strings.Contains(v, ":NoExecute") {
			return []Violation{{
				KataID: "ZC2000",
				Message: "`kubectl taint nodes … :NoExecute` evicts every non-tolerating " +
					"pod immediately — a typo on `--all` nodes empties the cluster. " +
					"Prefer `kubectl drain $NODE` or a `:NoSchedule` taint for " +
					"gentle drain.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}
	return nil
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC2001",
		Title:    "Warn on `unsetopt EVAL_LINENO` — `$LINENO` inside `eval` stops tracking source, stack traces go blank",
		Severity: SeverityWarning,
		Description: "On by default, Zsh's `EVAL_LINENO` keeps `$LINENO`, `$funcfiletrace`, and " +
			"`$funcstack` pointing at the line inside the `eval`ed string where the " +
			"error actually happened. Turning the option off (`unsetopt EVAL_LINENO` " +
			"or `setopt NO_EVAL_LINENO`) reverts to pre-Zsh-4.3 behaviour: `$LINENO` " +
			"collapses to the line that launched the `eval`, so every runtime error " +
			"inside a generated config, a lazy-loaded function, or a `compile`d string " +
			"reports the same line number and the stack trace loses every frame past " +
			"the eval. Keep the option on; if strict POSIX-matching line numbers are " +
			"needed inside one helper, scope with `emulate -LR sh` in that function.",
		Check: checkZC2001,
	})
}

func checkZC2001(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	var enabling bool
	switch ident.Value {
	case "setopt":
		enabling = true
	case "unsetopt":
		enabling = false
	default:
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := zc2001Canonical(arg.String())
		switch v {
		case "EVALLINENO":
			if !enabling {
				return zc2001Hit(cmd, "unsetopt EVAL_LINENO")
			}
		case "NOEVALLINENO":
			if enabling {
				return zc2001Hit(cmd, "setopt NO_EVAL_LINENO")
			}
		}
	}
	return nil
}

func zc2001Canonical(s string) string {
	out := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '_' || c == '-' {
			continue
		}
		if c >= 'a' && c <= 'z' {
			c -= 'a' - 'A'
		}
		out = append(out, c)
	}
	return string(out)
}

func zc2001Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC2001",
		Message: "`" + form + "` reverts `$LINENO` inside `eval` to the outer line — " +
			"errors in generated configs collapse to a single source line and " +
			"stack frames past `eval` vanish. Keep on; scope via `emulate -LR sh`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC2002",
		Title:    "Error on `crictl rmi -a` / `crictl rm -af` — wipes every image/container on the Kubernetes node",
		Severity: SeverityError,
		Description: "`crictl` talks directly to the node's CRI runtime (containerd, CRI-O), " +
			"below the kubelet and the cluster API. `crictl rmi -a` removes every " +
			"cached image including the ones currently backing running pods — the " +
			"kubelet must immediately re-pull from the registry, and image-pull rate " +
			"limits or network blips turn the node Unready. `crictl rm -af` force-" +
			"removes every container on the node, killing pods without running " +
			"PreStop hooks or honoring PodDisruptionBudget. Route maintenance through " +
			"`kubectl drain $NODE` + `kubectl delete pod --grace-period=30`; use " +
			"`crictl` at most on a cordoned, drained node with a documented recovery " +
			"plan.",
		Check: checkZC2002,
	})
}

func checkZC2002(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "crictl" {
		return nil
	}
	if len(cmd.Arguments) == 0 {
		return nil
	}
	sub := cmd.Arguments[0].String()
	if sub != "rmi" && sub != "rm" {
		return nil
	}
	hasAll := false
	hasForce := false
	for _, arg := range cmd.Arguments[1:] {
		v := arg.String()
		if v == "-a" || v == "--all" || v == "-af" || v == "-fa" {
			hasAll = true
		}
		if v == "-f" || v == "--force" || v == "-af" || v == "-fa" {
			hasForce = true
		}
	}
	if sub == "rmi" && hasAll {
		return zc2002Hit(cmd, "crictl rmi -a")
	}
	if sub == "rm" && hasAll && hasForce {
		return zc2002Hit(cmd, "crictl rm -af")
	}
	return nil
}

func zc2002Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC2002",
		Message: "`" + form + "` talks to the node CRI directly, under the kubelet — " +
			"images/containers backing running pods disappear, kubelet must " +
			"re-pull or re-run. Route through `kubectl drain`/`delete pod`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC2003",
		Title:    "Warn on `setopt KSH_ZERO_SUBSCRIPT` — `$arr[0]` stops aliasing the first element",
		Severity: SeverityWarning,
		Description: "Default Zsh treats `$arr[0]` as a quirk-compatibility alias for `$arr[1]` " +
			"— `arr=(a b c); echo $arr[0]` prints `a`, and `arr[0]=new` rewrites the " +
			"first element. `setopt KSH_ZERO_SUBSCRIPT` flips that to ksh semantics: " +
			"`$arr[0]` becomes a distinct slot (the element just before the " +
			"1-indexed head, which Zsh stores separately), so reads silently switch " +
			"to empty string and `arr[0]=new` no longer touches `$arr[1]`. Any Zsh " +
			"code that intentionally used `$arr[0]` as a shortcut breaks, and ported " +
			"Bash/ksh code that assumes 0-indexed access meets a split-world model. " +
			"Leave the option off; use `$arr[1]` explicitly when you want the first " +
			"element, and adopt `KSH_ARRAYS` scoped with `emulate -LR ksh` for " +
			"ksh-style code paths.",
		Check: checkZC2003,
	})
}

func checkZC2003(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	var enabling bool
	switch ident.Value {
	case "setopt":
		enabling = true
	case "unsetopt":
		enabling = false
	default:
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := zc2003Canonical(arg.String())
		switch v {
		case "KSHZEROSUBSCRIPT":
			if enabling {
				return zc2003Hit(cmd, "setopt KSH_ZERO_SUBSCRIPT")
			}
		case "NOKSHZEROSUBSCRIPT":
			if !enabling {
				return zc2003Hit(cmd, "unsetopt NO_KSH_ZERO_SUBSCRIPT")
			}
		}
	}
	return nil
}

func zc2003Canonical(s string) string {
	out := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '_' || c == '-' {
			continue
		}
		if c >= 'a' && c <= 'z' {
			c -= 'a' - 'A'
		}
		out = append(out, c)
	}
	return string(out)
}

func zc2003Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC2003",
		Message: "`" + form + "` stops aliasing `$arr[0]` to `$arr[1]` — every later " +
			"read of `$arr[0]` silently returns empty and `arr[0]=new` stops " +
			"updating the first element. Use `$arr[1]` explicitly.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}
