// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1819(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `xattr file` (read only)",
			input:    `xattr file`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `xattr -d com.apple.metadata:kMDLabel_xxx file` (unrelated xattr)",
			input:    `xattr -d com.apple.metadata:kMDLabel_xxx file`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `xattr -d com.apple.quarantine /Applications/MyApp.app`",
			input: `xattr -d com.apple.quarantine /Applications/MyApp.app`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1819",
					Message: "`xattr -d com.apple.quarantine` / `-cr` strips the macOS Gatekeeper quarantine — the binary runs with no signature / notarization check. Verify with `codesign --verify` and `spctl --assess --type execute` before stripping.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `xattr -cr $HOME/Downloads`",
			input: `xattr -cr $HOME/Downloads`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1819",
					Message: "`xattr -d com.apple.quarantine` / `-cr` strips the macOS Gatekeeper quarantine — the binary runs with no signature / notarization check. Verify with `codesign --verify` and `spctl --assess --type execute` before stripping.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1819")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
