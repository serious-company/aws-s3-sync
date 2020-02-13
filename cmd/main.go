package main

import (
	"fmt"
	"os"

	"github.com/directxman12/zapr"
	"github.com/fsnotify/fsnotify"
	"github.com/go-logr/logr"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/dafiti-group/aws-s3-sync-api/pkg/server"
	"github.com/dafiti-group/aws-s3-sync-api/pkg/sync"
)

var rootCmd = &cobra.Command{
	Use:  "aws-s3-sync",
	Long: `Aws S3 Sync`,
}

// Config
type Config struct {
	Log       logr.Logger
	Viper     *viper.Viper
	SyncPath  string `mapstructure:"sync_path"`
	AwsBucket string `mapstructure:"aws_bucket"`
	AwsRegion string `mapstructure:"aws_region"`
}

var config Config

func main() {
	// Setup Logger
	var log logr.Logger

	zapLog, err := zap.NewProduction()
	if err != nil {
		fmt.Printf("who watches the watchmen (%v)?", err)
		os.Exit(-1)
	}

	log = zapr.NewLogger(zapLog)
	config.Log = log

	// Add Commands
	addCommands()

	// Start Viper
	config.Viper = viper.New()

	config.Viper.SetConfigName("config")
	config.Viper.SetConfigType("yaml")
	config.Viper.AddConfigPath("/etc/")
	config.Viper.AddConfigPath("$HOME/")
	config.Viper.AddConfigPath(".")
	config.Viper.AddConfigPath("./hack")

	if err := config.Viper.ReadInConfig(); err != nil {
		log.Error(err, "Failed while reading config file")
		os.Exit(-1)
	}

	if err := config.Viper.Unmarshal(&config); err != nil {
		log.Error(err, "Failed while eunmarshaling configs")
		os.Exit(-1)
	}

	config.Viper.WatchConfig()

	if err := rootCmd.Execute(); err != nil {
		log.Error(err, "Failed while executing cmd")
		os.Exit(-1)
	}
}

//AddCommands adds child commands to the root command rootCmd.
func addCommands() {
	rootCmd.AddCommand(cmdServer)
	rootCmd.AddCommand(cmdDummy)
}

var cmdServer = &cobra.Command{
	Use:   "server",
	Short: "Web Server",
	Run: func(cmd *cobra.Command, args []string) {
		log := config.Log.WithName("server")

		// Try to sync
		if err := config.sync(); err != nil {
			log.Info("Sync Failed, Will wait for file to change")
		}

		// @TODO: Validation
		config.Viper.OnConfigChange(func(e fsnotify.Event) {
			if err := config.Viper.Unmarshal(&config); err != nil {
				log.Error(err, "Fatal error Unmarshaling")
				panic(fmt.Errorf("Fatal error Unmarshaling file: %s \n", err))
			}

			if err := config.sync(); err != nil {
				log.Error(err, "Fatal error syncing")
			}
		})

		s := &server.Server{
			Log: config.Log,
		}
		s.Initialize()
		s.Run(":3000")
	},
}

var cmdDummy = &cobra.Command{
	Use:   "dummy",
	Short: "Just exits",
	Run: func(cmd *cobra.Command, args []string) {
		log := config.Log.WithName("dummy")
		log.Info("Fine, will exit")
		os.Exit(0)
	},
}

func (c *Config) sync() error {
	log := c.Log.WithName("sync")
	s := sync.Sync{
		Bucket: c.AwsBucket,
		Region: c.AwsRegion,
		Path:   c.SyncPath,
		Log:    log,
	}
	log.Info("Start Sync")
	if err := s.AwsSync(); err != nil {
		log.Error(err, "Error while sync")
		return err
	}

	log.Info("Done")
	return nil
}
