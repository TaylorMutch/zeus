/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"log/slog"

	"github.com/TaylorMutch/zeus/pkg/telemetry"
	"github.com/spf13/cobra"
)

var tenantControllerCmd = &cobra.Command{
	Use:   "tenant-controller",
	Short: "Starts the zeus tenant controller addon server",
	Long:  `The tenant-controller addon is a service that exposes tenant information from LGTM backends for use in the zeus stack. See INSERT LINK for more information`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := runTenantController(
			cmd.Flag("server.addr").Value.String(),
			cmd.Flag("mimir.addr").Value.String(),
			cmd.Flag("mimir.metric").Value.String(),
			cmd.Flag("loki.addr").Value.String(),
			cmd.Flag("loki.metric").Value.String(),
			cmd.Flag("tempo.addr").Value.String(),
			cmd.Flag("tempo.metric").Value.String(),
		); err != nil {
			slog.Error("zeus tenant controller failed to run service correctly", "error", err)
		} else {
			slog.Info("zeus tenant controller shutdown complete, see you next time")
		}
	},
}

func init() {
	rootCmd.AddCommand(tenantControllerCmd)

	// Local flags
	tenantControllerCmd.Flags().String("server.addr", ":8080", "The address to run the tenant controller server on")
	tenantControllerCmd.Flags().String("mimir.addr", "http://mimir-distributor", "The address of the Mimir server which includes tenant information")
	tenantControllerCmd.Flags().String("mimir.metric", "TODO_METRIC", "The metric used to determine tenants from Mimir")
	tenantControllerCmd.Flags().String("loki.addr", "http://loki-distributor", "The address of the Loki server which includes tenant information")
	tenantControllerCmd.Flags().String("loki.metric", "TODO_METRIC", "The metric used to determine tenants from Loki")
	tenantControllerCmd.Flags().String("tempo.addr", "http://tempo-distributor", "The address of the Tempo server which includes tenant information")
	tenantControllerCmd.Flags().String("tempo.metric", "TODO_METRIC", "The metric used to determine tenants from Tempo")
}

func runTenantController(
	serverAddr,
	mimirAddr,
	mimirMetric,
	lokiAddr,
	lokiMetric,
	tempoAddr,
	tempoMetric string,
) error {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt)
	defer cancel()
	telemetryShutdown, err := telemetry.Init(ctx, "zeus-auth", serviceVersion)
	if err != nil {
		slog.Error("failed to setup telemetry", "error", err)
		os.Exit(1)
	}
	defer telemetryShutdown(ctx)

	// TODO - setup DNS discovery for mimir, loki, tempo

	return nil
}
