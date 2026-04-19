package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1914(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `curl https://api.example/resource`",
			input:    `curl https://api.example/resource`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `curl --resolve api:443:10.0.0.1 https://api/`",
			input:    `curl --resolve api:443:10.0.0.1 https://api/`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `curl https://api/ --doh-url=https://doh/dns-query`",
			input: `curl https://api/ --doh-url=https://doh/dns-query`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1914",
					Message: "`curl --doh-url` bypasses the host's resolver chain — `/etc/hosts`, `systemd-resolved`, split-horizon DNS — so the request lands at an IP the operator did not vet. Drop the flag or pair it with `--resolve` pinning.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `curl https://api/ --dns-servers=1.1.1.1`",
			input: `curl https://api/ --dns-servers=1.1.1.1`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1914",
					Message: "`curl --dns-servers` bypasses the host's resolver chain — `/etc/hosts`, `systemd-resolved`, split-horizon DNS — so the request lands at an IP the operator did not vet. Drop the flag or pair it with `--resolve` pinning.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1914")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
