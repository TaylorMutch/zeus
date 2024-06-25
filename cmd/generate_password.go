/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/TaylorMutch/zeus/pkg/auth"
	"github.com/TaylorMutch/zeus/pkg/storage"
	"github.com/TaylorMutch/zeus/pkg/telemetry"
	"github.com/spf13/cobra"
)

// generatePasswordCmd represents the generatePassword command
var generatePasswordCmd = &cobra.Command{
	Use:   "generate-password",
	Short: "Generate a zeus credentials",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if err := runGeneratePassword(
			cmd.Flag("tenant-id").Value.String(),
			cmd.Flag("objstore.config").Value.String(),
			cmd.Flag("objstore.config-file").Value.String(),
		); err != nil {
			slog.Error("error generating password: ", "error", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(generatePasswordCmd)
	generatePasswordCmd.Flags().String("tenant-id", "", "The tenant to generate a password for")
	generatePasswordCmd.MarkFlagRequired("tenant-id")
}

func runGeneratePassword(tenantID, objstoreConfigStr, objstoreConfigFile string) error {
	telemetry.InitLogging()

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

	password, err := auth.GenerateRandomString(24)
	if err != nil {
		return fmt.Errorf("failed to generate password: %w", err)
	}

	cred := auth.NewCredential(tenantID, password)
	byt, err := json.Marshal(cred)
	if err != nil {
		return fmt.Errorf("failed to marshal credential: %w", err)
	}

	slog.Info("generated credential", "tenant_id", tenantID, "username", cred.ID, "password", password)
	err = store.PutObject(context.Background(), fmt.Sprintf("%s/%s", auth.CredentialStoreStoragePrefix, cred.ID), byt)
	if err != nil {
		return fmt.Errorf("failed to store credential: %w", err)
	}
	slog.Info("credential stored successfully")
	return nil
}
