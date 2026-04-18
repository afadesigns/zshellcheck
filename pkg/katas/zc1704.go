package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1704",
		Title:    "Error on `aws ec2 authorize-security-group-ingress --cidr 0.0.0.0/0` — port open to the internet",
		Severity: SeverityError,
		Description: "`aws ec2 authorize-security-group-ingress --cidr 0.0.0.0/0` (or `::/0` for " +
			"IPv6) adds a rule that accepts the specified protocol/port from any source — " +
			"the exact shape shodan, automated login-probers, and every exploit-as-a-" +
			"service customer scans for. Restrict the source to the office CIDR, a VPN " +
			"range, or a named security-group (`--source-group sg-…`). If the workload " +
			"genuinely needs public access, front it with an ALB / API Gateway / CloudFront " +
			"with WAF — not a raw SG rule from a shell script.",
		Check: checkZC1704,
	})
}

func checkZC1704(node ast.Node) []Violation {
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

	if len(cmd.Arguments) < 2 {
		return nil
	}
	if cmd.Arguments[0].String() != "ec2" {
		return nil
	}
	if cmd.Arguments[1].String() != "authorize-security-group-ingress" {
		return nil
	}

	for i, arg := range cmd.Arguments[2:] {
		v := arg.String()
		if v != "--cidr" && v != "--cidr-ip" && v != "--cidr-ipv6" {
			continue
		}
		idx := i + 3
		if idx >= len(cmd.Arguments) {
			continue
		}
		cidr := cmd.Arguments[idx].String()
		if cidr == "0.0.0.0/0" || cidr == "::/0" {
			return []Violation{{
				KataID: "ZC1704",
				Message: "`aws ec2 authorize-security-group-ingress --cidr " + cidr +
					"` opens the port to the entire internet — scope to a known source CIDR " +
					"or `--source-group sg-…`.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}
	return nil
}
