// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1767(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `mongod --bind_ip 127.0.0.1`",
			input:    `mongod --bind_ip 127.0.0.1`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `mongod --bind_ip 0.0.0.0`",
			input: `mongod --bind_ip 0.0.0.0`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1767",
					Message: "`mongod --bind_ip 0.0.0.0` exposes MongoDB on every interface — 2017 ransomware-wave target. Bind to `127.0.0.1` or a private-network IP, enable `--auth`, firewall port 27017.",
					Line:    1,
					Column:  9,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1767")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
