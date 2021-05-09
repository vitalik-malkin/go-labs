package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:           "",
	Short:         "Very simple backup tool",
	SilenceUsage:  true,
	SilenceErrors: true,
}

func init() {
	// versionCmd
	rootCmd.AddCommand(versionCmd)

	// makeCmd
	rootCmd.AddCommand(makeCmd)
	makeCmd.Flags().AddFlagSet(makeCmdConfigVar.BindFlags(""))
}

func Exec() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Printf("=======\nERROR\n%s\n", err)
		os.Exit(1)
	}
}
