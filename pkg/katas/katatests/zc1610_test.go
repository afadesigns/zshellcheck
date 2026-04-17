package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1610(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — curl to temp path",
			input:    `curl -fsSL -o /tmp/download.tar.gz https://example.com/x.tar.gz`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — wget to user home",
			input:    `wget -O $HOME/.local/bin/tool https://example.com/tool`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — curl -o /etc/config",
			input: `curl -fsSL -o /etc/myapp/config.yaml https://example.com/config`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1610",
					Message: "`curl -o /etc/myapp/config.yaml` writes an HTTP response straight into a system path — a compromised URL replaces the target. Download to a temp file, verify, then `install` into place.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — wget -O /usr/local/bin/tool",
			input: `wget -O /usr/local/bin/tool https://example.com/tool`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1610",
					Message: "`wget -O /usr/local/bin/tool` writes an HTTP response straight into a system path — a compromised URL replaces the target. Download to a temp file, verify, then `install` into place.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — wget -O /lib/x.so",
			input: `wget -O /lib/evil.so https://example.com/evil.so`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1610",
					Message: "`wget -O /lib/evil.so` writes an HTTP response straight into a system path — a compromised URL replaces the target. Download to a temp file, verify, then `install` into place.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1610")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
