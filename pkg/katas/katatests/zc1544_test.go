package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1544(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — dnf install curl",
			input:    `dnf install curl`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — dnf copr enable user/repo",
			input: `dnf copr enable user/repo`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1544",
					Message: "`dnf copr enable` pulls from a single-contributor repo — no distro security team. Pin the build, verify key fingerprint, mirror internally.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — add-apt-repository ppa:user/repo",
			input: `add-apt-repository ppa:user/repo`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1544",
					Message: "`add-apt-repository` pulls from a single-contributor repo — no distro security team. Pin the build, verify key fingerprint, mirror internally.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1544")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
