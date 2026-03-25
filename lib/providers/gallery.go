package providers

import (
	"encoding/json"
	"fmt"
	"memba/lib/models"
	"memba/lib/utils"
	"net/http"
	"net/url"
)

type GalleryProvider struct{}

func (g GalleryProvider) Print(results []models.SearchResult) {
	var content string

	for _, res := range results {
		// Render the inner box for each result
		innerBox := utils.ResultStyle.Render(fmt.Sprintf("%s\n%s", res.Title, res.URL))
		content += innerBox + "\n"
	}

	// Wrap all results in the Provider's outer box
	fmt.Println(utils.ProviderStyle.Render(content))
}

func (g GalleryProvider) Name() string { return "Gallery" }

func (g GalleryProvider) Search(client *http.Client, query string, token string) ([]models.SearchResult, error) {
	url := fmt.Sprintf("https://gallery.csh.rit.edu/api/directory/get/%s", url.QueryEscape(query))
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
	var data []json.RawMessage
	json.NewDecoder(resp.Body).Decode(&data)
	if len(data) == 0 {
		return nil, nil
	}
	return []models.SearchResult{
		{
			Title:       "Profiles",
			Description: fmt.Sprintf("Gallery directory for %s", query),
			URL:         fmt.Sprintf("https://gallery.csh.rit.edu/view/dir/%s", query),
			Source:      g.Name(),
		},
	}, nil
}
