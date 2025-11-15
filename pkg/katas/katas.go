package katas

import (
	"reflect"

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

var KatasByNodeType = make(map[reflect.Type][]Kata)
var KatasByID = make(map[string]Kata)

func RegisterKata(nodeType reflect.Type, kata Kata) {
	KatasByNodeType[nodeType] = append(KatasByNodeType[nodeType], kata)
	KatasByID[kata.ID] = kata
}

func Check(node ast.Node, disabledKatas []string) []Violation {
	var violations []Violation
	nodeType := reflect.TypeOf(node)
	if katas, ok := KatasByNodeType[nodeType]; ok {
		for _, kata := range katas {
			if !isKataDisabled(kata.ID, disabledKatas) {
				violations = append(violations, kata.Check(node)...)
			}
		}
	}
	return violations
}

func isKataDisabled(kataID string, disabledKatas []string) bool {
	for _, disabledKata := range disabledKatas {
		if kataID == disabledKata {
			return true
		}
	}
	return false
}

func GetKata(id string) (Kata, bool) {
	kata, ok := KatasByID[id]
	return kata, ok
}