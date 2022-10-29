package cmd

import (
	"github.com/spf13/cobra"
)

var hostsCmd = &cobra.Command{
	Use:   "hosts",
	Short: "short",
	Long:  `longish`,
}

func init() {
	rootCmd.AddCommand(hostsCmd)
}
