// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1549(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — unzip to /tmp/stage",
			input:    `unzip foo.zip -d /tmp/stage`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — unzip -o to /opt/app",
			input:    `unzip -o foo.zip -d /opt/app`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — unzip -d /",
			input: `unzip foo.zip -d /`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1549",
					Message: "`unzip -d /` extracts into a system path — any archive entry overwrites matching system file. Stage, inspect, copy.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — unzip -o file -d /boot",
			input: `unzip -o foo.zip -d /boot`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1549",
					Message: "`unzip -d /boot` extracts into a system path — any archive entry overwrites matching system file. Stage, inspect, copy.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1549")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
