package redditpostsync

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	redditpostsync "github.com/pedrokunz/canoe_reddit_integration/internal/repository/reddit_post_sync"
	service "github.com/pedrokunz/canoe_reddit_integration/internal/service/reddit_post_sync"
)

func Execute(repository redditpostsync.Repository) func(c *gin.Context) {
	return func(c *gin.Context) {
		queriedSubreddits := c.Query("subreddits")
		if queriedSubreddits == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		subreddits := strings.Split(queriedSubreddits, ",")
		if len(subreddits) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		for _, subreddit := range subreddits {
			go func(subreddit string) {
				postReadLatestSyncedIDOutput, err := repository.PostReadLatestSyncedID(redditpostsync.PostReadLatestSyncedIDInput{
					Subreddit: subreddit,
				})
				if err != nil {
					log.Printf("Error reading latest synced ID for subreddit %s: %s\n", subreddit, err)
					return
				}

				serviceOutput, err := service.New().FetchRedditPosts(
					context.Background(),
					service.Input{
						Subreddit:          subreddit,
						LatestPostIDSynced: postReadLatestSyncedIDOutput.LatestPostIDSynced,
					},
				)
				if err != nil {
					log.Printf("Error fetching posts for subreddit %s: %s\n", subreddit, err)
					return
				}

				if len(serviceOutput.Posts) == 0 {
					log.Printf("No new posts for subreddit %s\n", subreddit)
					return
				}

				repositoryOutput, err := repository.PostCreate(redditpostsync.PostCreateInput{
					Posts: serviceOutput.Posts,
				})
				if err != nil {
					log.Printf("Error saving posts for subreddit %s: %s\n", subreddit, err)
					return
				}

				log.Printf("Successfully saved %d posts for subreddit %s\n", repositoryOutput.PostAmountCreated, subreddit)
			}(subreddit)
		}

		c.JSON(http.StatusOK, gin.H{"message": "Reddit posts sync started"})
	}
}
