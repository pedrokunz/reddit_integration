package main

import (
	"github.com/joho/godotenv"
	"github.com/pedrokunz/canoe_reddit_integration/api"
	"github.com/pedrokunz/canoe_reddit_integration/internal/repository"
	"github.com/pedrokunz/canoe_reddit_integration/internal/repository/login"
	"github.com/pedrokunz/canoe_reddit_integration/internal/repository/postgres/migration"
	redditpostsync "github.com/pedrokunz/canoe_reddit_integration/internal/repository/reddit_post_sync"
)

func init() {
	//load .env file
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	// Run migrations
	err = migration.Run()
	if err != nil {
		panic(err)
	}
}

func main() {
	repositories := repository.Repositories{
		Login:          login.Postgres(),
		RedditPostSync: redditpostsync.Postgres(),
	}

	api.New(repositories)
}
