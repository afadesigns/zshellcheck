package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1482(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — docker login --password-stdin",
			input:    `docker login --password-stdin -u user registry`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — docker push (not login)",
			input:    `docker push -p registry/image`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — docker login -p pass",
			input: `docker login -u user -p secretpass registry`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1482",
					Message: "`-p secretpass` puts the password in ps / /proc / history. Use `--password-stdin` piped from a secrets file or credential helper.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — docker login --password=pass",
			input: `docker login --password=secretpass -u user`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1482",
					Message: "`--password=secretpass` puts the password in ps / /proc / history. Use `--password-stdin` piped from a secrets file or credential helper.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — helm registry login -p",
			input: `helm registry login -u user -p secret example.com`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1482",
					Message: "`-p secret` puts the password in ps / /proc / history. Use `--password-stdin` piped from a secrets file or credential helper.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1482")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
