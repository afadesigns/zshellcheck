// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1580(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — go build",
			input:    `go build -o app ./cmd/app`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — go build -ldflags with version",
			input:    `go build -ldflags "-X main.Version=1.2.3" ./cmd/app`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — go build -ldflags with PASSWORD",
			input: `go build -ldflags "-X main.PASSWORD=hunter2" ./cmd/app`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1580",
					Message: "`go build -ldflags` injecting a secret bakes it into the binary. Read from os.Getenv / mounted secret file at runtime.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — go build -ldflags with API_KEY",
			input: `go build -ldflags "-X main.API_KEY=xyz" ./cmd/app`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1580",
					Message: "`go build -ldflags` injecting a secret bakes it into the binary. Read from os.Getenv / mounted secret file at runtime.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1580")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
