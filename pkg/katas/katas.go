package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

type Kata struct {
	ID          string
	Title       string
	Description string
	Check       func(node ast.Node) []Violation
}

type Violation struct {
	KataID  string
	Message string
	Line    int
	Column  int
}

type KatasRegistry struct {
	KatasByNodeType map[ast.NodeType][]Kata
	KatasByID       map[string]Kata
}

// Registry is the global instance of KatasRegistry that holds all registered katas.
var Registry = KatasRegistry{
	KatasByNodeType: make(map[ast.NodeType][]Kata),
	KatasByID:       make(map[string]Kata),
}

// RegisterKata registers a new kata with the global registry.
func RegisterKata(nodeType ast.NodeType, kata Kata) {
	Registry.KatasByNodeType[nodeType] = append(Registry.KatasByNodeType[nodeType], kata)
	Registry.KatasByID[kata.ID] = kata
}

// Check runs all applicable katas against a given AST node.
func (kr *KatasRegistry) Check(node ast.Node, disabledKatas []string) []Violation {
	var violations []Violation
	if katas, ok := kr.KatasByNodeType[node.Type()]; ok {
		for _, kata := range katas {
			if !kr.isKataDisabled(kata.ID, disabledKatas) {
				violations = append(violations, kata.Check(node)...)
			}
		}
	}
	return violations
}

// isKataDisabled checks if a kata is present in the list of disabled katas.
func (kr *KatasRegistry) isKataDisabled(kataID string, disabledKatas []string) bool {
	for _, disabledKata := range disabledKatas {
		if kataID == disabledKata {
			return true
		}
	}
	return false
}

// GetKata retrieves a kata by its ID from the global registry.
func (kr *KatasRegistry) GetKata(id string) (Kata, bool) {
	kata, ok := kr.KatasByID[id]
	return kata, ok
}