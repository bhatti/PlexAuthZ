package commands

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var dumpConfigCmd = &cobra.Command{
	Use:   "dump-config",
	Short: "Dumps the Config to Screen",
	Long:  `Used to debug the config`,
	Run: func(cmd *cobra.Command, args []string) {
		d, err := yaml.Marshal(&config)
		if err != nil {
			log.Fatalf("error: %v", err)
		}
		fmt.Printf("%s", string(d))
	},
}

func init() {
	rootCmd.AddCommand(dumpConfigCmd)
}
