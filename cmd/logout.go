package cmd

import (
	"memba/lib/auth"
	"time"

	"github.com/spf13/cobra"
	webview "github.com/webview/webview_go"
)

func init() {
	rootCmd.AddCommand(logoutCmd)
}

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Log out of CSH",
	Run: func(cmd *cobra.Command, args []string) {
		go auth.StartServer()
		w := webview.New(true)
		defer w.Destroy()

		w.SetTitle("CSH Auth Login")
		w.SetSize(480, 600, webview.HintNone)
		w.Navigate("http://localhost:8080/auth/logout")

		go func() {
			<-auth.LogoutDone
			time.Sleep(1 * time.Second) // Give the user a moment to see the "Success" message
			w.Terminate()
		}()

		w.Run()
	},
}
