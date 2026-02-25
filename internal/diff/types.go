package diff

import "html/template"

// LineType classifies a diff line.
type LineType int

const (
	Context LineType = iota
	Added
	Removed
)

// Line represents a single line in a diff hunk.
type Line struct {
	Type    LineType
	OldNum  int    // 0 if added
	NewNum  int    // 0 if removed
	Content string // raw text (no +/- prefix)
}

// Hunk is a contiguous group of diff lines.
type Hunk struct {
	OldStart int // starting line number on the old side
	NewStart int // starting line number on the new side
	OldCount int // number of old-side lines covered
	NewCount int // number of new-side lines covered
	Lines    []Line
}

// Changeset represents a single file's diff.
type Changeset struct {
	ID           int
	OldName      string
	NewName      string
	LinesAdded   int
	LinesRemoved int
	IsNew        bool
	IsDeleted    bool
	IsRenamed    bool
	IsBinary     bool
	Hunks        []Hunk
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
	OldNum     int           // 0 = no line number
	NewNum     int           // 0 = no line number
	OldClass   string        // "old", "old-full", or ""
	NewClass   string        // "new", "new-full", or ""
	OldContent template.HTML // syntax-highlighted HTML
	NewContent template.HTML // syntax-highlighted HTML
	IsContext  bool
}
