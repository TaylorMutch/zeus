/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// generatePasswordCmd represents the generatePassword command
var generatePasswordCmd = &cobra.Command{
	Use:   "generate-password",
	Short: "Generate a zeus credentials",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("generate-password called")
	},
}

func init() {
	rootCmd.AddCommand(generatePasswordCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// generatePasswordCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// generatePasswordCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
