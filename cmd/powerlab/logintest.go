package main

import (
	"fmt"

	"github.com/andrewjjenkins/powerlab/pkg/config"
	"github.com/andrewjjenkins/powerlab/pkg/server/megarac"
	"github.com/spf13/cobra"
)

var loginTestCmd = &cobra.Command{
	Use:  "logintest <hostname>",
	Args: cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		hostname := args[0]

		creds, err := config.LoadCredentials("")
		if err != nil {
			panic(err)
		}
		server, ok := creds.Servers[hostname]
		if !ok {
			panic(fmt.Errorf("cannot find %s in credentials", hostname))
		}

		a, err := megarac.NewApi(hostname, true)
		if err != nil {
			panic(err)
		}
		err = a.Login(server.Username, server.Password)
		if err != nil {
			panic(err)
		}
		fmt.Println("Logged in")
	},
}

func init() {
	rootCmd.AddCommand(loginTestCmd)
}
