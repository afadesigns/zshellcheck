package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1545(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — docker image prune",
			input:    `docker image prune`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — docker system prune (no -a / --volumes)",
			input:    `docker system prune -f`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — docker system prune -af --volumes",
			input: `docker system prune -af --volumes`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1545",
					Message: "`docker system prune` with `-a`/`--volumes` drops unused volumes — stopped stacks lose their databases. Scope the prune.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — docker volume prune -a",
			input: `docker volume prune -a`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1545",
					Message: "`docker volume prune` with `-a`/`--volumes` drops unused volumes — stopped stacks lose their databases. Scope the prune.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1545")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
