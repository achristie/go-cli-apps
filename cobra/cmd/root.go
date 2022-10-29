package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:     "pScan",
	Short:   "short",
	Long:    "longer",
	Version: "0.1",
}

var cfgFile string

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file")
	rootCmd.PersistentFlags().StringP("hosts-file", "f", "pScan.hosts", "pScan hosts file")

	versionTemplate := `{{printf "%s: %s - version %s\n" .Name .Short .Version}}`

	rootCmd.SetVersionTemplate(versionTemplate)
}

func initConfig() {
	return
}
