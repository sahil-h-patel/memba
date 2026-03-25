package providers

import (
	"fmt"
	"memba/lib/models"
	"memba/lib/utils"
	"net/http"
	"net/url"
)

type ProfilesProvider struct{}

func (p ProfilesProvider) Print(results []models.SearchResult) {
	var content string

	for _, res := range results {
		// Render the inner box for each result
		innerBox := utils.ResultStyle.Render(fmt.Sprintf("%s\n%s", res.Title, res.URL))
		content += innerBox + "\n"
	}

	// Wrap all results in the Provider's outer box
	fmt.Println(utils.ProviderStyle.Render(content))
}

func (p ProfilesProvider) Name() string { return "Profile" }

func (p ProfilesProvider) Search(client *http.Client, query string, token string) ([]models.SearchResult, error) {

	url := fmt.Sprintf("https://profiles.csh.rit.edu/search?q=%s&page=1", url.QueryEscape(query))
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
			Title:       "Profiles",
			Description: fmt.Sprintf("Profile for %s", query),
			URL:         fmt.Sprintf("https://profiles.csh.rit.edu/user/%s", query),
			Source:      p.Name(),
		},
	}, nil
}
