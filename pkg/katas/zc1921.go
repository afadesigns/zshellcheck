package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1921",
		Title:    "Warn on `systemctl kill -s KILL` / `--signal=SIGKILL` — skips `ExecStop=`, leaks resources",
		Severity: SeverityWarning,
		Description: "`systemctl kill UNIT -s KILL` (and `--signal=9` / `SIGKILL`) bypasses the " +
			"unit's `ExecStop=` sequence and the `TimeoutStopSec=` budget. Any lockfile, " +
			"socket, or shared-memory segment the service was supposed to unlink survives; the " +
			"next restart often fails with \"address already in use\" or a corrupt journal. " +
			"Default to `systemctl stop UNIT` (or `restart`) and let the stop sequence run. " +
			"Reserve `-s KILL` for a last-resort recovery path with a runbook attached.",
		Check: checkZC1921,
	})
}

func checkZC1921(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "systemctl" {
		return nil
	}
	if len(cmd.Arguments) == 0 || cmd.Arguments[0].String() != "kill" {
		return nil
	}

	for i, arg := range cmd.Arguments[1:] {
		v := arg.String()
		if strings.HasPrefix(v, "--signal=") {
			if zc1921IsHardSignal(strings.TrimPrefix(v, "--signal=")) {
				return zc1921Hit(cmd, v)
			}
		}
		if v == "-s" && i+2 < len(cmd.Arguments) {
			sig := cmd.Arguments[i+2].String()
			if zc1921IsHardSignal(sig) {
				return zc1921Hit(cmd, "-s "+sig)
			}
		}
	}
	return nil
}

func zc1921IsHardSignal(sig string) bool {
	switch strings.ToUpper(sig) {
	case "KILL", "SIGKILL", "9":
		return true
	}
	return false
}

func zc1921Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1921",
		Message: "`systemctl kill " + form + "` bypasses `ExecStop=` and " +
			"`TimeoutStopSec=` — lockfiles, sockets, and shm segments survive and the next " +
			"restart often fails with \"address already in use\". Use `systemctl stop` or " +
			"`restart` instead.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}
