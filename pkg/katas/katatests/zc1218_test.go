package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1218(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid useradd with nologin",
			input:    `useradd --system --shell /sbin/nologin myservice`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid regular useradd",
			input:    `useradd newuser`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid system useradd without nologin",
			input: `useradd -r myservice`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1218",
					Message: "Add `--shell /sbin/nologin` when creating system accounts with `useradd`. This prevents interactive login for service accounts.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1218")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
