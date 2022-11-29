package main;

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

var (
	dataBase *sql.DB
)

func initializeDb() (e error) {
	var err error
	
	dataBase, err = sql.Open("sqlite3", "./db/Users.db?cache=shared&mode=rwc"); if err != nil {
		return err
	}

	return nil
}