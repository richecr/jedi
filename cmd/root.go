package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "scan",
	Short: "A fast and flexible scanner",
	Long:  `Scan is a fast and flexible scanner that can be used for various purposes.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting scan...")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
