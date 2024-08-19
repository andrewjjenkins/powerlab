package main

import (
	"net/http"
	"time"

	"github.com/andrewjjenkins/powerlab/pkg/config"
	"github.com/andrewjjenkins/powerlab/pkg/serve"

	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use: "serve",
	Run: func(cmd *cobra.Command, args []string) {
		serverManager, err := config.LoadServers("")
		if err != nil {
			panic(err)
		}

		s := &http.Server{
			Addr:         "0.0.0.0:8080",
			WriteTimeout: 15 * time.Second,
			ReadTimeout:  15 * time.Second,
		}
		serve.Serve(s, serverManager, serveDevMode)

		s.ListenAndServe()
	},
}

var serveDevMode bool

func init() {
	serveCmd.PersistentFlags().BoolVar(&serveDevMode, "dev", false, "Proxy /ui/ to localhost:3000 for dev mode")
	rootCmd.AddCommand(serveCmd)
}
