package providers

import (
	"encoding/json"
	"fmt"
	"io"
	"memba/lib/models"
	"memba/lib/utils"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type WikiProvider struct {
	isAuthed bool
}

func (w *WikiProvider) AuthBot(client *http.Client) (err error) {
	if w.isAuthed {
		return nil
	}

	apiUrl := "https://wiki.csh.rit.edu/api.php"

	// get token for Wiki
	u, _ := url.Parse(apiUrl)
	params := u.Query()
	params.Set("action", "query")
	params.Set("meta", "tokens")
	params.Set("type", "login")
	params.Set("format", "json")
	u.RawQuery = params.Encode()
	tokenResp, err := client.Get(u.String())
	if err != nil {
		return err
	}

	var tokenData struct {
		Query struct {
			Tokens struct {
				LoginToken string `json:"logintoken"`
			} `json:"tokens"`
		} `json:"query"`
	}
	if err := json.NewDecoder(tokenResp.Body).Decode(&tokenData); err != nil {
		return fmt.Errorf("failed to decode token: %w", err)
	}
	if tokenData.Query.Tokens.LoginToken == "" {
		return fmt.Errorf("received empty login token from Wiki")
	}
	tokenResp.Body.Close()

	fmt.Printf("Token: %s\n", tokenData.Query.Tokens.LoginToken)
	fmt.Printf("Cookies after token req: %v\n", client.Jar.Cookies(u))

	// login req passing into token to wiki Bot
	loginData := url.Values{}
	loginData.Set("action", "login")
	loginData.Set("lgname", os.Getenv("WIKI_BOT_USER"))
	loginData.Set("lgpassword", os.Getenv("WIKI_BOT_PASSWORD"))
	loginData.Set("lgtoken", tokenData.Query.Tokens.LoginToken)
	loginData.Set("format", "json")
	// rawBody := fmt.Sprintf("action=login&lgname=%s&lgpassword=%s&lgtoken=%s&format=json",
	// 	url.QueryEscape(os.Getenv("WIKI_BOT_USER")),
	// 	url.QueryEscape(os.Getenv("WIKI_BOT_PASSWORD")),
	// 	tokenData.Query.Tokens.LoginToken)
	loginReq, _ := http.NewRequest("POST", apiUrl, strings.NewReader(loginData.Encode()))
	loginReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	loginResp, err := client.Do(loginReq)
	if err != nil {
		return err
	}
	defer loginResp.Body.Close()

	body, _ := io.ReadAll(loginResp.Body)
	fmt.Println("RAW LOGIN RESPONSE:", string(body))

	var loginResult struct {
		Login struct {
			Result string `json:"result"`
		} `json:"login"`
	}
	if err := json.Unmarshal(body, &loginResult); err != nil {
		return fmt.Errorf("failed to parse login result: %w", err)
	}
	if loginResult.Login.Result == "Success" {
		w.isAuthed = true
		return nil
	}
	return
}

func (w WikiProvider) Print(results []models.SearchResult) {
	var content string

	for _, res := range results {
		// Render the inner box for each result
		innerBox := utils.ResultStyle.Render(fmt.Sprintf("%s\n%s", res.Title, res.URL))
		content += innerBox + "\n"
	}

	// Wrap all results in the Provider's outer box
	fmt.Println(utils.ProviderStyle.Render(content))
}

func (w WikiProvider) Name() string { return "Wiki" }

func (w WikiProvider) Search(client *http.Client, query string, token string) ([]models.SearchResult, error) {
	// Actual api call here
	err := w.AuthBot(client)
	if err != nil {
		return nil, err
	}

	u, _ := url.Parse("https://wiki.csh.rit.edu/api.php")
	params := u.Query()
	params.Set("action", "query")
	params.Set("assert", "user")
	params.Set("list", "search")
	params.Set("srsearch", query)
	params.Set("srnamespace", "*")
	params.Set("srwhat", "text")
	params.Set("format", "json")

	u.RawQuery = params.Encode()
	fmt.Println(u.String())

	req, _ := http.NewRequest("GET", u.String(), nil)
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {

		// check if error.code == readapidenied or notloggedin
		// then reauthenticate and try again
		return nil, err
	}
	defer resp.Body.Close()

	type wikiItem struct {
		NS        int    `json:"ns"`
		Title     string `json:"title"`
		PageId    int    `json:"pageid"`
		Size      int    `json:"size"`
		WordCount int    `json:"wordcount"`
	}

	type wikiError struct {
		Code string `json:"code"`
		Info string `json:"info"`
	}

	type wikiResponse struct {
		// Errors is a list in MediaWiki API
		Errors []wikiError `json:"errors"`
		Query  struct {
			Search []wikiItem `json:"search"`
		} `json:"query"`
	}

	body, _ := io.ReadAll(resp.Body)

	var data wikiResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("could not parse JSON: %w", err)
	}

	if len(data.Errors) > 0 {
		errCode := data.Errors[0].Code
		if errCode == "readapidenied" || errCode == "notloggedin" {
			fmt.Println("Session expired, re-authenticating Wiki bot...")

			// Re-run your login logic
			err := w.AuthBot(client)
			if err != nil {
				return nil, fmt.Errorf("re-auth failed: %w", err)
			}

			// 3. RECURSION: Try the search again now that we are logged in
			return w.Search(client, query, token)
		}

		return nil, fmt.Errorf("Wiki API error: %s - %s", errCode, data.Errors[0].Info)
	}

	var finalResults []models.SearchResult
	for _, result := range data.Query.Search {
		finalResults = append(finalResults,
			models.SearchResult{
				Title:       result.Title,
				Description: fmt.Sprintf("%d count", result.WordCount),
				URL:         fmt.Sprintf("https://wiki.csh.rit.edu/index.php?curid=%d", result.PageId),
				Source:      w.Name(),
			},
		)
	}
	return finalResults, nil
}
