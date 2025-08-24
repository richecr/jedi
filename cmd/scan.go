package cmd

import (
	"fmt"
	"os"

	"github.com/richecr/jedi-scan/pkg/jediscan"
	"github.com/spf13/cobra"
)

const tempFilePattern = "gitleaks-*.json"

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan staged files for secrets using gitleaks",
	Long:  `Scan staged git files for secrets and sensitive information before committing.`,
	Run:   Scan,
}

func Scan(cmd *cobra.Command, args []string) {
	fmt.Println("Executing scan command...")

	jediScan, err := jediscan.NewJediScan(tempFilePattern)
	if err != nil {
		fmt.Printf("Error initializing scanner: %v\n", err)
		os.Exit(1)
	}

	err = jediScan.Scan()
	if err != nil {
		fmt.Println("Error during scan:", err)
		os.Exit(1)
	}

	fmt.Println("Scan completed successfully.")
	os.Exit(0)
}
