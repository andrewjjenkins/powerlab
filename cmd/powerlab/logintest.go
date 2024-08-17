package main

import (
	"fmt"

	"github.com/andrewjjenkins/powerlab/pkg/server/megarac"
	"github.com/spf13/cobra"
)

var loginTestCmd = &cobra.Command{
	Use:  "logintest <hostname> <username> <password>",
	Args: cobra.MatchAll(cobra.ExactArgs(3), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		hostname, username, password := args[0], args[1], args[2]
		a, err := megarac.NewApi(hostname, true)
		if err != nil {
			panic(err)
		}
		err = a.Login(username, password)
		if err != nil {
			panic(err)
		}
		fmt.Println("Logged in")
	},
}

func init() {
	rootCmd.AddCommand(loginTestCmd)
}
