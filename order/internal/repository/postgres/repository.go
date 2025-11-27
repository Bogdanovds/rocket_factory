package postgres

import (
	"database/sql"

	_ "github.com/lib/pq"
)

// Repository реализует интерфейс repository.Repository для PostgreSQL
type Repository struct {
	db *sql.DB
}

// NewRepository создаёт новый PostgreSQL репозиторий
func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}
