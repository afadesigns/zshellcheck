package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1501(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — docker compose up",
			input:    `docker compose up -d`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — docker-compose up",
			input: `docker-compose up -d`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1501",
					Message: "`docker-compose` is the deprecated Python V1 binary. Use `docker compose` (space-separated subcommand) for the bundled V2 plugin.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — docker-compose down",
			input: `docker-compose down --volumes`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1501",
					Message: "`docker-compose` is the deprecated Python V1 binary. Use `docker compose` (space-separated subcommand) for the bundled V2 plugin.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1501")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
