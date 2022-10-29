package cmd

import (
	"fmt"
	"io"
	"os"

	"achristie.net/cobra/scan"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:          "delet",
	Short:        "short",
	Long:         `longish`,
	Aliases:      []string{"d"},
	Args:         cobra.MinimumNArgs(1),
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		hostsFile, err := cmd.Flags().GetString("hosts-file")
		if err != nil {
			return err
		}

		return deleteAction(os.Stdout, hostsFile, args)
	},
}

func init() {
	hostsCmd.AddCommand(deleteCmd)
}

func deleteAction(out io.Writer, hostsFile string, args []string) error {
	hl := &scan.HostsList{}

	if err := hl.Load(hostsFile); err != nil {
		return err
	}

	for _, s := range args {
		if err := hl.Remove(s); err != nil {
			return err
		}

		fmt.Fprintln(out, "Deleted host:", s)
	}

	return hl.Save(hostsFile)
}
