package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1812",
		Title:    "Error on `aws ssm put-parameter --type SecureString --value SECRET` — plaintext in argv",
		Severity: SeverityError,
		Description: "`aws ssm put-parameter` stores the value as-is under the given parameter " +
			"name; the whole point of `--type SecureString` is that the value is sensitive. " +
			"Passing the plaintext with `--value SECRET` (or `--value=SECRET`) puts the " +
			"secret in argv where `ps`, `/proc/PID/cmdline`, shell history, and AWS CLI " +
			"debug logs (`--debug`) can read it. Pipe the value in from stdin with `--cli-" +
			"input-json file://param.json` (mode 0600) or use `aws secretsmanager " +
			"create-secret --secret-string file://secret` which supports `file://` in every " +
			"code path.",
		Check: checkZC1812,
	})
}

func checkZC1812(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "aws" {
		return nil
	}
	if len(cmd.Arguments) < 2 {
		return nil
	}
	if cmd.Arguments[0].String() != "ssm" {
		return nil
	}
	if cmd.Arguments[1].String() != "put-parameter" {
		return nil
	}

	hasSecureString := false
	hasInlineValue := false
	for i, arg := range cmd.Arguments[2:] {
		v := arg.String()
		switch {
		case v == "--type":
			if 2+i+1 < len(cmd.Arguments) && cmd.Arguments[2+i+1].String() == "SecureString" {
				hasSecureString = true
			}
		case v == "--type=SecureString":
			hasSecureString = true
		case v == "--value":
			if 2+i+1 < len(cmd.Arguments) {
				next := cmd.Arguments[2+i+1].String()
				if next != "" && !strings.HasPrefix(next, "file://") && !strings.HasPrefix(next, "-") {
					hasInlineValue = true
				}
			}
		case strings.HasPrefix(v, "--value="):
			val := strings.TrimPrefix(v, "--value=")
			if val != "" && !strings.HasPrefix(val, "file://") {
				hasInlineValue = true
			}
		}
	}
	if !hasSecureString || !hasInlineValue {
		return nil
	}
	return []Violation{{
		KataID: "ZC1812",
		Message: "`aws ssm put-parameter --type SecureString --value …` puts the " +
			"plaintext in argv — `ps` / `/proc/PID/cmdline` / history / CLI debug " +
			"logs can read it. Use `--cli-input-json file://…` (mode 0600) or the " +
			"`file://` form for `--value`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}
