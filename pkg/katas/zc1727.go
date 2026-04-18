package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1727",
		Title:    "Error on `curl/wget --proxy http://USER:PASS@HOST` — proxy credentials in argv",
		Severity: SeverityError,
		Description: "Embedding the proxy username and password in the URL passed to `--proxy` " +
			"(curl), `-x` (curl short form), or `--proxy-password=` (wget) lands the " +
			"credential in argv — visible in `ps`, `/proc/<pid>/cmdline`, shell history, " +
			"and CI logs. Configure the proxy through `~/.curlrc` / `~/.netrc` (chmod 600) " +
			"for curl, or `~/.wgetrc` for wget, so the secret never reaches the command " +
			"line.",
		Check: checkZC1727,
	})
}

func checkZC1727(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	switch ident.Value {
	case "curl":
		return zc1727Curl(cmd)
	case "wget":
		return zc1727Wget(cmd)
	}
	return nil
}

func zc1727Curl(cmd *ast.SimpleCommand) []Violation {
	prevProxy := false
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if prevProxy {
			if zc1727URLHasCreds(v) {
				return zc1727Hit(cmd, "curl --proxy "+v)
			}
			prevProxy = false
			continue
		}
		switch {
		case v == "--proxy" || v == "-x":
			prevProxy = true
		case strings.HasPrefix(v, "--proxy="):
			val := strings.TrimPrefix(v, "--proxy=")
			if zc1727URLHasCreds(val) {
				return zc1727Hit(cmd, "curl "+v)
			}
		case strings.HasPrefix(v, "-x"):
			val := v[2:]
			if zc1727URLHasCreds(val) {
				return zc1727Hit(cmd, "curl "+v)
			}
		}
	}
	return nil
}

func zc1727Wget(cmd *ast.SimpleCommand) []Violation {
	prevPwd := false
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if prevPwd {
			return zc1727Hit(cmd, "wget --proxy-password "+v)
		}
		switch {
		case v == "--proxy-password":
			prevPwd = true
		case strings.HasPrefix(v, "--proxy-password="):
			return zc1727Hit(cmd, "wget "+v)
		}
	}
	return nil
}

// zc1727URLHasCreds returns true when the URL contains a `userinfo` portion
// (text between `://` and the next `@` before any `/`).
func zc1727URLHasCreds(url string) bool {
	scheme := strings.Index(url, "://")
	if scheme < 0 {
		return false
	}
	rest := url[scheme+3:]
	at := strings.Index(rest, "@")
	if at < 0 {
		return false
	}
	if slash := strings.Index(rest, "/"); slash >= 0 && slash < at {
		return false
	}
	return true
}

func zc1727Hit(cmd *ast.SimpleCommand, what string) []Violation {
	return []Violation{{
		KataID: "ZC1727",
		Message: "`" + what + "` puts proxy credentials in argv — visible in `ps`, " +
			"`/proc`, history. Move them into `~/.curlrc` / `~/.netrc` (chmod 600) or " +
			"`~/.wgetrc`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}
