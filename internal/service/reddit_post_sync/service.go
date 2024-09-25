package redditpostsync

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/pedrokunz/canoe_reddit_integration/internal/domain"
)

type Service struct{}

type Input struct {
	Subreddit string
}

type Output struct {
	Posts []domain.Post
}

func New() Service {
	return Service{}
}

func (s Service) FetchRedditPosts(ctx context.Context, input Input) (*Output, error) {
	url := fmt.Sprintf("https://www.reddit.com/r/%s.json", input.Subreddit)

	// Create HTTP client
	client := &http.Client{}

	// Make request
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return nil, err
	}

	// Add headers
	req.Header.Set("User-Agent", "CanoeRedditIntegration:1.0 (by /u/Stunning_Commission4)")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error making request: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return nil, err
	}

	// Check response status
	if resp.StatusCode != http.StatusOK {
		log.Println(string(body))
		return nil, fmt.Errorf("failed to fetch posts: %s", resp.Status)
	}

	// Parse JSON response
	var redditResponse struct {
		Data struct {
			Children []struct {
				Data domain.Post `json:"data"`
			} `json:"children"`
		} `json:"data"`
	}

	err = json.Unmarshal(body, &redditResponse)
	if err != nil {
		log.Printf("Error parsing JSON response: %v", err)
		return nil, err
	}

	// Map to domain.Post
	var posts []domain.Post
	for _, child := range redditResponse.Data.Children {
		posts = append(posts, child.Data)
	}

	log.Printf("Fetched %d posts for subreddit %s", len(posts), input.Subreddit)

	return &Output{
		Posts: posts,
	}, nil
}
