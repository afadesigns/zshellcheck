// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1884(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `curl https://api.example/public`",
			input:    `curl https://api.example/public`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `curl -H \"Authorization: Bearer $T\" https://api.example/private`",
			input:    `curl -H "Authorization: Bearer $T" https://api.example/private`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `curl https://api/thing?apikey=abc`",
			input: `curl https://api.example/thing?apikey=abc`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1884",
					Message: "`curl https://api.example/thing?apikey=abc` carries `apikey...` in the URL query — logged by every proxy, CDN, and server access log along the path. Move credentials to `-H \"Authorization: Bearer \"$TOKEN\"` or a POST body.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `curl -X POST https://api.example/auth?token=xyz`",
			input: `curl -X POST https://api.example/auth?token=xyz`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1884",
					Message: "`curl https://api.example/auth?token=xyz` carries `token...` in the URL query — logged by every proxy, CDN, and server access log along the path. Move credentials to `-H \"Authorization: Bearer \"$TOKEN\"` or a POST body.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1884")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
