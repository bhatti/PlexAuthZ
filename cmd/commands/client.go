package commands

import (
	"github.com/spf13/cobra"
)

var clientCmd = &cobra.Command{
	Use:   "client for authz server",
	Short: "client for authz server",
	Long:  "client for authz server",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO
	},
}

func init() {
	cobra.OnInitialize(initConfig)

	clientCmd.PersistentFlags().StringVar(
		&cfgFile,
		"config",
		"",
		"config file")
	rootCmd.PersistentFlags().StringVar(
		&logFile,
		"logFile",
		"",
		"log file")
	rootCmd.AddCommand(clientCmd)
}
