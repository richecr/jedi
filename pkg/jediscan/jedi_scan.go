package jediscan

import (
	"errors"
	"fmt"
	"strings"

	"github.com/richecr/jedi-scan/internal/format"
	"github.com/richecr/jedi-scan/internal/git"
	"github.com/richecr/jedi-scan/internal/gitleaks"
	"github.com/richecr/jedi-scan/internal/model"
	"github.com/richecr/jedi-scan/internal/tempfile"
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

func (j *JediScan) Scan() error {
	files, err := git.GetStagedFiles()
	if err != nil {
		fmt.Println(format.FormatError("❌ Error getting staged files: " + err.Error()))
		return errors.New("Failed to get staged files")
	}

	fmt.Println(format.FormatFileList(files))

	exitCode := 0
	for _, file := range files {
		if strings.TrimSpace(file) == "" {
			continue
		}
		fmt.Println(format.FormatHeader("\n🔍 Verifying file: ") + format.FormatFileName(file))

		diff, err := git.GetDiffWithChangedLines(file)
		if err != nil {
			fmt.Println(format.FormatError("❌ Error getting diff for file: ") + format.FormatFileName(file))
			return errors.New("Failed to get diff for file: " + file)
		}

		tempFile, err := j.CreateTempFile()
		if err != nil {
			fmt.Printf("Error creating temporary file: %v", err)
			return errors.New("Failed to create temporary file")
		}
		err = gitleaks.Scan(diff, tempFile)

		if err != nil && err.Error() != "exit status 1" {
			fmt.Println(format.FormatError("❌ Error running gitleaks: ") + format.FormatFileName(file))
			return errors.New("Failed to run gitleaks on file: " + file)
		}

		leaks, err := gitleaks.UnmarshalResults(tempFile)
		tempFile.Remove()
		if err != nil {
			fmt.Println(format.FormatError("❌ Error parsing gitleaks JSON output: ") + format.FormatFileName(file))
			return errors.New("Failed to parse gitleaks JSON output for file: " + file)
		}

		changedLines := git.ExtractChangedLines(diff)

		if len(leaks) > 0 {
			exitCode = 1
			fmt.Println(format.FormatError("🚨 Secrets found in file: ") + format.FormatFileName(file))
			for _, leak := range leaks {
				lineNum := gitleaks.FindLineNumberForLeak(changedLines, leak)
				if lineNum == -1 {
					continue
				}
				fmt.Println(format.FormatIndented("• Rule: " + format.FormatWarning(leak.RuleID)))
				fmt.Println(format.FormatIndented("  Secret: " + leak.Secret))
				fmt.Println(format.FormatIndented("  Line: " + format.FormatLineNumber(lineNum)))
				fmt.Println(format.FormatIndented("  Match: " + leak.Match))
				fmt.Println(format.FormatIndented("  Description: " + leak.Description))
			}
		} else {
			fmt.Println(format.FormatSuccess("✅ No secrets found in file: ") + format.FormatFileName(file))
		}
	}

	if exitCode == 1 {
		fmt.Println(format.FormatError("🚨 Secrets found in staged files. Please review and fix them before committing."))
		return errors.New("Secrets found in staged files")
	}

	fmt.Println(format.FormatSuccess("✅ No secrets found in staged files. You can proceed with the commit."))
	return nil
}

func (j *JediScan) scanFile(file string) (bool, error) {
	fmt.Println(format.FormatHeader("\n🔍 Verifying file: ") + format.FormatFileName(file))

	diff, err := git.GetDiffWithChangedLines(file)
	if err != nil {
		fmt.Println(format.FormatError("❌ Error getting diff for file: ") + format.FormatFileName(file))
		return false, errors.New("Failed to get diff for file: " + file)
	}

	tempFile, err := j.CreateTempFile()
	if err != nil {
		fmt.Printf("Error creating temporary file: %v", err)
		return false, errors.New("Failed to create temporary file")
	}
	defer tempFile.Remove()

	err = gitleaks.Scan(diff, tempFile)
	if err != nil && err.Error() != "exit status 1" {
		fmt.Println(format.FormatError("❌ Error running gitleaks: ") + format.FormatFileName(file))
		return false, errors.New("Failed to run gitleaks on file: " + file)
	}

	leaks, err := gitleaks.UnmarshalResults(tempFile)
	if err != nil {
		fmt.Println(format.FormatError("❌ Error parsing gitleaks JSON output: ") + format.FormatFileName(file))
		return false, errors.New("Failed to parse gitleaks JSON output for file: " + file)
	}

	changedLines := git.ExtractChangedLines(diff)
	if len(leaks) > 0 {
		fmt.Println(format.FormatError("🚨 Secrets found in file: ") + format.FormatFileName(file))
		printLeaks(leaks, changedLines)
		return true, nil
	}

	fmt.Println(format.FormatSuccess("✅ No secrets found in file: ") + format.FormatFileName(file))
	return false, nil
}

func printLeaks(leaks []gitleaks.Leak, changedLines []model.ChangedLine) {
	for _, leak := range leaks {
		lineNum := gitleaks.FindLineNumberForLeak(changedLines, leak)
		if lineNum == -1 {
			continue
		}
		fmt.Println(format.FormatIndented("• Rule: " + format.FormatWarning(leak.RuleID)))
		fmt.Println(format.FormatIndented("  Secret: " + leak.Secret))
		fmt.Println(format.FormatIndented("  Line: " + format.FormatLineNumber(lineNum)))
		fmt.Println(format.FormatIndented("  Match: " + leak.Match))
		fmt.Println(format.FormatIndented("  Description: " + leak.Description))
	}
}

func GetInstance(pattern string) *JediScan {
	if jediScan == nil {
		jediScan, _ = NewJediScan(pattern)
	}

	return jediScan
}
