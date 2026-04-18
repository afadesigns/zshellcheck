package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1670(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — setsebool -P httpd_can_network_connect on (not in dangerous list)",
			input:    `setsebool -P httpd_can_network_connect on`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — setsebool without -P (session only)",
			input:    `setsebool httpd_execmem 1`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — setsebool -P httpd_execmem 1",
			input: `setsebool -P httpd_execmem 1`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1670",
					Message: "`setsebool -P httpd_execmem 1` persistently relaxes SELinux memory-protection policy — fix the binary instead (`execstack -c`, relabel with `chcon`, or change the domain).",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — setsebool -P allow_execstack on",
			input: `setsebool -P allow_execstack on`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1670",
					Message: "`setsebool -P allow_execstack on` persistently relaxes SELinux memory-protection policy — fix the binary instead (`execstack -c`, relabel with `chcon`, or change the domain).",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1670")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
