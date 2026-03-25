package lib

import (
	"fmt"
	"memba/lib/models"
	"memba/lib/providers"
	"net/http"
	"sync"
)

func Run(client *http.Client, query string, token string) {

	wiki := &providers.WikiProvider{}

	engine := models.Engine{
		Providers: []models.Provider{
			&providers.DiscourseProvider{},
			&providers.GalleryProvider{},
			&providers.ProfilesProvider{},
			&providers.QuotefaultProvider{},
			wiki,
		},
	}
	fmt.Println("Authenticating with CSH Wiki...")
	if err := wiki.AuthBot(client); err != nil {
		fmt.Printf("Warning: Wiki authentication failed: %v\n", err)
	}
	resultsChan := make(chan []models.SearchResult, len(engine.Providers))
	var wg sync.WaitGroup

	for _, provider := range engine.Providers {
		wg.Add(1)
		go func(p models.Provider) {
			defer wg.Done()
			results, err := p.Search(client, query, token)
			if err != nil {
				fmt.Printf("Error from %s: %v\n", p.Name(), err)
				return
			}
			resultsChan <- results
			p.Print(results)
		}(provider)
	}

	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	var allResults []models.SearchResult
	for res := range resultsChan {
		allResults = append(allResults, res...)
	}
	fmt.Println("All searches complete!")
}
