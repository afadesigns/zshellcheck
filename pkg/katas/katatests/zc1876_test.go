package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1876(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `cargo publish`",
			input:    `cargo publish`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `cargo publish --dry-run`",
			input:    `cargo publish --dry-run`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `cargo publish --allow-dirty`",
			input: `cargo publish --allow-dirty`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1876",
					Message: "`cargo publish --allow-dirty` uploads a tarball snapshot of the dirty working tree — debug prints and local-only patches end up on crates.io for a version that cannot be replaced. Commit first; `--dry-run` to rehearse.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1876")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
