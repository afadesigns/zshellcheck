// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1727(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `curl URL --proxy http://PROXY:8080` (no creds in URL)",
			input:    `curl URL --proxy http://PROXY:8080`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `curl URL` (no proxy)",
			input:    `curl URL`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `wget URL` (no proxy creds)",
			input:    `wget URL`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `curl URL --proxy http://USER:PASS@PROXY:8080`",
			input: `curl URL --proxy http://USER:PASS@PROXY:8080`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1727",
					Message: "`curl --proxy http://USER:PASS@PROXY:8080` puts proxy credentials in argv — visible in `ps`, `/proc`, history. Move them into `~/.curlrc` / `~/.netrc` (chmod 600) or `~/.wgetrc`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `curl URL -x http://USER:PASS@PROXY:8080`",
			input: `curl URL -x http://USER:PASS@PROXY:8080`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1727",
					Message: "`curl --proxy http://USER:PASS@PROXY:8080` puts proxy credentials in argv — visible in `ps`, `/proc`, history. Move them into `~/.curlrc` / `~/.netrc` (chmod 600) or `~/.wgetrc`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `wget URL --proxy-password=hunter2`",
			input: `wget URL --proxy-password=hunter2`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1727",
					Message: "`wget --proxy-password=hunter2` puts proxy credentials in argv — visible in `ps`, `/proc`, history. Move them into `~/.curlrc` / `~/.netrc` (chmod 600) or `~/.wgetrc`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1727")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
