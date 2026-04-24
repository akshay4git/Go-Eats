package unsplash

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
)

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func GetUnSplashImageURL(client HttpClient, menuItem string) string {
	imageUrl := "https://api.unsplash.com/search/photos/?page=1&query=" + url.QueryEscape(menuItem) + "&w=400&h=400"
	req, err := http.NewRequest("GET", imageUrl, nil)
	if err != nil {
		log.Printf("Failed to create request: %v", err)
		return ""
	}

	req.Header.Set("Authorization", fmt.Sprintf("%s %s", "Client-ID", os.Getenv("UNSPLASH_API_KEY")))
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to send request: %v", err)
		return ""
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response body: %v", err)
		return ""
	}

	var apiResponse UnSplash
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		log.Printf("UnSplash::Failed to decode JSON response: %v", err)
		return ""
	}

	// Guard against empty results
	if len(apiResponse.Results) == 0 {
		log.Printf("UnSplash::No results found for: %s", menuItem)
		return ""
	}

	return apiResponse.Results[0].Urls.Small
}