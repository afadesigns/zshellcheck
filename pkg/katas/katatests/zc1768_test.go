package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1768(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `sqlcmd -S server -U user -P` (prompt)",
			input:    `sqlcmd -S server -U user -P`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `sqlcmd -S server -E` (Windows auth, no password)",
			input:    `sqlcmd -S server -E`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `sqlcmd -S server -U user -P hunter2`",
			input: `sqlcmd -S server -U user -P hunter2`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1768",
					Message: "`sqlcmd -P hunter2` puts the SQL Server password in argv — visible in `ps`, `/proc`, history, SQL Server audit. Use `-P` with no arg (prompt) or `SQLCMDPASSWORD` env var.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `bcp mydb in data.csv -U user -P hunter2`",
			input: `bcp mydb in data.csv -U user -P hunter2`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1768",
					Message: "`bcp -P hunter2` puts the SQL Server password in argv — visible in `ps`, `/proc`, history, SQL Server audit. Use `-P` with no arg (prompt) or `SQLCMDPASSWORD` env var.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1768")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
