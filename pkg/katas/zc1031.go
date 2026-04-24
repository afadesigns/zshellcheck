package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.ShebangNode, Kata{
		ID:    "ZC1031",
		Title: "Use `#!/usr/bin/env zsh` for portability",
		Description: "Using `#!/usr/bin/env zsh` is more portable than `#!/bin/zsh` because it searches " +
			"for the `zsh` executable in the user's `PATH`.",
		Severity: SeverityInfo,
		Check:    checkZC1031,
		Fix:      fixZC1031,
	})
}

// fixZC1031 rewrites `#!/bin/zsh` to `#!/usr/bin/env zsh` in the
// shebang line. Span-aware: replaces the whole `#!/bin/zsh` run as a
// single edit at column 1, line 1.
func fixZC1031(node ast.Node, v Violation, source []byte) []FixEdit {
	shebang, ok := node.(*ast.Shebang)
	if !ok {
		return nil
	}
	if shebang.Path != "#!/bin/zsh" {
		return nil
	}
	return []FixEdit{{
		Line:    v.Line,
		Column:  v.Column,
		Length:  len("#!/bin/zsh"),
		Replace: "#!/usr/bin/env zsh",
	}}
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
				Level:   SeverityInfo,
			})
		}
	}

	return violations
}
