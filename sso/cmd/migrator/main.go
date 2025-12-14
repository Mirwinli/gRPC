package main

import (
	"errors"
	"flag"
	"fmt"

	//Библиотека для миграций
	"github.com/golang-migrate/migrate/v4"
	//Драйвер для виконання миграцій sqLite
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	// Драйвер для получення міграцій з файлів
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	var storagePath, migrationPath, migrationsTable string

	flag.StringVar(&storagePath, "storage-path", "", "path to store storage")
	flag.StringVar(&migrationPath, "migration-path", "", "path to store migrations")
	flag.StringVar(&migrationsTable, "migrations-table", "migrations", "path to store migrations table")
	flag.Parse()

	if storagePath == "" {
		panic("storage-path is required")
	}
	if migrationPath == "" {
		panic("migration-path is required")
	}

	m, err := migrate.New("file://"+migrationPath,
		fmt.Sprintf("sqlite3://%s&x-migrations-table=%s", storagePath, migrationsTable),
	)
	if err != nil {
		panic(err)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("no change")

			return
		}
		panic(err)
	}
	fmt.Println("successfully apply migrations")
}
