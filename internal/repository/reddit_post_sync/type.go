package redditpostsync

import "github.com/pedrokunz/canoe_reddit_integration/internal/domain"

type PostCreateInput struct {
	Posts []domain.Post
}

type PostCreateOutput struct {
	Posts []domain.Post
}

type Repository interface {
	PostCreate(input PostCreateInput) (PostCreateOutput, error)
}
