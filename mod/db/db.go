package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var instance *sql.DB

func GetInstance() *sql.DB {
	if instance == nil {
		instance, _ = sql.Open("sqlite3", "./sqlite-database.db")
	}
	return instance
}
