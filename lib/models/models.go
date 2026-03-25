package models

import "net/http"

type Engine struct {
	Providers []Provider
}

type SearchResult struct {
	Title       string
	Description string
	URL         string
	Source      string
}

type Provider interface {
	Name() string
	Print(results []SearchResult)
	Search(client *http.Client, query string, token string) ([]SearchResult, error)
}
