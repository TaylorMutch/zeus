/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/TaylorMutch/zeus/pkg/api"
	"github.com/TaylorMutch/zeus/pkg/auth"
	"github.com/TaylorMutch/zeus/pkg/storage"
	"github.com/TaylorMutch/zeus/pkg/telemetry"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

// authCmd represents the auth command
var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Starts the zeus auth addon server",
	Long:  `The auth addon is a service that provides authentication for the zeus stack. See INSERT LINK for more information`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := runAuth(); err != nil {
			slog.Error("zeus auth failed to run service")
		}
	},
}

func init() {
	rootCmd.AddCommand(authCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// authCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// authCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// runAuth runs the auth service
func runAuth() error {

	// Setup shutdown signals
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt)
	defer cancel()

	// Initialize telemetry
	telemetryShutdown, err := telemetry.Init(ctx, "zeus-auth", serviceVersion)
	if err != nil {
		slog.Error("failed to setup telemetry", "error", err)
		os.Exit(1)
	}
	defer telemetryShutdown(ctx)

	// Setup storage provider to sync authentication from
	// TODO - configuration for blob storage provider
	store, err := storage.NewObjectStore("zeus-auth", []byte(`type: FILESYSTEM
config:
  directory: ""
prefix: ""`))
	if err != nil {
		return fmt.Errorf("failed to setup storage provider: %w", err)
	}

	authstore, err := auth.NewObjectCredentialStore(store)
	if err != nil {
		return fmt.Errorf("failed to setup auth store: %w", err)
	}

	api := api.New()

	api.GET("/auth", func(c *gin.Context) {

		authstore.GetCredential("test")

	})

	// Setup a webserver to serve auth requests
	server := &http.Server{
		Addr:    ":8080", // TODO - configuration for port
		Handler: api,
	}

	// gracefulShutdown
	go gracefulShutdown(ctx, server, api)
	return server.ListenAndServe()
}
