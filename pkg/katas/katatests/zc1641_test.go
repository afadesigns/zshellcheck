package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1641(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — --from-file",
			input:    `kubectl create secret generic mysec --from-file=password=/run/secrets/pw`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — --from-env-file",
			input:    `kubectl create secret generic mysec --from-env-file=/run/secrets/env`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — --from-literal=password=X",
			input: `kubectl create secret generic mysec --from-literal=password=hunter2`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1641",
					Message: "`kubectl create secret --from-literal=password=hunter2` puts the secret in argv — visible via `ps`. Use `--from-file=KEY=PATH` / `--from-env-file=PATH`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — --docker-password=$PW",
			input: `kubectl create secret docker-registry reg --docker-password=$PW --docker-username=u`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1641",
					Message: "`kubectl create secret --docker-password=$PW` puts the secret in argv — visible via `ps`. Use `--from-file=KEY=PATH` / `--from-env-file=PATH`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1641")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
