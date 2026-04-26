// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1745(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `poetry publish --repository myrepo` (no password)",
			input:    `poetry publish --repository myrepo`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `twine upload dist/*` (token via env)",
			input:    `twine upload dist/*`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `poetry publish --username u --password hunter2`",
			input: `poetry publish --username u --password hunter2`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1745",
					Message: "`poetry publish --password hunter2` puts the registry password in argv — visible in `ps`, `/proc`, history. Use env vars (`POETRY_PYPI_TOKEN_<NAME>`, `TWINE_PASSWORD`) or a `0600` `~/.pypirc` file.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `twine upload -u u -p hunter2 dist/*`",
			input: `twine upload -u u -p hunter2 dist/*`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1745",
					Message: "`twine upload --password hunter2` puts the registry password in argv — visible in `ps`, `/proc`, history. Use env vars (`POETRY_PYPI_TOKEN_<NAME>`, `TWINE_PASSWORD`) or a `0600` `~/.pypirc` file.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `poetry publish --password=hunter2`",
			input: `poetry publish --password=hunter2`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1745",
					Message: "`poetry publish --password=hunter2` puts the registry password in argv — visible in `ps`, `/proc`, history. Use env vars (`POETRY_PYPI_TOKEN_<NAME>`, `TWINE_PASSWORD`) or a `0600` `~/.pypirc` file.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1745")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
