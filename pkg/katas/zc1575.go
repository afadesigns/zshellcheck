// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1575",
		Title:    "Error on `aws configure set aws_secret_access_key <value>` — secret on cmdline",
		Severity: SeverityError,
		Description: "`aws configure set aws_secret_access_key …` writes the secret access key " +
			"into `~/.aws/credentials` and leaves the raw value in `ps` / shell history until " +
			"the process exits. On a shared CI runner or a multi-user host, that window is " +
			"long enough for a co-tenant to snapshot the key. Use IAM-role-based auth (EC2 " +
			"instance profile, IRSA on EKS, OIDC from GitHub / GitLab) or read the value from " +
			"stdin / a 0600 file and let `aws configure` import it interactively.",
		Check: checkZC1575,
	})
}

func checkZC1575(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "aws" {
		return nil
	}

	args := make([]string, 0, len(cmd.Arguments))
	for _, a := range cmd.Arguments {
		args = append(args, a.String())
	}

	// aws configure set aws_secret_access_key VALUE
	for i := 0; i+3 < len(args); i++ {
		if args[i] == "configure" && args[i+1] == "set" {
			key := strings.ToLower(args[i+2])
			if key == "aws_secret_access_key" || key == "aws_session_token" ||
				strings.Contains(key, "secret") {
				return []Violation{{
					KataID: "ZC1575",
					Message: "`aws configure set " + args[i+2] + " …` puts the secret in " +
						"ps / history. Use IAM-role auth or import from stdin / 0600 file.",
					Line:   cmd.Token.Line,
					Column: cmd.Token.Column,
					Level:  SeverityError,
				}}
			}
		}
	}
	return nil
}
