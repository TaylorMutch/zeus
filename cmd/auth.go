/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// authCmd represents the auth command
var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Starts the zeus auth addon server",
	Long:  `The auth addon is a service that provides authentication for the zeus stack. See INSERT LINK for more information`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("auth called")

		// TODO setup a webserver to serve auth requests

		// TODO setup a storage provider to sync requests from

		// TODO setup telemetry to report metrics, traces
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
