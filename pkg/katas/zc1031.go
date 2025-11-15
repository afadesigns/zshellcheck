package katas

import (
	"reflect"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(reflect.TypeOf(&ast.Shebang{}), Kata{
		ID:          "ZC1031",
		Title:       "Use `#!/usr/bin/env zsh` for portability",
		Description: "Using `#!/usr/bin/env zsh` is more portable than `#!/bin/zsh` because it searches " +
			"for the `zsh` executable in the user's `PATH`.",
		Check:       checkZC1031,
	})
}

func checkZC1031(node ast.Node) []Violation {
	violations := []Violation{}

	if shebang, ok := node.(*ast.Shebang); ok {
		if shebang.Path == "#!/bin/zsh" {
			violations = append(violations, Violation{
				KataID:  "ZC1031",
				Message: "Use `#!/usr/bin/env zsh` for portability instead of `#!/bin/zsh`.",
				Line:    1,
				Column:  1,
			})
		}
	}

	return violations
}
