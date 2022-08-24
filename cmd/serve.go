package cmd

import (
	"strconv"

	"github.com/abhi-g80/chipku/server"
	"github.com/spf13/cobra"
)

var port int
var debug bool

func init() {
	RootCmd.AddCommand(serve)
	serve.Flags().IntVarP(&port, "port", "p", 8080, "port to serve on")
	serve.Flags().BoolVarP(&debug, "debug", "d", false, "print debug messages")
}

var serve = &cobra.Command{
	Use:   "serve",
	Short: "Start server",
	Long:  `Start the Chipku server`,
	Run: func(cmd *cobra.Command, args []string) {
		p := strconv.Itoa(port)
		d := debug
		server.Serve(p, d)
	},
}
