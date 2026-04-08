package postgres

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func Migrate(postgresDSN, migrationsPath string) error {
	m, err := migrate.New("file://"+migrationsPath, postgresDSN)
	if err != nil {
		return fmt.Errorf("Could not connect to postgres database: %s", err)

	}
	defer func() { _, _ = m.Close() }()
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("Could not run migrations: %s", err)
	}
	return nil
}
