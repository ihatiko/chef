package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "test",
	Short: "test",
	Long:  "test",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("test called %v\n", os.Args)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		_, err := fmt.Fprintf(os.Stderr, "Oops. An error while executing command '%s'\n", err)
		if err != nil {
			return
		}
		os.Exit(1)
	}
}

func main() {
	Execute()
}
