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

	"github.com/gin-gonic/gin"
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
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
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
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.zeus.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func gracefulShutdown(ctx context.Context, server *http.Server, api *gin.Engine) {
	// wait for signal
	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	slog.Info("gracefully shutting down the server")
	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error("unable to shutdown server", "error", err)
	}
	slog.Info("server shutdown complete, see you next time")
}
