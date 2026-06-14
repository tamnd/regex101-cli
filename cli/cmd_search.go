package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func (a *App) searchCmd() *cobra.Command {
	var flavor string
	cmd := &cobra.Command{
		Use:   "search <query>",
		Short: "Search regex patterns by keyword",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			query := args[0]
			if query == "" {
				return codeError(exitUsage, fmt.Errorf("query cannot be empty"))
			}
			limit := a.effectiveLimit(20)
			patterns, err := a.client.Search(cmd.Context(), query, flavor, limit)
			if err != nil {
				return mapFetchErr(err)
			}
			return a.renderOrEmpty(patterns, len(patterns))
		},
	}
	cmd.Flags().StringVar(&flavor, "flavor", "", "filter by flavor: javascript|python|pcre|pcre2|php|golang|java|ruby|rust|csharp")
	return cmd
}
