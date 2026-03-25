package cmd

import (
	"fmt"
	"memba/lib/auth"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "memba",
	Short: "A search tool to query through all of CSH services",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Skip token check if the user is literally trying to login or logout
		if cmd.Name() == "login" || cmd.Name() == "logout" || cmd.Name() == "help" {
			return
		}

		token, err := auth.GetToken()
		if err != nil || len(token) == 0 {
			fmt.Println("No valid token found. Please run: memba login")
			os.Exit(1) // Stop the search before it starts
		}
	},
}

var apiClient *http.Client

func Execute(client *http.Client) {
	apiClient = client
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
