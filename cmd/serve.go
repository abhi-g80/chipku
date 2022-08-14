package cmd

import (
	"strconv"

	"github.com/spf13/cobra"

	"github.com/abhi-g80/chipku/server"
)

var port int

func init() {
	rootCmd.AddCommand(serve)
	serve.Flags().IntVarP(&port, "port", "p", 8080, "port to serve on")
}

var serve = &cobra.Command{
	Use:   "serve",
	Short: "Start server",
	Long:  `Start the Chipku server`,
	Run: func(cmd *cobra.Command, args []string) {
		p := strconv.Itoa(port)
		server.Serve(p)
	},
}
