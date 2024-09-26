package redditpostsync

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/vartanbeno/go-reddit/v2/reddit"

	"github.com/pedrokunz/canoe_reddit_integration/internal/domain"
)

type Service struct {
	Client *reddit.Client
}

type Input struct {
	Subreddit          string
	LatestPostIDSynced string
}

type Output struct {
	Posts []domain.Post
}

func New() Service {
	credentials := reddit.Credentials{
		ID:       os.Getenv("REDDIT_CLIENT_ID"),
		Secret:   os.Getenv("REDDIT_CLIENT_SECRET"),
		Username: os.Getenv("REDDIT_USERNAME"),
		Password: os.Getenv("REDDIT_PASSWORD"),
	}
	client, _ := reddit.NewClient(credentials)

	return Service{
		Client: client,
	}
}

func (s Service) FetchRedditPosts(ctx context.Context, input Input) (*Output, error) {
	return s.fetchRedditPostsRecursive(ctx, input.Subreddit, "", input.LatestPostIDSynced, 100)
}

func (s Service) fetchRedditPostsRecursive(
	ctx context.Context,
	subreddit, after, latestPostIDSynced string,
	limit int,
) (*Output, error) {
	options := &reddit.ListOptions{
		After: after,
		Limit: limit,
	}

	newPosts, resp, err := s.Client.Subreddit.NewPosts(ctx, subreddit, options)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("reddit API returned status code %d", resp.StatusCode)
	}

	// Map to domain.Post
	var posts []domain.Post
	for _, post := range newPosts {
		if post.FullID == latestPostIDSynced {
			return &Output{
				Posts: posts,
			}, nil
		}

		posts = append(posts, domain.Post{
			Title:    post.Title,
			Author:   post.Author,
			Origin:   domain.Reddit,
			SyncedAt: post.Created.UTC(),
			Metadata: map[string]any{
				"reddit_id":  post.FullID,
				"reddit_url": post.URL,
				"subreddit":  subreddit,
			},
		})
	}

	log.Printf("Fetched %d posts for subreddit %s", len(posts), subreddit)

	// If there is an after value, fetch the next page
	if resp.After != "" {
		nextOutput, err := s.fetchRedditPostsRecursive(
			ctx,
			subreddit,
			resp.After,
			"",
			limit,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, nextOutput.Posts...)
	}

	return &Output{
		Posts: posts,
	}, nil
}
