package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1786",
		Title:    "Error on `mount.cifs ... -o password=SECRET` — SMB password in argv",
		Severity: SeverityError,
		Description: "Passing `password=` (or `pass=`) inside `mount.cifs` / `mount -t cifs` " +
			"options puts the SMB password in argv. Any local user who can read `ps`, " +
			"`/proc/PID/cmdline`, or process-accounting records gets the cleartext, and the " +
			"line also ends up in shell history and — if captured — in CI logs. Use a " +
			"`credentials=/etc/cifs-creds` file (`0600`, `username=` and `password=` lines), " +
			"the `$USER`/`$PASSWD` env vars `mount.cifs` reads when those options are " +
			"missing, or `pam_mount` for login-time mounts.",
		Check: checkZC1786,
	})
}

func checkZC1786(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	switch ident.Value {
	case "mount.cifs":
		// direct form
	case "mount":
		isCifs := false
		for i, arg := range cmd.Arguments {
			v := arg.String()
			if (v == "-t" || v == "--types") && i+1 < len(cmd.Arguments) {
				t := cmd.Arguments[i+1].String()
				if t == "cifs" || t == "smb3" {
					isCifs = true
					break
				}
			}
		}
		if !isCifs {
			return nil
		}
	default:
		return nil
	}

	for i, arg := range cmd.Arguments {
		v := arg.String()
		var opts string
		if strings.HasPrefix(v, "-o") && len(v) > 2 {
			opts = v[2:]
		} else if v == "-o" && i+1 < len(cmd.Arguments) {
			opts = cmd.Arguments[i+1].String()
		}
		if opts == "" {
			continue
		}
		opts = strings.Trim(opts, "\"'")
		if zc1786OptsHavePassword(opts) {
			return []Violation{{
				KataID: "ZC1786",
				Message: "`" + ident.Value + " ... password=…` leaks the SMB password " +
					"into argv / `ps` / `/proc/PID/cmdline`. Use `credentials=/path/" +
					"to/creds` (mode 0600) or `$PASSWD` env var instead.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}
	return nil
}

func zc1786OptsHavePassword(opts string) bool {
	for _, field := range strings.Split(opts, ",") {
		key, value, ok := strings.Cut(strings.TrimSpace(field), "=")
		if !ok {
			continue
		}
		switch strings.ToLower(key) {
		case "password", "pass", "password2":
			if value != "" {
				return true
			}
		}
	}
	return false
}
