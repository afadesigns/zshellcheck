// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1630(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — localhost bind",
			input:    `php -S 127.0.0.1:8000 -t public`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — php script (not -S)",
			input:    `php artisan migrate`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — php -S 0.0.0.0:8000",
			input: `php -S 0.0.0.0:8000`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1630",
					Message: "`php -S 0.0.0.0:8000` binds the dev server to every interface — unauthenticated access to the working directory. Use `127.0.0.1:PORT` locally, nginx / caddy for external exposure.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — php -S [::]:8080",
			input: `php -S [::]:8080 -t public`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1630",
					Message: "`php -S [::]:8080` binds the dev server to every interface — unauthenticated access to the working directory. Use `127.0.0.1:PORT` locally, nginx / caddy for external exposure.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1630")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
