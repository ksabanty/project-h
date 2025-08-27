package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file:", err)
		return
	}
	// get reddit access token
	accessToken := getAccessToken()
	fmt.Println("Access Token:", accessToken)

	subreddit := "soccer" // replace with your subreddit
	flair := "Media"      // replace with your flair
	posts := getPostsWithFlair(subreddit, flair, accessToken)
	fmt.Println("Posts:")
	for _, post := range posts {
		fmt.Println(post)
	}
}

// getPostsWithFlair fetches recent posts from a subreddit with a specific flair
func getPostsWithFlair(subreddit, flair, accessToken string) []string {
	userAgent := os.Getenv("USER_AGENT")
	token := strings.TrimSpace(accessToken)
	url := fmt.Sprintf("https://oauth.reddit.com/r/%s/new.json?q=flair_name%%3A\"%s\"&restrict_sr=on&limit=20", subreddit, flair)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil
	}
	req.Header.Set("Authorization", "bearer "+token)
	req.Header.Set("User-Agent", userAgent)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error getting response:", err)
		return nil
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return nil
	}

	// Parse JSON and filter by flair
	type Post struct {
		Title     string `json:"title"`
		Flair     string `json:"link_flair_text"`
		Permalink string `json:"permalink"`
	}
	type Listing struct {
		Data struct {
			Children []struct {
				Data Post `json:"data"`
			} `json:"children"`
		} `json:"data"`
	}

	var listing Listing
	fmt.Println("Response:", string(body))
	err = json.Unmarshal(body, &listing)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return nil
	}

	var results []string
	for _, child := range listing.Data.Children {
		if child.Data.Flair == flair {
			results = append(results, fmt.Sprintf("%s (https://reddit.com%s)", child.Data.Title, child.Data.Permalink))
		}
	}
	return results
}

// write a function that gets the access token from reddit
// write a function that gets the access token from reddit
func getAccessToken() string {
	url := "https://www.reddit.com/api/v1/access_token"
	clientID := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")
	userAgent := os.Getenv("USER_AGENT")
	data := "grant_type=client_credentials"

	req, err := http.NewRequest("POST", url, bytes.NewBufferString(data))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return ""
	}

	req.SetBasicAuth(clientID, clientSecret)
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error getting response:", err)
		return ""
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return ""
	}

	// Parse JSON response
	var result struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int    `json:"expires_in"`
		Scope       string `json:"scope"`
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println("Error parsing token JSON:", err)
		return ""
	}

	return result.AccessToken
}
