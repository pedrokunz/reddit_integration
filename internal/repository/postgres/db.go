package postgres

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq" // Import the PostgreSQL driver
)

func Connect() (*sql.DB, error) {
	dsn := getDataSourceName()
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	// Verify the connection
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func Close(db *sql.DB) {
	db.Close()
}

func getDataSourceName() string {
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbname := os.Getenv("DB_NAME")
	sslmode := os.Getenv("DB_SSLMODE")

	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		user,
		password,
		host,
		port,
		dbname,
		sslmode,
	)
}
