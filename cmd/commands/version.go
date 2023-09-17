package commands

import (
	"fmt"
	"github.com/bhatti/PlexAuthZ/internal/version"

	"github.com/spf13/cobra"
)

var (
	shortened = false
	// Version of latest change.
	Version = "dev"
	// Commit of latest change.
	Commit = "dirty"
	// Date of latest change.
	Date = ""
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "version",
	Long:  `Version will output the current build information`,
	Run: func(cmd *cobra.Command, _ []string) {
		v := version.New(Version, Commit, Date)
		var response string

		if shortened {
			response = v.ToYAML()
		} else {
			response = v.ToJSON()
		}

		fmt.Printf("%+v", response)
		return
	},
}

func init() {
	versionCmd.Flags().BoolVarP(
		&shortened,
		"short",
		"s",
		true,
		"Use shortened output for version information.")
	rootCmd.AddCommand(versionCmd)
}
