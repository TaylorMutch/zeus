/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"errors"
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

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Starts the zeus auth addon server",
	Long:  `The auth addon is a service that provides authentication for the zeus stack. See INSERT LINK for more information`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := runAuth(
			cmd.Flag("server.addr").Value.String(),
			cmd.Flag("objstore.config").Value.String(),
			cmd.Flag("objstore.config-file").Value.String(),
		); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				slog.Error("zeus auth failed to run service correctly", "error", err)
			}
		} else {
			slog.Info("auth shutdown complete, see you next time")
		}
	},
}

func init() {
	rootCmd.AddCommand(authCmd)

	// Local flags
	authCmd.Flags().String("server.addr", ":8080", "The address to run the auth server on")
}

func runAuth(serverAddr, objstoreConfigStr, objstoreConfigFile string) error {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt)
	defer cancel()
	telemetryShutdown, err := telemetry.Init(ctx, "zeus-auth", serviceVersion)
	if err != nil {
		slog.Error("failed to setup telemetry", "error", err)
		os.Exit(1)
	}
	defer telemetryShutdown(ctx)

	if objstoreConfigStr == "" && objstoreConfigFile == "" {
		return fmt.Errorf("either --objstore.config or --objstore.config-file must be provided")
	}
	objstoreConfig, err := storage.ReadObjstoreConfig(objstoreConfigStr, objstoreConfigFile)
	if err != nil {
		return fmt.Errorf("failed to read objstore config: %w", err)
	}

	store, err := storage.NewObjectStore("zeus-auth", objstoreConfig)
	if err != nil {
		return fmt.Errorf("failed to setup storage provider: %w", err)
	}

	authstore, err := auth.NewObjectCredentialStore(store)
	if err != nil {
		return fmt.Errorf("failed to setup auth store: %w", err)
	}

	api := api.New()
	api.GET("/auth", func(c *gin.Context) {
		user, pass, ok := c.Request.BasicAuth()
		if !ok {
			c.String(http.StatusUnauthorized, "unauthorized")
			return
		}

		cred, err := authstore.GetCredential(c.Request.Context(), user)
		if err != nil {
			if errors.Is(err, auth.CredentialDoesNotExistError) {
				c.String(http.StatusNotFound, "credential does not exist")
				return
			}
			c.String(http.StatusInternalServerError, "server error")
		}

		authorized := auth.DefaultPasswordFactory.VerifyPassword(pass, cred.Password.CipherText, cred.Password.Salt)
		if !authorized {
			c.String(http.StatusUnauthorized, "unauthorized")
			return
		}

		authstore.CacheCredential(user, cred)
		c.String(http.StatusOK, "ok")
	})

	server := &http.Server{
		Addr:    serverAddr,
		Handler: api,
	}
	go gracefulShutdown(ctx, server)
	return server.ListenAndServe()
}
