// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1952",
		Title:    "Error on `zfs set sync=disabled` — `fsync()` becomes a no-op, crash loses unflushed writes",
		Severity: SeverityError,
		Description: "`zfs set sync=disabled POOL/DATASET` turns `fsync()`, `O_SYNC`, and `O_DSYNC` " +
			"into no-ops on that dataset. PostgreSQL, MariaDB, etcd, and every application that " +
			"relies on fsync for durability will report success for writes that are still in the " +
			"ARC, so a panic or power cut loses minutes of committed transactions. The flag is " +
			"a benchmarking knob, not a production setting. Leave sync at `standard` and, if " +
			"latency is the concern, add a `log` vdev (SLOG) or tune " +
			"`zfs_txg_timeout` instead.",
		Check: checkZC1952,
	})
}

func checkZC1952(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "zfs" {
		return nil
	}
	if len(cmd.Arguments) < 2 || cmd.Arguments[0].String() != "set" {
		return nil
	}
	for _, arg := range cmd.Arguments[1:] {
		v := arg.String()
		if strings.HasPrefix(v, "sync=") {
			val := strings.TrimPrefix(v, "sync=")
			if val == "disabled" {
				return []Violation{{
					KataID: "ZC1952",
					Message: "`zfs set sync=disabled` turns `fsync()` into a no-op — DBs " +
						"(PostgreSQL/MariaDB/etcd) lose committed transactions on crash. Leave " +
						"sync at `standard`; use a SLOG vdev if latency is the concern.",
					Line:   cmd.Token.Line,
					Column: cmd.Token.Column,
					Level:  SeverityError,
				}}
			}
		}
	}
	return nil
}
