package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1687(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — snap install strict",
			input:    `snap install hello-world`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — snap list",
			input:    `snap list`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — snap install --classic",
			input: `snap install code --classic`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1687",
					Message: "`snap install --classic` drops AppArmor / cgroup / seccomp sandbox — find a strict snap or a distro-package alternative, or document why this specific snap needs it.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — snap install --devmode",
			input: `snap install pkg --devmode`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1687",
					Message: "`snap install --devmode` logs confinement violations instead of blocking — find a strict snap or a distro-package alternative, or document why this specific snap needs it.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1687")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
