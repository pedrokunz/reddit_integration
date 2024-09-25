package redditpostreadlist

import (
	"net/http"

	"github.com/gin-gonic/gin"
	redditpostreadlist "github.com/pedrokunz/canoe_reddit_integration/internal/repository/reddit_post_read_list"
)

func Execute(repository redditpostreadlist.Repository) func(c *gin.Context) {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "execute posts"})
	}
}
