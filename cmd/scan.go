package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "A fast and flexible scanner",
	Long:  `Scan is a fast and flexible scanner that can be used for various purposes.`,
	Run:   Scan,
}

func Scan(cmd *cobra.Command, args []string) {
	fmt.Println("Executing scan command...")
}
