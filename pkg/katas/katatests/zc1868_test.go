package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1868(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `gcloud config set compute/zone us-central1-a`",
			input:    `gcloud config set compute/zone us-central1-a`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `gcloud config set auth/disable_ssl_validation false`",
			input:    `gcloud config set auth/disable_ssl_validation false`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `gcloud config set auth/disable_ssl_validation true`",
			input: `gcloud config set auth/disable_ssl_validation true`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1868",
					Message: "`gcloud config set auth/disable_ssl_validation true` turns off TLS for every later `gcloud` call — service-account tokens and deploys become interceptable. Unset it; pin custom CAs via `core/custom_ca_certs_file`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1868")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
