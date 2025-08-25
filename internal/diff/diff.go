package diff

import (
	"strconv"
	"strings"
)

type LineType int

const (
	ContextLine LineType = iota
	AddedLine
	RemovedLine
	HeaderLine
)

type DiffLine struct {
	Type    LineType
	Number  int
	Content string
}

type Hunk struct {
	StartLine int
	Lines     []DiffLine
}

type FileDiff struct {
	FileName string
	Hunks    []Hunk
}

func ExtractRemovedLinesContent(lines []string) []string {
	var removedLines []string
	for _, line := range lines {
		if strings.HasPrefix(line, "-") && !strings.HasPrefix(line, "---") {
			content := strings.TrimPrefix(line, "-")
			normalized := NormalizeContent(content)
			if normalized != "" {
				removedLines = append(removedLines, normalized)
			}
		}
	}
	return removedLines
}

func NormalizeContent(content string) string {
	content = strings.ReplaceAll(content, "\\ No newline at end of file", "")
	return strings.TrimSpace(content)
}

func IsReformattedLine(content string, removedLines []string) bool {
	if content == "" {
		return false
	}
	for _, removedContent := range removedLines {
		if content == removedContent {
			return true
		}
	}
	return false
}

func IsContextLine(line string) bool {
	return !strings.HasPrefix(line, "@@") &&
		!strings.HasPrefix(line, "+++") &&
		!strings.HasPrefix(line, "---") &&
		!strings.HasPrefix(line, "diff") &&
		!strings.HasPrefix(line, "index") &&
		line != ""
}

func ParseUnifiedDiff(diff, file string) FileDiff {
	lines := strings.Split(diff, "\n")
	var hunks []Hunk
	var currentHunk *Hunk
	var currentLine int

	for _, line := range lines {
		if strings.HasPrefix(line, "@@") {
			if currentHunk != nil {
				hunks = append(hunks, *currentHunk)
			}
			parts := strings.Split(line, " ")
			if len(parts) >= 3 {
				newRange := strings.TrimPrefix(parts[2], "+")
				rangeParts := strings.Split(newRange, ",")
				start, _ := strconv.Atoi(rangeParts[0])
				currentLine = start
				currentHunk = &Hunk{StartLine: start}
			}
			continue
		}
		if currentHunk == nil {
			continue
		}
		if strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "+++") {
			currentHunk.Lines = append(currentHunk.Lines, DiffLine{Type: AddedLine, Number: currentLine, Content: strings.TrimPrefix(line, "+")})
			currentLine++
		} else if strings.HasPrefix(line, "-") && !strings.HasPrefix(line, "---") {
			currentHunk.Lines = append(currentHunk.Lines, DiffLine{Type: RemovedLine, Number: 0, Content: strings.TrimPrefix(line, "-")})
		} else {
			currentHunk.Lines = append(currentHunk.Lines, DiffLine{Type: ContextLine, Number: currentLine, Content: line})
			currentLine++
		}
	}

	if currentHunk != nil {
		hunks = append(hunks, *currentHunk)
	}
	return FileDiff{FileName: file, Hunks: hunks}
}

func (fd FileDiff) GetAddedLines() []DiffLine {
	var added []DiffLine
	for _, h := range fd.Hunks {
		for _, l := range h.Lines {
			if l.Type == AddedLine {
				added = append(added, l)
			}
		}
	}
	return added
}
