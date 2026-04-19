package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

var zc1757DangerousScopes = []string{
	"delete_repo",
	"admin:org",
	"admin:enterprise",
	"admin:public_key",
	"admin:org_hook",
	"site_admin",
	"admin:repo_hook",
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1757",
		Title:    "Warn on `gh auth refresh --scopes delete_repo|admin:*` — token escalated to destructive perms",
		Severity: SeverityWarning,
		Description: "`gh auth refresh --scopes <list>` (also `gh auth login --scopes`) rotates " +
			"the stored OAuth token with additional scopes. `delete_repo`, `admin:org`, " +
			"`admin:enterprise`, `admin:public_key`, and `admin:*_hook` give the token " +
			"permanent destructive perms that outlast the script that asked for them — a " +
			"compromised token now carries repo-deletion, org-membership, and SSH-key " +
			"manipulation rights. Request the minimum scope the task needs (`repo`, " +
			"`workflow`) and rotate the token off when the elevated operation completes.",
		Check: checkZC1757,
	})
}

func checkZC1757(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "gh" {
		return nil
	}
	if len(cmd.Arguments) < 2 {
		return nil
	}
	if cmd.Arguments[0].String() != "auth" {
		return nil
	}
	sub := cmd.Arguments[1].String()
	if sub != "refresh" && sub != "login" {
		return nil
	}

	prevScopes := false
	for _, arg := range cmd.Arguments[2:] {
		v := arg.String()
		if prevScopes {
			if hit := zc1757MatchScopes(v); hit != "" {
				return zc1757Hit(cmd, sub, hit)
			}
			prevScopes = false
			continue
		}
		switch {
		case v == "--scopes" || v == "-s":
			prevScopes = true
		case strings.HasPrefix(v, "--scopes="):
			if hit := zc1757MatchScopes(strings.TrimPrefix(v, "--scopes=")); hit != "" {
				return zc1757Hit(cmd, sub, hit)
			}
		}
	}
	return nil
}

func zc1757MatchScopes(list string) string {
	for _, scope := range strings.Split(list, ",") {
		scope = strings.TrimSpace(scope)
		for _, danger := range zc1757DangerousScopes {
			if scope == danger {
				return scope
			}
		}
	}
	return ""
}

func zc1757Hit(cmd *ast.SimpleCommand, sub, scope string) []Violation {
	return []Violation{{
		KataID: "ZC1757",
		Message: "`gh auth " + sub + " --scopes " + scope + "` escalates the token to " +
			"destructive privileges that outlast the script. Request the minimum " +
			"scope (`repo`, `workflow`) and rotate the token when done.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}
