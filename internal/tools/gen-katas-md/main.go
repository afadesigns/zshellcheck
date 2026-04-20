// Command gen-katas-md regenerates KATAS.md from the live kata registry.
//
//	go run ./internal/tools/gen-katas-md
//
// Writes to ./KATAS.md (relative to the module root) by default, or to the
// path passed via -o. Every kata's ID, title, severity, and description are
// sourced directly from the RegisterKata() calls in pkg/katas/zc*.go so
// the generated file cannot drift from the implementation.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/katas"
)

func main() {
	out := flag.String("o", "KATAS.md", "output path (default: KATAS.md at module root)")
	flag.Parse()

	registry := katas.Registry
	ids := make([]string, 0, len(registry.KatasByID))
	for id := range registry.KatasByID {
		ids = append(ids, id)
	}
	sort.Slice(ids, func(i, j int) bool {
		return idNum(ids[i]) < idNum(ids[j])
	})

	f, err := os.Create(*out)
	if err != nil {
		log.Fatalf("create %s: %v", *out, err)
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	defer w.Flush()

	count := len(ids)
	fmt.Fprintf(w, "# ZShellCheck Katas\n\n")
	fmt.Fprintf(w, "Auto-generated list of all %d implemented checks. Do not edit by hand — regenerate via `go run ./internal/tools/gen-katas-md`.\n\n", count)

	fmt.Fprintf(w, "## Summary\n\n")
	fmt.Fprintf(w, "| Severity | Count |\n| :--- | ---: |\n")
	sevCount := map[katas.Severity]int{}
	for _, id := range ids {
		sevCount[registry.KatasByID[id].Severity]++
	}
	for _, sev := range []katas.Severity{katas.SeverityError, katas.SeverityWarning, katas.SeverityInfo, katas.SeverityStyle} {
		fmt.Fprintf(w, "| `%s` | %d |\n", sev, sevCount[sev])
	}
	fmt.Fprintf(w, "| **total** | **%d** |\n\n", count)

	fmt.Fprintf(w, "## Table of Contents\n\n")
	for _, id := range ids {
		k := registry.KatasByID[id]
		fmt.Fprintf(w, "- [%s: %s](#%s)\n", k.ID, escapeTitle(k.Title), strings.ToLower(k.ID))
	}
	fmt.Fprintf(w, "\n---\n\n")

	for _, id := range ids {
		k := registry.KatasByID[id]
		fmt.Fprintf(w, "<a id=\"%s\"></a>\n", strings.ToLower(k.ID))
		fmt.Fprintf(w, "### %s — %s\n\n", k.ID, k.Title)
		fmt.Fprintf(w, "**Severity:** `%s`\n\n", k.Severity)
		fmt.Fprintf(w, "%s\n\n", k.Description)
		fmt.Fprintf(w, "Disable by adding `%s` to `disabled_katas` in `.zshellcheckrc`.\n\n", k.ID)
		fmt.Fprintf(w, "---\n\n")
	}
}

// idNum extracts the numeric portion of a kata ID ("ZC1234" -> 1234).
// Returns 0 for IDs that do not parse, which is acceptable for sort stability.
func idNum(id string) int {
	if len(id) < 3 {
		return 0
	}
	n, _ := strconv.Atoi(id[2:])
	return n
}

// escapeTitle drops markdown control characters that would corrupt a link label.
func escapeTitle(t string) string {
	replacer := strings.NewReplacer("[", "\\[", "]", "\\]", "|", "\\|")
	return replacer.Replace(t)
}
