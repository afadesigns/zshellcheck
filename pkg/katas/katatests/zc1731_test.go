// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1731(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `curl URL --data 'name=John'` (non-secret key)",
			input:    `curl URL --data 'name=John'`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `curl URL --data @secret.txt` (file reference)",
			input:    `curl URL --data @secret.txt`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `curl URL --data-binary @-` (stdin sentinel)",
			input:    `curl URL --data-binary @-`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `curl URL -d 'password=hunter2'`",
			input: `curl URL -d 'password=hunter2'`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1731",
					Message: "`curl -d 'password=hunter2'` puts secret-keyed POST body (`password=…`) in argv — visible in `ps`, `/proc`, history. Read the value from a file with `--data @PATH` or `--data-binary @-` piped from a secrets store.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `curl URL --data-urlencode 'token=ABC123'`",
			input: `curl URL --data-urlencode 'token=ABC123'`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1731",
					Message: "`curl --data-urlencode 'token=ABC123'` puts secret-keyed POST body (`token=…`) in argv — visible in `ps`, `/proc`, history. Read the value from a file with `--data @PATH` or `--data-binary @-` piped from a secrets store.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `wget URL --post-data='api_key=ABC123'`",
			input: `wget URL --post-data='api_key=ABC123'`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1731",
					Message: "`wget --post-data='api_key=ABC123'` puts secret-keyed POST body (`api_key=…`) in argv — visible in `ps`, `/proc`, history. Read the value from a file with `--data @PATH` or `--data-binary @-` piped from a secrets store.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1731")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
