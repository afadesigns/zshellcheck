package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1984",
		Title:    "Error on `sgdisk -Z` / `sgdisk -o` — erases the GPT partition table on the target disk",
		Severity: SeverityError,
		Description: "`sgdisk -Z $DISK` (`--zap-all`) wipes the primary GPT, the protective " +
			"MBR, and the backup GPT at the end of the device. `sgdisk -o $DISK` " +
			"(`--clear`) replaces the existing partition table with a fresh empty GPT. " +
			"Either command detaches every partition, LVM PV, LUKS container, and " +
			"filesystem header on the device — when the target variable resolves to a " +
			"wrong path (tab completion, `$DISK` defaulted to `/dev/sda`), the host " +
			"becomes unbootable. Require an `lsblk $DISK` + `blkid $DISK` preflight in " +
			"the script, route the action through `--pretend` (`-t`) first, and keep a " +
			"`sgdisk --backup=/root/$DISK.gpt $DISK` image before any zap.",
		Check: checkZC1984,
	})
}

func checkZC1984(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	// Parser caveat: `sgdisk --zap-all $DISK` mangles command name to
	// `zap-all` and leaves the `--` as a postfix on a prior expression.
	if ident.Value == "zap-all" {
		return zc1984Hit(cmd, "sgdisk --zap-all")
	}
	if ident.Value != "sgdisk" {
		return nil
	}
	for _, arg := range cmd.Arguments {
		v := arg.String()
		switch v {
		case "-Z":
			return zc1984Hit(cmd, "sgdisk -Z")
		case "-o":
			return zc1984Hit(cmd, "sgdisk -o")
		}
	}
	return nil
}

func zc1984Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1984",
		Message: "`" + form + "` erases the GPT on the target device — a wrong " +
			"`$DISK` detaches every partition/LVM/LUKS header and bricks boot. " +
			"`lsblk`/`blkid` preflight, `--backup` the old table, and test with " +
			"`-t`/`--pretend` first.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}
