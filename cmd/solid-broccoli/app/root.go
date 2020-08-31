package app

import (
	"fmt"
	"os"
	"runtime"

	sb "github.com/dstdfx/solid-broccoli/internal/app/solidbroccoli"
	"github.com/dstdfx/solid-broccoli/internal/pkg/config"
	"github.com/dstdfx/solid-broccoli/internal/pkg/log"
	"github.com/spf13/cobra"
)

const defaultCfgFile = "/etc/solid-broccoli/solid-broccoli.yaml"

var cfgFile string

// nolint
// Variables that are injected in build time.
var (
	buildGitCommit string
	buildGitTag    string
	buildDate      string
	buildCompiler  = runtime.Version()
)

// RootCmd represents the base command when called without any subcommands.
var RootCmd = &cobra.Command{
	Use:   "solid-broccoli",
	Short: "solid-broccoli represents a simple HTTP API service",
	Run: func(_ *cobra.Command, _ []string) {
		// Initialize application config and log.
		if _, err := os.Stat(cfgFile); err != nil {
			exitWithErr(fmt.Errorf("config file %s can't be read: %s", cfgFile, err))
		}
		if err := config.InitFromFile(cfgFile); err != nil {
			exitWithErr(err)
		}

		// Init logger
		logger, err := log.InitLogger(log.InitLoggerOpts{
			File:      config.Config.Log.File,
			UseStdout: config.Config.Log.UseStdout,
			Debug:     config.Config.Log.Debug,
		})
		if err != nil {
			exitWithErr(err)
		}

		// Start main routine
		if err := sb.StartService(logger, sb.StartOpts{Interrupt: make(chan os.Signal, 1)}); err != nil {
			exitWithErr(fmt.Errorf("error starting solidbroccoli app: %w", err))
		}
	},
}

func init() {
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config",
		defaultCfgFile, "path to application config")
}

// exitWithErr is a helper method to print errors in case of empty logger.
func exitWithErr(err error) {
	_, _ = fmt.Fprintf(os.Stderr, "application is exiting after error: %s\n", err)
	os.Exit(1)
}
