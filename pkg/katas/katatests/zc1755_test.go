package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1755(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `gcloud sql users create ... --prompt-for-password`",
			input:    `gcloud sql users create myuser --instance myinst --prompt-for-password`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `gcloud sql users list --instance myinst`",
			input:    `gcloud sql users list --instance myinst`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `gcloud sql users create ... --password PASS`",
			input: `gcloud sql users create myuser --instance myinst --password hunter2`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1755",
					Message: "`gcloud sql users create --password hunter2` puts the Cloud SQL password in argv — visible in `ps`, `/proc`, history, and Cloud Audit Logs. Use `--prompt-for-password` or call the SQL Admin API with a body file.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `gcloud sql users set-password ... --password=PASS`",
			input: `gcloud sql users set-password myuser --instance myinst --password=hunter2`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1755",
					Message: "`gcloud sql users set-password --password=hunter2` puts the Cloud SQL password in argv — visible in `ps`, `/proc`, history, and Cloud Audit Logs. Use `--prompt-for-password` or call the SQL Admin API with a body file.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1755")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
