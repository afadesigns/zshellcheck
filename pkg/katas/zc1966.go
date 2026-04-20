package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1966",
		Title:    "Error on `zpool import -f` / `zpool export -f` — forced ZFS pool op bypasses hostid/txg checks",
		Severity: SeverityError,
		Description: "`zpool import -f $POOL` force-imports a pool even when the on-disk " +
			"hostid differs — i.e. the pool is already imported on another host " +
			"(multipath/SAN, shared JBOD, HA cluster). The second import writes to the " +
			"same vdevs and silently corrupts the pool. `zpool export -f` skips the " +
			"graceful-flush path and detaches vdevs with in-flight txgs, which can lose " +
			"the tail of the ZIL. Export without `-f` after `zfs unmount -a`; import " +
			"without `-f` after verifying `zpool import` (no target) reports the pool " +
			"as `ONLINE` and the hostid matches.",
		Check: checkZC1966,
	})
}

func checkZC1966(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "zpool" || len(cmd.Arguments) < 2 {
		return nil
	}
	sub := cmd.Arguments[0].String()
	if sub != "import" && sub != "export" {
		return nil
	}
	for _, arg := range cmd.Arguments[1:] {
		v := arg.String()
		if v == "-f" || v == "--force" {
			return zc1966Hit(cmd, "zpool "+sub+" -f")
		}
	}
	return nil
}

func zc1966Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1966",
		Message: "`" + form + "` bypasses hostid/txg safety — forced import of a pool " +
			"already online elsewhere (SAN/HA) corrupts it; forced export drops in-flight " +
			"txgs. `zfs unmount -a` first, then plain `zpool export`/`import`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}
