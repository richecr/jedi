package jediscan

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/richecr/jedi-scan/internal/git"
	"github.com/richecr/jedi-scan/internal/gitleaks"
	"github.com/richecr/jedi-scan/internal/tempfile"
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

type JediScan struct {
	pattern string
}

var jediScan *JediScan

func NewJediScan(pattern string) (*JediScan, error) {
	return &JediScan{
		pattern: pattern,
	}, nil
}

func (j *JediScan) CreateTempFile() (*tempfile.TempFile, error) {
	tempFile, err := tempfile.NewTempFile(j.pattern)
	if err != nil {
		return nil, err
	}

	return tempFile, nil
}

func (j *JediScan) Scan() {
	files, err := git.GetStagedFiles()
	if err != nil {
		fmt.Println(errorStyle.Render("❌ Error getting staged files: ") + errorStyle.Render(err.Error()))
		os.Exit(1)
	}

	exitCode := 0
	for _, file := range files {
		if file != "" {
			fmt.Println(headerStyle.Render("\n🔍 Verifying file: ") + fileStyle.Render(file))

			diff, err := git.GetDiff(file)
			if err != nil {
				fmt.Println(errorStyle.Render("❌ Error getting diff for file: ") + fileStyle.Render(file))
				os.Exit(1)
			}

			tempFile, err := j.CreateTempFile()
			if err != nil {
				fmt.Printf("Error creating temporary file: %v", err)
				os.Exit(1)
			}
			err = gitleaks.Scan(diff, tempFile)

			// if err != nil {
			// 	fmt.Printf("Running gitleaks command: %s\n", err)
			// 	return fmt.Errorf("Error scanning file %s: %v", file, err)
			// }

			leaks, err := gitleaks.UnmarshalResults(tempFile)
			tempFile.Remove()
			if err != nil {
				fmt.Println(errorStyle.Render("❌ Error parsing gitleaks JSON output: ") + fileStyle.Render(file))
				os.Exit(1)
			}

			if len(leaks) > 0 {
				exitCode = 1
				fmt.Println(errorStyle.Render("🚨 Secrets found in file: ") + fileStyle.Render(file))
				for _, leak := range leaks {
					fmt.Println(indent.Render(fmt.Sprintf("• Rule: %s", warningStyle.Render(leak.RuleID))))
					fmt.Println(indent.Render(fmt.Sprintf("  Secret: %s", leak.Secret)))
					fmt.Println(indent.Render(fmt.Sprintf("  Line: %d-%d | Column: %d-%d", leak.StartLine, leak.EndLine, leak.StartColumn, leak.EndColumn)))
					fmt.Println(indent.Render(fmt.Sprintf("  Match: %s", leak.Match)))
					fmt.Println(indent.Render(fmt.Sprintf("  Description: %s", leak.Description)))
				}

			} else {
				fmt.Println(successStyle.Render("✅ No secrets found in file: ") + fileStyle.Render(file))
			}
		}
	}

	os.Exit(exitCode)
}

func GetInstance(pattern string) *JediScan {
	if jediScan == nil {
		jediScan, _ = NewJediScan(pattern)
	}

	return jediScan
}
