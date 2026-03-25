package providers

import (
	"fmt"
	"memba/lib/models"
	"memba/lib/utils"
	"net/http"
	"net/url"
)

type DiscourseProvider struct{}

func (d DiscourseProvider) Print(results []models.SearchResult) {
	var content string

	for _, res := range results {
		// Render the inner box for each result
		innerBox := utils.ResultStyle.Render(fmt.Sprintf("%s\n%s", res.Title, res.URL))
		content += innerBox + "\n"
	}

	// Wrap all results in the Provider's outer box
	fmt.Println(utils.ProviderStyle.Render(content))
}

func (d DiscourseProvider) Name() string { return "Discourse" }

func (d DiscourseProvider) Search(client *http.Client, query string, token string) ([]models.SearchResult, error) {

	url := fmt.Sprintf("https://discourse.csh.rit.edu/search?q=%s", url.QueryEscape(query))
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return []models.SearchResult{
		{
			Title:       "Discourse",
			Description: fmt.Sprintf("Discourse post for %s", query),
			URL:         fmt.Sprintf("https://discourse.csh.rit.edu/search?q=%s", query),
			Source:      d.Name(),
		},
	}, nil
}
