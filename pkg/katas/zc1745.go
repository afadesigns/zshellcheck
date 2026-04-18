package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1745",
		Title:    "Error on `poetry publish --password PASS` / `twine upload -p PASS` — registry secret in argv",
		Severity: SeverityError,
		Description: "Poetry's `publish --username USER --password PASS` and Twine's `upload " +
			"--username USER --password PASS` (or the short `-u`/`-p` forms) put the PyPI / " +
			"private-index password in argv — visible in `ps`, `/proc/<pid>/cmdline`, shell " +
			"history, and CI logs. Use the `POETRY_PYPI_TOKEN_<NAME>` / `TWINE_USERNAME` + " +
			"`TWINE_PASSWORD` environment variables (sourced from a secrets manager) or a " +
			"`~/.pypirc` file with `0600` perms so the credential never reaches the command " +
			"line.",
		Check: checkZC1745,
	})
}

func checkZC1745(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	var tool string
	switch ident.Value {
	case "poetry":
		if len(cmd.Arguments) == 0 || cmd.Arguments[0].String() != "publish" {
			return nil
		}
		tool = "poetry publish"
	case "twine":
		if len(cmd.Arguments) == 0 || cmd.Arguments[0].String() != "upload" {
			return nil
		}
		tool = "twine upload"
	default:
		return nil
	}

	prevPwd := false
	for _, arg := range cmd.Arguments[1:] {
		v := arg.String()
		if prevPwd {
			return zc1745Hit(cmd, tool, "--password "+v)
		}
		switch {
		case v == "--password" || v == "-p":
			prevPwd = true
		case strings.HasPrefix(v, "--password="):
			return zc1745Hit(cmd, tool, v)
		}
	}
	return nil
}

func zc1745Hit(cmd *ast.SimpleCommand, tool, what string) []Violation {
	return []Violation{{
		KataID: "ZC1745",
		Message: "`" + tool + " " + what + "` puts the registry password in argv — " +
			"visible in `ps`, `/proc`, history. Use env vars (`POETRY_PYPI_TOKEN_<" +
			"NAME>`, `TWINE_PASSWORD`) or a `0600` `~/.pypirc` file.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}
