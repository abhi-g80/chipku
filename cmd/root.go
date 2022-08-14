package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "chipku",
	Short: "Chipku is a no frill pastebin",
	Long: `A fast and reliable pastebin server.

May be used for sharing snippets with your loved ones and colleagues.
Partial documentation is available at http://github.com/abhi-go/chipku`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("And apparently, this is fun ?!")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
