package commands

import (
	"fmt"
	"github.com/bhatti/PlexAuthZ/internal/controller"
	cfg "github.com/bhatti/PlexAuthZ/internal/domain"
	"github.com/bhatti/PlexAuthZ/internal/factory"
	"github.com/bhatti/PlexAuthZ/internal/metrics"
	"github.com/bhatti/PlexAuthZ/internal/server"
	"github.com/bhatti/PlexAuthZ/internal/version"
	"github.com/bhatti/PlexAuthZ/internal/web"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
)

var cfgFile string

var config *cfg.Config
var logFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "platform-poc",
	Short: "Start an instance of Auth Server",
	Long:  "Start an instance of Auth Server",
	Run: func(cmd *cobra.Command, args []string) {
		_ = rootRun(cmd, args)
	},
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(
		&cfgFile,
		"config",
		"",
		"config file")

}

func rootRun(_ *cobra.Command, _ []string) error {
	// Start web server for health status
	webServer := web.NewDefaultWebServer(config)

	authService, _, err := factory.CreateAuthAdminService(config, metrics.New(), cfg.RootClientType, "")
	if err != nil {
		return err
	}
	err = controller.StartControllers(config, authService, webServer)
	if err != nil {
		return err
	}

	grpcServer, err := server.StartServers(config, authService)
	if err != nil {
		log.WithField("Error", err).Fatalf("could start for gRPC server")
		return err
	}

	go func() {
		err = grpcServer.Serve()
		if err != nil {
			log.WithField("Error", err).
				Fatal("failed to start GRPC Server")
		}
	}()

	return webServer.Start(config.HttpListenPort)
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(version, commit, date string) error {
	Version = version
	Commit = commit
	Date = date
	log.Infof("Starting root command...")
	if err := rootCmd.Execute(); err != nil {
		log.WithField("Error", err).
			Fatalf("Could not start server")

		return err
	}

	return nil
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	var err error
	if config, err = cfg.NewConfig(cfgFile); err != nil {
		panic(err)
	}
	config.Version = version.New(Version, Commit, Date)

	log.SetFormatter(&log.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})
	if logFile != "" {
		fmt.Printf("routing logs to %s\n", logFile)
		if f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666); err == nil {
			log.SetOutput(f)
			log.SetFormatter(&log.JSONFormatter{})
		}
	}
	log.SetLevel(log.InfoLevel)
}
