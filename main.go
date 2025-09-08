package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// TokenCache holds the access token and its expiry
type TokenCache struct {
	AccessToken string `json:"access_token"`
	Expiry      int64  `json:"expiry"`
}

type Post struct {
	Title     string `json:"title"`
	Flair     string `json:"link_flair_text"`
	Permalink string `json:"permalink"`
	Ups       int    `json:"ups"`
	IsVideo   bool   `json:"is_video"`
}

func main() {
	// Load .env file
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file:", err)
		return
	}

	accessToken := getOrCacheAccessToken()
	// fmt.Println("Access Token:", accessToken)

	// create an object that holds a list of subreddit and search_query pairs
	queriesFile, err := os.Open("assets/subreddit_queries.json")
	if err != nil {
		fmt.Println("Error opening queries file:", err)
		return
	}
	defer queriesFile.Close()

	var queries []struct {
		Subreddit   string `json:"subreddit"`
		SearchQuery string `json:"search_query"`
	}
	byteValue, _ := ioutil.ReadAll(queriesFile)
	json.Unmarshal(byteValue, &queries)

	var posts []Post

	for _, query := range queries {
		subreddit := query.Subreddit
		flair := strings.TrimPrefix(query.SearchQuery, "flair_name:")
		flair = strings.Trim(flair, "\"")
		fmt.Printf("Fetching posts from r/%s with flair '%s'\n", subreddit, flair)
		posts = append(posts, getPostsWithFlair(subreddit, flair, accessToken)...)
	}

	for _, post := range posts {
		// separate each of the posts on to a newline
		fmt.Printf("\033]8;;https://www.reddit.com%s\033\\%s\033]8;;\033\\\n", post.Permalink, post.Title)
	}
}

// getPostsWithFlair fetches recent posts from a subreddit with a specific flair
func getPostsWithFlair(subreddit, flair, accessToken string) []Post {

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

	type Listing struct {
		Data struct {
			Children []struct {
				Data Post `json:"data"`
			} `json:"children"`
		} `json:"data"`
	}

	var listing Listing
	// fmt.Println("Response:", string(body))
	err = json.Unmarshal(body, &listing)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return nil
	}

	var results []Post
	for _, child := range listing.Data.Children {
		if child.Data.Flair == flair && child.Data.IsVideo {
			results = append(results, Post{
				Title:     child.Data.Title,
				Flair:     child.Data.Flair,
				Permalink: child.Data.Permalink,
				Ups:       child.Data.Ups,
				IsVideo:   child.Data.IsVideo,
			})
		}
	}
	// sort the results by the number of ups
	sort.Slice(results, func(i, j int) bool {
		return results[i].Ups > results[j].Ups
	})
	return results
}

// getOrCacheAccessToken checks for a valid cached token, otherwise fetches and caches a new one
func getOrCacheAccessToken() string {
	const cacheFile = "token_cache.json"
	// Check for cached token
	if _, err := os.Stat(cacheFile); err == nil {
		data, err := ioutil.ReadFile(cacheFile)
		if err == nil {
			var cache TokenCache
			if json.Unmarshal(data, &cache) == nil {
				if cache.Expiry > time.Now().Unix() {
					return cache.AccessToken
				}
			}
		}
	}
	// No valid token, fetch new
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

	// Cache token with expiry
	cache := TokenCache{
		AccessToken: result.AccessToken,
		Expiry:      time.Now().Unix() + int64(result.ExpiresIn) - 30, // 30s buffer
	}
	cacheData, _ := json.Marshal(cache)
	ioutil.WriteFile(cacheFile, cacheData, 0644)
	return result.AccessToken
}
