// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC2000(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `kubectl taint nodes $NODE key=value:NoSchedule` (gentle)",
			input:    `kubectl taint nodes $NODE key=value:NoSchedule`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `kubectl drain $NODE`",
			input:    `kubectl drain $NODE`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `kubectl taint nodes $NODE key=value:NoExecute`",
			input: `kubectl taint nodes $NODE key=value:NoExecute`,
			expected: []katas.Violation{
				{
					KataID:  "ZC2000",
					Message: "`kubectl taint nodes … :NoExecute` evicts every non-tolerating pod immediately — a typo on `--all` nodes empties the cluster. Prefer `kubectl drain $NODE` or a `:NoSchedule` taint for gentle drain.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC2000")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC2001(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `setopt EVAL_LINENO` (default on)",
			input:    `setopt EVAL_LINENO`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `unsetopt NO_EVAL_LINENO`",
			input:    `unsetopt NO_EVAL_LINENO`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `unsetopt EVAL_LINENO`",
			input: `unsetopt EVAL_LINENO`,
			expected: []katas.Violation{
				{
					KataID:  "ZC2001",
					Message: "`unsetopt EVAL_LINENO` reverts `$LINENO` inside `eval` to the outer line — errors in generated configs collapse to a single source line and stack frames past `eval` vanish. Keep on; scope via `emulate -LR sh`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `setopt NO_EVAL_LINENO`",
			input: `setopt NO_EVAL_LINENO`,
			expected: []katas.Violation{
				{
					KataID:  "ZC2001",
					Message: "`setopt NO_EVAL_LINENO` reverts `$LINENO` inside `eval` to the outer line — errors in generated configs collapse to a single source line and stack frames past `eval` vanish. Keep on; scope via `emulate -LR sh`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC2001")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC2002(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `crictl ps`",
			input:    `crictl ps`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `crictl rmi $IMAGE_ID`",
			input:    `crictl rmi $IMAGE_ID`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `crictl rmi -a`",
			input: `crictl rmi -a`,
			expected: []katas.Violation{
				{
					KataID:  "ZC2002",
					Message: "`crictl rmi -a` talks to the node CRI directly, under the kubelet — images/containers backing running pods disappear, kubelet must re-pull or re-run. Route through `kubectl drain`/`delete pod`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `crictl rm -af`",
			input: `crictl rm -af`,
			expected: []katas.Violation{
				{
					KataID:  "ZC2002",
					Message: "`crictl rm -af` talks to the node CRI directly, under the kubelet — images/containers backing running pods disappear, kubelet must re-pull or re-run. Route through `kubectl drain`/`delete pod`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC2002")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC2003(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `unsetopt KSH_ZERO_SUBSCRIPT` (default)",
			input:    `unsetopt KSH_ZERO_SUBSCRIPT`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `setopt NO_KSH_ZERO_SUBSCRIPT`",
			input:    `setopt NO_KSH_ZERO_SUBSCRIPT`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `setopt KSH_ZERO_SUBSCRIPT`",
			input: `setopt KSH_ZERO_SUBSCRIPT`,
			expected: []katas.Violation{
				{
					KataID:  "ZC2003",
					Message: "`setopt KSH_ZERO_SUBSCRIPT` stops aliasing `$arr[0]` to `$arr[1]` — every later read of `$arr[0]` silently returns empty and `arr[0]=new` stops updating the first element. Use `$arr[1]` explicitly.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `unsetopt NO_KSH_ZERO_SUBSCRIPT`",
			input: `unsetopt NO_KSH_ZERO_SUBSCRIPT`,
			expected: []katas.Violation{
				{
					KataID:  "ZC2003",
					Message: "`unsetopt NO_KSH_ZERO_SUBSCRIPT` stops aliasing `$arr[0]` to `$arr[1]` — every later read of `$arr[0]` silently returns empty and `arr[0]=new` stops updating the first element. Use `$arr[1]` explicitly.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC2003")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
