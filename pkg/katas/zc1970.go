package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1970",
		Title:    "Warn on `losetup -P` / `kpartx -a` / `partprobe` on untrusted image — runs kernel partition parser",
		Severity: SeverityWarning,
		Description: "`losetup -P $LOOP $IMG`, `kpartx -av $IMG`, and `partprobe $LOOP` all " +
			"tell the kernel to rescan a block device's partition table and emit `/dev/" +
			"loopNpX` (or dm-N) entries. When the image comes from an untrusted source " +
			"— a customer-supplied VM disk, a downloaded installer, a forensic capture — " +
			"the scan runs MBR/GPT/LVM parsers over attacker-controlled bytes and has " +
			"historically triggered kernel CVEs (fsconfig heap overflow, ext4 mount " +
			"bugs). Do the inspection in a throwaway VM or an offline parser like " +
			"`fdisk -l $IMG` / `sfdisk --dump $IMG` that reads without kernel scan, and " +
			"only attach partitions with `losetup -P` after the layout is known-good.",
		Check: checkZC1970,
	})
}

func checkZC1970(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	switch ident.Value {
	case "losetup":
		for _, arg := range cmd.Arguments {
			v := arg.String()
			if v == "-P" || v == "--partscan" ||
				v == "-Pf" || v == "-fP" || v == "-rP" || v == "-Pr" {
				return zc1970Hit(cmd, "losetup -P")
			}
		}
	case "kpartx":
		for _, arg := range cmd.Arguments {
			v := arg.String()
			if v == "-a" || v == "-av" || v == "-va" || v == "-as" || v == "-sa" {
				return zc1970Hit(cmd, "kpartx -a")
			}
		}
	case "partprobe":
		return zc1970Hit(cmd, "partprobe")
	}
	return nil
}

func zc1970Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1970",
		Message: "`" + form + "` asks the kernel to parse the partition table of the " +
			"image — attacker-controlled bytes have tripped kernel CVEs. Use `fdisk " +
			"-l`/`sfdisk --dump` offline first, scan only known-good images.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}
