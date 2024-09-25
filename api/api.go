package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pedrokunz/canoe_reddit_integration/api/controller/login"
	redditpostreadlist "github.com/pedrokunz/canoe_reddit_integration/api/controller/reddit_post_read_list"
	redditpostsync "github.com/pedrokunz/canoe_reddit_integration/api/controller/reddit_post_sync"
	"github.com/pedrokunz/canoe_reddit_integration/api/middleware"
	"github.com/pedrokunz/canoe_reddit_integration/internal/repository"
)

func New(repositories repository.Repositories) {
	r := gin.Default()

	// Middleware
	r.Use(middleware.ErrorHandler())
	r.Use(middleware.RateLimiter())

	// TODO: get customers from database
	customers := []string{"afacf9fa-7516-468f-b048-ac4c0562aa3f"}

	r.GET("/ping", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"message": "pong"}) })
	r.POST("/login", login.Execute(repositories.Login))

	r.Use(middleware.Auth())

	// Routes
	v1 := r.Group("/v1")

	v1.POST(
		"/reddit/posts/sync",
		middleware.AuthorizeCustomers(customers...),
		middleware.AuthorizeRoles("admin"),
		middleware.AuthorizeAttribute("reddit.post.sync", "true"),
		redditpostsync.Execute(repositories.RedditPostSync),
	)

	v1.GET("/reddit/posts",
		middleware.AuthorizeCustomers(customers...),
		middleware.AuthorizeRoles("admin", "member", "reader"),
		middleware.AuthorizeAttribute("reddit.post.read_list", "true"),
		redditpostreadlist.Execute(repositories.RedditPostReadList))

	r.Run(":8080")
}
