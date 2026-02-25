package diff

import (
	"bytes"
	"html"
	"path/filepath"
	"strings"

	"github.com/alecthomas/chroma/v2"
	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
)

var (
	formatter = chromahtml.New(
		chromahtml.WithClasses(true),
		chromahtml.PreventSurroundingPre(true),
	)
	style = styles.Get("github")
)

// HighlightLines takes a filename and a slice of raw source lines,
// returning syntax-highlighted HTML for each line. The indices match 1:1.
func HighlightLines(filename string, lines []string) []string {
	lexer := lexers.Match(filename)
	if lexer == nil {
		lexer = lexers.Fallback
	}
	lexer = chroma.Coalesce(lexer)

	source := strings.Join(lines, "\n")
	iter, err := lexer.Tokenise(nil, source)
	if err != nil {
		return escapeLines(lines)
	}

	// Split tokens into per-line groups.
	allTokens := iter.Tokens()
	tokenLines := splitTokensByLine(allTokens)

	result := make([]string, len(lines))
	for i := range lines {
		if i < len(tokenLines) {
			var buf bytes.Buffer
			lineIter := chroma.Literator(tokenLines[i]...)
			if err := formatter.Format(&buf, style, lineIter); err != nil {
				result[i] = html.EscapeString(lines[i])
			} else {
				result[i] = buf.String()
			}
		} else {
			result[i] = html.EscapeString(lines[i])
		}
	}
	return result
}

// splitTokensByLine splits a flat token slice into per-line groups.
func splitTokensByLine(tokens []chroma.Token) [][]chroma.Token {
	var lines [][]chroma.Token
	var current []chroma.Token

	for _, tok := range tokens {
		parts := strings.Split(tok.Value, "\n")
		for j, part := range parts {
			if j > 0 {
				lines = append(lines, current)
				current = nil
			}
			if part != "" {
				current = append(current, chroma.Token{
					Type:  tok.Type,
					Value: part,
				})
			}
		}
	}
	if current != nil {
		lines = append(lines, current)
	}
	return lines
}

func escapeLines(lines []string) []string {
	out := make([]string, len(lines))
	for i, l := range lines {
		out[i] = html.EscapeString(l)
	}
	return out
}

// FileIcon returns a Font Awesome icon class for a filename.
func FileIcon(filename string) string {
	switch filepath.Ext(filename) {
	case ".go":
		return "fa-file-code-o"
	case ".js", ".ts", ".jsx", ".tsx":
		return "fa-file-code-o"
	case ".css", ".scss", ".less":
		return "fa-file-code-o"
	case ".html", ".htm", ".tmpl":
		return "fa-file-code-o"
	case ".md", ".txt", ".rst":
		return "fa-file-text"
	case ".json", ".yaml", ".yml", ".toml":
		return "fa-file-code-o"
	case ".png", ".jpg", ".jpeg", ".gif", ".svg", ".ico":
		return "fa-file-image-o"
	default:
		return "fa-file-text"
	}
}
