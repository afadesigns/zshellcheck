package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1053(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "grep -q in condition",
			input:    `if grep -q pattern file; then echo found; fi`,
			expected: []katas.Violation{},
		},
		{
			name:     "no grep in condition",
			input:    `if true; then echo yes; fi`,
			expected: []katas.Violation{},
		},
		{
			name:  "grep without -q in condition",
			input: `if grep pattern file; then echo found; fi`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1053",
					Message: "Silence `grep` output in conditions. Use `grep -q` or redirect to `/dev/null`.",
					Line:    1,
					Column:  4,
				},
			},
		},
		{
			name:     "grep -rq combined flags",
			input:    `if grep -rq pattern dir; then echo found; fi`,
			expected: []katas.Violation{},
		},
		{
			name:     "grep --quiet in condition",
			input:    `if grep --quiet pattern file; then echo found; fi`,
			expected: []katas.Violation{},
		},
		{
			name:  "egrep without -q in condition",
			input: `if egrep pattern file; then echo found; fi`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1053",
					Message: "Silence `grep` output in conditions. Use `grep -q` or redirect to `/dev/null`.",
					Line:    1,
					Column:  4,
				},
			},
		},
		{
			name:  "grep piped in condition",
			input: `if echo test | grep pattern; then echo found; fi`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1053",
					Message: "Silence `grep` output in conditions. Use `grep -q` or redirect to `/dev/null`.",
					Line:    1,
					Column:  16,
				},
			},
		},
		{
			name:  "grep without -q in while condition",
			input: `while grep pattern file; do echo loop; done`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1053",
					Message: "Silence `grep` output in conditions. Use `grep -q` or redirect to `/dev/null`.",
					Line:    1,
					Column:  7,
				},
			},
		},
		{
			name:     "grep -q in while condition",
			input:    `while grep -q pattern file; do echo loop; done`,
			expected: []katas.Violation{},
		},
		{
			name:     "grep redirected to /dev/null",
			input:    `if grep pattern file > /dev/null; then echo found; fi`,
			expected: []katas.Violation{},
		},
		{
			name:     "grep --silent in condition",
			input:    `if grep --silent pattern file; then echo found; fi`,
			expected: []katas.Violation{},
		},
		{
			name:  "negated grep without -q",
			input: `if ! grep pattern file; then echo not_found; fi`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1053",
					Message: "Silence `grep` output in conditions. Use `grep -q` or redirect to `/dev/null`.",
					Line:    1,
					Column:  6,
				},
			},
		},
		{
			name:  "fgrep without -q",
			input: `if fgrep pattern file; then echo found; fi`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1053",
					Message: "Silence `grep` output in conditions. Use `grep -q` or redirect to `/dev/null`.",
					Line:    1,
					Column:  4,
				},
			},
		},
		{
			name:  "zgrep without -q",
			input: `if zgrep pattern file; then echo found; fi`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1053",
					Message: "Silence `grep` output in conditions. Use `grep -q` or redirect to `/dev/null`.",
					Line:    1,
					Column:  4,
				},
			},
		},
		{
			name:  "grep in && condition without -q",
			input: `if true && grep pattern file; then echo found; fi`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1053",
					Message: "Silence `grep` output in conditions. Use `grep -q` or redirect to `/dev/null`.",
					Line:    1,
					Column:  12,
				},
			},
		},
		{
			name:     "grep left side of pipe is silenced",
			input:    `if grep pattern file | wc -l; then echo found; fi`,
			expected: []katas.Violation{},
		},
		{
			name:     "non-grep command in condition",
			input:    `if ls /tmp; then echo ok; fi`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1053")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
