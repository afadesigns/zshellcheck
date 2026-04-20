package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1982",
		Title:    "Error on `ipcrm -a` — removes every SysV IPC object, breaks Postgres/Oracle/shm apps",
		Severity: SeverityError,
		Description: "`ipcrm -a` deletes every System V shared-memory segment, semaphore set, " +
			"and message queue owned by the caller (or, as root, every object on the " +
			"host). Long-running services that rely on SysV IPC — PostgreSQL's shared " +
			"buffers, Oracle's SGA, the `sysv` session store in several RDBMS test " +
			"suites, shm-based mutexes in batch pipelines — lose their backing store " +
			"mid-transaction and either SIGSEGV or return `EINVAL` on the next access. " +
			"Scope the removal: `ipcrm -m ID`/`-s ID`/`-q ID` against the specific " +
			"identifier reported by `ipcs -a`, after confirming no running process " +
			"attached to it.",
		Check: checkZC1982,
	})
}

func checkZC1982(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "ipcrm" {
		return nil
	}
	for _, arg := range cmd.Arguments {
		if arg.String() == "-a" {
			return []Violation{{
				KataID: "ZC1982",
				Message: "`ipcrm -a` deletes every SysV shm/sem/mqueue object — " +
					"Postgres/Oracle/shm-based services lose their backing store " +
					"mid-transaction. Scope with `-m`/`-s`/`-q` on the specific ID " +
					"after checking `ipcs -a`.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}
	return nil
}
