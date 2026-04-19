package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1834",
		Title:    "Error on `tc qdisc … root netem loss 100%` — hard blackhole on a live interface",
		Severity: SeverityError,
		Description: "`tc qdisc add/replace dev IFACE root netem loss 100%` (also `corrupt 100%` " +
			"or `duplicate 100%` with no buffering) installs a Linux kernel qdisc that " +
			"drops every outbound packet on the named interface. Running this on the " +
			"interface that carries your SSH session is the canonical way to lock " +
			"yourself out of a remote host — the `tc` command returns success, the kernel " +
			"happily applies the rule, and the next TCP segment ACK never arrives. Even on " +
			"the console it halts any process that depends on the interface. Stage netem " +
			"experiments on a secondary interface, wrap them in `at now + 5 minutes` (or a " +
			"`timeout … tc qdisc del …` recovery trap) so a partial failure does not leave " +
			"the link permanently black-holed.",
		Check: checkZC1834,
	})
}

func checkZC1834(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "tc" {
		return nil
	}
	args := cmd.Arguments
	if len(args) < 3 {
		return nil
	}
	if args[0].String() != "qdisc" {
		return nil
	}
	action := args[1].String()
	if action != "add" && action != "replace" && action != "change" {
		return nil
	}
	// Walk remaining args looking for `netem <mode> 100%`.
	for i := 2; i+2 < len(args); i++ {
		if args[i].String() != "netem" {
			continue
		}
		for j := i + 1; j+1 < len(args); j++ {
			mode := args[j].String()
			if mode != "loss" && mode != "corrupt" && mode != "duplicate" {
				continue
			}
			val := args[j+1].String()
			if val == "100%" || val == "100" {
				return []Violation{{
					KataID: "ZC1834",
					Message: "`tc qdisc " + action + " … netem " + mode + " 100%` " +
						"black-holes every packet on the target interface — remote SSH " +
						"dies instantly. Stage on a secondary dev or wrap in a timed " +
						"recovery (`at now + N minutes … tc qdisc del …`).",
					Line:   cmd.Token.Line,
					Column: cmd.Token.Column,
					Level:  SeverityError,
				}}
			}
		}
	}
	return nil
}
