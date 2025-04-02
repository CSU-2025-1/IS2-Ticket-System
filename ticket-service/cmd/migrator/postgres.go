package migrator

import (
	"context"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/pressly/goose"
)

func MigratePostgres(ctx context.Context, conStr string) error {
	db, err := goose.OpenDBWithDriver("postgres", conStr)
	if err != nil {
		return fmt.Errorf("migrator: %w", err)
	}
	defer db.Close()

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("goose dialect: %w", err)
	}

	migrationsDir := "./migrations"
	if err := goose.Up(db, migrationsDir); err != nil {
		return fmt.Errorf("goose up: %w", err)
	}

	fmt.Println("Migrations applied successfully!")
	return nil
}
