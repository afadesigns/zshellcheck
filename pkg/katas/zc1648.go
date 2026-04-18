package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1648",
		Title:    "Error on `cp /dev/null /var/log/...` / `truncate -s 0 /var/log/...` — audit-log wipe",
		Severity: SeverityError,
		Description: "Replacing a file under `/var/log/` with `/dev/null` or truncating it to " +
			"size zero erases audit evidence: failed login attempts from `auth.log`, sudo " +
			"usage from `sudo.log`, kernel audit trail from `audit/audit.log`, console " +
			"history from `wtmp` / `btmp`. Scripts that do this during \"cleanup\" are almost " +
			"always misusing logrotate (which handles rotation safely via a `create` stage) " +
			"or deliberately covering tracks. Use `logrotate -f /etc/logrotate.d/...` for " +
			"rotation, `journalctl --vacuum-time=...` for journald.",
		Check: checkZC1648,
	})
}

func checkZC1648(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	switch ident.Value {
	case "cp":
		if len(cmd.Arguments) < 2 {
			return nil
		}
		if cmd.Arguments[0].String() != "/dev/null" {
			return nil
		}
		dest := cmd.Arguments[1].String()
		if strings.HasPrefix(dest, "/var/log/") {
			return zc1648Hit(cmd, "cp /dev/null "+dest)
		}
	case "truncate":
		var zeroSize bool
		var target string
		for i, arg := range cmd.Arguments {
			v := arg.String()
			if v == "-s" && i+1 < len(cmd.Arguments) && cmd.Arguments[i+1].String() == "0" {
				zeroSize = true
			}
			if strings.HasPrefix(v, "/var/log/") {
				target = v
			}
		}
		if zeroSize && target != "" {
			return zc1648Hit(cmd, "truncate -s 0 "+target)
		}
	}
	return nil
}

func zc1648Hit(cmd *ast.SimpleCommand, desc string) []Violation {
	return []Violation{{
		KataID: "ZC1648",
		Message: "`" + desc + "` wipes an audit log — use `logrotate -f` or " +
			"`journalctl --vacuum-time=...` instead.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}
