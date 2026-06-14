package cli

import (
	"github.com/spf13/cobra"
)

func (a *App) listCmd() *cobra.Command {
	var flavor string
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List top regex patterns by upvotes",
		RunE: func(cmd *cobra.Command, _ []string) error {
			limit := a.effectiveLimit(20)
			patterns, err := a.client.List(cmd.Context(), flavor, limit)
			if err != nil {
				return mapFetchErr(err)
			}
			return a.renderOrEmpty(patterns, len(patterns))
		},
	}
	cmd.Flags().StringVar(&flavor, "flavor", "", "filter by flavor: javascript|python|pcre|pcre2|php|golang|java|ruby|rust|csharp")
	return cmd
}
