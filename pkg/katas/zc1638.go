package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

var zc1638SecretArgs = []string{
	"password", "passwd", "secret", "token",
	"apikey", "api_key", "api-key",
	"accesskey", "access_key", "access-key",
	"privatekey", "private_key", "private-key",
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1638",
		Title:    "Error on `docker/podman build --build-arg SECRET=VALUE` — secret baked into image layer",
		Severity: SeverityError,
		Description: "`--build-arg KEY=VALUE` values land in the image metadata that `docker " +
			"history` (and the analogous podman / buildah tooling) read back from the layer. " +
			"Even if the Dockerfile only uses the arg to export as a build-time env var, the " +
			"literal value is cached in the layer forever. A key-shaped name (`password`, " +
			"`secret`, `token`, `apikey`, `access_key`, `private_key`) with a concrete value " +
			"embeds that secret in every image pulled. Use BuildKit secrets " +
			"(`--secret id=mysecret,src=path`) or a multi-stage build where the secret stays " +
			"in a discarded stage.",
		Check: checkZC1638,
	})
}

func checkZC1638(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	switch ident.Value {
	case "docker", "podman", "buildah":
	default:
		return nil
	}
	if len(cmd.Arguments) == 0 || cmd.Arguments[0].String() != "build" {
		return nil
	}

	for i, arg := range cmd.Arguments {
		if arg.String() != "--build-arg" {
			continue
		}
		if i+1 >= len(cmd.Arguments) {
			continue
		}
		pair := cmd.Arguments[i+1].String()
		eq := strings.Index(pair, "=")
		if eq < 0 {
			continue
		}
		key := strings.ToLower(pair[:eq])
		for _, s := range zc1638SecretArgs {
			if strings.Contains(key, s) {
				return []Violation{{
					KataID: "ZC1638",
					Message: "`" + ident.Value + " build --build-arg " + pair + "` bakes " +
						"the secret into the image layer metadata. Use `--secret " +
						"id=NAME,src=PATH` (BuildKit) or a multi-stage build.",
					Line:   cmd.Token.Line,
					Column: cmd.Token.Column,
					Level:  SeverityError,
				}}
			}
		}
	}
	return nil
}
