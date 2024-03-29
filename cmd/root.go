package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// RootCmd the root command for chipku
var RootCmd = &cobra.Command{
	Use:   "chipku",
	Short: "Chipku is a no frill pastebin",
	Long: `A fast and reliable pastebin server.

May be used for sharing snippets with your loved ones and colleagues.
Partial documentation is available at http://github.com/abhi-g80/chipku`,
	Run: func(cmd *cobra.Command, args []string) {
		// print cmd help if no serve subcommand isn't invoked
		err := cmd.Help()
		if err != nil {
			os.Exit(1)
		}
	},
}

// Execute try to run the root command
func Execute() error {
	if err := RootCmd.Execute(); err != nil {
		return fmt.Errorf("root cmd: %v", err)
	}
	return nil
}
