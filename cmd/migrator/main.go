package main

import (
	"errors"
	"flag"
	"fmt"

	"github.com/golang-migrate/migrate/v4"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	var storagePath, migrationsPath, migrationsTable, direction string

	flag.StringVar(&storagePath, "storage-path", "", "path to storage")
	flag.StringVar(&migrationsPath, "migrations-path", "", "path to migrations")
	flag.StringVar(&migrationsTable, "migrations-table", "migrations", "name of migrations table")
	flag.StringVar(&direction, "direction", "", "up or down migrations")
	flag.Parse()

	if storagePath == "" {
		panic("storage path is requires")
	}
	if migrationsPath == "" {
		panic("migrations path is required")
	}

	// postgres://postgres:postgres@localhost:5432/grpc_auth?sslmode=disable?x-migrations-table=migrationsTable
	// storagePath: --storage-path=postgres:postgres@localhost:5432/grpc_auth?sslmode=disable
	m, err := migrate.New(
		"file://"+migrationsPath,
		fmt.Sprintf("postgres://%s&x-migrations-table=%s", storagePath, migrationsTable),
	)
	if err != nil {
		panic(err)
	}

	switch direction {
	case "up":
		if err := m.Up(); err != nil {
			if errors.Is(err, migrate.ErrNoChange) {
				fmt.Println("no migrations to apply")

				return
			}
			panic(err)
		}
		fmt.Println("migrations applied successfully")
	case "down":
		if err := m.Down(); err != nil {
			if errors.Is(err, migrate.ErrNoChange) {
				fmt.Println("no migrations to apply")

				return
			}
			panic(err)
		}
		fmt.Println("migrations canceled successfully")
	}

}
