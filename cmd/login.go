package cmd

import (
	"fmt"
	"memba/lib/auth"
	"time"

	"github.com/spf13/cobra"
	webview "github.com/webview/webview_go"
)

func init() {
	rootCmd.AddCommand(loginCmd)
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Log in to CSH",
	Run: func(cmd *cobra.Command, args []string) {
		token, err := auth.GetToken()
		if err != nil || len(token) == 0 {
			fmt.Printf("No token found, launching login...")
			Login()

			token, err = auth.GetToken()
			if err != nil {
				fmt.Println("Login failed or was cancelled.")
				return
			}
		}
	},
}

func Login() {
	go auth.StartServer()
	w := webview.New(true)
	defer w.Destroy()

	w.SetTitle("CSH Auth Login")
	w.SetSize(480, 600, webview.HintNone)
	w.Navigate("http://localhost:8080/auth/login")

	go func() {
		<-auth.LoginDone
		time.Sleep(1 * time.Second)
		w.Terminate()
	}()

	w.Run()
}
