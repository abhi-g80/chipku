package cmd

import (
	// "fmt"

	"github.com/spf13/cobra"

  "github.com/abhi-g80/chipku/server"
)

var port string

func init() {
	rootCmd.AddCommand(serve)
    serve.Flags().StringVarP(&port, "port", "p", "8080", "port to serve on")
}


var serve = &cobra.Command{
	Use:   "serve",
	Short: "Start server",
	Long:  `Start the pastebin server`,
	Run: func(cmd *cobra.Command, args []string) {
		// fmt.Printf("Chipku - a no frill pastebin v%s\n", Version)
        server.Serve(port)
	},
}

// serve.Flags().IntVarP(&echoTimes, "times", "t", 1, "times to echo the input")
