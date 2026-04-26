// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1069(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid local in function",
			input:    `my_func() { local x=1; }`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid typeset global",
			input:    `typeset x=1`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid local global",
			input: `local x=1`,
			expected: []katas.Violation{
				{
					KataID: "ZC1069",
					Message: "`local` can only be used inside functions. " +
						"Use `typeset`, `declare`, or just assignment for global variables.",
					Line:   1,
					Column: 1,
				},
			},
		},
		{
			name:  "invalid local in if block (global)",
			input: `if true; then local x=1; fi`,
			expected: []katas.Violation{
				{
					KataID: "ZC1069",
					Message: "`local` can only be used inside functions. " +
						"Use `typeset`, `declare`, or just assignment for global variables.",
					Line:   1,
					Column: 15,
				},
			},
		},
		{
			name:     "valid local in nested function",
			input:    `outer() { inner() { local x=1; }; }`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid local in subshell (global)",
			input: `( local x=1 )`,
			expected: []katas.Violation{
				{
					KataID: "ZC1069",
					Message: "`local` can only be used inside functions. " +
						"Use `typeset`, `declare`, or just assignment for global variables.",
					Line:   1,
					Column: 3,
				},
			},
		},
		{
			name:     "valid local in function keyword",
			input:    "function myfunc { local x=1; }",
			expected: []katas.Violation{},
		},
		{
			name:  "invalid local in while loop (global)",
			input: `while true; do local x=1; done`,
			expected: []katas.Violation{
				{
					KataID: "ZC1069",
					Message: "`local` can only be used inside functions. " +
						"Use `typeset`, `declare`, or just assignment for global variables.",
					Line:   1,
					Column: 16,
				},
			},
		},
		{
			name:  "invalid local in for loop (global)",
			input: `for i in a b c; do local x=1; done`,
			expected: []katas.Violation{
				{
					KataID: "ZC1069",
					Message: "`local` can only be used inside functions. " +
						"Use `typeset`, `declare`, or just assignment for global variables.",
					Line:   1,
					Column: 20,
				},
			},
		},
		{
			name:  "invalid local in case (global)",
			input: "case $x in\na) local y=1;;\nesac",
			expected: []katas.Violation{
				{
					KataID: "ZC1069",
					Message: "`local` can only be used inside functions. " +
						"Use `typeset`, `declare`, or just assignment for global variables.",
					Line:   2,
					Column: 4,
				},
			},
		},
		{
			name:     "regular echo command",
			input:    `echo hello world`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1069")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
