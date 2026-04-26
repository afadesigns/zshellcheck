// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1782(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `flatpak remote-add flathub https://flathub.org/repo/flathub.flatpakrepo`",
			input:    `flatpak remote-add flathub https://flathub.org/repo/flathub.flatpakrepo`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `flatpak install flathub org.gimp.GIMP`",
			input:    `flatpak install flathub org.gimp.GIMP`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `flatpak remote-add --no-gpg-verify local /srv/repo`",
			input: `flatpak remote-add --no-gpg-verify local /srv/repo`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1782",
					Message: "`flatpak remote-add --no-gpg-verify` disables signature verification — updates from this remote are accepted with only HTTPS as identity. Sign the repo (`ostree gpg-sign`) and import the key with `--gpg-import=KEYFILE`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `flatpak remote-modify --gpg-verify=false myrepo`",
			input: `flatpak remote-modify --gpg-verify=false myrepo`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1782",
					Message: "`flatpak remote-modify --gpg-verify=false` disables signature verification — updates from this remote are accepted with only HTTPS as identity. Sign the repo (`ostree gpg-sign`) and import the key with `--gpg-import=KEYFILE`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1782")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
