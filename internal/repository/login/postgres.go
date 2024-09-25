package login

import (
	"context"
	"encoding/json"
	"log"

	"github.com/lib/pq"
	"github.com/pedrokunz/canoe_reddit_integration/internal/domain"
	"github.com/pedrokunz/canoe_reddit_integration/internal/repository/postgres"
)

type postgresRepository struct{}

func Postgres() Repository {
	return postgresRepository{}
}

func (p postgresRepository) UserRead(ctx context.Context, input LoginRepositoryUserReadInput) (LoginRepositoryUserReadOutput, error) {
	db, err := postgres.Connect()
	if err != nil {
		return LoginRepositoryUserReadOutput{}, err
	}
	defer db.Close()

	query := `
		SELECT 
			users.id, 
			users.email, 
			users.password, 
			users.roles, 
			users.attributes, 
			users.customer_id
		FROM users
		INNER JOIN customers ON users.customer_id = customers.id
		WHERE users.email = $1 
		AND customers.id = $2
		AND users.deleted_at IS NULL
		AND customers.deleted_at IS NULL
		LIMIT 1`

	args := []interface{}{input.Email, input.CustomerID}
	row := db.QueryRowContext(ctx, query, args...)

	log.Printf("Query: %s, args: %v", query, args)

	var user domain.User
	var attributes []byte
	var roles pq.StringArray

	err = row.Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&roles,
		&attributes,
		&user.Customer.ID,
	)
	if err != nil {
		return LoginRepositoryUserReadOutput{}, err
	}

	err = json.Unmarshal(attributes, &user.Attributes)
	if err != nil {
		return LoginRepositoryUserReadOutput{}, err
	}

	user.Roles = []string(roles)

	return LoginRepositoryUserReadOutput{
		User: user,
	}, nil
}
