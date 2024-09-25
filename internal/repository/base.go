package repository

import (
	"github.com/pedrokunz/canoe_reddit_integration/internal/repository/login"
	redditpostreadlist "github.com/pedrokunz/canoe_reddit_integration/internal/repository/reddit_post_read_list"
	redditpostsync "github.com/pedrokunz/canoe_reddit_integration/internal/repository/reddit_post_sync"
)

type Repositories struct {
	Login              login.Repository
	RedditPostSync     redditpostsync.Repository
	RedditPostReadList redditpostreadlist.Repository
}
