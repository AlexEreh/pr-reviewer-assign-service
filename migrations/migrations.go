package migrations

import (
	"database/sql"
	"embed"
	"fmt"

	"github.com/pressly/goose/v3"

	"pr-reviewer-assign-service/pkg/app"
)

//go:embed sql/*.sql
var Res embed.FS

func MigrateUp(_ *app.App, db *sql.DB) error {
	goose.SetBaseFS(Res)

	err := goose.SetDialect("postgres")
	if err != nil {
		return fmt.Errorf("failed to set postgres dialect: %w", err)
	}

	err = goose.Up(db, "sql")
	if err != nil {
		return fmt.Errorf("failed to up migrations: %w", err)
	}
	return nil
}
