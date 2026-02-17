package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var (
	// Global flags
	apiKey       string
	outputFormat string
	outputFile   string
	verbose      bool
	quiet        bool
	dryRun       bool
	listCommands bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ahrefs",
	Short: "Ahrefs API CLI - AI agent-friendly interface to Ahrefs API v3",
	Long: `Ahrefs CLI is a command-line interface for the Ahrefs API v3.

Designed for maximum discoverability and ease of use by AI coding agents.
Every command supports --help, --describe, and --list-fields for introspection.

Authentication:
  Set API key via --api-key flag or AHREFS_API_KEY environment variable.
  Or use 'ahrefs config set-key <key>' to persist in config file.

Output Formats:
  json (default), yaml, csv, table

Examples:
  # Get domain rating
  ahrefs site-explorer domain-rating --target example.com

  # List all available commands
  ahrefs --list-commands

  # Get structured command metadata
  ahrefs site-explorer backlinks --describe`,
	Version: "0.1.0",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Handle --list-commands at root level
		if listCommands {
			return printCommandList(cmd.Root())
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// If --list-commands was specified, it was already handled in PersistentPreRunE
		if listCommands {
			return nil
		}
		// Otherwise show help
		return cmd.Help()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Global flags available to all commands
	rootCmd.PersistentFlags().StringVar(&apiKey, "api-key", os.Getenv("AHREFS_API_KEY"), "Ahrefs API key (or set AHREFS_API_KEY env var)")
	rootCmd.PersistentFlags().StringVar(&outputFormat, "format", "json", "Output format: json, yaml, csv, table")
	rootCmd.PersistentFlags().StringVarP(&outputFile, "output", "o", "", "Output file (default: stdout)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output (show request/response details)")
	rootCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "Quiet mode (errors only)")
	rootCmd.PersistentFlags().BoolVar(&dryRun, "dry-run", false, "Validate request without executing")

	// Root-level flags
	rootCmd.Flags().BoolVar(&listCommands, "list-commands", false, "List all available commands as JSON")
}

// AddCommands adds all subcommands to root
func AddCommands(commands ...*cobra.Command) {
	rootCmd.AddCommand(commands...)
}

// CommandInfo represents metadata about a command for introspection
type CommandInfo struct {
	Name        string        `json:"name"`
	Use         string        `json:"use"`
	Short       string        `json:"short"`
	Long        string        `json:"long"`
	Subcommands []CommandInfo `json:"subcommands,omitempty"`
	Flags       []FlagInfo    `json:"flags,omitempty"`
	Examples    string        `json:"examples,omitempty"`
}

type FlagInfo struct {
	Name      string `json:"name"`
	Shorthand string `json:"shorthand,omitempty"`
	Usage     string `json:"usage"`
	DefValue  string `json:"default,omitempty"`
	Required  bool   `json:"required"`
}

// printCommandList outputs all available commands as JSON
func printCommandList(cmd *cobra.Command) error {
	info := buildCommandInfo(cmd)
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(info); err != nil {
		return fmt.Errorf("failed to encode command list: %w", err)
	}
	return nil
}

// buildCommandInfo recursively builds command metadata
func buildCommandInfo(cmd *cobra.Command) CommandInfo {
	info := CommandInfo{
		Name:     cmd.Name(),
		Use:      cmd.Use,
		Short:    cmd.Short,
		Long:     cmd.Long,
		Examples: cmd.Example,
	}

	// Add flags
	cmd.Flags().VisitAll(func(flag *pflag.Flag) {
		flagInfo := FlagInfo{
			Name:      flag.Name,
			Shorthand: flag.Shorthand,
			Usage:     flag.Usage,
			DefValue:  flag.DefValue,
		}
		// Check if required
		if requiredAnnotation, ok := flag.Annotations["required"]; ok && len(requiredAnnotation) > 0 {
			flagInfo.Required = true
		}
		info.Flags = append(info.Flags, flagInfo)
	})

	// Add subcommands recursively
	for _, subcmd := range cmd.Commands() {
		if !subcmd.Hidden {
			info.Subcommands = append(info.Subcommands, buildCommandInfo(subcmd))
		}
	}

	return info
}

// GetGlobalFlags returns the current global flag values
func GetGlobalFlags() GlobalFlags {
	return GlobalFlags{
		APIKey:       apiKey,
		OutputFormat: outputFormat,
		OutputFile:   outputFile,
		Verbose:      verbose,
		Quiet:        quiet,
		DryRun:       dryRun,
	}
}

// GlobalFlags holds all global flag values
type GlobalFlags struct {
	APIKey       string
	OutputFormat string
	OutputFile   string
	Verbose      bool
	Quiet        bool
	DryRun       bool
}
