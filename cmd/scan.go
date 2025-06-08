package cmd

import (
	"fmt"
	"os"

	"github.com/richecr/jedi-scan/pkg/jediscan"
	"github.com/spf13/cobra"
)

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "A fast and flexible scanner",
	Long:  `Scan is a fast and flexible scanner that can be used for various purposes.`,
	Run:   Scan,
}

var jediScan *jediscan.JediScan
var tempFilePattern = "gitleaks-*.json"

func init() {
	jediScan = jediscan.GetInstance(tempFilePattern)
}

func Scan(cmd *cobra.Command, args []string) {
	fmt.Println("Executing scan command...")
	err := jediScan.Scan()
	if err != nil {
		fmt.Println("Error during scan:", err)
		os.Exit(1)
	}

	fmt.Println("Scan completed successfully.")
	os.Exit(0)
}
