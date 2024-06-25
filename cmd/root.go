/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var (
	// See https://www.digitalocean.com/community/tutorials/using-ldflags-to-set-version-information-for-go-applications
	// for setting ldflags for dynamic version
	serviceVersion = "development"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "zeus",
	Short: "Addons for Grafana LGTM",
	Long:  `Zeus is a collection of add-ons for Grafana LGTM stack that help facilitate operating an OSS observability stack`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().String("objstore.config", "", "The configuration for the object store used by zeus")
	rootCmd.PersistentFlags().String("objstore.config-file", "", "The configuration file for the object store used by zeus")

	// Local flags
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func gracefulShutdown(ctx context.Context, server *http.Server) {
	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	slog.Info("gracefully shutting down the service")
	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error("unable to shutdown http server", "error", err)
	}
}
