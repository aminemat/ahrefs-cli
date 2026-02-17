package siteexplorer

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/amine/ahrefs-cli/cmd"
	"github.com/amine/ahrefs-cli/internal/config"
	"github.com/amine/ahrefs-cli/pkg/client"
	"github.com/amine/ahrefs-cli/pkg/models"
	"github.com/amine/ahrefs-cli/pkg/output"
	"github.com/spf13/cobra"
)

// NewSiteExplorerCmd creates the site-explorer command
func NewSiteExplorerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "site-explorer",
		Short: "Site Explorer API endpoints",
		Long: `Access Site Explorer data including domain rating, backlinks,
referring domains, anchors, organic keywords, and more.`,
		Aliases: []string{"se"},
	}

	cmd.AddCommand(newDomainRatingCmd())
	cmd.AddCommand(newBacklinksCmd())
	cmd.AddCommand(newBacklinksStatsCmd())

	return cmd
}

func newDomainRatingCmd() *cobra.Command {
	var (
		target string
		mode   string
		date   string
	)

	cmd := &cobra.Command{
		Use:   "domain-rating",
		Short: "Get domain rating for a target",
		Long: `Get the domain rating (DR) for a domain or URL.

Domain Rating is a metric that shows the strength of a website's backlink profile
on a logarithmic scale from 0 to 100, with the latter being the strongest.`,
		Example: `  # Get domain rating for a domain
  ahrefs site-explorer domain-rating --target example.com

  # Get domain rating for a specific URL
  ahrefs site-explorer domain-rating --target example.com/page --mode exact

  # Get historical domain rating
  ahrefs site-explorer domain-rating --target example.com --date 2024-01-01`,
		RunE: func(cobraCmd *cobra.Command, args []string) error {
			return runDomainRating(target, mode, date)
		},
	}

	cmd.Flags().StringVar(&target, "target", "", "Target domain or URL (required)")
	cmd.Flags().StringVar(&mode, "mode", "domain", "Mode: exact, domain, prefix, subdomains")
	cmd.Flags().StringVar(&date, "date", "", "Date for historical data (YYYY-MM-DD)")

	cmd.MarkFlagRequired("target")

	return cmd
}

func newBacklinksStatsCmd() *cobra.Command {
	var (
		target string
		mode   string
		date   string
	)

	cmd := &cobra.Command{
		Use:   "backlinks-stats",
		Short: "Get backlinks statistics",
		Long:  "Get aggregated statistics about backlinks for a target.",
		Example: `  # Get backlinks stats for a domain
  ahrefs site-explorer backlinks-stats --target example.com

  # Get stats for a specific URL
  ahrefs site-explorer backlinks-stats --target example.com/page --mode exact`,
		RunE: func(cobraCmd *cobra.Command, args []string) error {
			return runBacklinksStats(target, mode, date)
		},
	}

	cmd.Flags().StringVar(&target, "target", "", "Target domain or URL (required)")
	cmd.Flags().StringVar(&mode, "mode", "domain", "Mode: exact, domain, prefix, subdomains")
	cmd.Flags().StringVar(&date, "date", "", "Date for historical data (YYYY-MM-DD)")

	cmd.MarkFlagRequired("target")

	return cmd
}

func newBacklinksCmd() *cobra.Command {
	var (
		target string
		mode   string
		limit  int
		offset int
		sel    string
		where  string
	)

	cmd := &cobra.Command{
		Use:   "backlinks",
		Short: "Get backlinks for a target",
		Long:  "List backlinks pointing to a target domain or URL.",
		Example: `  # Get backlinks for a domain
  ahrefs site-explorer backlinks --target example.com --limit 100

  # Get specific fields
  ahrefs site-explorer backlinks --target example.com \
    --select url_from,domain_rating,anchor --limit 50

  # Filter backlinks
  ahrefs site-explorer backlinks --target example.com \
    --where 'domain_rating>50' --limit 100`,
		RunE: func(cobraCmd *cobra.Command, args []string) error {
			return runBacklinks(target, mode, limit, offset, sel, where)
		},
	}

	cmd.Flags().StringVar(&target, "target", "", "Target domain or URL (required)")
	cmd.Flags().StringVar(&mode, "mode", "domain", "Mode: exact, domain, prefix, subdomains")
	cmd.Flags().IntVar(&limit, "limit", 100, "Maximum number of results")
	cmd.Flags().IntVar(&offset, "offset", 0, "Offset for pagination")
	cmd.Flags().StringVar(&sel, "select", "", "Comma-separated list of fields to return")
	cmd.Flags().StringVar(&where, "where", "", "Filter expression (Ahrefs filter syntax)")

	cmd.MarkFlagRequired("target")

	return cmd
}

func runDomainRating(target, mode, date string) error {
	flags := cmd.GetGlobalFlags()

	// Get API key
	apiKey := flags.APIKey
	if apiKey == "" {
		apiKey = config.GetAPIKey()
	}
	if apiKey == "" {
		return fmt.Errorf("API key required. Set via --api-key flag, AHREFS_API_KEY env var, or 'ahrefs config set-key'")
	}

	// Create client
	c := client.NewClient(client.Config{
		APIKey: apiKey,
	})

	// Build request params
	params := url.Values{}
	params.Set("target", target)
	params.Set("mode", mode)
	if date != "" {
		params.Set("date", date)
	}

	if flags.DryRun {
		fmt.Printf("✓ Valid request. Would call: GET %s/site-explorer/domain-rating?%s\n",
			client.BaseURL, params.Encode())
		return nil
	}

	if flags.Verbose {
		fmt.Printf("Requesting: GET /site-explorer/domain-rating?%s\n", params.Encode())
	}

	// Make request
	resp, err := c.Get(context.Background(), "/site-explorer/domain-rating", params)
	if err != nil {
		w, _ := output.NewWriter(flags.OutputFormat, flags.OutputFile)
		w.WriteError(err)
		return err
	}

	// Parse response
	var result models.DomainRatingResponse
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	// Output result
	w, err := output.NewWriter(flags.OutputFormat, flags.OutputFile)
	if err != nil {
		return err
	}
	defer w.Close()

	return w.WriteSuccess(result, &resp.Meta)
}

func runBacklinksStats(target, mode, date string) error {
	flags := cmd.GetGlobalFlags()

	apiKey := flags.APIKey
	if apiKey == "" {
		apiKey = config.GetAPIKey()
	}
	if apiKey == "" {
		return fmt.Errorf("API key required")
	}

	c := client.NewClient(client.Config{
		APIKey: apiKey,
	})

	params := url.Values{}
	params.Set("target", target)
	params.Set("mode", mode)
	if date != "" {
		params.Set("date", date)
	}

	if flags.DryRun {
		fmt.Printf("✓ Valid request. Would call: GET %s/site-explorer/backlinks-stats?%s\n",
			client.BaseURL, params.Encode())
		return nil
	}

	if flags.Verbose {
		fmt.Printf("Requesting: GET /site-explorer/backlinks-stats?%s\n", params.Encode())
	}

	resp, err := c.Get(context.Background(), "/site-explorer/backlinks-stats", params)
	if err != nil {
		w, _ := output.NewWriter(flags.OutputFormat, flags.OutputFile)
		w.WriteError(err)
		return err
	}

	var result models.BacklinksStatsResponse
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	w, err := output.NewWriter(flags.OutputFormat, flags.OutputFile)
	if err != nil {
		return err
	}
	defer w.Close()

	return w.WriteSuccess(result, &resp.Meta)
}

func runBacklinks(target, mode string, limit, offset int, sel, where string) error {
	flags := cmd.GetGlobalFlags()

	apiKey := flags.APIKey
	if apiKey == "" {
		apiKey = config.GetAPIKey()
	}
	if apiKey == "" {
		return fmt.Errorf("API key required")
	}

	c := client.NewClient(client.Config{
		APIKey: apiKey,
	})

	params := url.Values{}
	params.Set("target", target)
	params.Set("mode", mode)
	params.Set("limit", fmt.Sprintf("%d", limit))
	if offset > 0 {
		params.Set("offset", fmt.Sprintf("%d", offset))
	}
	if sel != "" {
		params.Set("select", sel)
	}
	if where != "" {
		params.Set("where", where)
	}

	if flags.DryRun {
		fmt.Printf("✓ Valid request. Would call: GET %s/site-explorer/backlinks?%s\n",
			client.BaseURL, params.Encode())
		return nil
	}

	if flags.Verbose {
		fmt.Printf("Requesting: GET /site-explorer/backlinks?%s\n", params.Encode())
	}

	resp, err := c.Get(context.Background(), "/site-explorer/backlinks", params)
	if err != nil {
		w, _ := output.NewWriter(flags.OutputFormat, flags.OutputFile)
		w.WriteError(err)
		return err
	}

	var result models.BacklinksResponse
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	w, err := output.NewWriter(flags.OutputFormat, flags.OutputFile)
	if err != nil {
		return err
	}
	defer w.Close()

	return w.WriteSuccess(result, &resp.Meta)
}
