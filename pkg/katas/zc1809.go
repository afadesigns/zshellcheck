package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1809",
		Title:    "Error on `gsutil rm -r gs://…` / `gsutil rb -f gs://…` — bulk GCS deletion",
		Severity: SeverityError,
		Description: "`gsutil rm -r gs://bucket/prefix` and `gsutil rm -rf gs://bucket` delete " +
			"every object under the prefix — with `-m` (parallel) they do it faster than any " +
			"undo window. `gsutil rb -f gs://bucket` removes the bucket after force-deleting " +
			"the contents. Neither soft-deletes; Object Versioning can help only if it is " +
			"turned on in advance, and `gsutil rb` leaves no retention grace. Preview with " +
			"`gsutil ls`, enable Object Versioning or retention locks before the fact, and " +
			"prefer narrower `gsutil rm gs://bucket/specific-object` calls.",
		Check: checkZC1809,
	})
}

func checkZC1809(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "gsutil" {
		return nil
	}

	subIdx := -1
	for i, arg := range cmd.Arguments {
		v := arg.String()
		if v == "rm" || v == "rb" {
			subIdx = i
			break
		}
	}
	if subIdx == -1 {
		return nil
	}
	sub := cmd.Arguments[subIdx].String()

	hasDestFlag := false
	for _, arg := range cmd.Arguments[subIdx+1:] {
		v := arg.String()
		if sub == "rm" && (v == "-r" || v == "-R" || v == "-rf" || v == "-fr" || v == "--recursive") {
			hasDestFlag = true
			break
		}
		if sub == "rb" && (v == "-f" || v == "--force") {
			hasDestFlag = true
			break
		}
	}
	if !hasDestFlag {
		return nil
	}
	return []Violation{{
		KataID: "ZC1809",
		Message: "`gsutil " + sub + "` with recursive/force deletes every matching " +
			"GCS object (or the bucket itself). Preview with `gsutil ls`, enable " +
			"Object Versioning / retention locks ahead of time, and prefer narrower " +
			"object-level `gsutil rm` calls.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}
