package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1982(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `ipcs -a` (list)",
			input:    `ipcs -a`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `ipcrm -m $SHMID` (scoped)",
			input:    `ipcrm -m $SHMID`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `ipcrm -a`",
			input: `ipcrm -a`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1982",
					Message: "`ipcrm -a` deletes every SysV shm/sem/mqueue object — Postgres/Oracle/shm-based services lose their backing store mid-transaction. Scope with `-m`/`-s`/`-q` on the specific ID after checking `ipcs -a`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1982")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
