package login

import (
	"context"

	"github.com/pedrokunz/canoe_reddit_integration/internal/domain"
)

type LoginRepositoryUserReadInput struct {
	Email      string
	CustomerID string
}

type LoginRepositoryUserReadOutput struct {
	User domain.User
}

type Repository interface {
	UserRead(ctx context.Context, input LoginRepositoryUserReadInput,
	) (LoginRepositoryUserReadOutput, error)
}
