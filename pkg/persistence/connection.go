package persistence

import (
	"automsg/pkg/config"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func NewConnection(config config.RootConfig) (*sql.DB, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Database.Host, config.Database.Port, config.Database.User,
		config.Database.Password, config.Database.Name, config.Database.SSLMode,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to the database: %v", err)
	}

	return db, nil
}
