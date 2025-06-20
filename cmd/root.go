package cmd

import (
	"context"
	"os"

	"github.com/charmbracelet/fang"
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
	opts := fang.WithVersion("0.1.0")

	if err := fang.Execute(context.TODO(), rootCmd, opts); err != nil {
		os.Exit(1)
	}
}
