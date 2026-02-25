// Package phui provides Go view structs mirroring Phabricator's PHUI component system.
// Each struct uses a builder pattern and produces HTML via Render().
package phui

import (
	"html/template"
	"strings"
)

// esc HTML-escapes a string.
func esc(s string) string {
	return template.HTMLEscapeString(s)
}

// classes joins non-empty class names with a space.
func classes(cc ...string) string {
	var out []string
	for _, c := range cc {
		if c != "" {
			out = append(out, c)
		}
	}
	return strings.Join(out, " ")
}

// attr returns ` key="value"` with escaping, or empty string if value is empty.
func attr(key, value string) string {
	if value == "" {
		return ""
	}
	return ` ` + key + `="` + esc(value) + `"`
}
