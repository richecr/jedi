package gitleaks

import (
	"encoding/json"
	"os/exec"
	"strings"

	"github.com/richecr/jedi-scan/internal/tempfile"
)

func Scan(diff string, tempFile *tempfile.TempFile) error {
	cmd := exec.Command("gitleaks", "stdin", "-f=json", "-r="+tempFile.GetFileName())
	cmd.Stdin = strings.NewReader(diff)

	err := cmd.Run()
	return err
}

func UnmarshalResults(tempFile *tempfile.TempFile) ([]Leak, error) {
	data, err := tempFile.ReadFile()
	if err != nil {
		return nil, err
	}

	var leaks []Leak
	err = json.Unmarshal(data, &leaks)
	if err != nil {
		return nil, err
	}

	return leaks, nil
}
