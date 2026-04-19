package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1796(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `pg_restore -d mydb backup.dump`",
			input:    `pg_restore -d mydb backup.dump`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `pg_restore --list backup.dump` (TOC only)",
			input:    `pg_restore --list backup.dump`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `pg_restore -c -d mydb backup.dump`",
			input: `pg_restore -c -d mydb backup.dump`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1796",
					Message: "`pg_restore -c` drops every object in the target DB before recreating from the archive — stale or wrong-target dump silently loses data. Restore into a fresh DB (`createdb new && pg_restore -d new`), or snapshot first.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1796")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
