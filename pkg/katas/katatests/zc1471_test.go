package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1471(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — kubectl get pods",
			input:    `kubectl get pods`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — helm install nginx bitnami/nginx",
			input:    `helm install nginx bitnami/nginx`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — kubectl get pods --insecure-skip-tls-verify",
			input: `kubectl get pods --insecure-skip-tls-verify`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1471",
					Message: "`--insecure-skip-tls-verify` turns off API-server certificate verification — MITM steals every secret. Fix the CA bundle instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — helm install --insecure-skip-tls-verify=true foo",
			input: `helm install --insecure-skip-tls-verify=true foo`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1471",
					Message: "`--insecure-skip-tls-verify` turns off API-server certificate verification — MITM steals every secret. Fix the CA bundle instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — oc login --insecure-skip-tls-verify=true",
			input: `oc login --insecure-skip-tls-verify=true https://cluster`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1471",
					Message: "`--insecure-skip-tls-verify` turns off API-server certificate verification — MITM steals every secret. Fix the CA bundle instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1471")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
