package cmd

import (
	"fmt"

	"github.com/abhi-g80/chipku/server"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(version)
}

var version = &cobra.Command{
	Use:   "version",
	Short: "Print version",
	Long:  `Print version number for Chipku`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s\n", server.Version)
	},
}
