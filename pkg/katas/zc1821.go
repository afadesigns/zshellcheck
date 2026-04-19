package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

var zc1821DiskutilDestructive = map[string]string{
	"eraseDisk":       "reformats the whole disk",
	"eraseVolume":     "reformats the volume",
	"secureErase":     "overwrites every block, no undo",
	"zeroDisk":        "writes zeros across the whole disk",
	"randomDisk":      "writes random bytes across the whole disk",
	"reformat":        "reformats the volume in place",
	"eraseCD":         "erases the optical disc",
	"erasePartitions": "removes every partition on the disk",
	"partitionDisk":   "rewrites the partition table",
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1821",
		Title:    "Error on `diskutil eraseDisk` / `secureErase` / `partitionDisk` — macOS storage reformat",
		Severity: SeverityError,
		Description: "The `diskutil` subcommands `eraseDisk`, `eraseVolume`, `secureErase`, " +
			"`zeroDisk`, `randomDisk`, `reformat`, `erasePartitions`, and `partitionDisk` all " +
			"rewrite disk or volume state with no Time Machine snapshot or APFS " +
			"preservation. A wrong `/dev/diskN` (especially after a reboot that reordered " +
			"the BSD names) erases the wrong drive, and the only recovery is an offline " +
			"backup. Always pair the call with a typed confirmation, resolve the target by " +
			"`diskutil info -plist` / mount-point rather than by index, and run " +
			"`diskutil list` right before the destructive call.",
		Check: checkZC1821,
	})
}

func checkZC1821(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "diskutil" {
		return nil
	}
	if len(cmd.Arguments) == 0 {
		return nil
	}
	sub := cmd.Arguments[0].String()
	note, ok := zc1821DiskutilDestructive[sub]
	if !ok {
		return nil
	}
	return []Violation{{
		KataID: "ZC1821",
		Message: "`diskutil " + sub + "` " + note + ". Resolve the target by " +
			"`diskutil info -plist` / mount-point (not by index), run " +
			"`diskutil list` immediately before, and require a typed " +
			"confirmation.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}
