package main

import (
	"net/http"
	"time"

	"github.com/andrewjjenkins/powerlab/pkg/serve"

	"github.com/golang/glog"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use: "serve",
	Run: func(cmd *cobra.Command, args []string) {
		s := &http.Server{
			Addr:         "0.0.0.0:8080",
			WriteTimeout: 15 * time.Second,
			ReadTimeout:  15 * time.Second,
		}
		serve.Serve(s)

		glog.Fatal(s.ListenAndServe())
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
