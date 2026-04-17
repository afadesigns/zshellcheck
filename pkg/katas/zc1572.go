package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1572",
		Title:    "Warn on `docker run -e PASSWORD=<value>` — secret in container env / inspect",
		Severity: SeverityWarning,
		Description: "Passing a secret through `docker run -e NAME=value` puts it in the output " +
			"of `docker inspect`, the container's `/proc/1/environ` (readable by anything that " +
			"shares the PID namespace), and the shell history of whoever launched the " +
			"container. Use `--env-file` with 0600 perms, a secret-mount `--secret` via " +
			"BuildKit / Swarm, or mount a tmpfs file the container reads at runtime.",
		Check: checkZC1572,
	})
}

var secretEnvPrefixes = []string{
	"PASSWORD", "PASSWD", "PASS",
	"SECRET", "SECRET_KEY", "API_KEY",
	"TOKEN", "AUTH_TOKEN", "ACCESS_TOKEN",
	"PRIVATE_KEY", "DB_PASSWORD", "DB_PASS",
	"AWS_SECRET", "AWS_SECRET_ACCESS_KEY",
	"GITHUB_TOKEN", "GH_TOKEN", "NPM_TOKEN",
}

func checkZC1572(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "docker" && ident.Value != "podman" && ident.Value != "nerdctl" {
		return nil
	}

	if len(cmd.Arguments) == 0 || cmd.Arguments[0].String() != "run" {
		return nil
	}

	var prevE bool
	for _, arg := range cmd.Arguments[1:] {
		v := arg.String()
		if prevE {
			prevE = false
			name, value, ok := strings.Cut(v, "=")
			if !ok || value == "" {
				// `-e NAME` (value from caller env) is fine.
				continue
			}
			if looksLikeSecret(name) {
				return []Violation{{
					KataID: "ZC1572",
					Message: "`-e " + name + "=<value>` writes the secret into `docker " +
						"inspect` and `/proc/1/environ`. Use `--env-file` 0600 or " +
						"`--secret`.",
					Line:   cmd.Token.Line,
					Column: cmd.Token.Column,
					Level:  SeverityWarning,
				}}
			}
		}
		if v == "-e" || v == "--env" {
			prevE = true
			continue
		}
		// Joined form: --env=NAME=value
		if strings.HasPrefix(v, "--env=") {
			inner := v[len("--env="):]
			name, value, ok := strings.Cut(inner, "=")
			if ok && value != "" && looksLikeSecret(name) {
				return []Violation{{
					KataID: "ZC1572",
					Message: "`--env=" + name + "=<value>` writes the secret into `docker " +
						"inspect` and `/proc/1/environ`. Use `--env-file` 0600 or " +
						"`--secret`.",
					Line:   cmd.Token.Line,
					Column: cmd.Token.Column,
					Level:  SeverityWarning,
				}}
			}
		}
	}
	return nil
}

func looksLikeSecret(name string) bool {
	up := strings.ToUpper(name)
	for _, p := range secretEnvPrefixes {
		if up == p || strings.Contains(up, p) {
			return true
		}
	}
	return false
}
