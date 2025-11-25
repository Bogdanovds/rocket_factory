package migrator

import (
	"database/sql"
	"embed"
	"fmt"
	"log"

	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

// Migrator отвечает за применение миграций к базе данных
type Migrator struct {
	db *sql.DB
}

// New создаёт новый экземпляр Migrator
func New(db *sql.DB) *Migrator {
	return &Migrator{db: db}
}

// Up применяет все pending миграции
func (m *Migrator) Up(migrationsDir string) error {
	goose.SetBaseFS(nil)

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set dialect: %w", err)
	}

	if err := goose.Up(m.db, migrationsDir); err != nil {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	log.Println("✅ Migrations applied successfully")
	return nil
}

// UpEmbed применяет миграции из встроенных файлов
func (m *Migrator) UpEmbed() error {
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set dialect: %w", err)
	}

	if err := goose.Up(m.db, "migrations"); err != nil {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	log.Println("✅ Migrations applied successfully")
	return nil
}

// Down откатывает последнюю миграцию
func (m *Migrator) Down(migrationsDir string) error {
	goose.SetBaseFS(nil)

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set dialect: %w", err)
	}

	if err := goose.Down(m.db, migrationsDir); err != nil {
		return fmt.Errorf("failed to rollback migration: %w", err)
	}

	log.Println("✅ Migration rolled back successfully")
	return nil
}

// Status выводит статус миграций
func (m *Migrator) Status(migrationsDir string) error {
	goose.SetBaseFS(nil)

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set dialect: %w", err)
	}

	if err := goose.Status(m.db, migrationsDir); err != nil {
		return fmt.Errorf("failed to get migration status: %w", err)
	}

	return nil
}

