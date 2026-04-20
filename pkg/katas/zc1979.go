package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1979",
		Title:    "Warn on `setopt HIST_FCNTL_LOCK` — `fcntl()` lock on NFS `$HISTFILE` stalls or deadlocks",
		Severity: SeverityWarning,
		Description: "Off by default, Zsh serialises writes to `$HISTFILE` with its own " +
			"lock-file dance next to the history. `setopt HIST_FCNTL_LOCK` switches to " +
			"POSIX `fcntl()` advisory locking — which is the safer primitive on local " +
			"filesystems, but on NFS homes the lock is proxied through `rpc.lockd` and " +
			"a single hung client or rebooted NFS server leaves every other shell " +
			"blocked the next time it tries to write history. The interactive shell " +
			"appears frozen on prompt return, and scripts that source user rc files " +
			"hang in `zshaddhistory`. Keep the option off on NFS homes; only turn it on " +
			"when `$HISTFILE` lives on a local filesystem (ext4, xfs, btrfs, zfs local " +
			"pool) that implements `fcntl()` without network round-trips.",
		Check: checkZC1979,
	})
}

func checkZC1979(node ast.Node) []Violation {
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
		v := zc1979Canonical(arg.String())
		switch v {
		case "HISTFCNTLLOCK":
			if enabling {
				return zc1979Hit(cmd, "setopt HIST_FCNTL_LOCK")
			}
		case "NOHISTFCNTLLOCK":
			if !enabling {
				return zc1979Hit(cmd, "unsetopt NO_HIST_FCNTL_LOCK")
			}
		}
	}
	return nil
}

func zc1979Canonical(s string) string {
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

func zc1979Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1979",
		Message: "`" + form + "` routes `$HISTFILE` locking through POSIX `fcntl()` — " +
			"on NFS home directories a hung `rpc.lockd` freezes every other shell " +
			"at the next prompt. Keep off; enable only when `$HISTFILE` is on a " +
			"local fs.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}
