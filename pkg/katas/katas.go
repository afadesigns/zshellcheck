package katas

import (
	"fmt"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

// Severity defines the severity of a violation.
type Severity string

const (
	SeverityError   Severity = "error"
	SeverityWarning Severity = "warning"
	SeverityInfo    Severity = "info"
	SeverityStyle   Severity = "style"
)

// Violation represents a found violation in the code.
type Violation struct {
	KataID  string
	Message string
	Line    int
	Column  int
	Level   Severity
}

// FixEdit is a single text replacement applied by the auto-fixer.
// Coordinates are 1-based; Length is the number of bytes of source to
// replace starting at Line:Column. Replace may be empty (deletion) and
// may span multiple lines.
type FixEdit struct {
	Line    int
	Column  int
	Length  int
	Replace string
}

// Kata represents a single linting rule. Fix is optional; when non-nil
// the auto-fixer invokes it with the AST node, the violation, and the
// full file source (byte slice) so the fix can inspect a span around
// the violation before producing edits. Katas with no safe
// deterministic fix leave Fix nil and the fixer skips them.
type Kata struct {
	ID          string
	Title       string
	Description string
	Severity    Severity
	Check       func(node ast.Node) []Violation
	Fix         func(node ast.Node, v Violation, source []byte) []FixEdit
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
		kata.Severity = SeverityWarning
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

// FixesFor invokes the Fix function of the kata that produced the
// given violation, if any, and returns the resulting edits. An empty
// slice means the kata has no deterministic auto-fix; callers should
// skip the violation in fix mode.
func (kr *KatasRegistry) FixesFor(node ast.Node, v Violation, source []byte) []FixEdit {
	kata, ok := kr.KatasByID[v.KataID]
	if !ok || kata.Fix == nil {
		return nil
	}
	return kata.Fix(node, v, source)
}

// CheckAndFix runs Check for every kata registered against the node
// type and, when a kata also declares a Fix, invokes that Fix for each
// emitted violation. Returns the violations (including ones without a
// fix) and the concatenated edits. Use this from the CLI fix mode so
// each node is visited exactly once.
func (kr *KatasRegistry) CheckAndFix(node ast.Node, disabledKatas []string, source []byte) ([]Violation, []FixEdit) {
	var violations []Violation
	var edits []FixEdit
	key := fmt.Sprintf("%T", node)
	katasForNode, ok := kr.KatasByType[key]
	if !ok {
		return nil, nil
	}
	for _, kata := range katasForNode {
		skip := false
		for _, d := range disabledKatas {
			if d == kata.ID {
				skip = true
				break
			}
		}
		if skip {
			continue
		}
		vs := kata.Check(node)
		for i := range vs {
			if vs[i].Level == "" {
				vs[i].Level = kata.Severity
			}
			if kata.Fix != nil {
				edits = append(edits, kata.Fix(node, vs[i], source)...)
			}
		}
		violations = append(violations, vs...)
	}
	return violations, edits
}

// DefaultKatasRegistry is the default registry.
var DefaultKatasRegistry = NewKatasRegistry()

// RegisterKata registers a new Kata with the default registry.
func RegisterKata(nodeType ast.Node, kata Kata) {
	DefaultKatasRegistry.RegisterKata(nodeType, kata)
}

// Registry is the global registry.
var Registry = DefaultKatasRegistry
