package cmd

import (
	"fmt"
	_ "fmt"
	_ "os"
	"runtime"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:          "version",
	Short:        "Shows version info",
	SilenceUsage: true,
	RunE:         versionCmdExec,
}

func versionCmdExec(cmd *cobra.Command, args []string) error {
	fmt.Printf("Version: %s (%s, %s)\n", "v1.0", runtime.GOOS, runtime.GOARCH)
	return nil
}
