package postgres

import (
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log/slog"
	"support/internal/config"
)

func InitDB() (*sqlx.DB, error) {
	dbInfo := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.AppConfig.DB.Host,
		config.AppConfig.DB.Port,
		config.AppConfig.DB.User,
		config.AppConfig.DB.Pass,
		config.AppConfig.DB.Name,
	)
	db, err := sqlx.Connect("postgres", dbInfo)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func MakeMigrations(up bool) error {
	dbLine := fmt.Sprintf("postgres://%s:%s@db:%s/%s?sslmode=disable",
		config.AppConfig.DB.User,
		config.AppConfig.DB.Pass,
		config.AppConfig.DB.Port,
		config.AppConfig.DB.Name,
	)
	m, err := migrate.New("file://migrations", dbLine)
	if err != nil {
		return err
	}

	if up {
		err = m.Up()
		if err != nil {
			if errors.Is(err, migrate.ErrNoChange) {
				slog.Info("No new migrations to apply")
				return nil
			}
			return err
		}
	} else {
		err = m.Down()
		if err != nil {
			return err
		}
	}
	return nil
}
