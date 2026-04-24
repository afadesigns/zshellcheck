package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1803(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `mysql --ssl-mode=VERIFY_IDENTITY -h db -u u`",
			input:    `mysql --ssl-mode=VERIFY_IDENTITY -h db -u u`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `psql postgresql://u@db/mydb`",
			input:    `psql postgresql://u@db/mydb`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `mysql --skip-ssl -h db -u u`",
			input: `mysql --skip-ssl -h db -u u`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1803",
					Message: "`mysql --skip-ssl` disables TLS — login handshake and queries travel in plaintext. Use `--ssl-mode=VERIFY_IDENTITY` (MySQL) / `sslmode=verify-full` (psql) with a pinned CA.",
					Line:    1,
					Column:  8,
				},
			},
		},
		{
			name:  "invalid — `psql \"host=db sslmode=disable user=u\"`",
			input: `psql "host=db sslmode=disable user=u"`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1803",
					Message: "`psql host=db sslmode=disable user=u` disables TLS — login handshake and queries travel in plaintext. Use `--ssl-mode=VERIFY_IDENTITY` (MySQL) / `sslmode=verify-full` (psql) with a pinned CA.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1803")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
