package cli

import (
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/fornellas/tasmota_exporter/cli/server"
	"github.com/fornellas/tasmota_exporter/cli/version"
	"github.com/fornellas/tasmota_exporter/log"
)

var ExitFunc func(int) = func(code int) { os.Exit(code) }

var logLevelStr string
var defaultLogLevelStr = "info"
var forceColor bool
var defaultForceColor = false

var Cmd = &cobra.Command{
	Use:   "tasmota_exporter",
	Short: "Prometheus exporter for Tasmota sensors",
	Args:  cobra.NoArgs,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if forceColor {
			color.NoColor = false
		}
		cmd.SetContext(log.SetLoggerValue(
			cmd.Context(), cmd.OutOrStderr(), logLevelStr, ExitFunc,
		))
	},
	Run: func(cmd *cobra.Command, args []string) {
		logger := log.GetLogger(cmd.Context())
		if err := cmd.Help(); err != nil {
			logger.Fatal(err)
		}
	},
}

func init() {
	Cmd.PersistentFlags().StringVarP(
		&logLevelStr, "log-level", "l", defaultLogLevelStr,
		"Logging level",
	)
	Cmd.PersistentFlags().BoolVarP(
		&forceColor, "force-color", "", defaultForceColor,
		"Force colored output",
	)

	Cmd.AddCommand(server.Cmd)
	Cmd.AddCommand(version.Cmd)
}
