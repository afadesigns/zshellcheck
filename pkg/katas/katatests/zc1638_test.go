package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1638(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — non-secret build-arg",
			input:    `docker build --build-arg VERSION=1.0 -t app .`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — BuildKit secret",
			input:    `docker build --secret id=dbpass,src=/run/secrets/db -t app .`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — docker build --build-arg PASSWORD=s3cret",
			input: `docker build --build-arg PASSWORD=s3cret -t app .`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1638",
					Message: "`docker build --build-arg PASSWORD=s3cret` bakes the secret into the image layer metadata. Use `--secret id=NAME,src=PATH` (BuildKit) or a multi-stage build.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — podman build --build-arg API_KEY=$KEY",
			input: `podman build --build-arg API_KEY=$KEY -t app .`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1638",
					Message: "`podman build --build-arg API_KEY=$KEY` bakes the secret into the image layer metadata. Use `--secret id=NAME,src=PATH` (BuildKit) or a multi-stage build.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1638")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
