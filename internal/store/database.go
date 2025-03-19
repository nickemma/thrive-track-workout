package store

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"io/fs"
)

func Open() (*sql.DB, error) {
	conn, err := sql.Open("pgx", "host=localhost user=root password=postgres dbname=postgres port=5432 sslmode=disable")

	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	fmt.Println("Connected to database")
	return conn, nil
}

func MigrateFs(db *sql.DB, migrationFs fs.FS, dir string) error {

}

func Migrate(db *sql.DB, dir string) error {
	err := goose.SetDialect("postgres")

	if err != nil {
		return fmt.Errorf("error setting postgres dialect: %w", err)
	}

	err = goose.Up(db, dir)
	if err != nil {
		return fmt.Errorf("error performing Up migrations: %w", err)
	}
	return nil
}
