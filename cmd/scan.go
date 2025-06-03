package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

type Leak struct {
	RuleID      string   `json:"RuleID"`
	Description string   `json:"Description"`
	StartLine   int      `json:"StartLine"`
	EndLine     int      `json:"EndLine"`
	StartColumn int      `json:"StartColumn"`
	EndColumn   int      `json:"EndColumn"`
	Match       string   `json:"Match"`
	Secret      string   `json:"Secret"`
	File        string   `json:"File"`
	Commit      string   `json:"Commit"`
	Author      string   `json:"Author"`
	Email       string   `json:"Email"`
	Date        string   `json:"Date"`
	Message     string   `json:"Message"`
	Tags        []string `json:"Tags"`
	Fingerprint string   `json:"Fingerprint"`
}

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "A fast and flexible scanner",
	Long:  `Scan is a fast and flexible scanner that can be used for various purposes.`,
	Run:   Scan,
}

func Scan(cmd *cobra.Command, args []string) {
	fmt.Println("Executing scan command...")
	command := exec.Command("git", "diff", "--name-only", "--staged")
	filesOutput, err := command.CombinedOutput()
	if err != nil {
		log.Fatal("Error on getting staged files:", err)
	}

	files := strings.Split(string(filesOutput), "\n")

	for _, file := range files {
		if file != "" {
			fmt.Printf("Verfifying file: %s\n", file)
			err := checkSecretsInFile(file)
			if err != nil {
				log.Println("Error checking secrets in file:", file, err)
			}
		}
	}
}

func checkSecretsInFile(file string) error {
	cmdDiff := exec.Command("git", "diff", "--staged", file)
	diffOutput, err := cmdDiff.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Error in getting diff for file %s: %v", file, err)
	}

	tempFile, err := os.CreateTemp("", "gitleaks-*.json")
	if err != nil {
		return fmt.Errorf("Error creating temporary file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	cmdGitleaks := exec.Command("gitleaks", "stdin", "-f=json", "-r="+tempFile.Name())
	cmdGitleaks.Stdin = strings.NewReader(string(diffOutput))

	err = cmdGitleaks.Run()

	jsonData, err := os.ReadFile(tempFile.Name())
	if err != nil {
		return fmt.Errorf("Error reading temporary file %s: %v", tempFile.Name(), err)
	}

	var leaks []Leak
	err = json.Unmarshal(jsonData, &leaks)
	if err != nil {
		return fmt.Errorf("Error parsing gitleaks JSON output: %v", err)
	}

	if len(leaks) > 0 {
		fmt.Printf("Secrets found in file: %s:\n", file)
		for _, leak := range leaks {
			fmt.Printf("Rule: %s | Secret: %s | Line: %d-%d | Column: %d-%d\n",
				leak.RuleID, leak.Secret, leak.StartLine, leak.EndLine, leak.StartColumn, leak.EndColumn)
			fmt.Printf("Description: %s\n", leak.Description)
			fmt.Printf("Match: %s\n", leak.Match)
		}
	} else {
		fmt.Printf("No secrets found in file: %s.\n", file)
	}

	return nil
}
