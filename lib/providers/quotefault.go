package providers

import (
	"encoding/json"
	"fmt"
	"memba/lib/models"
	"memba/lib/utils"
	"net/http"
	"net/url"
)

type QuotefaultProvider struct{}

func (q QuotefaultProvider) Print(results []models.SearchResult) {
	var content string

	for _, res := range results {
		// Render the inner box for each result
		innerBox := utils.ResultStyle.Render(fmt.Sprintf("%s\n%s", res.Title, res.URL))
		content += innerBox + "\n"
	}

	// Wrap all results in the Provider's outer box
	fmt.Println(utils.ProviderStyle.Render(content))
}

func (q QuotefaultProvider) Name() string { return "Wiki" }

func (q QuotefaultProvider) Search(client *http.Client, query string, token string) ([]models.SearchResult, error) {
	url := fmt.Sprintf("https://quotefault.csh.rit.edu/api/quotes?lt=0&q=%s&hidden=false", url.QueryEscape(query))
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Bad response: %d %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	type quoteJson struct {
		Submitter struct {
			Cn  string `json:"cn"`
			Uid string `json:"uid"`
		} `json:"submitter"`
		Timestamp string `json:"timestamp"`
		Shards    struct {
			Body    string `json:"body"`
			Speaker struct {
				Cn  string `json:"cn"`
				Uid string `json:"uid"`
			} `json:"speaker"`
		}
		Id    int `json:"size"`
		Score int `json:"wordcount"`
	}

	var data struct {
		Query struct {
			Search []quoteJson `json:"search"`
		} `json:"query"`
	}

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, fmt.Errorf("could not parse JSON: %w", err)
	}

	var finalResults []models.SearchResult
	for _, result := range data.Query.Search {
		finalResults = append(finalResults,
			models.SearchResult{
				Title:       result.Shards.Body,
				Description: fmt.Sprintf("Submitted by %s, Speaker %s", result.Submitter.Cn, result.Shards.Speaker.Cn),
				URL:         fmt.Sprintf("https://quotefault.csh.rit.edu/storage?q=%s", query),
				Source:      q.Name(),
			},
		)
	}
	return finalResults, nil
}
