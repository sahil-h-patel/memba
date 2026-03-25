package cmd

import (
	"fmt"
	engine "memba/lib"
	"memba/lib/auth"
	"memba/lib/providers"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(searchCmd)
	searchCmd.AddCommand(wikiSearchCmd)
	searchCmd.AddCommand(discourseSearchCmd)
	searchCmd.AddCommand(gallerySearchCmd)
	searchCmd.AddCommand(profilesSearchCmd)
	searchCmd.AddCommand(quotefaultSearchCmd)
}

var searchCmd = &cobra.Command{
	Use:   "search [source] [query]",
	Short: "search through ALL CSH services",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			return
		}
		query := args[0]
		token, err := auth.GetToken()
		if err != nil {
			fmt.Println("Login failed or was cancelled.")
			return
		}
		engine.Run(apiClient, query, string(token))
	},
}

var wikiSearchCmd = &cobra.Command{
	Use:   "search [source] [query]",
	Short: "search through ALL CSH services",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			return
		}
		query := args[0]
		token, err := auth.GetToken()
		if err != nil {
			fmt.Println("Login failed or was cancelled.")
			return
		}
		wiki := providers.WikiProvider{}
		wiki.Search(apiClient, query, string(token))
	},
}

var discourseSearchCmd = &cobra.Command{
	Use:   "search [source] [query]",
	Short: "search through ALL CSH services",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			return
		}
		query := args[0]
		token, err := auth.GetToken()
		if err != nil {
			fmt.Println("Login failed or was cancelled.")
			return
		}
		wiki := providers.DiscourseProvider{}
		wiki.Search(apiClient, query, string(token))
	},
}

var gallerySearchCmd = &cobra.Command{
	Use:   "search [source] [query]",
	Short: "search through ALL CSH services",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			return
		}
		query := args[0]
		token, err := auth.GetToken()
		if err != nil {
			fmt.Println("Login failed or was cancelled.")
			return
		}
		wiki := providers.GalleryProvider{}
		wiki.Search(apiClient, query, string(token))
	},
}

var profilesSearchCmd = &cobra.Command{
	Use:   "search [source] [query]",
	Short: "search through ALL CSH services",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			return
		}
		query := args[0]
		token, err := auth.GetToken()
		if err != nil {
			fmt.Println("Login failed or was cancelled.")
			return
		}
		wiki := providers.ProfilesProvider{}
		wiki.Search(apiClient, query, string(token))
	},
}

var quotefaultSearchCmd = &cobra.Command{
	Use:   "search [source] [query]",
	Short: "search through ALL CSH services",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			return
		}
		query := args[0]
		token, err := auth.GetToken()
		if err != nil {
			fmt.Println("Login failed or was cancelled.")
			return
		}
		wiki := providers.QuotefaultProvider{}
		wiki.Search(apiClient, query, string(token))
	},
}
