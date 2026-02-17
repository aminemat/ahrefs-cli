package config

import (
	"fmt"

	"github.com/amine/ahrefs-cli/internal/config"
	"github.com/spf13/cobra"
)

// NewConfigCmd creates the config command
func NewConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage CLI configuration",
		Long:  "Manage configuration settings for the Ahrefs CLI, including API key storage.",
	}

	cmd.AddCommand(newSetKeyCmd())
	cmd.AddCommand(newShowCmd())
	cmd.AddCommand(newValidateCmd())

	return cmd
}

func newSetKeyCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "set-key <api-key>",
		Short: "Set the Ahrefs API key",
		Long:  "Save the Ahrefs API key to the configuration file (~/.ahrefsrc).",
		Args:  cobra.ExactArgs(1),
		Example: `  # Set API key
  ahrefs config set-key sk_your_api_key_here`,
		RunE: func(cmd *cobra.Command, args []string) error {
			apiKey := args[0]

			cfg := &config.Config{
				APIKey: apiKey,
			}

			if err := config.Save(cfg); err != nil {
				return fmt.Errorf("failed to save config: %w", err)
			}

			fmt.Println("API key saved successfully")
			return nil
		},
	}
}

func newShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "show",
		Short: "Show current configuration",
		Long:  "Display the current configuration settings.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			if cfg.APIKey == "" {
				fmt.Println("No API key configured")
				fmt.Println("Set one with: ahrefs config set-key <your-key>")
			} else {
				// Mask the API key
				masked := maskAPIKey(cfg.APIKey)
				fmt.Printf("API Key: %s\n", masked)
			}

			return nil
		},
	}
}

func newValidateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "validate",
		Short: "Validate API key",
		Long:  "Test if the configured API key is valid by making a test API request.",
		RunE: func(cmd *cobra.Command, args []string) error {
			apiKey := config.GetAPIKey()
			if apiKey == "" {
				return fmt.Errorf("no API key configured. Use 'ahrefs config set-key <key>'")
			}

			fmt.Println("API key validation not yet implemented")
			fmt.Println("Will test with a lightweight API request in the future")
			return nil
		},
	}
}

func maskAPIKey(key string) string {
	if len(key) <= 8 {
		return "****"
	}
	return key[:4] + "****" + key[len(key)-4:]
}
