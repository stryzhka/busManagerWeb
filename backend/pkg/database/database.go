package database

import (
	"database/sql"
	"fmt"
)

func NewPostgresDatabase(connStr string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	fmt.Println("success connection")
	return db, err
}
