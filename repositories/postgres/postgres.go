package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/GlebMoskalev/go-event-bot/config"
	"github.com/GlebMoskalev/go-event-bot/repositories"
	_ "github.com/lib/pq"
	"log/slog"
)

type postgres struct {
	db  *sql.DB
	log *slog.Logger
}

func dsn(cfg config.Postgres) string {
	return fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name,
	)
}

func New(ctx context.Context, cfg config.Postgres, logger *slog.Logger) (repositories.DB, error) {
	db, err := sql.Open("postgres", dsn(cfg))

	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping the database")
	}
	return &postgres{db: db, log: logger}, err
}

func (p *postgres) Close() error {
	err := p.db.Close()
	if err != nil {
		return fmt.Errorf("can not close database: %v", err)
	}
	return nil
}
