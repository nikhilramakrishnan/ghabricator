package diff

import (
	"encoding/json"
	"fmt"
	"html/template"
	"strings"

	"github.com/nikhilr/ghabricator/internal/remarkup"
)

// InlineComment represents an inline comment to be rendered within the diff table.
type InlineComment struct {
	ID        int64
	Author    string
	AvatarURL string
	Body      string
	Path      string // file path (used as changesetID)
	Line      int
	Side      string // "LEFT" or "RIGHT"
}

// InlineCommentMeta builds the full data-meta map expected by DiffInline.js bindToRow.
func InlineCommentMeta(c InlineComment) map[string]any {
	contentState := map[string]any{
		"text":           c.Body,
		"suggestionText": nil,
		"hasSuggestion":  false,
	}
	return map[string]any{
		"id":                  c.ID,
		"phid":                fmt.Sprintf("GHCMT-%d", c.ID),
		"on_right":            c.Side == "RIGHT",
		"number":              c.Line,
		"length":              0,
		"isNewFile":           c.Side == "RIGHT",
		"changesetID":         c.Path,
		"isDraft":             false,
		"isFixed":             false,
		"isGhost":             false,
		"isSynthetic":         false,
		"isDraftDone":         false,
		"isEditing":           false,
		"replyToCommentPHID":  nil,
		"snippet":             truncateSnippet(c.Body),
		"menuItems":           defaultMenuItems(),
		"documentEngineKey":   nil,
		"startOffset":         nil,
		"endOffset":           nil,
		"canSuggestEdit":      false,
		"state": map[string]any{
			"initial":   contentState,
			"committed": contentState,
			"active":    contentState,
		},
	}
}

func truncateSnippet(body string) string {
	if len(body) <= 80 {
		return body
	}
	return body[:80] + "..."
}

func defaultMenuItems() []map[string]any {
	return []map[string]any{
		{"action": "reply", "label": "Reply", "icon": "fa-reply", "key": "r"},
		{"action": "quote", "label": "Quote", "icon": "fa-quote-left"},
		{"action": "edit", "label": "Edit", "icon": "fa-pencil", "key": "e"},
		{"action": "delete", "label": "Delete", "icon": "fa-times", "key": "x"},
		{"action": "collapse", "label": "Collapse", "icon": "fa-compress"},
	}
}

// ChangesetMeta is the metadata structure for a changeset, stored in Javelin's
// data store and referenced via data-meta pointer on the DOM element.
type ChangesetMeta struct {
	Left        string `json:"left"`
	Right       string `json:"right"`
	RenderURI   string `json:"renderURI"`
	Ref         string `json:"ref"`
	Loaded      bool   `json:"loaded"`
	DisplayPath string `json:"displayPath"`
	Icon        string `json:"icon"`
}

// BuildChangesetMeta creates the metadata for a changeset.
func BuildChangesetMeta(cs Changeset) ChangesetMeta {
	path := cs.DisplayPath()
	return ChangesetMeta{
		Left:        path,
		Right:       path,
		RenderURI:   "/api/changeset/render",
		Ref:         path,
		Loaded:      true,
		DisplayPath: path,
		Icon:        FileIcon(path),
	}
}

// RenderChangeset produces the full Phabricator-style HTML for a single changeset.
// metaRef is the Javelin metadata pointer (e.g., "0_3") for the data-meta attribute.
// comments are inline comments to be rendered within the diff table at their respective lines.
func RenderChangeset(cs Changeset, metaRef string, comments []InlineComment) template.HTML {
	hunkRows := buildRowsByHunk(cs)
	path := cs.DisplayPath()
	icon := FileIcon(path)

	// Index comments by (line, side) for quick lookup after each row.
	type lineKey struct {
		line int
		side string
	}
	commentsByKey := make(map[lineKey][]InlineComment)
	for _, c := range comments {
		k := lineKey{line: c.Line, side: c.Side}
		commentsByKey[k] = append(commentsByKey[k], c)
	}

	var b strings.Builder

	// Outer wrapper — data-meta is a Javelin pointer, not inline JSON.
	// data-lines-added/removed used by collapse JS to decide auto-collapse.
	fmt.Fprintf(&b, `<div class="differential-changeset" id="diff-C%d" data-sigil="differential-changeset" data-meta="%s" data-lines-added="%d" data-lines-removed="%d">`, cs.ID, metaRef, cs.LinesAdded, cs.LinesRemoved)
	fmt.Fprintf(&b, `<a name="C%d" class="differential-inline-comment-anchor"></a>`, cs.ID)

	// File header — sticky bar with collapse toggle and line stats (moodboard style)
	b.WriteString(`<div class="mood-changeset-header" data-sigil="changeset-header">`)
	b.WriteString(`<span class="phui-icon-view phui-font-fa fa-chevron-down changeset-collapse-toggle"></span>`)
	fmt.Fprintf(&b, `<span class="phui-icon-view phui-font-fa %s" style="opacity:0.5"></span>`, icon)
	fmt.Fprintf(&b, `<span class="differential-changeset-path-name" data-sigil="changeset-header-path-name">%s</span>`, template.HTMLEscapeString(path))
	b.WriteString(`<span class="stats">`)
	if cs.LinesAdded > 0 {
		fmt.Fprintf(&b, `<span class="add-stat">+%d</span>`, cs.LinesAdded)
	}
	if cs.LinesRemoved > 0 {
		if cs.LinesAdded > 0 {
			b.WriteString(` `)
		}
		fmt.Fprintf(&b, `<span class="del-stat">-%d</span>`, cs.LinesRemoved)
	}
	b.WriteString(`</span>`)
	b.WriteString(`</div>`)

	// Diff table
	b.WriteString(`<div class="changeset-view-content" data-sigil="changeset-view-content">`)
	b.WriteString(`<table class="differential-diff remarkup-code PhabricatorMonospaced diff-2up chroma" data-sigil="differential-diff intercept-copy">`)
	b.WriteString(`<colgroup>`)
	b.WriteString(`<col class="num" style="width:4em" /><col class="left"/>`)
	b.WriteString(`<col class="num" style="width:4em" /><col class="copy"/>`)
	b.WriteString(`<col class="right"/><col class="cov"/>`)
	b.WriteString(`</colgroup>`)

	for hunkIdx, hunk := range cs.Hunks {
		// "Show more" separator at the top of the first hunk if it doesn't start at line 1.
		if hunkIdx == 0 && hunk.NewStart > 1 {
			renderShowMore(&b, fmt.Sprintf("Context above (lines 1\u2013%d)", hunk.NewStart-1), path, 1, hunk.NewStart-1, cs.ID)
		}

		// "Show more" separator between hunks.
		if hunkIdx > 0 {
			prevHunk := cs.Hunks[hunkIdx-1]
			prevEndNew := prevHunk.NewStart + prevHunk.NewCount
			gapLines := hunk.NewStart - prevEndNew
			if gapLines > 0 {
				renderShowMore(&b, fmt.Sprintf("Show %d more lines", gapLines), path, prevEndNew, hunk.NewStart-1, cs.ID)
			}
		}

		rows := hunkRows[hunkIdx]
		for _, row := range rows {
			renderRow(&b, cs.ID, row)

			// Insert inline comments that target this row's lines.
			if row.NewNum > 0 {
				if cmts, ok := commentsByKey[lineKey{line: row.NewNum, side: "RIGHT"}]; ok {
					for _, c := range cmts {
						renderInlineCommentRow(&b, c)
					}
				}
			}
			if row.OldNum > 0 {
				if cmts, ok := commentsByKey[lineKey{line: row.OldNum, side: "LEFT"}]; ok {
					for _, c := range cmts {
						renderInlineCommentRow(&b, c)
					}
				}
			}
		}

		// "Context below" separator after the last hunk.
		if hunkIdx == len(cs.Hunks)-1 {
			lastEnd := hunk.NewStart + hunk.NewCount
			renderShowMore(&b, "Context below", path, lastEnd, lastEnd+20, cs.ID)
		}
	}

	b.WriteString(`</table>`)
	b.WriteString(`</div>`) // changeset-view-content
	b.WriteString(`</div>`) // differential-changeset

	return template.HTML(b.String())
}

// renderShowMore writes a clickable "show more lines" separator row.
func renderShowMore(b *strings.Builder, label, path string, startLine, endLine, csID int) {
	b.WriteString(`<tr class="show-more">`)
	b.WriteString(`<th class="num"></th>`)
	fmt.Fprintf(b, `<td class="show-more-content" colspan="5" data-action="context-expand" data-path="%s" data-start="%d" data-end="%d" data-cs="%d" `+
		`style="text-align:center;padding:6px;background:rgba(55,55,55,.04);color:#6b748c;font-size:12px;cursor:pointer">`+
		`<span class="phui-icon-view phui-font-fa fa-ellipsis-h mrs"></span>`+
		`%s</td>`,
		template.HTMLEscapeString(path), startLine, endLine, csID,
		template.HTMLEscapeString(label))
	b.WriteString(`</tr>`)
}

// RenderContextRows produces <tr> elements for expanded context lines.
// fileContent is the full file text, startLine/endLine are 1-indexed inclusive.
func RenderContextRows(path, fileContent string, startLine, endLine, csID int) string {
	lines := strings.Split(fileContent, "\n")
	if startLine < 1 {
		startLine = 1
	}
	if endLine > len(lines) {
		endLine = len(lines)
	}
	if startLine > endLine {
		return ""
	}

	// Highlight the context slice
	slice := lines[startLine-1 : endLine]
	highlighted := HighlightLines(path, slice)

	var b strings.Builder
	for i, hl := range highlighted {
		lineNum := startLine + i
		fmt.Fprintf(&b, `<tr>`)
		fmt.Fprintf(&b, `<td class="n" data-n="%d" id="C%dOL%d"></td>`, lineNum, csID, lineNum)
		fmt.Fprintf(&b, `<td data-copy-mode="copy-l">%s</td>`, hl)
		fmt.Fprintf(&b, `<td class="n" data-n="%d" id="C%dNL%d"></td>`, lineNum, csID, lineNum)
		fmt.Fprintf(&b, `<td class="copy"></td>`)
		fmt.Fprintf(&b, `<td colspan="2" data-copy-mode="copy-r">%s</td>`, hl)
		b.WriteString(`</tr>`)
	}
	return b.String()
}

// renderInlineCommentRow writes a single inline comment as a <tr> inside the diff table.
// Styled to match the moodboard design with avatar, author, and action buttons.
func renderInlineCommentRow(b *strings.Builder, c InlineComment) {
	metaJSON, _ := json.Marshal(InlineCommentMeta(c))

	b.WriteString(`<tr class="inline" data-sigil="inline-row">`)
	b.WriteString(`<td colspan="6">`)
	fmt.Fprintf(b, `<div class="differential-inline-comment inline-comment-element mood-inline-comment" data-sigil="differential-inline-comment" data-meta='%s'>`,
		template.HTMLEscapeString(string(metaJSON)))

	// Header with avatar
	b.WriteString(`<div class="mood-inline-header">`)
	if c.AvatarURL != "" {
		fmt.Fprintf(b, `<img src="%s" style="width:20px;height:20px;border-radius:3px">`, template.HTMLEscapeString(c.AvatarURL))
	}
	fmt.Fprintf(b, `<strong>%s</strong>`, template.HTMLEscapeString(c.Author))
	b.WriteString(`</div>`)

	// Body
	b.WriteString(`<div class="mood-inline-body">`)
	b.WriteString(`<div class="phabricator-remarkup">`)
	b.WriteString(remarkup.Render(c.Body))
	b.WriteString(`</div></div>`)

	// Action buttons
	b.WriteString(`<div class="mood-inline-actions">`)
	b.WriteString(`<a class="inline-action-reply"><span class="phui-icon-view phui-font-fa fa-reply"></span> Reply</a>`)
	b.WriteString(`<a class="inline-action-done"><span class="phui-icon-view phui-font-fa fa-check"></span> Done</a>`)
	b.WriteString(`</div>`)

	b.WriteString(`</div>`)
	b.WriteString(`</td></tr>`)
}

func renderRow(b *strings.Builder, csID int, row DiffRow) {
	b.WriteString(`<tr>`)

	// Old line number — include change class like the PHP renderer
	if row.OldNum > 0 {
		cls := "n"
		if row.OldClass != "" {
			cls = row.OldClass + " n"
		}
		fmt.Fprintf(b, `<td class="%s" data-n="%d" id="C%dOL%d"></td>`, cls, row.OldNum, csID, row.OldNum)
	} else {
		b.WriteString(`<td class="n"></td>`)
	}

	// Old content — always include data-copy-mode
	if row.OldClass != "" {
		fmt.Fprintf(b, `<td class="%s" data-copy-mode="copy-l">%s</td>`, row.OldClass, row.OldContent)
	} else {
		fmt.Fprintf(b, `<td data-copy-mode="copy-l">%s</td>`, row.OldContent)
	}

	// New line number — include change class
	if row.NewNum > 0 {
		cls := "n"
		if row.NewClass != "" {
			cls = row.NewClass + " n"
		}
		fmt.Fprintf(b, `<td class="%s" data-n="%d" id="C%dNL%d"></td>`, cls, row.NewNum, csID, row.NewNum)
	} else {
		b.WriteString(`<td class="n"></td>`)
	}

	// Copy column
	b.WriteString(`<td class="copy"></td>`)

	// New content — colspan=2 merges with coverage column (no separate cov cell)
	if row.NewClass != "" {
		fmt.Fprintf(b, `<td class="%s" colspan="2" data-copy-mode="copy-r">%s</td>`, row.NewClass, row.NewContent)
	} else {
		fmt.Fprintf(b, `<td colspan="2" data-copy-mode="copy-r">%s</td>`, row.NewContent)
	}

	b.WriteString(`</tr>`)
}

// buildRowsByHunk converts a Changeset into rows grouped by hunk for the two-up view.
// It pairs removed+added lines as modifications when they appear consecutively.
func buildRowsByHunk(cs Changeset) [][]DiffRow {
	// Collect all lines from all hunks, highlighting each side.
	oldLines, newLines := collectSides(cs)
	oldHL := HighlightLines(cs.DisplayPath(), oldLines)
	newHL := HighlightLines(cs.DisplayPath(), newLines)

	oldIdx, newIdx := 0, 0
	result := make([][]DiffRow, len(cs.Hunks))

	for hunkIdx, hunk := range cs.Hunks {
		var rows []DiffRow
		// Group consecutive removed/added lines for pairing.
		i := 0
		for i < len(hunk.Lines) {
			line := hunk.Lines[i]
			switch line.Type {
			case Context:
				rows = append(rows, DiffRow{
					OldNum:     line.OldNum,
					NewNum:     line.NewNum,
					OldContent: template.HTML(oldHL[oldIdx]),
					NewContent: template.HTML(newHL[newIdx]),
					IsContext:  true,
				})
				oldIdx++
				newIdx++
				i++

			case Removed:
				// Collect consecutive removed lines.
				var removed []int
				for i < len(hunk.Lines) && hunk.Lines[i].Type == Removed {
					removed = append(removed, i)
					i++
				}
				// Collect consecutive added lines that follow.
				var added []int
				for i < len(hunk.Lines) && hunk.Lines[i].Type == Added {
					added = append(added, i)
					i++
				}
				// Pair them up row by row.
				maxPairs := max(len(removed), len(added))
				for j := 0; j < maxPairs; j++ {
					row := DiffRow{}
					if j < len(removed) {
						rl := hunk.Lines[removed[j]]
						row.OldNum = rl.OldNum
						row.OldContent = template.HTML(oldHL[oldIdx])
						oldIdx++
						if j < len(added) {
							row.OldClass = "old" // paired with added on this row
						} else {
							row.OldClass = "old old-full" // no counterpart
						}
					}
					if j < len(added) {
						al := hunk.Lines[added[j]]
						row.NewNum = al.NewNum
						row.NewContent = template.HTML(newHL[newIdx])
						newIdx++
						if j < len(removed) {
							row.NewClass = "new" // paired with removed on this row
						} else {
							row.NewClass = "new new-full" // no counterpart
						}
					}
					rows = append(rows, row)
				}

			case Added:
				// Added lines not preceded by removed.
				rows = append(rows, DiffRow{
					NewNum:     line.NewNum,
					NewClass:   "new new-full",
					NewContent: template.HTML(newHL[newIdx]),
				})
				newIdx++
				i++
			}
		}
		result[hunkIdx] = rows
	}
	return result
}

// collectSides extracts the old-side and new-side line contents in order,
// for use with the syntax highlighter.
func collectSides(cs Changeset) (oldLines, newLines []string) {
	for _, hunk := range cs.Hunks {
		for _, line := range hunk.Lines {
			switch line.Type {
			case Context:
				oldLines = append(oldLines, line.Content)
				newLines = append(newLines, line.Content)
			case Removed:
				oldLines = append(oldLines, line.Content)
			case Added:
				newLines = append(newLines, line.Content)
			}
		}
	}
	return
}

// RenderFileTree produces a Phabricator-style file tree for the left sidebar.
// Uses diff-tree-view CSS classes that are part of differential.pkg.css.
func RenderFileTree(changesets []Changeset) template.HTML {
	var b strings.Builder

	b.WriteString(`<div class="diff-tree-view">`)

	for _, cs := range changesets {
		path := cs.DisplayPath()
		icon := FileIcon(path)

		// Show just the filename, with directory prefix dimmed.
		dir, file := splitPath(path)

		fmt.Fprintf(&b, `<div class="diff-tree-path diff-tree-path-changeset" data-changeset-id="C%d">`, cs.ID)
		b.WriteString(`<div class="diff-tree-path-indent">`)
		fmt.Fprintf(&b, `<div class="diff-tree-path-icon diff-tree-path-icon-kind"><span class="phui-icon-view phui-font-fa %s"></span></div>`, icon)
		b.WriteString(`<div class="diff-tree-path-name">`)
		fmt.Fprintf(&b, `<a href="#C%d" style="text-decoration:none;color:inherit">`, cs.ID)
		// Status indicator dot
		switch {
		case cs.IsNew || (cs.OldName == "" && cs.NewName != ""):
			b.WriteString(`<span style="color:#2EA86B;margin-right:4px">●</span>`)
		case cs.IsDeleted || cs.NewName == "" || cs.NewName == "/dev/null":
			b.WriteString(`<span style="color:#C0392B;margin-right:4px">●</span>`)
		default:
			b.WriteString(`<span style="color:#E8A83E;margin-right:4px">●</span>`)
		}
		if dir != "" {
			fmt.Fprintf(&b, `<span style="opacity:0.65">%s</span>`, template.HTMLEscapeString(dir))
		}
		b.WriteString(template.HTMLEscapeString(file))
		b.WriteString(`</a>`)
		b.WriteString(`</div>`)

		// Inline count badge
		if cs.LinesAdded > 0 || cs.LinesRemoved > 0 {
			fmt.Fprintf(&b, `<div class="diff-tree-path-inlines diff-tree-path-inlines-visible">`)
			if cs.LinesAdded > 0 {
				fmt.Fprintf(&b, `<span class="green" style="color:#2EA86B;background:rgba(46,168,107,.1);padding:0 4px;border-radius:2px;font-size:11px">+%d</span>`, cs.LinesAdded)
			}
			if cs.LinesRemoved > 0 {
				if cs.LinesAdded > 0 {
					b.WriteString(` `)
				}
				fmt.Fprintf(&b, `<span class="red" style="color:#C0392B;background:rgba(192,57,43,.1);padding:0 4px;border-radius:2px;font-size:11px">-%d</span>`, cs.LinesRemoved)
			}
			b.WriteString(`</div>`)
		}

		b.WriteString(`</div>`) // diff-tree-path-indent
		b.WriteString(`</div>`) // diff-tree-path
	}

	b.WriteString(`</div>`)

	return template.HTML(b.String())
}

// splitPath splits a file path into directory prefix and filename.
func splitPath(path string) (dir, file string) {
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '/' {
			return path[:i+1], path[i+1:]
		}
	}
	return "", path
}
