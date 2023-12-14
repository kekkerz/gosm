package cmd

import (
	"github.com/spf13/cobra"
	"os"
)

var Profile string

var rootCmd = &cobra.Command{
	Use: "gosm",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&Profile, "profile", "p", "", "AWS profile")
}
