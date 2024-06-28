/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/TaylorMutch/zeus/pkg/api"
	"github.com/TaylorMutch/zeus/pkg/enums"
	"github.com/TaylorMutch/zeus/pkg/telemetry"
	"github.com/spf13/cobra"
)

// alertmanagerLabelProxyCmd represents the alertmanagerLabelProxy command
var alertmanagerLabelProxyCmd = &cobra.Command{
	Use:   "alertmanager-label-proxy",
	Short: "Starts the zeus alertmanager label proxy addon server",
	Long: `The alertmanager label proxy addon is a service that attaches a tenant label to incoming alerts and passes them onto a downstream alertmanager instance.` +
		` See INSERT LINK for more information`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := runAlertmanagerLabelProxy(
			cmd.Flag("server.addr").Value.String(),
			cmd.Flag("alertmanager.addr").Value.String(),
			cmd.Flag("alerts.proxy-label").Value.String(),
			cmd.Flag("alerts.tenant-id-header").Value.String(),
		); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				slog.Error("zeus alertmanager-label-proxy failed to run service correctly", "error", err)
			}
		} else {
			slog.Info("zeus alertmanager label proxy shutdown complete, see you next time")
		}
	},
}

func init() {
	rootCmd.AddCommand(alertmanagerLabelProxyCmd)

	// Local flags
	alertmanagerLabelProxyCmd.Flags().String("server.addr", ":8080", "The address to run the alertmanager-label-proxy server on")
	alertmanagerLabelProxyCmd.Flags().String("alertmanager.addr", "http://alertmanager:9093", "The address of the alertmanager server")
	alertmanagerLabelProxyCmd.Flags().String("alerts.proxy-label", "zeus_tenant", "The label to add to alerts which indicates the tenant the alert originated from")
	alertmanagerLabelProxyCmd.Flags().String("alerts.tenant-id-header", "X-Scope-Orgid", "The header which contains the tenant ID")
	//alertmanagerLabelProxyCmd.Flags().Bool("alerts.tenant-id-header-required", false, "If the org-id header is required to be present in the request")
}

func runAlertmanagerLabelProxy(
	serverAddr,
	alertmanagerAddr,
	proxyLabel,
	tenantIDHeader string,
) error {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt)
	defer cancel()
	telemetryShutdown, err := telemetry.Init(ctx, enums.AlertmanagerLabelProxyServiceName, serviceVersion)
	if err != nil {
		slog.Error("failed to setup telemetry", "error", err)
		os.Exit(1)
	}
	defer telemetryShutdown(ctx)
	api, err := api.NewAlertmanagerLabelProxy(enums.AlertmanagerLabelProxyServiceName, alertmanagerAddr, proxyLabel, tenantIDHeader)
	if err != nil {
		return err
	}
	server := &http.Server{
		Addr:    serverAddr,
		Handler: api,
	}
	go gracefulShutdown(ctx, server)
	return server.ListenAndServe()
}
