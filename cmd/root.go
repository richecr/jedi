package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(scanCmd)
}

var rootCmd = &cobra.Command{
	Use:   "jedi-scan",
	Short: "JediScan is a security tool that uses gitleaks to scan staged code for sensitive information ",
	Long:  `JediScan is a security tool that uses gitleaks to scan staged code for sensitive information like API keys and tokens, ensuring your code stays secure before being committed.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
