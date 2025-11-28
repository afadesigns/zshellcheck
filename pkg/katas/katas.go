package katas

import (
	"fmt"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

// Severity defines the severity of a violation.
type Severity string

const (
	Error   Severity = "Error"
	Warning Severity = "Warning"
	Info    Severity = "Info"
)

// Violation represents a found violation in the code.
type Violation struct {
	KataID  string
	Message string
	Line    int
	Column  int
	Level   Severity
}

// Kata represents a single linting rule.
type Kata struct {
	ID          string
	Title       string
	Description string
	Severity    Severity
	Check       func(node ast.Node) []Violation
}

// KatasRegistry is a registry for all available Katas.
type KatasRegistry struct {
	KatasByType map[string][]Kata
	KatasByID   map[string]Kata
}

// NewKatasRegistry creates a new KatasRegistry.
func NewKatasRegistry() *KatasRegistry {
	return &KatasRegistry{
		KatasByType: make(map[string][]Kata),
		KatasByID:   make(map[string]Kata),
	}
}

// RegisterKata registers a new Kata.
func (kr *KatasRegistry) RegisterKata(nodeType ast.Node, kata Kata) {
	if kata.Severity == "" {
		kata.Severity = Warning
	}
	key := fmt.Sprintf("%T", nodeType)
	kr.KatasByType[key] = append(kr.KatasByType[key], kata)
	kr.KatasByID[kata.ID] = kata
}

// GetKata returns a Kata by its ID.
func (kr *KatasRegistry) GetKata(id string) (Kata, bool) {
	kata, ok := kr.KatasByID[id]
	return kata, ok
}

// KatasByNodeType returns all registered Katas grouped by node type.
func (kr *KatasRegistry) KatasByNodeType() map[string][]Kata {
	return kr.KatasByType
}

func (kr *KatasRegistry) Check(node ast.Node, disabledKatas []string) []Violation {
	var violations []Violation
	key := fmt.Sprintf("%T", node)
	if katasForNode, ok := kr.KatasByType[key]; ok {
		for _, kata := range katasForNode {
			// Check if disabled
			disabled := false
			for _, d := range disabledKatas {
				if d == kata.ID {
					disabled = true
					break
				}
			}
			if !disabled {
				vs := kata.Check(node)
				for i := range vs {
					if vs[i].Level == "" {
						vs[i].Level = kata.Severity
					}
				}
				violations = append(violations, vs...)
			}
		}
	}
	return violations
}

// DefaultKatasRegistry is the default registry.
var DefaultKatasRegistry = NewKatasRegistry()

// RegisterKata registers a new Kata with the default registry.
func RegisterKata(nodeType ast.Node, kata Kata) {
	DefaultKatasRegistry.RegisterKata(nodeType, kata)
}

// Registry is the global registry.
var Registry = DefaultKatasRegistry
