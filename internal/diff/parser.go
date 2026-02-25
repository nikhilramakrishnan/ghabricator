package diff

import (
	"strings"

	godiff "github.com/sourcegraph/go-diff/diff"
)

// ParseDiff parses a unified diff string into structured Changesets.
func ParseDiff(raw string) ([]Changeset, error) {
	fileDiffs, err := godiff.ParseMultiFileDiff([]byte(raw))
	if err != nil {
		return nil, err
	}

	changesets := make([]Changeset, 0, len(fileDiffs))
	for i, fd := range fileDiffs {
		cs := Changeset{
			ID:      i + 1,
			OldName: cleanPath(fd.OrigName),
			NewName: cleanPath(fd.NewName),
		}

		if cs.OldName == "/dev/null" {
			cs.IsNew = true
		}
		if cs.NewName == "/dev/null" {
			cs.IsDeleted = true
		}
		if !cs.IsNew && !cs.IsDeleted && cs.OldName != cs.NewName {
			cs.IsRenamed = true
		}
		for _, ext := range fd.Extended {
			if strings.Contains(ext, "Binary files") || strings.Contains(ext, "GIT binary patch") {
				cs.IsBinary = true
				break
			}
		}

		for _, h := range fd.Hunks {
			hunk := parseHunk(h)
			for _, l := range hunk.Lines {
				switch l.Type {
				case Added:
					cs.LinesAdded++
				case Removed:
					cs.LinesRemoved++
				}
			}
			cs.Hunks = append(cs.Hunks, hunk)
		}

		changesets = append(changesets, cs)
	}
	return changesets, nil
}

func parseHunk(h *godiff.Hunk) Hunk {
	bodyLines := strings.Split(string(h.Body), "\n")
	oldNum := int(h.OrigStartLine)
	newNum := int(h.NewStartLine)

	hunk := Hunk{
		OldStart: int(h.OrigStartLine),
		NewStart: int(h.NewStartLine),
		OldCount: int(h.OrigLines),
		NewCount: int(h.NewLines),
	}

	for _, raw := range bodyLines {
		if raw == "" {
			continue
		}
		prefix := raw[0]
		content := ""
		if len(raw) > 1 {
			content = raw[1:]
		}

		switch prefix {
		case '+':
			hunk.Lines = append(hunk.Lines, Line{
				Type:    Added,
				NewNum:  newNum,
				Content: content,
			})
			newNum++
		case '-':
			hunk.Lines = append(hunk.Lines, Line{
				Type:    Removed,
				OldNum:  oldNum,
				Content: content,
			})
			oldNum++
		case ' ':
			hunk.Lines = append(hunk.Lines, Line{
				Type:    Context,
				OldNum:  oldNum,
				NewNum:  newNum,
				Content: content,
			})
			oldNum++
			newNum++
		case '\\':
			// "\ No newline at end of file" — skip
		}
	}
	return hunk
}

// cleanPath strips the a/ or b/ prefix from diff paths.
func cleanPath(p string) string {
	if strings.HasPrefix(p, "a/") || strings.HasPrefix(p, "b/") {
		return p[2:]
	}
	return p
}
