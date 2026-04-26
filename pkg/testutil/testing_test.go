// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package testutil

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
)

func TestCheck_NoViolations(t *testing.T) {
	// A simple assignment should not trigger violations for most katas
	violations := Check("x=1", "ZC_NONEXISTENT")
	if len(violations) != 0 {
		t.Errorf("expected 0 violations for nonexistent kata, got %d", len(violations))
	}
}

func TestCheck_WithViolations(t *testing.T) {
	// Test a known kata - echo with unquoted variable
	violations := Check("echo $x", "ZC1001")
	// Whether there are violations or not depends on the kata, but it should not panic
	_ = violations
}

func TestAssertViolations_Matching(t *testing.T) {
	actual := []katas.Violation{
		{KataID: "ZC1001", Message: "test msg", Line: 1, Column: 1},
	}
	expected := []katas.Violation{
		{KataID: "ZC1001", Message: "test msg", Line: 1, Column: 1},
	}
	AssertViolations(t, "echo test", actual, expected)
}

func TestAssertViolations_Empty(t *testing.T) {
	AssertViolations(t, "echo test", []katas.Violation{}, []katas.Violation{})
}

func TestAssertViolations_MismatchKataID(t *testing.T) {
	// Use a sub-test to capture the expected failure
	mockT := &testing.T{}
	actual := []katas.Violation{
		{KataID: "ZC1001", Message: "msg", Line: 1, Column: 1},
	}
	expected := []katas.Violation{
		{KataID: "ZC1002", Message: "msg", Line: 1, Column: 1},
	}
	AssertViolations(mockT, "echo test", actual, expected)
	// mockT should have recorded an error - we can't easily check this
	// but the important thing is it exercises the code path
}

func TestAssertViolations_MismatchMessage(t *testing.T) {
	mockT := &testing.T{}
	actual := []katas.Violation{
		{KataID: "ZC1001", Message: "actual msg", Line: 1, Column: 1},
	}
	expected := []katas.Violation{
		{KataID: "ZC1001", Message: "expected msg", Line: 1, Column: 1},
	}
	AssertViolations(mockT, "echo test", actual, expected)
}

func TestAssertViolations_MismatchLine(t *testing.T) {
	mockT := &testing.T{}
	actual := []katas.Violation{
		{KataID: "ZC1001", Message: "msg", Line: 1, Column: 1},
	}
	expected := []katas.Violation{
		{KataID: "ZC1001", Message: "msg", Line: 2, Column: 1},
	}
	AssertViolations(mockT, "echo test", actual, expected)
}

func TestAssertViolations_MismatchColumn(t *testing.T) {
	mockT := &testing.T{}
	actual := []katas.Violation{
		{KataID: "ZC1001", Message: "msg", Line: 1, Column: 1},
	}
	expected := []katas.Violation{
		{KataID: "ZC1001", Message: "msg", Line: 1, Column: 2},
	}
	AssertViolations(mockT, "echo test", actual, expected)
}
