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

const (
	ErrGetStagedFiles = "failed to get staged files"
	ErrGetDiff        = "failed to get diff for file"
	ErrCreateTempFile = "failed to create temporary file"
	ErrRunGitleaks    = "failed to run gitleaks on file"
	ErrParseGitleaks  = "failed to parse gitleaks JSON output for file"
	ErrSecretsFound   = "secrets found in staged files"
)

type JediScan struct {
	pattern string
}

func NewJediScan(pattern string) (*JediScan, error) {
	if pattern == "" {
		return nil, errors.New("pattern cannot be empty")
	}
	return &JediScan{
		pattern: pattern,
	}, nil
}

func (j *JediScan) CreateTempFile() (*tempfile.TempFile, error) {
	return tempfile.NewTempFile(j.pattern)
}

func (j *JediScan) Scan() error {
	files, err := git.GetStagedFiles()
	if err != nil {
		fmt.Println(format.FormatError("❌ Error getting staged files: " + err.Error()))
		return errors.New(ErrGetStagedFiles)
	}

	fmt.Println(format.FormatFileList(files))

	hasSecrets := false
	for _, file := range files {
		if strings.TrimSpace(file) == "" {
			continue
		}

		foundSecrets, err := j.scanFile(file)
		if err != nil {
			fmt.Println(format.FormatError("❌ Error scanning file: ") + format.FormatFileName(file))
			return err
		}

		if foundSecrets {
			hasSecrets = true
		}
	}

	if hasSecrets {
		fmt.Println(format.FormatError("🚨 Secrets found in staged files. Please review and fix them before committing."))
		return errors.New(ErrSecretsFound)
	}

	fmt.Println(format.FormatSuccess("✅ No secrets found in staged files. You can proceed with the commit."))
	return nil
}

func (j *JediScan) scanFile(file string) (bool, error) {
	fmt.Println(format.FormatHeader("\n🔍 Verifying file: ") + format.FormatFileName(file))

	originalDiff, err := git.GetOriginalDiff(file)
	if err != nil {
		fmt.Println(format.FormatError("❌ Error getting original diff for file: ") + format.FormatFileName(file))
		return false, fmt.Errorf("%s: %s", ErrGetDiff, file)
	}

	filteredDiff, err := git.GetDiffWithChangedLines(originalDiff, file)
	if err != nil {
		fmt.Println(format.FormatError("❌ Error getting diff for file: ") + format.FormatFileName(file))
		return false, fmt.Errorf("%s: %s", ErrGetDiff, file)
	}

	tempFile, err := j.CreateTempFile()
	if err != nil {
		fmt.Printf("Error creating temporary file: %v", err)
		return false, fmt.Errorf("%s: %w", ErrCreateTempFile, err)
	}
	defer tempFile.Remove()

	err = gitleaks.Scan(filteredDiff, tempFile)
	if err != nil && err.Error() != "exit status 1" {
		fmt.Println(format.FormatError("❌ Error running gitleaks: ") + format.FormatFileName(file))
		return false, fmt.Errorf("%s: %s", ErrRunGitleaks, file)
	}

	leaks, err := gitleaks.UnmarshalResults(tempFile)
	if err != nil {
		fmt.Println(format.FormatError("❌ Error parsing gitleaks JSON output: ") + format.FormatFileName(file))
		return false, fmt.Errorf("%s: %s", ErrParseGitleaks, file)
	}

	changedLines := git.ExtractTrulyNewLines(originalDiff)

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
