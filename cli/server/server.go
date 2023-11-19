package server

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/fornellas/tasmota_exporter/log"
	"github.com/fornellas/tasmota_exporter/server"
)

var addr string

var Cmd = &cobra.Command{
	Use:   "server",
	Short: "Start HTTP server to handle probe requests.",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()

		logger := log.GetLogger(ctx)

		srv := server.NewServer(addr)

		go func() {
			sig := make(chan os.Signal, 1)
			signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
			<-sig

			logger.Info("Shutting down...")
			if err := srv.Shutdown(ctx); err != nil {
				logger.Errorf("Shutdown request failed: %v", err)
			}
		}()

		logger.Infof("Starting server on %s", addr)
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			logger.Fatalf("Server error: %v", err)
		}
		logger.Info("Exiting")
	},
}

func init() {
	Cmd.Flags().StringVarP(
		&addr, "address", "", ":8244",
		"TCP address for the server to listen on.",
	)
}
