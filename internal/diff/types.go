package diff

import (
	"encoding/json"
	"html/template"
)

// LineType classifies a diff line.
type LineType int

const (
	Context LineType = iota
	Added
	Removed
)

func (lt LineType) String() string {
	switch lt {
	case Context:
		return "context"
	case Added:
		return "added"
	case Removed:
		return "removed"
	default:
		return "unknown"
	}
}

func (lt LineType) MarshalJSON() ([]byte, error) {
	return json.Marshal(lt.String())
}

// Line represents a single line in a diff hunk.
type Line struct {
	Type    LineType `json:"type"`
	OldNum  int      `json:"oldNum"`  // 0 if added
	NewNum  int      `json:"newNum"`  // 0 if removed
	Content string   `json:"content"` // raw text (no +/- prefix)
}

// Hunk is a contiguous group of diff lines.
type Hunk struct {
	OldStart int    `json:"oldStart"` // starting line number on the old side
	NewStart int    `json:"newStart"` // starting line number on the new side
	OldCount int    `json:"oldCount"` // number of old-side lines covered
	NewCount int    `json:"newCount"` // number of new-side lines covered
	Lines    []Line `json:"lines"`
}

// Changeset represents a single file's diff.
type Changeset struct {
	ID           int    `json:"id"`
	OldName      string `json:"oldName"`
	NewName      string `json:"newName"`
	LinesAdded   int    `json:"linesAdded"`
	LinesRemoved int    `json:"linesRemoved"`
	IsNew        bool   `json:"isNew"`
	IsDeleted    bool   `json:"isDeleted"`
	IsRenamed    bool   `json:"isRenamed"`
	IsBinary     bool   `json:"isBinary"`
	Hunks        []Hunk `json:"hunks"`
}

// DisplayPath returns the best path to show for this changeset.
func (c *Changeset) DisplayPath() string {
	if c.NewName != "" && c.NewName != "/dev/null" {
		return c.NewName
	}
	return c.OldName
}

// DiffRow is a single row in the two-up diff table.
type DiffRow struct {
	OldNum     int           `json:"oldNum"`     // 0 = no line number
	NewNum     int           `json:"newNum"`     // 0 = no line number
	OldClass   string        `json:"oldClass"`   // "old", "old-full", or ""
	NewClass   string        `json:"newClass"`   // "new", "new-full", or ""
	OldContent template.HTML `json:"oldContent"` // syntax-highlighted HTML
	NewContent template.HTML `json:"newContent"` // syntax-highlighted HTML
	IsContext  bool          `json:"isContext"`
}
