package git

import (
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
	cmd := exec.Command("git", "diff", "--staged", "-U0", file)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	diff := string(output)
	return diff, nil
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
		}
	}

	return changedLines
}
