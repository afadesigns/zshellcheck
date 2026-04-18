package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1624(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — az login (interactive)",
			input:    `az login`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — az login --service-principal with federated token",
			input:    `az login --service-principal -u appid -t tenantid`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — az login -p $SECRET",
			input: `az login --service-principal -u appid -p $SECRET`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1624",
					Message: "`az login -p` puts the SP password in argv — visible in `ps` / `/proc/<pid>/cmdline`. Use federated-token OIDC, managed identity, or `AZURE_PASSWORD` via a protected env var.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — az login --password literal",
			input: `az login --password hunter2 -u user`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1624",
					Message: "`az login --password` puts the SP password in argv — visible in `ps` / `/proc/<pid>/cmdline`. Use federated-token OIDC, managed identity, or `AZURE_PASSWORD` via a protected env var.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1624")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
