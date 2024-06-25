/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/TaylorMutch/zeus/pkg/auth"
	"github.com/TaylorMutch/zeus/pkg/storage"
	"github.com/spf13/cobra"
)

// generatePasswordCmd represents the generatePassword command
var generatePasswordCmd = &cobra.Command{
	Use:   "generate-password",
	Short: "Generate a zeus credentials",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("generate-password called")
		if err := runGeneratePassword(); err != nil {
			fmt.Println("error generating password: ", err)
		}
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

func runGeneratePassword() error {
	// TODO - parse CLI flag for tenant, object store
	store, err := storage.NewObjectStore("zeus-auth", []byte(`type: FILESYSTEM
config:
  directory: "/tmp/zeus-auth"
prefix: ""`))
	if err != nil {
		return fmt.Errorf("failed to setup storage provider: %w", err)
	}

	// Generate a password
	tenant := "taylor-123" // TODO - parse CLI flag
	password, err := auth.GenerateRandomString(24)
	if err != nil {
		return fmt.Errorf("failed to generate password: %w", err)
	}

	cred := auth.NewCredential(tenant, password)
	username := cred.ID

	// Store the credential
	byt, err := json.Marshal(cred)
	if err != nil {
		return fmt.Errorf("failed to marshal credential: %w", err)
	}

	fmt.Printf("Generated credential for tenant %s:\n%s\n%s\n", tenant, username, password)
	err = store.PutObject(context.Background(), fmt.Sprintf("%s/%s", auth.CredentialStoreStoragePrefix, username), byt)
	if err != nil {
		return fmt.Errorf("failed to store credential: %w", err)
	}
	fmt.Println("Credential stored successfully")
	return nil
}
