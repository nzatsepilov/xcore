package cmd

import (
	"github.com/spf13/cobra"
)

var infoCmd = &cobra.Command{
	Use:       "info",
	ValidArgs: []string{"info"},
	Run: func(cmd *cobra.Command, args []string) {
		println("info command")
	},
}
