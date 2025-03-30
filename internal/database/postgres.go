package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

func InitDatabase(host, port, user, name, password string) (*sql.DB, error) {
	connStr := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable",
		user, password, host, port, name,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("error pinging database: %v", err)
	}
	return db, err
}
