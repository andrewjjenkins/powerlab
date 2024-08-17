package main

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "powerlab",
	Short: "Control server power state using IPMI",
}

func main() {
	rootCmd.Execute()
}
