package migrator

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/egor-zakharov/go-musthave-diploma-tpl/db"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

type migrator struct {
	db *sql.DB
}

func New(db *sql.DB) Migrator {
	return &migrator{db: db}
}

func (m *migrator) Run() error {
	driver, err := iofs.New(db.MigrationsFS, "migrations")
	if err != nil {
		return err
	}

	d, err := postgres.WithInstance(m.db, &postgres.Config{})
	if err != nil {
		return err
	}

	instance, err := migrate.NewWithInstance("iofs", driver, "psql_db", d)
	if err != nil {
		return err
	}

	defer func() {
		_, _ = instance.Close()
	}()

	err = instance.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("unable to apply migrations %v", err)
	}

	return nil
}
