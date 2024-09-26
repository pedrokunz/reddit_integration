package redditpostsync

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/pedrokunz/canoe_reddit_integration/internal/repository/postgres"
)

type postgresRepository struct{}

func Postgres() Repository {
	return postgresRepository{}
}

func (p postgresRepository) PostReadLatestSyncedID(input PostReadLatestSyncedIDInput) (PostReadLatestSyncedIDOutput, error) {
	db, err := postgres.Connect()
	if err != nil {
		return PostReadLatestSyncedIDOutput{}, err
	}
	defer db.Close()

	var latestPostIDSynced string
	err = db.QueryRow(`
		SELECT posts.metadata->>'reddit_id' AS latest_post_id_synced
		FROM posts
		WHERE posts.origin = 'reddit' AND posts.metadata->>'subreddit' = $1
		ORDER BY posts.synced_at DESC
		LIMIT 1
	`, input.Subreddit).Scan(&latestPostIDSynced)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return PostReadLatestSyncedIDOutput{}, err
	}

	return PostReadLatestSyncedIDOutput{
		LatestPostIDSynced: latestPostIDSynced,
	}, nil
}

func (p postgresRepository) PostCreate(input PostCreateInput) (PostCreateOutput, error) {
	db, err := postgres.Connect()
	if err != nil {
		return PostCreateOutput{}, err
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return PostCreateOutput{}, err
	}
	defer tx.Rollback()

	// Prepare bulk insert for posts table
	postValueStrings := make([]string, 0, len(input.Posts))
	postValueArgs := make([]interface{}, 0, len(input.Posts)*5)

	for i, post := range input.Posts {
		metadataJSON, err := json.Marshal(post.Metadata)
		if err != nil {
			return PostCreateOutput{}, err
		}

		postValueStrings = append(postValueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d)", i*5+1, i*5+2, i*5+3, i*5+4, i*5+5))
		postValueArgs = append(postValueArgs, post.Title, post.Author, post.Origin, post.SyncedAt, metadataJSON)
	}

	// Insert into posts table
	postStmt := fmt.Sprintf(`
		INSERT INTO posts (title, author, origin, synced_at, metadata)
		VALUES %s
		ON CONFLICT (metadata) DO NOTHING
	`, strings.Join(postValueStrings, ","))
	result, err := tx.Exec(postStmt, postValueArgs...)
	if err != nil {
		return PostCreateOutput{}, err
	}

	insertedCount, err := result.RowsAffected()
	if err != nil {
		return PostCreateOutput{}, err
	}

	err = tx.Commit()
	if err != nil {
		return PostCreateOutput{}, err
	}

	return PostCreateOutput{
		PostAmountCreated: int(insertedCount),
	}, nil
}
