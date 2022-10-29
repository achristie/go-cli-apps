package cmd

import (
	"fmt"
	"io"
	"os"

	"achristie.net/cobra/scan"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:          "add",
	Short:        "short",
	Long:         `longish`,
	Aliases:      []string{"a"},
	Args:         cobra.MinimumNArgs(1),
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		hostsFile, err := cmd.Flags().GetString("hosts-file")
		if err != nil {
			return err
		}

		return addAction(os.Stdout, hostsFile, args)
	},
}

func init() {
	hostsCmd.AddCommand(addCmd)
}

func addAction(out io.Writer, hostsFile string, args []string) error {
	hl := &scan.HostsList{}

	if err := hl.Load(hostsFile); err != nil {
		return err
	}

	for _, s := range args {
		if err := hl.Add(s); err != nil {
			return err
		}

		fmt.Fprintln(out, "Added host:", s)
	}

	return hl.Save(hostsFile)
}
