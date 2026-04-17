package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1424",
		Title:    "Dangerous: `mkfs.*` / `mkfs -t` — formats a filesystem, destroys data",
		Severity: SeverityError,
		Description: "`mkfs.ext4 /dev/sda1`, `mkfs.xfs /dev/...`, `mkfs -t ...` all destroy the " +
			"existing filesystem on the target device. A typo on the target path reformats the " +
			"wrong disk. Validate the device path, use `blkid` / `lsblk` first, and consider a " +
			"confirmation prompt.",
		Check: checkZC1424,
	})
}

func checkZC1424(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	// mkfs, mkfs.ext4, mkfs.xfs, mkfs.btrfs, mkfs.vfat, etc.
	name := ident.Value
	if name == "mkfs" || strings.HasPrefix(name, "mkfs.") || name == "mke2fs" ||
		name == "mkswap" || name == "wipefs" {
		return []Violation{{
			KataID: "ZC1424",
			Message: "`" + name + "` formats / wipes a device — destroys data. Validate the " +
				"target with `lsblk` / `blkid` first, and consider an interactive confirmation.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityError,
		}}
	}

	return nil
}
