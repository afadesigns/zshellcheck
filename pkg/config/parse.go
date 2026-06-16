// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package config

import (
	"fmt"
	"strconv"
	"strings"
)

// Parse reads a ZShellCheck configuration from its YAML-subset format.
// The schema is flat: scalar `key: value` pairs plus the `disabled_katas`
// sequence (a block list of `- ZC####` items or an inline `[ZC####, …]`).
// It is implemented without a third-party YAML dependency to keep the
// binary dependency-free, and accepts the documented format: `#` comments,
// single/double quotes, and standard escapes inside double quotes.
//
// Parse returns an error on a structurally malformed line — a key with no
// value separator, an empty key, or a sequence item outside a list — so a
// genuinely broken config is reported rather than silently ignored.
func Parse(data []byte) (Config, error) {
	var cfg Config
	inList := false
	for n, raw := range strings.Split(string(data), "\n") {
		line := stripComment(raw)
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		if strings.HasPrefix(trimmed, "- ") || trimmed == "-" {
			if !inList {
				return cfg, fmt.Errorf("config: line %d: list item outside a sequence", n+1)
			}
			if item := unquote(strings.TrimSpace(trimmed[1:])); item != "" {
				cfg.DisabledKatas = append(cfg.DisabledKatas, item)
			}
			continue
		}
		key, val, ok := strings.Cut(trimmed, ":")
		key = strings.TrimSpace(key)
		val = strings.TrimSpace(val)
		if !ok || key == "" {
			return cfg, fmt.Errorf("config: line %d: expected `key: value`", n+1)
		}
		inList = false
		if key == "disabled_katas" {
			inList = parseListValue(&cfg, val)
			continue
		}
		if err := assignScalar(&cfg, key, unquote(val)); err != nil {
			return cfg, fmt.Errorf("config: line %d: %w", n+1, err)
		}
	}
	return cfg, nil
}

// parseListValue handles the value after `disabled_katas:`. An empty
// value opens a block sequence (reported by the true return); `[]` is the
// empty inline list; `[a, b]` is a populated inline list; anything else is
// a single bare item.
func parseListValue(cfg *Config, val string) (blockOpen bool) {
	switch {
	case val == "":
		return true
	case val == "[]":
		return false
	case strings.HasPrefix(val, "[") && strings.HasSuffix(val, "]"):
		for _, item := range strings.Split(val[1:len(val)-1], ",") {
			if item := unquote(strings.TrimSpace(item)); item != "" {
				cfg.DisabledKatas = append(cfg.DisabledKatas, item)
			}
		}
	default:
		cfg.DisabledKatas = append(cfg.DisabledKatas, unquote(val))
	}
	return false
}

// assignScalar sets the field named by key. Unknown keys are ignored for
// forward compatibility; only a malformed boolean is an error.
func assignScalar(cfg *Config, key, val string) error {
	switch key {
	case "error_color":
		cfg.ErrorColor = val
	case "warning_color":
		cfg.WarningColor = val
	case "info_color":
		cfg.InfoColor = val
	case "id_color":
		cfg.IDColor = val
	case "title_color":
		cfg.TitleColor = val
	case "message_color":
		cfg.MessageColor = val
	case "line_color":
		cfg.LineColor = val
	case "column_color":
		cfg.ColumnColor = val
	case "no_color":
		b, err := strconv.ParseBool(val)
		if err != nil {
			return fmt.Errorf("invalid boolean for no_color: %q", val)
		}
		cfg.NoColor = b
	case "verbose":
		b, err := strconv.ParseBool(val)
		if err != nil {
			return fmt.Errorf("invalid boolean for verbose: %q", val)
		}
		cfg.Verbose = b
	}
	return nil
}

// stripComment removes a trailing `#` comment. A `#` starts a comment only
// at line start or after whitespace, and never inside quotes, so values
// like a quoted `"#ff0000"` survive.
func stripComment(line string) string {
	inSingle, inDouble := false, false
	for i := 0; i < len(line); i++ {
		switch c := line[i]; {
		case inSingle:
			if c == '\'' {
				inSingle = false
			}
		case inDouble:
			if c == '"' {
				inDouble = false
			}
		case c == '\'':
			inSingle = true
		case c == '"':
			inDouble = true
		case c == '#' && (i == 0 || line[i-1] == ' ' || line[i-1] == '\t'):
			return line[:i]
		}
	}
	return line
}

// unquote strips a matching pair of surrounding quotes. Double-quoted
// values have their escapes decoded; single-quoted values are literal.
func unquote(s string) string {
	if len(s) < 2 {
		return s
	}
	switch {
	case s[0] == '"' && s[len(s)-1] == '"':
		return decodeEscapes(s[1 : len(s)-1])
	case s[0] == '\'' && s[len(s)-1] == '\'':
		return s[1 : len(s)-1]
	}
	return s
}

// simpleEscapes maps a single character after `\` to its byte value.
// ESC is reachable through `\e` and `\xNN` — the two forms a double-quoted
// YAML scalar uses for ANSI colour codes.
var simpleEscapes = map[byte]byte{
	'n': '\n', 't': '\t', 'r': '\r', 'e': 0x1b, '0': 0, '\\': '\\', '"': '"',
}

// decodeEscapes resolves the escape sequences a double-quoted YAML scalar
// may carry. An unrecognised escape is left verbatim.
func decodeEscapes(s string) string {
	if !strings.ContainsRune(s, '\\') {
		return s
	}
	var b strings.Builder
	for i := 0; i < len(s); i++ {
		if s[i] != '\\' || i+1 >= len(s) {
			b.WriteByte(s[i])
			continue
		}
		i++
		if s[i] == 'x' && i+2 < len(s) {
			if v, err := strconv.ParseUint(s[i+1:i+3], 16, 8); err == nil {
				b.WriteByte(byte(v))
				i += 2
				continue
			}
		}
		if v, ok := simpleEscapes[s[i]]; ok {
			b.WriteByte(v)
			continue
		}
		b.WriteByte('\\')
		b.WriteByte(s[i])
	}
	return b.String()
}
