package siteexplorer

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/aminemat/ahrefs-cli/cmd"
	"github.com/aminemat/ahrefs-cli/internal/config"
	"github.com/aminemat/ahrefs-cli/pkg/client"
	"github.com/aminemat/ahrefs-cli/pkg/models"
	"github.com/aminemat/ahrefs-cli/pkg/output"
	"github.com/spf13/cobra"
)

// newAnchorsCmd creates the anchors command
func newAnchorsCmd() *cobra.Command {
	var (
		target  string
		mode    string
		limit   int
		offset  int
		sel     string
		where   string
		orderBy string
	)

	c := &cobra.Command{
		Use:   "anchors",
		Short: "Get anchor text distribution",
		Long:  "List anchor texts used in backlinks pointing to the target.",
		Example: `  # Get anchor texts for a domain
  ahrefs site-explorer anchors --target example.com --limit 100

  # Get anchor texts with backlink count
  ahrefs site-explorer anchors --target example.com \
    --select anchor,backlinks,refdomains --limit 50`,
		RunE: func(cobraCmd *cobra.Command, args []string) error {
			return runAnchors(target, mode, limit, offset, sel, where, orderBy)
		},
	}

	c.Flags().StringVar(&target, "target", "", "Target domain or URL (required)")
	c.Flags().StringVar(&mode, "mode", "domain", "Mode: exact, domain, prefix, subdomains")
	c.Flags().IntVar(&limit, "limit", 100, "Maximum number of results")
	c.Flags().IntVar(&offset, "offset", 0, "Offset for pagination")
	c.Flags().StringVar(&sel, "select", "", "Comma-separated list of fields to return")
	c.Flags().StringVar(&where, "where", "", "Filter expression (Ahrefs filter syntax)")
	c.Flags().StringVar(&orderBy, "order-by", "", "Sort order (e.g., backlinks:desc)")

	c.MarkFlagRequired("target")

	return c
}

func runAnchors(target, mode string, limit, offset int, sel, where, orderBy string) error {
	flags := cmd.GetGlobalFlags()

	apiKey := flags.APIKey
	if apiKey == "" {
		apiKey = config.GetAPIKey()
	}
	if apiKey == "" {
		return fmt.Errorf("API key required")
	}

	c := client.NewClient(client.Config{APIKey: apiKey})

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
	if orderBy != "" {
		params.Set("order_by", orderBy)
	}

	if flags.DryRun {
		fmt.Printf("✓ Valid request. Would call: GET %s/site-explorer/anchors?%s\n",
			client.BaseURL, params.Encode())
		return nil
	}

	if flags.Verbose {
		fmt.Printf("Requesting: GET /site-explorer/anchors?%s\n", params.Encode())
	}

	resp, err := c.Get(context.Background(), "/site-explorer/anchors", params)
	if err != nil {
		w, _ := output.NewWriter(flags.OutputFormat, flags.OutputFile)
		w.WriteError(err)
		return err
	}

	var result models.AnchorsResponse
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

// newOrganicKeywordsCmd creates the organic-keywords command
func newOrganicKeywordsCmd() *cobra.Command {
	var (
		target  string
		mode    string
		limit   int
		offset  int
		sel     string
		where   string
		orderBy string
		country string
	)

	c := &cobra.Command{
		Use:   "organic-keywords",
		Short: "Get organic keywords",
		Long:  "List organic keywords that the target ranks for in search engines.",
		Example: `  # Get organic keywords for a domain
  ahrefs site-explorer organic-keywords --target example.com --limit 100

  # Get keywords for a specific country
  ahrefs site-explorer organic-keywords --target example.com \
    --country us --limit 50

  # Get high-traffic keywords
  ahrefs site-explorer organic-keywords --target example.com \
    --where 'traffic>100' --order-by traffic:desc --limit 100`,
		RunE: func(cobraCmd *cobra.Command, args []string) error {
			return runOrganicKeywords(target, mode, limit, offset, sel, where, orderBy, country)
		},
	}

	c.Flags().StringVar(&target, "target", "", "Target domain or URL (required)")
	c.Flags().StringVar(&mode, "mode", "domain", "Mode: exact, domain, prefix, subdomains")
	c.Flags().IntVar(&limit, "limit", 100, "Maximum number of results")
	c.Flags().IntVar(&offset, "offset", 0, "Offset for pagination")
	c.Flags().StringVar(&sel, "select", "", "Comma-separated list of fields to return")
	c.Flags().StringVar(&where, "where", "", "Filter expression (Ahrefs filter syntax)")
	c.Flags().StringVar(&orderBy, "order-by", "", "Sort order (e.g., traffic:desc)")
	c.Flags().StringVar(&country, "country", "", "Country code (e.g., us, gb, de)")

	c.MarkFlagRequired("target")

	return c
}

func runOrganicKeywords(target, mode string, limit, offset int, sel, where, orderBy, country string) error {
	flags := cmd.GetGlobalFlags()

	apiKey := flags.APIKey
	if apiKey == "" {
		apiKey = config.GetAPIKey()
	}
	if apiKey == "" {
		return fmt.Errorf("API key required")
	}

	c := client.NewClient(client.Config{APIKey: apiKey})

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
	if orderBy != "" {
		params.Set("order_by", orderBy)
	}
	if country != "" {
		params.Set("country", country)
	}

	if flags.DryRun {
		fmt.Printf("✓ Valid request. Would call: GET %s/site-explorer/organic-keywords?%s\n",
			client.BaseURL, params.Encode())
		return nil
	}

	if flags.Verbose {
		fmt.Printf("Requesting: GET /site-explorer/organic-keywords?%s\n", params.Encode())
	}

	resp, err := c.Get(context.Background(), "/site-explorer/organic-keywords", params)
	if err != nil {
		w, _ := output.NewWriter(flags.OutputFormat, flags.OutputFile)
		w.WriteError(err)
		return err
	}

	var result models.OrganicKeywordsResponse
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

// newTopPagesCmd creates the top-pages command
func newTopPagesCmd() *cobra.Command {
	var (
		target  string
		mode    string
		limit   int
		offset  int
		sel     string
		where   string
		orderBy string
		country string
	)

	c := &cobra.Command{
		Use:   "top-pages",
		Short: "Get top pages by organic traffic",
		Long:  "List pages that receive the most organic search traffic.",
		Example: `  # Get top pages for a domain
  ahrefs site-explorer top-pages --target example.com --limit 100

  # Get top pages in a specific country
  ahrefs site-explorer top-pages --target example.com \
    --country us --limit 50

  # Get top pages with specific fields
  ahrefs site-explorer top-pages --target example.com \
    --select url,traffic,keywords --limit 100`,
		RunE: func(cobraCmd *cobra.Command, args []string) error {
			return runTopPages(target, mode, limit, offset, sel, where, orderBy, country)
		},
	}

	c.Flags().StringVar(&target, "target", "", "Target domain or URL (required)")
	c.Flags().StringVar(&mode, "mode", "domain", "Mode: exact, domain, prefix, subdomains")
	c.Flags().IntVar(&limit, "limit", 100, "Maximum number of results")
	c.Flags().IntVar(&offset, "offset", 0, "Offset for pagination")
	c.Flags().StringVar(&sel, "select", "", "Comma-separated list of fields to return")
	c.Flags().StringVar(&where, "where", "", "Filter expression (Ahrefs filter syntax)")
	c.Flags().StringVar(&orderBy, "order-by", "", "Sort order (e.g., traffic:desc)")
	c.Flags().StringVar(&country, "country", "", "Country code (e.g., us, gb, de)")

	c.MarkFlagRequired("target")

	return c
}

func runTopPages(target, mode string, limit, offset int, sel, where, orderBy, country string) error {
	flags := cmd.GetGlobalFlags()

	apiKey := flags.APIKey
	if apiKey == "" {
		apiKey = config.GetAPIKey()
	}
	if apiKey == "" {
		return fmt.Errorf("API key required")
	}

	c := client.NewClient(client.Config{APIKey: apiKey})

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
	if orderBy != "" {
		params.Set("order_by", orderBy)
	}
	if country != "" {
		params.Set("country", country)
	}

	if flags.DryRun {
		fmt.Printf("✓ Valid request. Would call: GET %s/site-explorer/top-pages?%s\n",
			client.BaseURL, params.Encode())
		return nil
	}

	if flags.Verbose {
		fmt.Printf("Requesting: GET /site-explorer/top-pages?%s\n", params.Encode())
	}

	resp, err := c.Get(context.Background(), "/site-explorer/top-pages", params)
	if err != nil {
		w, _ := output.NewWriter(flags.OutputFormat, flags.OutputFile)
		w.WriteError(err)
		return err
	}

	var result models.TopPagesResponse
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

// newBrokenBacklinksCmd creates the broken-backlinks command
func newBrokenBacklinksCmd() *cobra.Command {
	var (
		target  string
		mode    string
		limit   int
		offset  int
		sel     string
		where   string
		orderBy string
	)

	c := &cobra.Command{
		Use:   "broken-backlinks",
		Short: "Get broken backlinks",
		Long:  "List backlinks pointing to non-existing pages (404 errors) on the target.",
		Example: `  # Get broken backlinks for a domain
  ahrefs site-explorer broken-backlinks --target example.com --limit 100

  # Get broken backlinks sorted by domain rating
  ahrefs site-explorer broken-backlinks --target example.com \
    --order-by domain_rating:desc --limit 50`,
		RunE: func(cobraCmd *cobra.Command, args []string) error {
			return runBrokenBacklinks(target, mode, limit, offset, sel, where, orderBy)
		},
	}

	c.Flags().StringVar(&target, "target", "", "Target domain or URL (required)")
	c.Flags().StringVar(&mode, "mode", "domain", "Mode: exact, domain, prefix, subdomains")
	c.Flags().IntVar(&limit, "limit", 100, "Maximum number of results")
	c.Flags().IntVar(&offset, "offset", 0, "Offset for pagination")
	c.Flags().StringVar(&sel, "select", "", "Comma-separated list of fields to return")
	c.Flags().StringVar(&where, "where", "", "Filter expression (Ahrefs filter syntax)")
	c.Flags().StringVar(&orderBy, "order-by", "", "Sort order (e.g., domain_rating:desc)")

	c.MarkFlagRequired("target")

	return c
}

func runBrokenBacklinks(target, mode string, limit, offset int, sel, where, orderBy string) error {
	flags := cmd.GetGlobalFlags()

	apiKey := flags.APIKey
	if apiKey == "" {
		apiKey = config.GetAPIKey()
	}
	if apiKey == "" {
		return fmt.Errorf("API key required")
	}

	c := client.NewClient(client.Config{APIKey: apiKey})

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
	if orderBy != "" {
		params.Set("order_by", orderBy)
	}

	if flags.DryRun {
		fmt.Printf("✓ Valid request. Would call: GET %s/site-explorer/broken-backlinks?%s\n",
			client.BaseURL, params.Encode())
		return nil
	}

	if flags.Verbose {
		fmt.Printf("Requesting: GET /site-explorer/broken-backlinks?%s\n", params.Encode())
	}

	resp, err := c.Get(context.Background(), "/site-explorer/broken-backlinks", params)
	if err != nil {
		w, _ := output.NewWriter(flags.OutputFormat, flags.OutputFile)
		w.WriteError(err)
		return err
	}

	var result models.BrokenBacklinksResponse
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

// newLinkedDomainsCmd creates the linked-domains command
func newLinkedDomainsCmd() *cobra.Command {
	var (
		target  string
		mode    string
		limit   int
		offset  int
		sel     string
		where   string
		orderBy string
	)

	c := &cobra.Command{
		Use:   "linked-domains",
		Short: "Get linked domains",
		Long:  "List domains that the target links out to.",
		Example: `  # Get linked domains for a domain
  ahrefs site-explorer linked-domains --target example.com --limit 100

  # Filter by domain rating
  ahrefs site-explorer linked-domains --target example.com \
    --where 'domain_rating>50' --order-by domain_rating:desc --limit 50`,
		RunE: func(cobraCmd *cobra.Command, args []string) error {
			return runLinkedDomains(target, mode, limit, offset, sel, where, orderBy)
		},
	}

	c.Flags().StringVar(&target, "target", "", "Target domain or URL (required)")
	c.Flags().StringVar(&mode, "mode", "domain", "Mode: exact, domain, prefix, subdomains")
	c.Flags().IntVar(&limit, "limit", 100, "Maximum number of results")
	c.Flags().IntVar(&offset, "offset", 0, "Offset for pagination")
	c.Flags().StringVar(&sel, "select", "", "Comma-separated list of fields to return")
	c.Flags().StringVar(&where, "where", "", "Filter expression (Ahrefs filter syntax)")
	c.Flags().StringVar(&orderBy, "order-by", "", "Sort order (e.g., domain_rating:desc)")

	c.MarkFlagRequired("target")

	return c
}

func runLinkedDomains(target, mode string, limit, offset int, sel, where, orderBy string) error {
	flags := cmd.GetGlobalFlags()

	apiKey := flags.APIKey
	if apiKey == "" {
		apiKey = config.GetAPIKey()
	}
	if apiKey == "" {
		return fmt.Errorf("API key required")
	}

	c := client.NewClient(client.Config{APIKey: apiKey})

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
	if orderBy != "" {
		params.Set("order_by", orderBy)
	}

	if flags.DryRun {
		fmt.Printf("✓ Valid request. Would call: GET %s/site-explorer/linked-domains?%s\n",
			client.BaseURL, params.Encode())
		return nil
	}

	if flags.Verbose {
		fmt.Printf("Requesting: GET /site-explorer/linked-domains?%s\n", params.Encode())
	}

	resp, err := c.Get(context.Background(), "/site-explorer/linked-domains", params)
	if err != nil {
		w, _ := output.NewWriter(flags.OutputFormat, flags.OutputFile)
		w.WriteError(err)
		return err
	}

	var result models.LinkedDomainsResponse
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

// newMetricsCmd creates the metrics command
func newMetricsCmd() *cobra.Command {
	var (
		target  string
		mode    string
		sel     string
		country string
	)

	c := &cobra.Command{
		Use:   "metrics",
		Short: "Get site metrics overview",
		Long:  "Get organic and paid traffic metrics for a target.",
		Example: `  # Get metrics for a domain
  ahrefs site-explorer metrics --target example.com

  # Get metrics for a specific country
  ahrefs site-explorer metrics --target example.com --country us`,
		RunE: func(cobraCmd *cobra.Command, args []string) error {
			return runMetrics(target, mode, sel, country)
		},
	}

	c.Flags().StringVar(&target, "target", "", "Target domain or URL (required)")
	c.Flags().StringVar(&mode, "mode", "domain", "Mode: exact, domain, prefix, subdomains")
	c.Flags().StringVar(&sel, "select", "", "Comma-separated list of fields to return")
	c.Flags().StringVar(&country, "country", "", "Country code (e.g., us, gb, de)")

	c.MarkFlagRequired("target")

	return c
}

func runMetrics(target, mode, sel, country string) error {
	flags := cmd.GetGlobalFlags()

	apiKey := flags.APIKey
	if apiKey == "" {
		apiKey = config.GetAPIKey()
	}
	if apiKey == "" {
		return fmt.Errorf("API key required")
	}

	c := client.NewClient(client.Config{APIKey: apiKey})

	params := url.Values{}
	params.Set("target", target)
	params.Set("mode", mode)
	if sel != "" {
		params.Set("select", sel)
	}
	if country != "" {
		params.Set("country", country)
	}

	if flags.DryRun {
		fmt.Printf("✓ Valid request. Would call: GET %s/site-explorer/metrics?%s\n",
			client.BaseURL, params.Encode())
		return nil
	}

	if flags.Verbose {
		fmt.Printf("Requesting: GET /site-explorer/metrics?%s\n", params.Encode())
	}

	resp, err := c.Get(context.Background(), "/site-explorer/metrics", params)
	if err != nil {
		w, _ := output.NewWriter(flags.OutputFormat, flags.OutputFile)
		w.WriteError(err)
		return err
	}

	var result models.MetricsResponse
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

// newMetricsHistoryCmd creates the metrics-history command
func newMetricsHistoryCmd() *cobra.Command {
	var (
		target   string
		mode     string
		sel      string
		country  string
		dateFrom string
		dateTo   string
	)

	c := &cobra.Command{
		Use:   "metrics-history",
		Short: "Get historical metrics",
		Long:  "Get historical organic and paid traffic metrics for a target.",
		Example: `  # Get metrics history for a domain
  ahrefs site-explorer metrics-history --target example.com

  # Get metrics history for a specific date range
  ahrefs site-explorer metrics-history --target example.com \
    --date-from 2024-01-01 --date-to 2024-12-31

  # Get metrics history for a specific country
  ahrefs site-explorer metrics-history --target example.com --country us`,
		RunE: func(cobraCmd *cobra.Command, args []string) error {
			return runMetricsHistory(target, mode, sel, country, dateFrom, dateTo)
		},
	}

	c.Flags().StringVar(&target, "target", "", "Target domain or URL (required)")
	c.Flags().StringVar(&mode, "mode", "domain", "Mode: exact, domain, prefix, subdomains")
	c.Flags().StringVar(&sel, "select", "", "Comma-separated list of fields to return")
	c.Flags().StringVar(&country, "country", "", "Country code (e.g., us, gb, de)")
	c.Flags().StringVar(&dateFrom, "date-from", "", "Start date (YYYY-MM-DD)")
	c.Flags().StringVar(&dateTo, "date-to", "", "End date (YYYY-MM-DD)")

	c.MarkFlagRequired("target")

	return c
}

func runMetricsHistory(target, mode, sel, country, dateFrom, dateTo string) error {
	flags := cmd.GetGlobalFlags()

	apiKey := flags.APIKey
	if apiKey == "" {
		apiKey = config.GetAPIKey()
	}
	if apiKey == "" {
		return fmt.Errorf("API key required")
	}

	c := client.NewClient(client.Config{APIKey: apiKey})

	params := url.Values{}
	params.Set("target", target)
	params.Set("mode", mode)
	if sel != "" {
		params.Set("select", sel)
	}
	if country != "" {
		params.Set("country", country)
	}
	if dateFrom != "" {
		params.Set("date_from", dateFrom)
	}
	if dateTo != "" {
		params.Set("date_to", dateTo)
	}

	if flags.DryRun {
		fmt.Printf("✓ Valid request. Would call: GET %s/site-explorer/metrics-history?%s\n",
			client.BaseURL, params.Encode())
		return nil
	}

	if flags.Verbose {
		fmt.Printf("Requesting: GET /site-explorer/metrics-history?%s\n", params.Encode())
	}

	resp, err := c.Get(context.Background(), "/site-explorer/metrics-history", params)
	if err != nil {
		w, _ := output.NewWriter(flags.OutputFormat, flags.OutputFile)
		w.WriteError(err)
		return err
	}

	var result models.MetricsHistoryResponse
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

// newPagesByTrafficCmd creates the pages-by-traffic command
func newPagesByTrafficCmd() *cobra.Command {
	var (
		target  string
		mode    string
		limit   int
		offset  int
		sel     string
		where   string
		orderBy string
		country string
	)

	c := &cobra.Command{
		Use:   "pages-by-traffic",
		Short: "Get pages sorted by traffic",
		Long:  "List pages sorted by organic search traffic.",
		Example: `  # Get pages by traffic for a domain
  ahrefs site-explorer pages-by-traffic --target example.com --limit 100

  # Get pages by traffic for a specific country
  ahrefs site-explorer pages-by-traffic --target example.com \
    --country us --limit 50`,
		RunE: func(cobraCmd *cobra.Command, args []string) error {
			return runPagesByTraffic(target, mode, limit, offset, sel, where, orderBy, country)
		},
	}

	c.Flags().StringVar(&target, "target", "", "Target domain or URL (required)")
	c.Flags().StringVar(&mode, "mode", "domain", "Mode: exact, domain, prefix, subdomains")
	c.Flags().IntVar(&limit, "limit", 100, "Maximum number of results")
	c.Flags().IntVar(&offset, "offset", 0, "Offset for pagination")
	c.Flags().StringVar(&sel, "select", "", "Comma-separated list of fields to return")
	c.Flags().StringVar(&where, "where", "", "Filter expression (Ahrefs filter syntax)")
	c.Flags().StringVar(&orderBy, "order-by", "", "Sort order (e.g., traffic:desc)")
	c.Flags().StringVar(&country, "country", "", "Country code (e.g., us, gb, de)")

	c.MarkFlagRequired("target")

	return c
}

func runPagesByTraffic(target, mode string, limit, offset int, sel, where, orderBy, country string) error {
	flags := cmd.GetGlobalFlags()

	apiKey := flags.APIKey
	if apiKey == "" {
		apiKey = config.GetAPIKey()
	}
	if apiKey == "" {
		return fmt.Errorf("API key required")
	}

	c := client.NewClient(client.Config{APIKey: apiKey})

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
	if orderBy != "" {
		params.Set("order_by", orderBy)
	}
	if country != "" {
		params.Set("country", country)
	}

	if flags.DryRun {
		fmt.Printf("✓ Valid request. Would call: GET %s/site-explorer/pages-by-traffic?%s\n",
			client.BaseURL, params.Encode())
		return nil
	}

	if flags.Verbose {
		fmt.Printf("Requesting: GET /site-explorer/pages-by-traffic?%s\n", params.Encode())
	}

	resp, err := c.Get(context.Background(), "/site-explorer/pages-by-traffic", params)
	if err != nil {
		w, _ := output.NewWriter(flags.OutputFormat, flags.OutputFile)
		w.WriteError(err)
		return err
	}

	var result models.PagesByTrafficResponse
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

// newBestByLinksCmd creates the best-by-links command
func newBestByLinksCmd() *cobra.Command {
	var (
		target  string
		mode    string
		limit   int
		offset  int
		sel     string
		where   string
		orderBy string
	)

	c := &cobra.Command{
		Use:   "best-by-links",
		Short: "Get best pages by backlinks",
		Long:  "List pages sorted by the number of backlinks they receive.",
		Example: `  # Get best pages by links for a domain
  ahrefs site-explorer best-by-links --target example.com --limit 100

  # Get pages with most referring domains
  ahrefs site-explorer best-by-links --target example.com \
    --order-by refdomains:desc --limit 50`,
		RunE: func(cobraCmd *cobra.Command, args []string) error {
			return runBestByLinks(target, mode, limit, offset, sel, where, orderBy)
		},
	}

	c.Flags().StringVar(&target, "target", "", "Target domain or URL (required)")
	c.Flags().StringVar(&mode, "mode", "domain", "Mode: exact, domain, prefix, subdomains")
	c.Flags().IntVar(&limit, "limit", 100, "Maximum number of results")
	c.Flags().IntVar(&offset, "offset", 0, "Offset for pagination")
	c.Flags().StringVar(&sel, "select", "", "Comma-separated list of fields to return")
	c.Flags().StringVar(&where, "where", "", "Filter expression (Ahrefs filter syntax)")
	c.Flags().StringVar(&orderBy, "order-by", "", "Sort order (e.g., backlinks:desc)")

	c.MarkFlagRequired("target")

	return c
}

func runBestByLinks(target, mode string, limit, offset int, sel, where, orderBy string) error {
	flags := cmd.GetGlobalFlags()

	apiKey := flags.APIKey
	if apiKey == "" {
		apiKey = config.GetAPIKey()
	}
	if apiKey == "" {
		return fmt.Errorf("API key required")
	}

	c := client.NewClient(client.Config{APIKey: apiKey})

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
	if orderBy != "" {
		params.Set("order_by", orderBy)
	}

	if flags.DryRun {
		fmt.Printf("✓ Valid request. Would call: GET %s/site-explorer/best-by-links?%s\n",
			client.BaseURL, params.Encode())
		return nil
	}

	if flags.Verbose {
		fmt.Printf("Requesting: GET /site-explorer/best-by-links?%s\n", params.Encode())
	}

	resp, err := c.Get(context.Background(), "/site-explorer/best-by-links", params)
	if err != nil {
		w, _ := output.NewWriter(flags.OutputFormat, flags.OutputFile)
		w.WriteError(err)
		return err
	}

	var result models.BestByLinksResponse
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
