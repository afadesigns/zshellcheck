// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1816",
		Title:    "Warn on `docker/podman commit` — produces un-reproducible image, bakes in runtime state",
		Severity: SeverityWarning,
		Description: "`docker commit CONTAINER IMAGE` (and the podman / nerdctl equivalents) " +
			"snapshots a running container's filesystem into a new image. There is no " +
			"Dockerfile, so the build is not reproducible; the snapshot inherits whatever " +
			"`/tmp` scratch, shell history, environment variables, and — frequently — " +
			"credentials the container held at that moment; and the resulting image's layer " +
			"metadata records only the container id, not what was actually installed. Build " +
			"from a `Dockerfile` (or `docker buildx build`) so the image can be regenerated " +
			"from source, and use `docker commit` only for one-off rescue work on a local " +
			"image you are about to discard.",
		Check: checkZC1816,
	})
}

func checkZC1816(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "docker" && ident.Value != "podman" && ident.Value != "nerdctl" {
		return nil
	}
	if len(cmd.Arguments) == 0 || cmd.Arguments[0].String() != "commit" {
		return nil
	}
	return []Violation{{
		KataID: "ZC1816",
		Message: "`" + ident.Value + " commit` snapshots a running container — no " +
			"Dockerfile trail, runtime env / `/tmp` scratch / shell history get baked " +
			"in, and the layer metadata does not record what was installed. Build from " +
			"a `Dockerfile` instead.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}
