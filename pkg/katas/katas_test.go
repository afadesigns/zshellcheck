package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/ast"
	"github.com/afadesigns/zshellcheck/pkg/token"
)

func TestKatas(t *testing.T) {
	if len(Registry.KatasByID) == 0 {
		t.Errorf("Registry is empty")
	}
}

func TestNewKatasRegistry(t *testing.T) {
	kr := NewKatasRegistry()
	if kr.KatasByType == nil {
		t.Error("KatasByType should not be nil")
	}
	if kr.KatasByID == nil {
		t.Error("KatasByID should not be nil")
	}
}

func TestKatasRegistry_RegisterAndGetKata(t *testing.T) {
	kr := NewKatasRegistry()

	kata := Kata{
		ID:    "ZC_TEST_001",
		Title: "Test kata",
		Check: func(node ast.Node) []Violation {
			return []Violation{{KataID: "ZC_TEST_001", Message: "found"}}
		},
	}

	kr.RegisterKata(ast.IdentifierNode, kata)

	// GetKata should find it
	got, ok := kr.GetKata("ZC_TEST_001")
	if !ok {
		t.Fatal("expected to find registered kata")
	}
	if got.ID != "ZC_TEST_001" {
		t.Errorf("expected ID=ZC_TEST_001, got %s", got.ID)
	}
	// Default severity should be applied
	if got.Severity != SeverityWarning {
		t.Errorf("expected default severity SeverityWarning, got %s", got.Severity)
	}

	// GetKata should not find unknown
	_, ok = kr.GetKata("ZC_NONEXISTENT")
	if ok {
		t.Error("expected not to find nonexistent kata")
	}

	// KatasByNodeType should return the registry map
	byType := kr.KatasByNodeType()
	if len(byType) == 0 {
		t.Error("expected non-empty KatasByNodeType")
	}
}

func TestKatasRegistry_RegisterWithExplicitSeverity(t *testing.T) {
	kr := NewKatasRegistry()

	kata := Kata{
		ID:       "ZC_TEST_002",
		Title:    "Error kata",
		Severity: SeverityError,
		Check:    func(node ast.Node) []Violation { return nil },
	}

	kr.RegisterKata(ast.IdentifierNode, kata)

	got, ok := kr.GetKata("ZC_TEST_002")
	if !ok {
		t.Fatal("expected to find registered kata")
	}
	if got.Severity != SeverityError {
		t.Errorf("expected severity SeverityError, got %s", got.Severity)
	}
}

func TestKatasRegistry_Check(t *testing.T) {
	kr := NewKatasRegistry()

	kr.RegisterKata(ast.IdentifierNode, Kata{
		ID:    "ZC_CHK_001",
		Title: "Check test",
		Check: func(node ast.Node) []Violation {
			return []Violation{
				{KataID: "ZC_CHK_001", Message: "violation found", Line: 1, Column: 1},
			}
		},
	})

	kr.RegisterKata(ast.IdentifierNode, Kata{
		ID:    "ZC_CHK_002",
		Title: "Check test 2",
		Check: func(node ast.Node) []Violation {
			return []Violation{
				{KataID: "ZC_CHK_002", Message: "another violation"},
			}
		},
	})

	node := &ast.Identifier{
		Token: token.Token{Type: token.IDENT, Literal: "x"},
		Value: "x",
	}

	// Check with no disabled katas
	violations := kr.Check(node, nil)
	if len(violations) != 2 {
		t.Fatalf("expected 2 violations, got %d", len(violations))
	}
	// Default severity should be applied when violation has no level
	if violations[0].Level != SeverityWarning {
		t.Errorf("expected default level SeverityWarning, got %s", violations[0].Level)
	}

	// Check with one disabled kata
	violations = kr.Check(node, []string{"ZC_CHK_001"})
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation with disabled kata, got %d", len(violations))
	}
	if violations[0].KataID != "ZC_CHK_002" {
		t.Errorf("expected ZC_CHK_002, got %s", violations[0].KataID)
	}

	// Check with node type that has no katas
	intNode := &ast.IntegerLiteral{
		Token: token.Token{Type: token.INT, Literal: "42"},
		Value: 42,
	}
	violations = kr.Check(intNode, nil)
	if len(violations) != 0 {
		t.Errorf("expected 0 violations for unregistered node type, got %d", len(violations))
	}
}

func TestKatasRegistry_CheckPreservesExplicitLevel(t *testing.T) {
	kr := NewKatasRegistry()

	kr.RegisterKata(ast.IdentifierNode, Kata{
		ID:       "ZC_LVL_001",
		Title:    "Level test",
		Severity: SeverityWarning,
		Check: func(node ast.Node) []Violation {
			return []Violation{
				{KataID: "ZC_LVL_001", Message: "has level", Level: SeverityError},
				{KataID: "ZC_LVL_001", Message: "no level"},
			}
		},
	})

	node := &ast.Identifier{
		Token: token.Token{Type: token.IDENT, Literal: "x"},
		Value: "x",
	}

	violations := kr.Check(node, nil)
	if len(violations) != 2 {
		t.Fatalf("expected 2 violations, got %d", len(violations))
	}
	// First violation has explicit SeverityError level - should be preserved
	if violations[0].Level != SeverityError {
		t.Errorf("expected explicit SeverityError level preserved, got %s", violations[0].Level)
	}
	// Second violation has no level - should get kata's default (SeverityWarning)
	if violations[1].Level != SeverityWarning {
		t.Errorf("expected default SeverityWarning level, got %s", violations[1].Level)
	}
}
