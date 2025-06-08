package format

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

var (
	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#00BFFF"))

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FF00")).
			Bold(true)

	warningStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFA500"))

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF4136")).
			Bold(true)

	fileStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7D56F4")).
			Italic(true)

	indent = lipgloss.NewStyle().
		MarginLeft(2)
)

func FormatHeader(text string) string {
	return headerStyle.Render(text)
}

func FormatSuccess(text string) string {
	return successStyle.Render(text)
}

func FormatWarning(text string) string {
	return warningStyle.Render(text)
}

func FormatError(text string) string {
	return errorStyle.Render(text)
}

func FormatFileName(fileName string) string {
	return fileStyle.Render(fileName)
}

func FormatIndented(text string) string {
	return indent.Render(text)
}

func FormatFileList(files []string) string {
	if len(files) == 0 {
		return "No files found"
	}

	formattedFiles := make([]string, len(files))
	for i, file := range files {
		formattedFiles[i] = FormatFileName(file)
	}

	return "Files:\n" + indent.Render(fmt.Sprintf("- %s", formattedFiles))
}

func FormatLineNumber(lineNumber int) string {
	if lineNumber < 0 {
		return "N/A"
	}
	return fmt.Sprintf("Line %d", lineNumber)
}
