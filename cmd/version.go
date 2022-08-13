package cmd

import (
	"fmt"

    "github.com/abhi-g80/chipku/server"
	"github.com/spf13/cobra"
)

// const Version = "0.1.1"

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version",
	Long:  `Print version number for Chipku`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Chipku - a no frill pastebin v%s\n", server.Version)
	},
}
