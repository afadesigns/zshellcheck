package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1206(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid crontab file",
			input:    `crontab /tmp/cron.txt`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid crontab -l",
			input:    `crontab -l`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid crontab -e",
			input: `crontab -e`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1206",
					Message: "Avoid `crontab -e` in scripts — it opens an interactive editor. Use `crontab file` or `echo '...' | crontab -` for programmatic cron management.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1206")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
