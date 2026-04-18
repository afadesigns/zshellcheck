package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1704(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — scoped CIDR",
			input:    `aws ec2 authorize-security-group-ingress --group-id sg-123 --protocol tcp --port 22 --cidr 10.0.0.0/8`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — describe-security-groups (different subcommand)",
			input:    `aws ec2 describe-security-groups`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — --cidr 0.0.0.0/0",
			input: `aws ec2 authorize-security-group-ingress --group-id sg-123 --protocol tcp --port 22 --cidr 0.0.0.0/0`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1704",
					Message: "`aws ec2 authorize-security-group-ingress --cidr 0.0.0.0/0` opens the port to the entire internet — scope to a known source CIDR or `--source-group sg-…`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — --cidr-ipv6 ::/0",
			input: `aws ec2 authorize-security-group-ingress --group-id sg-123 --ip-permissions proto --cidr-ipv6 ::/0`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1704",
					Message: "`aws ec2 authorize-security-group-ingress --cidr ::/0` opens the port to the entire internet — scope to a known source CIDR or `--source-group sg-…`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1704")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
