package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1701(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — dpkg -l (list)",
			input:    `dpkg -l`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — apt install from repo",
			input:    `apt install mypkg`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — dpkg -i local.deb",
			input: `dpkg -i local.deb`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1701",
					Message: "`dpkg -i FILE.deb` runs the package without signature verification — `sha256sum -c` or `debsig-verify` the file first, or install via `apt install` from a signed repo.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — dpkg -i /tmp/download.deb",
			input: `dpkg -i /tmp/download.deb`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1701",
					Message: "`dpkg -i FILE.deb` runs the package without signature verification — `sha256sum -c` or `debsig-verify` the file first, or install via `apt install` from a signed repo.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1701")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
