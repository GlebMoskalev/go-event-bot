package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/GlebMoskalev/go-event-bot/configs"
	"github.com/GlebMoskalev/go-event-bot/pkg/logger"
	"github.com/GlebMoskalev/go-event-bot/repositories"
	_ "github.com/lib/pq"
	"log/slog"
)

type postgres struct {
	dbBot   *sql.DB
	dbStaff *sql.DB
	log     *slog.Logger
}

func dsn(user, password, host, port, name string) string {
	return fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable",
		user, password, host, port, name,
	)
}

func New(ctx context.Context, botCfg config.BotPostgres, staffCfg config.StaffPostgres, log *slog.Logger) (repositories.DB, error) {
	log = logger.SetupLogger(log, "repository", "New")
	log.Info("initializing database connections")

	botDb, err := sql.Open("postgres", dsn(botCfg.User, botCfg.Password, botCfg.Host, botCfg.Port, botCfg.Name))

	if err != nil {
		log.Error("failed to open bot database", "error", err)
		return nil, fmt.Errorf("failed to open bot database: %v", err)
	}

	if err = botDb.Ping(); err != nil {
		log.Error("failed to ping bot database", "error", err)
		return nil, fmt.Errorf("failed to ping the bot database: %v", err)
	}

	staffDb, err := sql.Open("postgres", dsn(staffCfg.User, staffCfg.Password, staffCfg.Host, staffCfg.Port, staffCfg.Name))
	if err != nil {
		log.Error("failed to open staff database", "error", err)
		return nil, fmt.Errorf("failed to open staff database: %v", err)
	}

	if err = staffDb.Ping(); err != nil {
		log.Error("failed to ping staff database", "error", err)
		return nil, fmt.Errorf("failed to ping the staff database: %v", err)
	}

	log.Info("database connections established successfully")
	return &postgres{dbBot: botDb, dbStaff: staffDb, log: log}, err
}

func (p *postgres) Close() error {
	log := logger.SetupLogger(p.log, "repository", "Close")
	log.Info("closing database connections")

	err := p.dbBot.Close()
	if err != nil {
		log.Error("failed to close bot database", "error", err)
		return fmt.Errorf("can not close database: %v", err)
	}
	err = p.dbStaff.Close()
	if err != nil {
		log.Error("failed to close staff database", "error", err)
		return fmt.Errorf("can not close database: %v", err)
	}

	log.Info("database connections closed successfully")
	return nil
}
