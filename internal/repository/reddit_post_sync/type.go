package redditpostsync

import "github.com/pedrokunz/canoe_reddit_integration/internal/domain"

type PostReadLatestSyncedIDInput struct {
	Subreddit string
}

type PostReadLatestSyncedIDOutput struct {
	LatestPostIDSynced string
}

type PostCreateInput struct {
	Posts []domain.Post
}

type PostCreateOutput struct {
	PostAmountCreated int
}

type Repository interface {
	PostReadLatestSyncedID(input PostReadLatestSyncedIDInput) (PostReadLatestSyncedIDOutput, error)
	PostCreate(input PostCreateInput) (PostCreateOutput, error)
}
