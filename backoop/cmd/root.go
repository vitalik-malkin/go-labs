package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:          "",
	Short:        "Very simple backup tool",
	SilenceUsage: true,
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func Exec() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Printf("=======\nERROR\n%s\n", err)
		os.Exit(1)
	}
}
