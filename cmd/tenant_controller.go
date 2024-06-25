/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var tenantControllerCmd = &cobra.Command{
	Use:   "tenant-controller",
	Short: "Starts the zeus tenant controller addon server",
	Long:  `The tenant-controller addon is a service that exposes tenant information from LGTM backends for use in the zeus stack. See INSERT LINK for more information`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("tenantController called")
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
