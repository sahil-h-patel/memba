package main

import (
	"memba/cmd"
	"net/http"
	"net/http/cookiejar"
	"time"
)

func main() {
	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar:     jar,
		Timeout: 10 * time.Second,
	}
	cmd.Execute(client)
}
