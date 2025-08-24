package git

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/richecr/jedi-scan/internal/model"
)

func GetStagedFiles() ([]string, error) {
	cmd := exec.Command("git", "diff", "--name-only", "--staged")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	files := strings.Split(string(output), "\n")
	return files, nil
}

func GetDiff(file string) (string, error) {
	cmd := exec.Command("git", "diff", "--staged", file)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	return string(output), nil
}

func GetDiffWithChangedLines(file string) (string, error) {
	diff, err := getDiffRaw(file)
	if err != nil {
		return "", err
	}

	return filterDiffForNewLines(diff, file), nil
}

func GetOriginalDiff(file string) (string, error) {
	return getDiffRaw(file)
}

func getDiffRaw(file string) (string, error) {
	cmd := exec.Command("git", "diff", "--staged", "--no-prefix", "-U0", file)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func filterDiffForNewLines(diff, file string) string {
	lines := strings.Split(diff, "\n")
	var result []string

	result = append(result, "diff --git a/"+file+" b/"+file)
	result = append(result, "index 0000000..1111111 100644")
	result = append(result, "--- a/"+file)
	result = append(result, "+++ b/"+file)

	var removedLines []string
	var addedLines []string

	for _, line := range lines {
		if strings.HasPrefix(line, "-") && !strings.HasPrefix(line, "---") {
			content := strings.TrimPrefix(line, "-")
			content = strings.ReplaceAll(content, "\\ No newline at end of file", "")
			content = strings.TrimSpace(content)
			if content != "" {
				removedLines = append(removedLines, content)
			}
		} else if strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "+++") {
			content := strings.TrimPrefix(line, "+")
			content = strings.ReplaceAll(content, "\\ No newline at end of file", "")
			content = strings.TrimSpace(content)
			if content != "" {
				addedLines = append(addedLines, "+"+content)
			} else {
				addedLines = append(addedLines, line)
			}
		}
	}

	var trulyNewLines []string
	for _, addedLine := range addedLines {
		addedContent := strings.TrimPrefix(addedLine, "+")
		addedContent = strings.TrimSpace(addedContent)

		isReformat := false
		for _, removedContent := range removedLines {
			if addedContent == removedContent {
				isReformat = true
				break
			}
		}

		if !isReformat {
			trulyNewLines = append(trulyNewLines, addedLine)
		}
	}

	if len(trulyNewLines) > 0 {
		result = append(result, fmt.Sprintf("@@ -0,0 +1,%d @@", len(trulyNewLines)))
		result = append(result, trulyNewLines...)
	}

	return strings.Join(result, "\n")
}

func ExtractTrulyNewLines(diff string) []model.ChangedLine {
	lines := strings.Split(diff, "\n")
	var changedLines []model.ChangedLine
	var currentLine int

	removedLines := extractRemovedLinesContent(lines)
	hasReformattedLines := len(removedLines) > 0

	for _, line := range lines {
		if strings.HasPrefix(line, "@@") {
			currentLine = parseHunkStartLine(line)
		} else if strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "+++") {
			content := strings.TrimPrefix(line, "+")
			originalContent := content

			normalizedContent := normalizeContent(content)

			if !isReformattedLine(normalizedContent, removedLines) {
				lineNumber := currentLine

				if hasReformattedLines && currentLine > 1 {
					lineNumber = currentLine - 1
				}

				changedLines = append(changedLines, model.ChangedLine{
					Number:  lineNumber,
					Content: originalContent,
				})
			}
			currentLine++
		} else if strings.HasPrefix(line, "-") && !strings.HasPrefix(line, "---") {
			continue
		} else if isContextLine(line) {
			currentLine++
		}
	}

	return changedLines
}

func extractRemovedLinesContent(lines []string) []string {
	var removedLines []string
	for _, line := range lines {
		if strings.HasPrefix(line, "-") && !strings.HasPrefix(line, "---") {
			content := strings.TrimPrefix(line, "-")
			normalized := normalizeContent(content)
			if normalized != "" {
				removedLines = append(removedLines, normalized)
			}
		}
	}
	return removedLines
}

func parseHunkStartLine(line string) int {
	parts := strings.Split(line, " ")
	if len(parts) >= 3 {
		newRange := strings.TrimPrefix(parts[2], "+")
		rangeParts := strings.Split(newRange, ",")
		if start, err := strconv.Atoi(rangeParts[0]); err == nil {
			return start
		}
	}
	return 1
}

func normalizeContent(content string) string {
	content = strings.ReplaceAll(content, "\\ No newline at end of file", "")
	return strings.TrimSpace(content)
}

func isReformattedLine(content string, removedLines []string) bool {
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

func isContextLine(line string) bool {
	return !strings.HasPrefix(line, "@@") &&
		!strings.HasPrefix(line, "+++") &&
		!strings.HasPrefix(line, "---") &&
		!strings.HasPrefix(line, "diff") &&
		!strings.HasPrefix(line, "index") &&
		line != ""
}

func ExtractChangedLines(diff string) []model.ChangedLine {
	lines := strings.Split(diff, "\n")
	var changedLines []model.ChangedLine
	var currentLine int

	for _, line := range lines {
		if strings.HasPrefix(line, "@@") {
			parts := strings.Split(line, " ")
			if len(parts) >= 3 {
				newRange := parts[2]
				newRange = strings.TrimPrefix(newRange, "+")
				rangeParts := strings.Split(newRange, ",")
				start, _ := strconv.Atoi(rangeParts[0])
				currentLine = start
			}
		} else if strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "+++") {
			content := strings.TrimPrefix(line, "+")
			if strings.TrimSpace(content) != "" {
				changedLines = append(changedLines, model.ChangedLine{
					Number:  currentLine,
					Content: content,
				})
			}
			currentLine++
		} else if strings.HasPrefix(line, "-") && !strings.HasPrefix(line, "---") {
			continue
		} else if !strings.HasPrefix(line, "@@") && !strings.HasPrefix(line, "+++") && !strings.HasPrefix(line, "---") && !strings.HasPrefix(line, "diff") && !strings.HasPrefix(line, "index") && line != "" {
			currentLine++
		}
	}

	return changedLines
}
