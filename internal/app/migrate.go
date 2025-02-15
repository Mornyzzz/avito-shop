//go:build migrate

package app

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/golang-migrate/migrate/v4"
	// migrate tools
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const (
	_defaultAttempts = 20
	_defaultTimeout  = time.Second
)

func init() {
	const op = "app.migrate"
	databaseURL, ok := os.LookupEnv("PG_URL")
	if !ok || len(databaseURL) == 0 {
		log.Fatalf(op, "environment variable not declared: PG_URL")
	}

	databaseURL += "?sslmode=disable"

	var (
		attempts = _defaultAttempts
		err      error
		m        *migrate.Migrate
	)

	for attempts > 0 {
		m, err = migrate.New("file://migrations", databaseURL)
		if err == nil {
			break
		}

		log.Printf(op, "postgres is trying to connect, attempts left: %d", attempts)
		time.Sleep(_defaultTimeout)
		attempts--
	}

	if err != nil {
		log.Fatalf("%w: %s", op, err)
	}

	err = m.Up()
	defer m.Close()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("%w: %s", op, err)
	}

	if errors.Is(err, migrate.ErrNoChange) {
		log.Printf(op, "no change")
		return
	}

	log.Printf(op, "up success")
}
